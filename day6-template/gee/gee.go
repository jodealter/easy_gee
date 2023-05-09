package gee

import (
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"
	"time"
)

type HandlerFunc func(ctx *Context)

// RouterGroup 路由分组
type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc //support middleware
	parent      *RouterGroup  //support nesting
	engine      *Engine       //all group share an engine instance
}

// Engine 引擎，是最重要的，是gee对外提供的接口
// 这里加一个router是把router分离出去了，这样比较解耦合
// engine的作用是进行包装，对get，post，等的包装
type Engine struct {
	*RouterGroup
	router       *router
	groups       []*RouterGroup
	htmlTemplate *template.Template
	funcMap      template.FuncMap
}

func (e *Engine) SetFuncMap(funcmap template.FuncMap) {
	e.funcMap = funcmap
}
func (e *Engine) LoadHTMLGlob(pattern string) {
	e.htmlTemplate = template.Must(template.New("").Funcs(e.funcMap).ParseGlob(pattern))
}
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}
func (group *RouterGroup) addRouter(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Router + %4s - %s", method, pattern)
	group.engine.router.addRouter(method, pattern, handler)
}
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRouter("GET", pattern, handler)
}
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRouter("POST", pattern, handler)
}
func (e *Engine) addRouter(method string, pattern string, handlerFunc HandlerFunc) {
	e.router.addRouter(method, pattern, handlerFunc)
}
func (e *Engine) GET(pattern string, handlerFunc HandlerFunc) {
	e.addRouter("GET", pattern, handlerFunc)
}
func (e *Engine) POST(pattern string, handlerFunc HandlerFunc) {
	e.addRouter("POST", pattern, handlerFunc)
}
func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var middles []HandlerFunc
	for _, group := range e.groups {
		if strings.HasPrefix(r.URL.Path, group.prefix) {
			middles = append(middles, group.middlewares...)
		}
	}
	c := newContext(w, r)
	c.handlers = middles
	c.engine = e
	e.router.handle(c)
}
func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}
func Logger() HandlerFunc {
	return func(ctx *Context) {
		t := time.Now()
		ctx.Next()
		log.Printf("[%d] %s in %v", ctx.StatusCode, ctx.Req.RequestURI, time.Since(t))
	}
}
func (group *RouterGroup) Use(middle ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middle...)
}
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.prefix, relativePath)
	fileserver := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(ctx *Context) {
		file := ctx.Param("filepath")
		if _, err := fs.Open(file); err != nil {
			ctx.Status(http.StatusNotFound)
		}
		fileserver.ServeHTTP(ctx.Write, ctx.Req)
	}
}
func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	group.GET(urlPattern, handler)
}
