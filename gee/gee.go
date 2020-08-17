package gee

import (
	"html/template"
	"log"
	"net/http"
	"path"
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
	router        *router
	groups        []*RouterGroup     // 存储所有分组
	htmlTemplates *template.Template // 用于 html 渲染
	funcMap       template.FuncMap
}

// 初始化 engine
func New() *Engine {
	engine := &Engine{router: NewRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// 设置 funcMap
func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

// 渲染函数
func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
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
	// 给 ctx 的 engine 赋值
	c.engine = engine
	// 将 middleWares 赋值给 Context 中的 handlers
	c.handlers = middleWares
	engine.router.handle(c)
}

// 创建静态文件处理 handler
func (group *RouterGroup) CreateStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(ctx *Context) {
		file := ctx.Param("filepath")
		// 判断文件是否存在 or 是否有权限处理文件
		if _, err := fs.Open(file); err != nil {
			ctx.Status(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(ctx.Writer, ctx.Request)
	}
}

// 将硬盘上的 root 路径映射到路由 relativePath 上
func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.CreateStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	// 注册 handler
	group.GET(urlPattern, handler)
}
