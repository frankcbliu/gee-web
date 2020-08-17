package main

import (
	"fmt"
	"gee"
	"html/template"
	"net/http"
	"time"
)

type student struct {
	Name string
	Age  int8
}

func formatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func main() {
	r := gee.New()
	r.Use(gee.Logger())
	r.SetFuncMap(template.FuncMap{
		"formatAsDate": formatAsDate,
	})
	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./static")

	api := r.Group("/api")
	{
		api.GET("/", func(c *gee.Context) {
			c.HTML(http.StatusOK, "css.tmpl", nil)
		})
		api.GET("/speak", func(c *gee.Context) {
			// expect /speak?name=geek
			c.String(http.StatusOK, "hello %s, you are at %s\n", c.Query("name"), c.Path)
		})
		stu1 := &student{Name: "geek", Age: 20}
		stu2 := &student{Name: "frank", Age: 22}
		api.GET("/students", func(c *gee.Context) {
			c.HTML(http.StatusOK, "arr.tmpl", gee.H{
				"title":  "gee",
				"stuArr": [2]*student{stu1, stu2},
			})
		})
		api.GET("/date", func(c *gee.Context) {
			c.HTML(http.StatusOK, "custom_func.tmpl", gee.H{
				"title": "gee",
				"now":   time.Now(),
			})
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
