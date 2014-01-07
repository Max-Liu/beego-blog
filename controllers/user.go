package controllers

import (
	"github.com/astaxie/beego"
 	// spew "github.com/davecgh/go-spew/spew"
 	"github.com/garyburd/redigo/redis"
 	"log"
)

type UserController struct {
	beego.Controller
}

type User struct {
	// Id  int `form:"-"`
	Email    string `form:"email"`
	Password string `form:"password"`
}

var flash = beego.NewFlash()

func (this *UserController) Login_api() {
	user := User{}
	if err := this.ParseForm(&user); err != nil {
		beego.Info(err)
	} else {
		c, err := redis.Dial("tcp", ":6379")
		if err != nil{
			log.Println(err)
		}
		defer c.Close()
		reply,_:=c.Do("get", "user:password")
		if user.Password == string(reply.([]uint8)) {
			this.SetSession("login", true)
			this.Ctx.Redirect(302, "/blog/new")
		} else {
			flash.Error("error password")
			flash.Store(&this.Controller)
			this.Ctx.Redirect(302, "/user/login")
		}
	}
}

func (this *UserController) Login() {
	this.TplNames = "login.html"
	if this.GetSession("login") == true{
		this.Redirect("/blog/new", 302)
	}
	beego.ReadFromRequest(&this.Controller)
	this.Render()
}

func (this *UserController) Logout() {
	this.DelSession("login")
	this.Redirect("/user/home", 302)
}


