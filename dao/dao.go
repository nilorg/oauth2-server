package dao

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/nilorg/oauth2-server/models"
)

var (
	db     *gorm.DB
	Client = &client{}
	User   = &user{}
)

func init() {
	var err error
	db, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	// Migrate the schema
	db.AutoMigrate(&models.Client{}, &models.User{})
}
