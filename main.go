package main

import (
	"blog/controllers"
	"log"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/garyburd/redigo/redis"
	// spew "github.com/davecgh/go-spew/spew"
)

var checkDatabase = func(ctx *context.Context) {

	if ctx.Input.Url() != "/blog/password" && ctx.Input.Url() != "/blog/setpass" {
		c, err := redis.Dial("tcp", ":6379")
		defer c.Close()
		if err != nil {
			log.Println(err)
		}
		reply, _ := c.Do("GET", "user:password")
		if reply == nil {
			ctx.Redirect(302, "/blog/password")
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
	beego.Router("/blog/", &controllers.BlogController{}, "get:Blog")
	beego.Run()
}
