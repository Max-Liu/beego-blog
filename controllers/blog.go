package controllers

import (
	"github.com/astaxie/beego"
	spew "github.com/davecgh/go-spew/spew"
	"github.com/garyburd/redigo/redis"
	"log"
	"strconv"
	"time"
	// "fmt"
)

type Blog struct {
	Id          int64
	Title       string `form:"title"`
	Content     string `form:"content"`
	TimeCreated int64
}

type BlogController struct {
	beego.Controller
}

func (this *BlogController) New() {
	this.TplNames = "blog/edit.html"
	this.Layout = "layout/layout.html"
	this.Data["css"] = `<link rel="stylesheet" type="text/css" href="/static/css/bootstrap-wysihtml5.css"></link>
    <link href="/static/css/blog.css" rel="stylesheet">`
    this.Data["js"] = `<script src="/static/js/wysihtml5-0.3.0.min.js"></script>
    <script src="/static/js/jquery-1.7.1.min.js"></script>
    <script src="/static/js/bootstrap.min.js"></script>
    <script src="/static/js/bootstrap-wysihtml5.min.js"></script>
    <script type="text/javascript">
    $('#textarea').wysihtml5();
    </script>`
	this.Render()
}

func (this *BlogController) Get() {
	this.Layout = "layout/layout.html"
	this.TplNames = "blog/list.html"
	this.Data["css"]= `<link href="/static/css/list.css" rel="stylesheet">`

	if this.GetSession("login") == true {
		c, err := redis.Dial("tcp", ":6379")
		if err != nil {
			log.Println(err)
		}
		postIndexList, _ := redis.Values(c.Do("LRANGE", "post:list", 0, 10))

		var post Blog
		var postList []Blog

		for _, v := range postIndexList {
			r, _ := redis.Values(c.Do("HGETALL", "post:"+string(v.([]uint8))))
			redis.ScanStruct(r, &post)
			postList = append(postList, post)
		}

		this.Data["blogList"] = postList
		this.Render()
	} else {
		flash.Error("need login")
		flash.Store(&this.Controller)
		this.Redirect("/user/login", 302)
	}
}

func (this *BlogController) Post() {
	blog := Blog{}
	if err := this.ParseForm(&blog); err != nil {
		beego.Info(err)
	} else {
		c, err := redis.Dial("tcp", ":6379")
		defer c.Close()
		if err != nil {
			log.Println(err)
		}
		post_count, _ := c.Do("INCR", "post:count")
		blog.Id = post_count.(int64)
		blog.TimeCreated = time.Now().Unix()
		spew.Dump(blog)
		c.Send("LPUSH", "post:list", blog.Id)
		c.Send("HMSET", redis.Args{}.Add("post:"+strconv.FormatInt(post_count.(int64), 36)).AddFlat(&blog)...)
		c.Flush()
		r, _ := c.Receive()
		spew.Dump(r)
		if err != nil {
			spew.Dump(err)
		} else {
			this.Redirect("/", 302)
		}
	}
}

func (this *BlogController) Artical() {
	this.Layout = "layout/layout.html"
	this.TplNames = "blog/index.html"
	blogId := this.Input().Get("id")
	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		log.Println(err)
	}
	var blog Blog
	reply, _ := c.Do("HGETALL", "post:"+blogId)
	redis.ScanStruct(reply.([]interface{}), &blog)

	this.Data["blog"] = blog
	this.Render()
}
