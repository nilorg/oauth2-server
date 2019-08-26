package models

import "github.com/jinzhu/gorm"

// Client ...
type Client struct {
	gorm.Model
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}
