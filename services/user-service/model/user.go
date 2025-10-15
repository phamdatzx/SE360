package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserName string `json:"username" gorm:"uniqueIndex"`
	Password string `json:"password"`
	Name     string `json:"name"`
}
