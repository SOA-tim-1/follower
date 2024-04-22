package service

import (
	"follower/dtos"
)

type IUserService interface {
	CheckConnection()
	WriteUser(user *dtos.UserDto) error
	FindById(id int64) (*dtos.UserDto, error)
	CreateFollowConnection(firstId int64, secondId int64) error
	GetFollows(id int64) (*[]int64, error)
	GetFollowers(id int64) (*[]int64, error)
	GetSuggestionsForUser(id int64) (*[]int64, error)
	CheckIfFollowingConnectionExist(id1, id2 int64) (bool, error)
	DeleteFollowConnection(firstId int64, secondId int64) error
}
