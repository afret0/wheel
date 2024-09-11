package router

import (
	"github.com/gin-gonic/gin"
	"sync"
)

var mu = new(sync.Mutex)

type Group struct {
	path   string
	router map[string]*item
	use    []gin.HandlerFunc
}

func (g *Group) Use(middleware ...gin.HandlerFunc) {
	g.use = append(g.use, middleware...)
}

func (g *Group) POST(relativePath string, handle HandleFuncWrap, middleware ...gin.HandlerFunc) {
	mu.Lock()
	defer mu.Unlock()
	g.router[relativePath] = &item{
		method:     "POST",
		handle:     handle,
		path:       relativePath,
		middleware: middleware,
	}
}

func (g *Group) GET(relativePath string, handle HandleFuncWrap, middleware ...gin.HandlerFunc) {
	mu.Lock()
	defer mu.Unlock()

	g.router[relativePath] = &item{
		method:     "GET",
		handle:     handle,
		path:       relativePath,
		middleware: middleware,
	}
}
