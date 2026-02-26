package framework

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"mime"
	"net/http"
	"strings"
)

type Context struct {
	W         http.ResponseWriter
	R         *http.Request
	templates *template.Template
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	var tmpl *template.Template
	if t, ok := r.Context().Value("framework_templates").(*template.Template); ok {
		tmpl = t
	}
	return &Context{
		W:         w,
		R:         r,
		templates: tmpl,
	}
}

func (c *Context) Status(code int) {
	c.W.WriteHeader(code)
}

func (c *Context) JSON(code int, v any) error {
	c.W.Header().Set("Content-Type", "application/json")
	c.W.WriteHeader(code)
	return json.NewEncoder(c.W).Encode(v)
}

func (c *Context) Param(key string) string {
	val, ok := c.R.Context().Value(key).(string)
	if !ok {
		return ""
	}
	return val
}

func (c *Context) BindJSON(v any) error {
	return json.NewDecoder(c.R.Body).Decode(v)
}

func (c *Context) BindJSONStrict(v any) error {
	ct := c.R.Header.Get("Content-Type")
	if ct != "" {
		mediaType, _, err := mime.ParseMediaType(ct)
		if err != nil || mediaType != "application/json" {
			return &HTTPError{Code: http.StatusUnsupportedMediaType, Message: "Content-Type must be application/json"}
		}
	} else {
		return &HTTPError{Code: http.StatusUnsupportedMediaType, Message: "Content-Type header is missing"}
	}

	c.R.Body = http.MaxBytesReader(c.W, c.R.Body, 1048576)
	dec := json.NewDecoder(c.R.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(v)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var maxBytesError *http.MaxBytesError

		msg := "Invalid JSON"
		code := http.StatusBadRequest

		switch {
		case errors.As(err, &syntaxError):
			msg = fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
		case errors.As(err, &unmarshalTypeError):
			msg = fmt.Sprintf("Invalid value for field %q (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
		case errors.As(err, &maxBytesError):
			msg = "Request body must not be larger than 1MB"
			code = http.StatusRequestEntityTooLarge
		case strings.HasPrefix(err.Error(), "json: unknown field"):
			msg = fmt.Sprintf("Unknown field in JSON body: %s", strings.TrimPrefix(err.Error(), "json: unknown field "))
		case errors.Is(err, io.EOF):
			msg = "Request body must not be empty"
		default:
			msg = err.Error()
		}
		return &HTTPError{Code: code, Message: msg}
	}
	return nil
}

type HTTPError struct {
	Code    int
	Message string
}

func (e *HTTPError) Error() string {
	return e.Message
}

func (c *Context) Query(key string) string {
	return c.R.URL.Query().Get(key)
}

func (c *Context) HTML(code int, name string, data any) {
	if c.templates == nil {
		http.Error(c.W, "Templates not loaded", http.StatusInternalServerError)
		return
	}
	var buf bytes.Buffer
	if err := c.templates.ExecuteTemplate(&buf, name, data); err != nil {
		log.Printf("Template error (%s): %v", name, err)
		http.Error(c.W, "Template Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	c.W.Header().Set("Content-Type", "text/html; charset=utf-8")
	c.W.WriteHeader(code)
	buf.WriteTo(c.W)
}
