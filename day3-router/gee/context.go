package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// 比较实用的结构，构建html，json时使用
type H map[string]interface{}

// context 上下文，其实是对元语句的包装，不然每次都要写很小力度的元语句，会很费劲，字段包括了几个常用的
type Context struct {
	Write      http.ResponseWriter //响应
	Req        *http.Request       //请求
	Path       string              //路径，即请求资源的路径
	Method     string              //方法，get post还是什么的
	StatusCode int                 //状态码
	Params     map[string]string   //参数形如不确定的 /hello/:filepath/kds,冒号filepath是不确定的，可以匹配这个字段的所有
}

// 创建新的上下文
func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Write:  w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
	}
}

// 返回key对应的参数，即上边filepath对应的  比如/hello/go/kds  中的go
func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}
func (c *Context) PostFrom(key string) string {
	return c.Req.FormValue(key)
}

// 获得url中的query（询问）,与上边差不多
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// 设置状态码
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Write.WriteHeader(code)
}

// 设置头部信息
func (c *Context) SetHeader(key, value string) {
	c.Write.Header().Set(key, value)
}

// 下边都是一些格式
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
