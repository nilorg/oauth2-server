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
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/nilorg/oauth2"
	"github.com/nilorg/oauth2-server/models"
)

var (
	redisClient *redis.Client
	db          *gorm.DB
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

	db, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	// Migrate the schema
	db.AutoMigrate(&models.Client{}, &models.User{})
	fmt.Println(".....")
}

func main() {

	oauth2Server := oauth2.NewServer()
	oauth2Server.VerifyClient = func(clientID string) (basic *oauth2.ClientBasic, err error) {
		var client *models.Client
		err = db.Where("client_id = ?", clientID).First(client).Error
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
	oauth2Server.VerifyPassword = func(username, password string, scope []string) (err error) {
		user := &models.User{}
		err = db.Where("username = ?", username).First(user).Error
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
	oauth2Server.VerifyCredentialsScope = func(clientID string, scope []string) (err error) {
		client := &models.Client{}
		err = db.Where("client_id = ?", clientID).First(client).Error
		if err != nil {
			err = oauth2.ErrAccessDenied
			return
		}
		// scope 额外处理,这里不做处理
		return
	}
	oauth2Server.VerifyAuthorization = func(clientID, redirectUri string, scope []string) (err error) {
		client := &models.Client{}
		err = db.Where("client_id = ?", clientID).First(client).Error
		if err != nil {
			err = oauth2.ErrAccessDenied
			return
		}
		// redirectUri 额外处理,这里不做处理
		// scope 额外处理,这里不做处理
		return
	}
	oauth2Server.GenerateCode = func(clientID, redirectUri string, scope []string) (code string, err error) {
		client := &models.Client{}
		err = db.Where("client_id = ?", clientID).First(client).Error
		if err != nil {
			err = oauth2.ErrAccessDenied
			return
		}
		// redirectUri 额外处理,这里不做处理
		// scope 额外处理,这里不做处理
		code = oauth2.RandomCode()
		value := &oauth2.CodeValue{
			ClientID:    clientID,
			RedirectUri: redirectUri,
			Scope:       scope,
		}
		redisClient.Set(fmt.Sprintf("oauth2_code_%s", code), value, time.Minute)
		return
	}
	oauth2Server.VerifyCode = func(code, clientID, redirectUri string) (value *oauth2.CodeValue, err error) {
		value = &oauth2.CodeValue{}
		err = redisClient.Get(fmt.Sprintf("oauth2_code_%s", code)).Scan(value)
		if err != nil {
			err = oauth2.ErrAccessDenied
			return
		}
		if value.ClientID != clientID || value.RedirectUri != redirectUri {
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
			})
		})
		oauth2Group.POST("/login", func(ctx *gin.Context) {
			username := ctx.PostForm("username")
			password := ctx.PostForm("password")
			if username == "haha" && password == "haha" {
				session := sessions.Default(ctx)
				session.Set("oauth2_current_user", "haha")
				session.Save()

				//uri, _ := url.Parse(ctx.Query("login_redirect_uri"))
				//query := ctx.Request.URL.Query()
				//query.Del("login_redirect_uri")
				//uri.RawQuery = query.Encode()
				//ctx.Redirect(302, uri.String())

				ctx.Redirect(302, ctx.Query("login_redirect_uri"))
			} else {
				ctx.HTML(http.StatusOK, "login.tmpl", gin.H{
					"message": "登录错误",
					"title":   "登录",
				})
			}
		})
		oauth2Group.GET("/authorize", func(ctx *gin.Context) {
			clientID := ctx.Query("client_id")
			session := sessions.Default(ctx)
			currentAccount := session.Get("oauth2_current_user")
			if currentAccount == nil {
				redirectUri, _ := url.Parse("/oauth2/login")
				query := url.Values{}
				query.Set("client_id", clientID)
				query.Set("login_redirect_uri", ctx.Request.URL.String())
				redirectUri.RawQuery = query.Encode()
				ctx.Redirect(302, redirectUri.String())
			} else {
				ctx.HTML(http.StatusOK, "authorize.tmpl", gin.H{
					"title": "授权页面",
				})
			}
		})
		oauth2Group.POST("/authorize", func(ctx *gin.Context) {
			oauth2Server.HandleAuthorize(ctx.Writer, ctx.Request)
			// 模拟请求客户端
		})
	}
	r.Run() // listen and serve on 0.0.0.0:8080
}
