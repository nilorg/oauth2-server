package main

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	store := cookie.NewStore([]byte("secret"))
	r := gin.Default()
	r.Use(sessions.Sessions("mysession", store))
	r.LoadHTMLGlob("templates/*")
	oauth2Group := r.Group("/oauth2")
	{
		oauth2Group.GET("/login", func(ctx *gin.Context) {
			redirectUri := ctx.Query("login_redirect_uri")
			uri, _ := url.Parse(redirectUri)

			query := ctx.Request.URL.Query()
			query.Del("login_redirect_uri")

			uri.RawQuery = query.Encode()
			ctx.HTML(http.StatusOK, "login.tmpl", gin.H{
				"title":              "登录",
				"login_redirect_uri": uri.String(),
			})
		})
		oauth2Group.POST("/login", func(ctx *gin.Context) {
			username := ctx.PostForm("username")
			password := ctx.PostForm("password")
			uri := ctx.PostForm("login_redirect_uri")
			if username == "haha" && password == "haha" {
				session := sessions.Default(ctx)
				session.Set("oauth2_current_user", "haha")
				session.Save()

				ctx.Redirect(302, uri)
			} else {
				ctx.HTML(http.StatusOK, "login.tmpl", gin.H{
					"message":            "登录错误",
					"title":              "登录",
					"login_redirect_uri": uri,
				})
			}
		})
		oauth2Group.GET("/authorize", func(ctx *gin.Context) {
			clientID := ctx.Query("client_id")
			session := sessions.Default(ctx)
			currentAccount := session.Get("oauth2_current_user")
			if currentAccount == nil {
				uri := *ctx.Request.URL
				uri.RawQuery = uri.Query().Encode()
				ctx.Redirect(302, fmt.Sprintf("/oauth2/login?client_id=%s&login_redirect_uri=%s", clientID, uri.String()))
			} else {
				redirectURI := ctx.Query("redirect_uri")
				responseType := ctx.Query("response_type")
				state := ctx.Query("state")
				scope := ctx.Query("scope")
				ctx.HTML(http.StatusOK, "authorize.tmpl", gin.H{
					"title":         "授权页面",
					"redirect_uri":  redirectURI,
					"response_type": responseType,
					"state":         state,
					"scope":         scope,
				})
			}
		})
		oauth2Group.POST("/authorize", func(ctx *gin.Context) {
			redirectURI := ctx.PostForm("redirect_uri")
			responseType := ctx.PostForm("response_type")
			state := ctx.PostForm("state")
			scope := ctx.PostForm("scope")

			fmt.Printf("redirectURI: %s \n", redirectURI)
			fmt.Printf("responseType: %s \n", responseType)
			fmt.Printf("state: %s \n", state)
			fmt.Printf("scope: %s \n", scope)

			// 模拟请求客户端
		})
	}
	r.Run() // listen and serve on 0.0.0.0:8080
}
