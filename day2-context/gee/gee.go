package gee

import "net/http"

type HandlerFunc func(ctx *Context)

// 引擎，是最重要的，是gee对外提供的接口
// 这里加一个router是把router分离出去了，这样比较解耦合
// engine的作用是进行包装，对get，post，等的包装
type Engine struct {
	router *router
}

func New() *Engine {
	return &Engine{router: newRouter()}
}
func (e *Engine) addRouter(method string, pattern string, handlerFunc HandlerFunc) {
	e.router.addRouter(method, pattern, handlerFunc)
}
func (e *Engine) Get(pattern string, handlerFunc HandlerFunc) {
	e.addRouter("GET", pattern, handlerFunc)
}
func (e *Engine) POST(pattern string, handlerFunc HandlerFunc) {
	e.addRouter("POST", pattern, handlerFunc)
}
func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := newContext(w, r)
	e.router.handle(c)
}
func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}
