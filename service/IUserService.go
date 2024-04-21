package service

import "follower/dtos"

type IUserService interface {
	Create(userDto *dtos.UserDto) (*dtos.UserDto, error)
}
