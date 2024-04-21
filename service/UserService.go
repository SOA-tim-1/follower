package service

import (
	"follower/dtos"
	"follower/model"
	"follower/repo"

	"github.com/rafiulgits/go-automapper"
)

type UserService struct {
	UserRepo repo.IUserRepository
}

func (service *UserService) Create(userDto *dtos.UserDto) (*dtos.UserDto, error) {
	var user model.User
	automapper.Map(userDto, &user)

	createdUser, err := service.UserRepo.CreateUser(&user)
	if err != nil {
		return nil, err
	}

	// Map the created tour back to DTO
	createdUserDto := dtos.UserDto{}
	automapper.Map(&createdUser, &createdUserDto)

	return &createdUserDto, nil
}
