package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/nilorg/oauth2-server/controller"
	"github.com/nilorg/oauth2-server/middleware"
)

func main() {
	store := cookie.NewStore([]byte("secret"))
	r := gin.Default()
	r.Use(sessions.Sessions("mysession", store))
	r.LoadHTMLGlob("templates/*")
	oauth2Group := r.Group("/oauth2")
	{
		oauth2Group.GET("/login", controller.LoginPage)
		oauth2Group.POST("/login", controller.Login)
		oauth2Group.GET("/authorize", middleware.AuthRequired, controller.AuthorizePage)
		oauth2Group.POST("/authorize", middleware.AuthRequired, controller.Authorize)
	}
	r.Run() // listen and serve on 0.0.0.0:8080
}
