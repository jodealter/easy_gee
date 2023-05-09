package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}
type Context struct {
	Write      http.ResponseWriter
	Req        *http.Request
	Path       string
	Method     string
	StatusCode int
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Write:  w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
	}
}

// 下边两个差不多，都是获取url中的参数
func (c *Context) PostFrom(key string) string {
	return c.Req.FormValue(key)
}
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// 设置状态码
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Write.WriteHeader(code)
}
func (c *Context) SetHeader(key, value string) {
	c.Write.Header().Set(key, value)
}
func (c *Context) String(Code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(Code)
	c.Write.Write([]byte(fmt.Sprintf(format, values...)))
}
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Write)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Write, err.Error(), 500)
	}
}
func (c *Context) Data(code int, data string) {
	c.Status(code)
	c.Write.Write([]byte(data))
}
func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Write.Write([]byte(html))
}
