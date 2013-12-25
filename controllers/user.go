package controllers

import (
	"github.com/astaxie/beego"
 	spew "github.com/davecgh/go-spew/spew"
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
		if user.Email == "forevervmax@gmail.com" {
			this.SetSession("login", true)
			this.Ctx.Redirect(302, "/user/home")
		} else {
			flash.Error("error password")
			flash.Store(&this.Controller)
			this.Ctx.Redirect(302, "/user/home")
		}
	}
}

func (this *UserController) Login() {
	this.TplNames = "login.html"
	if this.GetSession("login") == true{
		this.Redirect("/blog/list", 302)
	}
	beego.ReadFromRequest(&this.Controller)
	this.Render()
}

func (this *UserController) Logout() {
	this.DelSession("login")
	this.Redirect("/user/home", 302)
}

func (this *UserController) Home(){
	this.TplNames = "home.html"
	if this.GetSession("login") == true{
		this.Render()
	}else{
		spew.Dump(this.GetSession("login"))
		flash.Error("need login")
		flash.Store(&this.Controller)
		this.Redirect("/user/login", 302)
	}
}
