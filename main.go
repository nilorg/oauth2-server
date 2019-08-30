package main

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/nilorg/oauth2"
	"github.com/nilorg/oauth2-server/dao"
	"github.com/nilorg/oauth2-server/middleware"
	"github.com/nilorg/oauth2-server/models"
	"golang.org/x/net/publicsuffix"
)

var (
	redisClient *redis.Client
)

func init() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := redisClient.Ping().Result()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(pong)
}

func main() {

	oauth2Server := oauth2.NewServer()
	oauth2Server.VerifyClient = func(clientID string) (basic *oauth2.ClientBasic, err error) {
		var client *models.Client
		client, err = dao.Client.SelectByID(clientID)
		if err != nil {
			err = oauth2.ErrUnauthorizedClient
			return
		}
		basic = &oauth2.ClientBasic{
			ID:     client.ClientID,
			Secret: client.ClientSecret,
		}
		return
	}
	oauth2Server.VerifyPassword = func(username, password string) (openID string, err error) {
		var user *models.User
		user, err = dao.User.SelectByUsername(username)
		if err != nil {
			err = oauth2.ErrAccessDenied
			return
		}
		if user.Username != username || user.Password != password {
			err = oauth2.ErrAccessDenied
		}
		// scope 额外处理,这里不做处理
		return
	}

	oauth2Server.VerifyAuthorization = func(clientID, redirectUri string) (err error) {
		// var client *models.Client
		// client, err = dao.Client.SelectByID(clientID)
		// if err != nil {
		// 	err = oauth2.ErrAccessDenied
		// 	return
		// }
		// redirectUri 额外处理,这里不做处理
		// scope 额外处理,这里不做处理
		return
	}
	oauth2Server.GenerateCode = func(clientID, openID, redirectURI string, scope []string) (code string, err error) {
		// var client models.Client
		// client, err = dao.Client.SelectByID(clientID)
		// if err != nil {
		// 	err = oauth2.ErrAccessDenied
		// 	return
		// }
		// redirectUri 额外处理,这里不做处理
		// scope 额外处理,这里不做处理
		code = oauth2.RandomCode()
		value := &oauth2.CodeValue{
			ClientID:    clientID,
			RedirectURI: redirectURI,
			Scope:       scope,
		}
		redisClient.Set(fmt.Sprintf("oauth2_code_%s", code), value, time.Minute)
		return
	}
	oauth2Server.VerifyCode = func(code, clientID, redirectURI string) (value *oauth2.CodeValue, err error) {
		value = &oauth2.CodeValue{}
		err = redisClient.Get(fmt.Sprintf("oauth2_code_%s", code)).Scan(value)
		if err != nil {
			err = oauth2.ErrAccessDenied
			return
		}
		if value.ClientID != clientID || value.RedirectURI != redirectURI {
			err = oauth2.ErrAccessDenied
			return
		}
		// scope 额外处理,这里不做处理
		return
	}
	oauth2Server.Init()

	store := cookie.NewStore([]byte("secret"))
	r := gin.Default()
	r.Use(sessions.Sessions("mysession", store))
	r.LoadHTMLGlob("templates/*")
	oauth2Group := r.Group("/oauth2")
	{
		oauth2Group.GET("/login", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "login.tmpl", gin.H{
				"title": "登录",
				"error": GetErrorMessage(ctx),
			})
		})
		oauth2Group.POST("/login", func(ctx *gin.Context) {
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
			session.Set("oauth2_current_user", user)
			_ = session.Save()

			ctx.Redirect(302, ctx.Query("login_redirect_uri"))
		})
		oauth2Group.GET("/authorize", middleware.AuthRequired, func(ctx *gin.Context) {
			clientID := ctx.Query("client_id")
			var err error
			var client *models.Client
			client, err = dao.Client.SelectByID(clientID)
			if err != nil {
				_ = SetErrorMessage(ctx, err.Error())
				ctx.Redirect(http.StatusFound, ctx.Request.RequestURI)
				return
			}
			uri := *ctx.Request.URL
			query := uri.Query()
			queryRedirectURI := query.Get(oauth2.RedirectURIKey)
			if queryRedirectURI == "" {
				query.Set(oauth2.RedirectURIKey, client.RedirectURI)
				uri.RawQuery = query.Encode()
			} else {
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
			}

			session := sessions.Default(ctx)
			currentAccount := session.Get("current_user")
			if currentAccount == nil {
				redirectURI, _ := url.Parse("/oauth2/login")
				redirectURIQuery := url.Values{}
				redirectURIQuery.Set("client_id", clientID)
				redirectURIQuery.Set("login_redirect_uri", uri.String())
				redirectURI.RawQuery = redirectURIQuery.Encode()
				ctx.Redirect(302, redirectURI.String())
			} else {
				ctx.HTML(http.StatusOK, "authorize.tmpl", gin.H{
					"title": "授权页面",
				})
			}
		})
		oauth2Group.POST("/authorize", middleware.AuthRequired, func(ctx *gin.Context) {
			session := sessions.Default(ctx)
			currentAccount := session.Get("current_user")
			rctx := oauth2.NewOpenIDContext(ctx.Request.Context(), currentAccount.(string))
			req := ctx.Request.WithContext(rctx)
			// 模拟请求客户端
			oauth2Server.HandleAuthorize(ctx.Writer, req)
		})
	}
	r.Run() // listen and serve on 0.0.0.0:8080
}

// SetErrorMessage sets a error message,
func SetErrorMessage(ctx *gin.Context, msg string) error {
	session := sessions.Default(ctx)
	session.Set("error_message", msg)
	return session.Save()
}

// GetErrorMessage returns the first error message
func GetErrorMessage(ctx *gin.Context) string {
	session := sessions.Default(ctx)
	value := session.Get("error_message")
	if value != nil {
		session.Delete("error_message")
		_ = session.Save()
		return value.(string)
	}
	return ""
}
