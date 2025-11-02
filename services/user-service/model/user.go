package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserName string    `json:"username" gorm:"uniqueIndex"`
	Password string    `json:"password"`
	Name     string    `json:"name"`
	Role     string    `json:"role"`
}

// Hook tự động sinh UUID trước khi tạo record
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}
