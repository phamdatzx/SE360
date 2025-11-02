package service

import (
	"user-service/dto"
	customError "user-service/error"
	"user-service/model"
	"user-service/repository"
	"user-service/utils"
)

type UserService interface {
	RegisterUser(user model.User) (dto.UserResponse, error)
	Login(request dto.LoginRequest) (dto.LoginResponse, error)
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

	//hash password
	user.Password, _ = utils.HashPassword(user.Password)
	user.Role = "USER"
	err = s.repo.CreateUser(&user)
	//map to dto
	resultDto := dto.UserResponse{user.ID.String(), user.UserName, user.Name}

	return resultDto, err
}

func (s *userService) Login(request dto.LoginRequest) (dto.LoginResponse, error) {
	//get user model
	user, err := s.repo.GetUserByUsername(request.Username)
	if err != nil {
		return dto.LoginResponse{}, customError.NewAppErrorWithErr(401, "username not found", err)
	}

	//check password
	if utils.CheckPasswordHash(request.Password, user.Password) { //correct info
		token, err := utils.GenerateToken(user.ID.String(), user.UserName, user.Role)
		if err != nil {
			return dto.LoginResponse{}, err
		}
		return dto.LoginResponse{AccessToken: token}, nil
	} else { //incorrect info
		return dto.LoginResponse{}, customError.NewAppError(401, "incorrect password")
	}
}
