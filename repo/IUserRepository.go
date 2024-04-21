package repo

import (
	"context"
	"follower/model"
)

type IUserRepository interface {
	CheckConnection()
	WriteUser(user *model.User) error
	FindById(id int64) (model.User, error)
	CreateFollowConnection(firstId int64, secondId int64) error
	GetFollows(id int64) (*[]int64, error)
	GetFollowers(id int64) (*[]int64, error)
	GetSuggestionsForUser(id int64) (*[]int64, error)
	DropAll() error
	CloseDriverConnection(ctx context.Context)
}
