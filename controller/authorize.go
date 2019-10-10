package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/nilorg/oauth2"
	"github.com/nilorg/oauth2-server/dao"
	"github.com/nilorg/oauth2-server/models"
	"golang.org/x/net/publicsuffix"
)

// AuthorizePage 授权页面
func AuthorizePage(ctx *gin.Context) {
	errMsg := GetErrorMessage(ctx)
	if errMsg != "" {
		ctx.HTML(http.StatusOK, "authorize.tmpl", gin.H{
			"title": "授权页面",
			"error": GetErrorMessage(ctx),
		})
		return
	}

	clientID := ctx.Query("client_id")
	var err error
	var client *models.Client
	client, err = dao.Client.SelectByID(clientID)
	if err != nil {
		fmt.Println("查询客户端错误.....")
		_ = SetErrorMessage(ctx, err.Error())
		ctx.Redirect(http.StatusFound, ctx.Request.RequestURI)
		return
	}
	// uri := *ctx.Request.URL
	// query := uri.Query()
	query := ctx.Request.URL.Query()
	queryRedirectURI := query.Get(oauth2.RedirectURIKey)
	if queryRedirectURI == "" {
		// query.Set(oauth2.RedirectURIKey, client.RedirectURI)
		// uri.RawQuery = query.Encode()
		queryRedirectURI = client.RedirectURI
	}
	// 判断重定向顶级域名是否和数据库中的顶级域名相等
	var qrLevelDomain string
	qrLevelDomain, err = publicsuffix.EffectiveTLDPlusOne(queryRedirectURI)
	if err != nil {
		_ = SetErrorMessage(ctx, err.Error())
		ctx.Redirect(http.StatusFound, ctx.Request.RequestURI)
		return
	}
	var dbLevelDomain string
	dbLevelDomain, err = publicsuffix.EffectiveTLDPlusOne(client.RedirectURI)
	if err != nil {
		_ = SetErrorMessage(ctx, err.Error())
		ctx.Redirect(http.StatusFound, ctx.Request.RequestURI)
		return
	}
	if qrLevelDomain != dbLevelDomain {
		_ = SetErrorMessage(ctx, "重定向域名不符合后台配置规范")
		ctx.Redirect(http.StatusFound, ctx.Request.RequestURI)
		return
	}
	ctx.HTML(http.StatusOK, "authorize.tmpl", gin.H{
		"title": "授权页面",
		"error": nil,
	})
	return
}

// Authorize authorize post
func Authorize(ctx *gin.Context) {
	session := sessions.Default(ctx)
	currentAccount := session.Get("current_user")
	cu := currentAccount.(*models.User)
	rctx := oauth2.NewOpenIDContext(ctx.Request.Context(), fmt.Sprint(cu.ID))
	req := ctx.Request.WithContext(rctx)
	// 模拟请求客户端
	oauth2Server.HandleAuthorize(ctx.Writer, req)
}
