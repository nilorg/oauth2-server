package controller

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/nilorg/oauth2-server/dao"
)

func LoginPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "login.tmpl", gin.H{
		"title": "登录",
		"error": GetErrorMessage(ctx),
	})
}

func Login(ctx *gin.Context) {
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")

	user, err := dao.User.SelectByUsername(username)
	if err != nil {
		_ = SetErrorMessage(ctx, err.Error())
		ctx.Redirect(http.StatusFound, ctx.Request.RequestURI)
		return
	}
	if user.Username != username || user.Password != password {
		_ = SetErrorMessage(ctx, "账号密码不正确")
		ctx.Redirect(http.StatusFound, ctx.Request.RequestURI)
		return
	}
	session := sessions.Default(ctx)
	session.Set("current_user", user)
	err = session.Save()
	if err != nil {
		_ = SetErrorMessage(ctx, err.Error())
		ctx.Redirect(http.StatusFound, ctx.Request.RequestURI)
		return
	}
	ctx.Redirect(302, ctx.Query("login_redirect_uri"))
}
