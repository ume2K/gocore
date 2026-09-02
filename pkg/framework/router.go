package framework

import (
	"context"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type HandlerFunc func(*Context)

type Middleware func(http.Handler) http.Handler

type Router struct {
	*RouteGroup
	roots          map[string]*node
	staticHandlers map[string]http.Handler
	templates      *template.Template
	funcMap        template.FuncMap
	notFound       http.Handler
}

func NewRouter() *Router {
	r := &Router{
		roots:          make(map[string]*node),
		staticHandlers: make(map[string]http.Handler),
		funcMap:        make(template.FuncMap),
	}

	r.RouteGroup = &RouteGroup{
		prefix:      "/",
		middlewares: []Middleware{},
		router:      r,
	}
	return r
}

func (r *Router) SetFuncMap(funcs template.FuncMap) {
	for name, fn := range funcs {
		r.funcMap[name] = fn
	}
}

func (r *Router) SetNotFound(handler http.Handler) {
	r.notFound = handler
}

func (r *Router) handle(method, path string, handler http.Handler) {
	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &node{}
	}

	r.roots[method].insert(method, path, handler)
}

func (r *Router) handleNotFound(w http.ResponseWriter, req *http.Request) {
	if r.notFound != nil {
		r.notFound.ServeHTTP(w, req)
	} else {
		http.NotFound(w, req)
	}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for prefix, handler := range r.staticHandlers {
		if strings.HasPrefix(req.URL.Path, prefix) {
			handler.ServeHTTP(w, req)
			return
		}
	}

	root, ok := r.roots[req.Method]
	if !ok {
		r.handleNotFound(w, req)
		return
	}

	pathParts := parsePath(req.URL.Path)
	node, params := root.search(pathParts)

	if node == nil || node.handler == nil {
		r.handleNotFound(w, req)
		return
	}

	ctx := req.Context()
	for k, v := range params {
		ctx = context.WithValue(ctx, k, v)
	}
	if r.templates != nil {
		ctx = context.WithValue(ctx, "framework_templates", r.templates)
	}

	handler := node.handler.(http.Handler)
	handler.ServeHTTP(w, req.WithContext(ctx))
}

func (r *Router) Static(urlPath, rootDir string) {
	if !strings.HasPrefix(urlPath, "/") {
		urlPath = "/" + urlPath
	}
	fs := http.StripPrefix(urlPath, http.FileServer(http.Dir(rootDir)))
	r.staticHandlers[urlPath] = fs
}

func (r *Router) LoadHTMLGlob(pattern string) {
	if r.templates == nil {
		r.templates = template.New("").Funcs(r.funcMap)
	}

	rootDir := filepath.Clean(pattern)
	rootDir = strings.TrimSuffix(rootDir, "*.html")
	rootDir = strings.TrimSuffix(rootDir, "*")
	rootDir = strings.TrimSuffix(rootDir, string(filepath.Separator))
	if rootDir == "" {
		rootDir = "."
	}

	err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && strings.HasSuffix(d.Name(), ".html") {
			relPath, err := filepath.Rel(rootDir, path)
			if err != nil {
				return err
			}

			templateName := filepath.ToSlash(relPath)
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			_, err = r.templates.New(templateName).Parse(string(content))
			if err != nil {
				return err
			}

			if templateName != d.Name() {
				_, err = r.templates.New(d.Name()).Parse(string(content))
				if err != nil {
					return err
				}
			}
		}
		return nil
	})

	if err != nil {
		panic(err)
	}
}
