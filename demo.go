package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
	"path"
	"time"
)

func m1(context *gin.Context) {
	start := time.Now()
	context.Next()
	cost := time.Since(start)
	fmt.Printf("[COST] %s %v\n", context.Request.URL, cost)
}

func main() {
	r := gin.Default()

	// 全局注册中间件
	r.Use(m1)

	// 加载静态文件
	r.Static("/xxx", "./statics")
	// gin框架中给模板添加自定义函数
	r.SetFuncMap(template.FuncMap{
		"safe": func(str string) template.HTML {
			return template.HTML(str)
		},
	})

	// r.LoadHTMLFiles("templates/index.tmpl")
	r.LoadHTMLGlob("templates/**/*")

	r.GET("/posts/index", func(context *gin.Context) {
		context.HTML(http.StatusOK, "posts/index.tmpl", gin.H{
			"title": "posts/index.tmpl",
		})
	})
	r.GET("/users/index", func(context *gin.Context) {
		context.HTML(http.StatusOK, "users/index.tmpl", gin.H{
			"title": "<a href=\"www.baidu.com\">baidu</a>",
		})
	})
	r.GET("/json1", func(context *gin.Context) {
		data := map[string]interface{}{
			"name":    "张三",
			"message": "hello world",
			"age":     18,
		}
		context.JSON(http.StatusOK, data)
	})
	r.GET("/json2", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{"name": "张三", "message": "hello world", "age": 18})
	})
	r.GET("json3", func(context *gin.Context) {
		type msg struct {
			Name    string `json:"name"`
			Message string
			Age     int
		}
		var data = msg{
			Name:    "张三",
			Message: "Hello golang",
			Age:     18,
		}
		context.JSON(http.StatusOK, data)
	})

	// 获取Query参数
	r.GET("/web", func(context *gin.Context) {
		//name := context.Query("name") // 获取参数
		//name := context.DefaultQuery("name", "somebody") // 取不到用默认值
		name, ok := context.GetQuery("name")
		if !ok {
			name = "somebody"
		}
		context.JSON(http.StatusOK, gin.H{
			"name": name,
		})
	})

	// 获取Form参数
	r.POST("/login", func(context *gin.Context) {
		username := context.PostForm("name")
		password := context.PostForm("password")
		context.JSON(http.StatusOK, gin.H{
			"username": username,
			"password": password,
		})
	})

	// 获取路径参数
	r.GET("/:name/:age", func(context *gin.Context) {
		name := context.Param("name")
		age := context.Param("age")
		context.JSON(http.StatusOK, gin.H{
			"name": name,
			"age":  age,
		})
	})

	// 参数绑定
	r.GET("/user", func(context *gin.Context) {
		type UserInfo struct {
			Username string `form:"username"`
			Password string `form:"password"`
		}
		var u UserInfo
		err := context.ShouldBind(&u)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		} else {
			fmt.Printf("%#v\n", u)
			context.JSON(http.StatusOK, gin.H{
				"status": "ok",
			})
		}
	})

	// 接收上传文件
	r.POST("/upload", func(context *gin.Context) {
		file, err := context.FormFile("f1")
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		} else {
			dst := path.Join("./", file.Filename)
			context.SaveUploadedFile(file, dst)
			context.JSON(http.StatusOK, gin.H{
				"status": "ok",
			})
		}
	})

	// 重定向
	r.GET("/xxx", func(context *gin.Context) {
		context.Redirect(http.StatusMovedPermanently, "http://www.baidu.com")
	})

	// 请求转发
	r.GET("/a", func(context *gin.Context) {
		context.Request.URL.Path = "/b"
		r.HandleContext(context)
	})
	r.GET("/b", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"message": "b",
		})
	})

	// Any
	r.Any("/xyz", func(context *gin.Context) {
		switch context.Request.Method {
		case "GET":
			context.JSON(http.StatusOK, gin.H{"method": "GET"})
		case "POST":
			context.JSON(http.StatusOK, gin.H{"method": "POST"})
		case "DELETE":
			context.JSON(http.StatusOK, gin.H{"method": "DELETE"})
			// ...
		}
	})

	// NoRoute
	r.NoRoute(func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"msg": "hell world",
		})
	})

	// 路由组
	userGroup := r.Group("/user")
	userGroup.GET("/xx", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{"msg": "ok"})
	})
	userGroup.GET("/oo", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{"msg": "ok"})
	})

	r.Run(":9090")
}
