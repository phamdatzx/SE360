package repository

import (
	"user-service/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *model.User) error
	CheckUserExists(username string) (bool, error)
	GetUserByUsername(username string) (*model.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) CheckUserExists(username string) (bool, error) {
	var count int64
	err := r.db.Model(&model.User{}).Where("user_name = ?", username).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *userRepository) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, "user_name = ?", username).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
