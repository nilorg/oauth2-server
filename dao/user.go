package dao

import (
	"github.com/nilorg/oauth2-server/models"
)

type user struct {
}

func (*user) SelectByUsername(username string) (mu *models.User, err error) {
	var dbResult models.User
	err = db.Where("username = ?", username).First(&dbResult).Error
	if err != nil {
		return
	}
	mu = &dbResult
	return
}
