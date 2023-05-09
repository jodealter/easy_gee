package gee

import (
	"fmt"
	"net/http"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request)
type Engine struct {
	router map[string]HandlerFunc
}

func New() *Engine {
	return &Engine{router: make(map[string]HandlerFunc)}
}
func (e *Engine) addRouter(method string, pattern string, handlerFunc HandlerFunc) {
	key := method + "-" + pattern
	e.router[key] = handlerFunc
}
func (e *Engine) Get(pattern string, handlerFunc HandlerFunc) {
	e.addRouter("GET", pattern, handlerFunc)
}
func (e *Engine) POST(pattern string, handlerFunc HandlerFunc) {
	e.addRouter("POST", pattern, handlerFunc)
}
func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}
func (e *Engine) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	key := request.Method + "-" + request.URL.Path
	if handler, ok := e.router[key]; ok {
		handler(writer, request)
	} else {
		fmt.Fprintf(writer, "404 NOT FOUND:%q\n", request.URL.Path)
	}
}
