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
	Writer    http.ResponseWriter
	Request   *http.Request
	templates *template.Template
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	var tmpl *template.Template
	if t, ok := r.Context().Value("framework_templates").(*template.Template); ok {
		tmpl = t
	}
	return &Context{
		Writer:    w,
		Request:   r,
		templates: tmpl,
	}
}

func (context *Context) Status(code int) {
	context.Writer.WriteHeader(code)
}

func (context *Context) JSON(code int, v any) error {
	context.Writer.Header().Set("Content-Type", "application/json")
	context.Writer.WriteHeader(code)
	return json.NewEncoder(context.Writer).Encode(v)
}

func (context *Context) Param(key string) string {
	val, ok := context.Request.Context().Value(key).(string)
	if !ok {
		return ""
	}
	return val
}

func (context *Context) BindJSON(v any) error {
	return json.NewDecoder(context.Request.Body).Decode(v)
}

func (context *Context) BindJSONStrict(v any) error {
	contentType := context.Request.Header.Get("Content-Type")
	if contentType != "" {
		mediaType, _, err := mime.ParseMediaType(contentType)
		if err != nil || mediaType != "application/json" {
			return &HTTPError{Code: http.StatusUnsupportedMediaType, Message: "Content-Type must be application/json"}
		}
	} else {
		return &HTTPError{Code: http.StatusUnsupportedMediaType, Message: "Content-Type header is missing"}
	}

	context.Request.Body = http.MaxBytesReader(context.Writer, context.Request.Body, 1048576)
	dec := json.NewDecoder(context.Request.Body)
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

func (context *Context) Query(key string) string {
	return context.Request.URL.Query().Get(key)
}

func (context *Context) HTML(code int, name string, data any) {
	if context.templates == nil {
		http.Error(context.Writer, "Templates not loaded", http.StatusInternalServerError)
		return
	}
	var buf bytes.Buffer
	if err := context.templates.ExecuteTemplate(&buf, name, data); err != nil {
		log.Printf("Template error (%s): %v", name, err)
		http.Error(context.Writer, "Template Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	context.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	context.Writer.WriteHeader(code)
	buf.WriteTo(context.Writer)
}
