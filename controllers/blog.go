package controllers

import (
	"github.com/astaxie/beego"
	spew "github.com/davecgh/go-spew/spew"
	"strconv"
	"github.com/garyburd/redigo/redis"
	"time"
)

type Blog struct {
	Id int64
	Title       string `form:"title"`
	Content     string `form:"content"`
	TimeCreated int64
}

type BlogController struct {
	beego.Controller
}

func (this *BlogController) New() {
	this.TplNames = "blog/edit.html"
	this.Render()
}

func (this *BlogController) Post() {
	blog := Blog{}
	if err := this.ParseForm(&blog); err != nil {
		beego.Info(err)
	} else {
		c,err :=redis.Dial("tcp", ":6379")
		if err!=nil{
			spew.Dump(err)
		}
		post_count,_ := c.Do("INCR","post:count")
		blog.Id = post_count.(int64)
		blog.TimeCreated = time.Now().Unix()
		spew.Dump(blog)
		c.Send("LPUSH","post:list",blog.Id)
		c.Send("HMSET",redis.Args{}.Add("post:"+strconv.FormatInt(post_count.(int64),36)).AddFlat(&blog)...)
		c.Flush()
		r,_ :=c.Receive()
		spew.Dump(r)
		if err != nil{
			spew.Dump(err)
		}else{
			this.Redirect("/user/home", 302)
		}
	}
}

func (this *BlogController) List() {
	this.Ctx.WriteString("list")
}
