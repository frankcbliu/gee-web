package gee

import (
	"log"
	"net/http"
	"strings"
)

// 定义一个处理函数
type HandlerFunc func(*Context)

type RouterGroup struct {
	prefix      string
	middleWares []HandlerFunc // 支持中间件
	parent      *RouterGroup  // 支持多级分组
	engine      *Engine       // 全局共用一个 engine 实例
}

type Engine struct {
	*RouterGroup
	router *router
	groups []*RouterGroup // 存储所有分组
}

// 初始化 engine
func New() *Engine {
	engine := &Engine{router: NewRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// 创建分组
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

func (group *RouterGroup) addRoute(method string, pattern string, handler HandlerFunc) {
	pattern = group.prefix + pattern
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

// 添加 Get 请求
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

// 添加 POST 请求
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

// 启动 HTTP 服务
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// 将中间件应用到某个 group
func (group *RouterGroup) Use(middleWares ...HandlerFunc) {
	group.middleWares = append(group.middleWares, middleWares...)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middleWares []HandlerFunc
	// 将分组中的中间件加入到 middleWares 中
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middleWares = append(middleWares, group.middleWares...)
		}
	}

	c := NewContext(w, req)
	// 将 middleWares 赋值给 Context 中的 handlers
	c.handlers = middleWares
	engine.router.handle(c)
}
