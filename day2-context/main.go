package main

import (
	"go_test/easy_gee/day2-context/gee"
)

func main() {
	r := gee.New()
	r.Get("/", func(ctx *gee.Context) {
		ctx.HTML(200, "<h1>hello world</h1>")
	})
	r.Get("/hello", func(ctx *gee.Context) {
		ctx.String(200, "hello %s ,you are at %s\n", ctx.Query("name"), ctx.Path)
	})
	r.POST("/login", func(ctx *gee.Context) {
		ctx.JSON(200, gee.H{
			"username": ctx.PostFrom("username"),
			"kds":      ctx.Query("username"),
			"password": ctx.PostFrom("password"),
		})
	})
	r.Run(":9999")
}
