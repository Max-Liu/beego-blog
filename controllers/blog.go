package controllers

import (
	"log"
	"github.com/astaxie/beego"
	spew "github.com/davecgh/go-spew/spew"
	"github.com/garyburd/redigo/redis"

	// "log"
	//
	"strconv"

	"time"
)

type Blog struct {
	Id          int64
	Title       string `form:"title"`
	Slug        string `form:"slug"`
	Content     string `form:"content"`
	TimeCreated int64
}

type BlogController struct {
	beego.Controller
}

func (this *BlogController) New() {
	if this.GetSession("login") != true {
		this.Redirect("/user/login", 302)
	}
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

func (this *BlogController) Password() {
	this.TplNames = "blog/setpass.html"
	this.Layout = "layout/layout.html"
	this.Render()
}

func (this *BlogController) SetPass() {
	password := this.Input().Get("password")
	if password != "" {
		log.Println(password)
		c, err := redis.Dial("tcp", ":6379")
		if err != nil {
			log.Println(err)
		}
		defer c.Close()

		_, err = c.Do("SET", "user:password", password)
		if err != nil {
			log.Println(err)
		}
		this.Redirect("/blog/new", 302)
	}

}
func (this *BlogController) Home() {
	this.TplNames = "index.html"
	this.Layout = "layout/layout.html"
	this.Data["css"] = `    <link rel="stylesheet" href="/static/css/home/normalize.css">
    <link rel="stylesheet" href="/static/css/home/main.css">
    <link rel="stylesheet" href="/static/css/home/bootstrap.css">
    <script src="/static/js/home/vendor/modernizr-2.7.1.min.js"></script>`

	this.Data["js"] = `<script src="//ajax.googleapis.com/ajax/libs/jquery/1.10.2/jquery.min.js"></script>
	<script>window.jQuery || document.write('<script src="js/vendor/jquery-1.10.2.min.js"><\/script>')</script>
	<script src="/static/js/home/plugins.js"></script>
	<script src="/static/js/home/bootstrap.js"></script>
	<script src="/static/js/home/main.js"></script>`

	this.Data["page"] = "home"

	this.Render()
}

func (this *BlogController) Blogroute() {

	route := this.Ctx.Input.Param(":route")

	switch route {
	case "home":
		this.Home()
	case "new":
		this.New()
	case "post":
		this.Post()
	case "setpass":
		this.SetPass()
	case "password":
		this.Password()
	default:
		this.Layout = "layout/layout.html"
		this.TplNames = "blog/index.html"

		c, err := redis.Dial("tcp", ":6379")
		if err != nil {
			log.Println(err)
		}

		blogId, _ := redis.String(c.Do("HGET", "slug.to.id", route))

		var blog Blog
		reply, _ := c.Do("HGETALL", "post:"+blogId)
		redis.ScanStruct(reply.([]interface{}), &blog)

		this.Data["blog"] = blog
		this.Render()
	}

}

func (this *BlogController) Blog() {
	this.Layout = "layout/layout.html"
	this.TplNames = "blog/list.html"

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
		slug := this.Input().Get("slug")
		post_count, _ := redis.Int(c.Do("GET", "post:count"))

		reply, _ := redis.Bool(c.Do("HSETNX", "slug.to.id", slug, post_count+1))
		if reply == true {
			post_count, _ := c.Do("INCR", "post:count")
			blog.Id = post_count.(int64)
			blog.TimeCreated = time.Now().Unix()
			c.Send("LPUSH", "post:list", blog.Id)
			c.Send("HMSET", redis.Args{}.Add("post:"+strconv.FormatInt(post_count.(int64), 36)).AddFlat(&blog)...)
			c.Flush()
			if err != nil {
				spew.Dump(err)
			} else {
				this.Redirect("/", 302)
			}
		}

	}
}
