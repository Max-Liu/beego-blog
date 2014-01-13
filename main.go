package main

import (
	"blog/controllers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/garyburd/redigo/redis"
	// spew "github.com/davecgh/go-spew/spew"
)

var checkDatabase = func(ctx *context.Context) {

	if ctx.Input.Url() != "/blog/password" && ctx.Input.Url() != "/blog/setpass" && ctx.Input.Url() != "/err/db" {
		c, err := redis.Dial("tcp", ":6379")
		if err != nil {
			ctx.Redirect(302, "/err/db")
		} else {
			defer c.Close()
			reply, _ := c.Do("GET", "user:password")
			if reply == nil {
				ctx.Redirect(302, "/blog/password")
			}
		}
	}
}

func main() {
	beego.SessionProvider = "redis"
	beego.SessionSavePath = "127.0.0.1:6379"

	beego.InsertFilter("*", 1, checkDatabase)

	beego.AutoRouter(&controllers.UserController{})

	beego.Router(`/blog/:route([\w-]+)`, &controllers.BlogController{}, "post:Blogroute;get:Blogroute")
	beego.Router(`/edit/:route([\w-]+)`, &controllers.EditController{}, "post:Editroute;get:Editroute")
	beego.Router("/", &controllers.BlogController{}, "get:Home")
	beego.Router("/err/db", &controllers.BlogController{}, "get:Errdb")
	beego.Router("/blog/", &controllers.BlogController{}, "get:Blog")

	beego.Run()
}
