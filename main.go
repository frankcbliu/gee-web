package main

import (
	"gee"
	"net/http"
)

func main() {
	r := gee.New()
	r.Use(gee.Logger())
	r.GET("/hello", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee!</h1>\n")
	})

	api := r.Group("/api")
	{
		api.GET("/", func(c *gee.Context) {
			c.HTML(http.StatusOK, "<h1>Hello Api.</h1>\n")
		})
		api.GET("/speak", func(c *gee.Context) {
			// expect /speak?name=geek
			c.String(http.StatusOK, "hello %s, you are at %s\n", c.Query("name"), c.Path)
		})
	}

	auth := r.Group("/auth")
	auth.Use(gee.OnlyForAuth())
	{
		auth.GET("/hello/:name", func(c *gee.Context) {
			// expect /hello/geek
			c.String(http.StatusOK, "hello %s, you are at %s\n", c.Param("name"), c.Path)
		})
		auth.GET("/assets/*filepath", func(c *gee.Context) {
			c.JSON(http.StatusOK, gee.H{"filepath": c.Param("filepath")})
		})
	}
	err := r.Run(":8080")
	if err != nil { // 启动服务失败
		panic("start server error!")
	}
}
