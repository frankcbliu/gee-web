package gee

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	// origin objects
	Writer  http.ResponseWriter
	Request *http.Request
	// 请求信息
	Method string
	Path   string
	Params map[string]string
	// 返回信息
	StatusCode int
	// 中间件信息
	handlers []HandlerFunc
	index    int
	// engine 指针
	engine *Engine
}

func NewContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer:  w,
		Request: req,
		Method:  req.Method,
		Path:    req.URL.Path,
		index:   -1, // 初始化为 -1
	}
}

// 依次调用当前 Context 中所有的中间件
func (c *Context) Next() {
	c.index++
	size := len(c.handlers)
	// 这里遍历所有 handler，是因为不是所有 handler 都会手动调用 c.Next()
	// 对于只作用于请求前的 handler，可以省略 c.Next()
	for ; c.index < size; c.index++ {
		c.handlers[c.index](c)
	}
}

// 返回失败信息
func (c *Context) Fail(code int, err string) {
	c.index = len(c.handlers)
	c.JSON(code, H{"message": err})
}

// 获取 Form 参数
func (c *Context) PostForm(key string) string {
	return c.Request.FormValue(key)
}

// 获取 Query 参数
func (c *Context) Query(key string) string {
	return c.Request.URL.Query().Get(key)
}

// 设置状态码
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// 设置请求头
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

// 返回 format 字符串
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	_, err := c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
	if err != nil {
		panic(err)
	}
}

// 返回 JSON
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		panic(err)
	}
}

// 返回 Data
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	_, err := c.Writer.Write(data)
	if err != nil {
		panic(err)
	}
}

// 返回 HTML
func (c *Context) HTML(code int, name string, data interface{}) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	if err := c.engine.htmlTemplates.ExecuteTemplate(c.Writer, name, data); err != nil {
		c.Fail(http.StatusInternalServerError, err.Error())
	}
}

// 根据参数获取值
func (c *Context) Param(key string) string {
	value, ok := c.Params[key]
	if !ok {
		log.Printf("Find Param key: %v error!", key)
	}
	return value
}
