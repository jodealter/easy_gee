package gee

import (
	"log"
	"net/http"
)

// router负责路由功能，添加router或者 把url中path对应的处理方式
// 目前只有添加路由和执行路由
type router struct {
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{handlers: map[string]HandlerFunc{}}
}
func (r *router) addRouter(method string, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - %s", method, pattern)
	key := method + "-" + pattern
	r.handlers[key] = handler
}
func (r *router) handle(c *Context) {
	key := c.Method + "-" + c.Path
	if handler, ok := r.handlers[key]; ok {
		handler(c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND :%s\n", c.Path)
	}
}
