package framework

import (
	"context"
	"html/template"
	"net/http"
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
	// 1. Static Files
	for prefix, handler := range r.staticHandlers {
		if strings.HasPrefix(req.URL.Path, prefix) {
			handler.ServeHTTP(w, req)
			return
		}
	}

	// 2. Trie Lookup
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

	// 3. Context
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
	_, err := r.templates.ParseGlob(pattern)
	if err != nil {
		panic(err)
	}
}
