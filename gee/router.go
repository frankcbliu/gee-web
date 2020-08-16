package gee

import (
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node       // 存储根节点
	handlers map[string]HandlerFunc // 存储节点到 handleFunc 的映射
}

func NewRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// 解析 pattern
// 如果路由是 /static/*filepath，则 results = [static]
// 如果路由是 /static/:name/doc，则 results = [static, :name, doc]
func parsePattern(pattern string) []string {
	parts := strings.Split(pattern, "/")

	results := make([]string, 0)

	for _, part := range parts {
		if part != "" {
			results = append(results, part)
			if part[0] == '*' {
				break
			}
		}
	}
	return results
}

// 往路由表中添加路由
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)
	key := method + "-" + pattern

	// 获取根节点
	_, ok := r.roots[method]
	if !ok { // 不存在则创建
		r.roots[method] = &node{}
	}
	// 往根节点中插入
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

// 从路由表中查询节点
func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path) // 获取 parts 数组
	params := make(map[string]string)

	root, ok := r.roots[method] // 获取根节点
	if !ok {
		return nil, nil
	}

	n := root.search(searchParts, 0)

	if n != nil { // 匹配节点非空
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 { // *开头，且不只有*
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}

// 处理路由函数
func (r *router) handle(c *Context) {
	// 获取节点和从路由中解析出来的参数
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil { // 节点是否存在，是判断路由是否存在的依据
		c.Params = params
		key := c.Method + "-" + n.pattern
		if handler, ok := r.handlers[key]; ok {
			c.handlers = append(c.handlers, handler)
		} else {
			c.handlers = append(c.handlers, func(context *Context) {
				c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
			})
		}
	}
	c.Next()
}
