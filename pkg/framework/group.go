package framework

import (
	"net/http"
	"path"
)

type RouteGroup struct {
	prefix      string
	middlewares []Middleware
	router      *Router
}

func (g *RouteGroup) Group(prefix string) *RouteGroup {
	mwCopy := make([]Middleware, len(g.middlewares))
	copy(mwCopy, g.middlewares)
	return &RouteGroup{
		prefix:      path.Join(g.prefix, prefix),
		middlewares: mwCopy,
		router:      g.router,
	}
}

func (g *RouteGroup) Use(mw Middleware) {
	g.middlewares = append(g.middlewares, mw)
}

func (g *RouteGroup) Add(method, pathStr string, handler HandlerFunc) {
	fullPath := path.Join(g.prefix, pathStr)
	var finalHandler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := NewContext(w, r)
		handler(c)
	})
	for i := len(g.middlewares) - 1; i >= 0; i-- {
		finalHandler = g.middlewares[i](finalHandler)
	}
	g.router.handle(method, fullPath, finalHandler)
}

func (g *RouteGroup) GET(path string, handler HandlerFunc)  { g.Add("GET", path, handler) }
func (g *RouteGroup) POST(path string, handler HandlerFunc) { g.Add("POST", path, handler) }
func (g *RouteGroup) PUT(path string, handler HandlerFunc)  { g.Add("PUT", path, handler) }
func (g *RouteGroup) DELETE(path string, handler HandlerFunc) { g.Add("DELETE", path, handler) }
