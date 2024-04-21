package repo

import "follower/model"

type IUserRepository interface {
	CreateUser(user *model.User) (model.User, error)
}
