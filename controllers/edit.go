package controllers

import (
	"log"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/davecgh/go-spew/spew"
	"github.com/garyburd/redigo/redis"
)

type EditController struct {
	beego.Controller
}

func (this *EditController) Editroute() {

	route := this.Ctx.Input.Param(":route")

	switch route {
	case "excute":
		this.Excute()
	default:
		this.Layout = "layout/layout.html"
		this.TplNames = "blog/edit.html"
		this.TplNames = "blog/edit.html"
		this.Layout = "layout/layout.html"
		this.Data["css"] = `<link rel="stylesheet" type="text/css" href="/static/css/bootstrap-wysihtml5.css"></link>
    <link href="/static/css/blog_new.css" rel="stylesheet">`
		this.Data["js"] = `<script src="/static/js/wysihtml5-0.3.0.min.js"></script>
    <script src="/static/js/bootstrap-wysihtml5.min.js"></script>
    <script type="text/javascript">
    $('#textarea').wysihtml5({
	"font-styles": true, //Font styling, e.g. h1, h2, etc. Default true
	"emphasis": true, //Italics, bold, etc. Default true
	"lists": true, //(Un)ordered lists, e.g. Bullets, Numbers. Default true
	"html": true, //Button which allows you to edit the generated HTML. Default false
	"link": true, //Button to insert a link. Default true
	"image": true, //Button to insert an image. Default true,
	"color": true //Button to change color of font  
})</script>`

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

func (this *EditController) Excute() {
	blog := Blog{}
	if err := this.ParseForm(&blog); err != nil {
		beego.Info(err)
	} else {
		c, err := redis.Dial("tcp", ":6379")
		defer c.Close()

		if err != nil {
			log.Println(err)
		}

		reply, _ := redis.Bool(c.Do("HGET", "slug.to.id", blog.Slug))
		id, _ := redis.Int64(c.Do("HGET", "slug.to.id", blog.Slug))

		if reply == true {
			blog.Id = id
			c.Do("HMSET", redis.Args{}.Add("post:"+strconv.Itoa(int(blog.Id))).AddFlat(&blog)...)
			if err != nil {
				spew.Dump(err)
			} else {
				this.Redirect("/edit/"+blog.Slug, 302)
			}
		}

	}
}
