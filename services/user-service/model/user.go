package model

import (
	"user-service/dto"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey;column:id"`
	Username string    `json:"username" gorm:"uniqueIndex;column:username"`
	Password string    `json:"password" gorm:"column:password"`
	Name     string    `json:"name" gorm:"column:name"`
	Email    string    `json:"email" gorm:"column:email"`
	Role     string    `json:"role" gorm:"column:role"`
	IsActive bool      `json:"is_active" gorm:"column:is_active;default:false"`
	IsVerify bool      `json:"is_verify" gorm:"column:is_verify;default:false"`
	IsBanned bool      `json:"is_banned" gorm:"column:is_banned;default:false"`
}

// Hook tự động sinh UUID trước khi tạo record
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}

func NewUser(request dto.RegisterRequest) *User {
	return &User{
		Username: request.Username,
		Password: request.Password, // nếu có mã hóa password thì xử lý trước khi gán
		Name:     request.Name,
		Email:    request.Email,
		Role:     request.Role,
	}
}
