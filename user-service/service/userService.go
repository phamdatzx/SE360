package service

import (
	"user-service/dto"
	customError "user-service/error"
	"user-service/model"
	"user-service/repository"
)

type UserService interface {
	RegisterUser(user model.User) (dto.UserResponse, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{repo: userRepo}
}

func (s *userService) RegisterUser(user model.User) (dto.UserResponse, error) {
	exist, err := s.repo.CheckUserExists(user.UserName)
	if err != nil {
		return dto.UserResponse{}, err
	}
	if exist {
		return dto.UserResponse{}, customError.NewAppError(409, "username has already existed")
	}

	err = s.repo.CreateUser(&user)
	//map to dto
	resultDto := dto.UserResponse{user.ID, user.UserName, user.Name}

	return resultDto, err
}
