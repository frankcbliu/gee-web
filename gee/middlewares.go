package gee

import (
	"log"
	"time"
)

// 统一日志打印中间件
func Logger() HandlerFunc {
	return func(ctx *Context) {
		// 开始计时
		t := time.Now()
		// 继续执行请求处理
		ctx.Next()
		// 计算耗时并打印日志
		log.Printf("[%d] %s in %v", ctx.StatusCode, ctx.Request.RequestURI, time.Since(t))
	}
}

// 仅用于 auth 路由的测试中间件
func OnlyForAuth() HandlerFunc {
	return func(c *Context) {
		// 开始计时
		t := time.Now()
		// 假设当前服务返回出错
		c.Fail(500, "Internal Server Error")
		// 计算耗时并打印日志
		log.Printf("[%d] %s in %v for group Auth", c.StatusCode, c.Request.RequestURI, time.Since(t))
	}
}
