package dao

import (
	"github.com/nilorg/oauth2-server/models"
)

type client struct {
}

func (*client) SelectByID(clientID string) (mc *models.Client, err error) {
	var dbResult models.Client
	err = db.Where("client_id = ?", clientID).First(&dbResult).Error
	if err != nil {
		return
	}
	mc = &dbResult
	return
}
