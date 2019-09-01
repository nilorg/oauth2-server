package controller

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/nilorg/oauth2"
	"github.com/nilorg/oauth2-server/dao"
	"github.com/nilorg/oauth2-server/models"
	"github.com/nilorg/oauth2-server/module/store"
	"github.com/nilorg/oauth2-server/utils"
)

var (
	oauth2Server *oauth2.Server
	sourceScope  = []string{
		"read_write",
		"read",
		"write",
	}
)

func init() {
	oauth2Server = oauth2.NewServer()
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
		return
	}

	oauth2Server.VerifyRedirectURI = func(clientID, redirectURI string) (err error) {
		var client *models.Client
		client, err = dao.Client.SelectByID(clientID)
		if err != nil {
			err = oauth2.ErrAccessDenied
			return
		}
		if strings.Index(redirectURI, client.RedirectURI) == -1 {
			err = oauth2.ErrInvalidRedirectURI
		}
		return
	}
	oauth2Server.GenerateCode = func(clientID, openID, redirectURI string, scope []string) (code string, err error) {
		code = oauth2.RandomCode()
		value := &oauth2.CodeValue{
			ClientID:    clientID,
			OpenID:      openID,
			RedirectURI: redirectURI,
			Scope:       scope,
		}
		store.RedisClient.Set(fmt.Sprintf("oauth2_code_%s", code), value, time.Minute)
		return
	}
	oauth2Server.VerifyCode = func(code, clientID, redirectURI string) (value *oauth2.CodeValue, err error) {
		value = &oauth2.CodeValue{}
		err = store.RedisClient.Get(fmt.Sprintf("oauth2_code_%s", code)).Scan(value)
		if err != nil {
			err = oauth2.ErrAccessDenied
			return
		}
		if value.ClientID != clientID || value.RedirectURI != redirectURI {
			err = oauth2.ErrAccessDenied
			return
		}
		if !utils.ScopeIsEqual(value.Scope, sourceScope) {
			err = oauth2.ErrInvalidScope
		}
		return
	}
	oauth2Server.VerifyScope = func(scope []string) (err error) {
		if !utils.ScopeIsSubset(scope, sourceScope) {
			err = oauth2.ErrInvalidScope
		}
		return
	}
	oauth2Server.Init()

}

// SetErrorMessage set a error message
func SetErrorMessage(ctx *gin.Context, msg string) error {
	session := sessions.Default(ctx)
	session.Set("error_message", msg)
	return session.Save()
}

// GetErrorMessage return the first error message
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
