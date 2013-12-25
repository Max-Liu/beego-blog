package main

import (
	"blog/controllers"
	"github.com/astaxie/beego"
)

func main() {
	beego.SessionProvider = "redis"
	beego.SessionSavePath = "127.0.0.1:6379"
	beego.AutoRouter(&controllers.UserController{})
	beego.AutoRouter(&controllers.BlogController{})
	beego.Run()
}

