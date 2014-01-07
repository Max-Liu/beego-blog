package main

import (
	"blog/controllers"
	"github.com/astaxie/beego"
	"log"
	"github.com/astaxie/beego/context"
	"github.com/garyburd/redigo/redis"
	// spew "github.com/davecgh/go-spew/spew"
)

var checkDatabase = func (ctx *context.Context){

	if ctx.Input.Url() != "/blog/setpass"{
		c, err := redis.Dial("tcp", ":6379")
		defer c.Close()
		if err != nil {
			log.Println(err)
		}
		reply,_ :=c.Do("GET","user:password")
		if reply == nil{
			ctx.Redirect(302,"/blog/setpass")
		}
	}
}

func main() {
	beego.SessionProvider = "redis"
	beego.SessionSavePath = "127.0.0.1:6379"
	beego.InsertFilter("*",1,checkDatabase )
	beego.AutoRouter(&controllers.UserController{})
	beego.AutoRouter(&controllers.BlogController{})
	beego.Router("/", &controllers.BlogController{},"get:Home")
	beego.Router("/blog/", &controllers.BlogController{},"get:Blog")
	beego.Run()
}

