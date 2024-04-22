package service

import (
	"fmt"
	"follower/dtos"
	"follower/model"
	"follower/repo"

	"github.com/rafiulgits/go-automapper"
)

type UserService struct {
	UserRepo repo.IUserRepository
}

func (service *UserService) CheckConnection() {
	service.UserRepo.CheckConnection()
}

func (service *UserService) WriteUser(userDto *dtos.UserDto) error {
	var user model.User
	automapper.Map(userDto, &user)
	err := service.UserRepo.WriteUser(&user)
	if err != nil {
		return err
	}

	return nil
}

func (service *UserService) FindById(id int64) (*dtos.UserDto, error) {
	user, err := service.UserRepo.FindById(id)
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("menu item with id %d not found", id))
	}

	userDto := dtos.UserDto{}
	automapper.Map(user, &userDto)

	return &userDto, nil

}

func (service *UserService) CreateFollowConnection(firstId int64, secondId int64) error {
	firstExist, err := service.CheckIfExist(firstId)
	if err != nil {
		return err
	}
	secondExist, err := service.CheckIfExist(secondId)
	if err != nil {
		return err
	}

	if !firstExist {
		user := model.User{ID: firstId}
		service.UserRepo.WriteUser(&user)
	}

	if !secondExist {
		user := model.User{ID: secondId}
		service.UserRepo.WriteUser(&user)
	}

	connectionExist, error := service.UserRepo.CheckIfFollowingConnectionExist(firstId, secondId)
	if error != nil {
		return error
	}

	if connectionExist {
		return nil
	} else {
		return service.UserRepo.CreateFollowConnection(firstId, secondId)
	}

}

func (service *UserService) DeleteFollowConnection(firstId int64, secondId int64) error {
	return service.UserRepo.DeleteFollowConnection(firstId, secondId)
}

func (service *UserService) GetFollows(id int64) (*[]int64, error) {
	return service.UserRepo.GetFollows(id)
}

func (service *UserService) GetFollowers(id int64) (*[]int64, error) {
	return service.UserRepo.GetFollowers(id)
}

func (service *UserService) GetSuggestionsForUser(id int64) (*[]int64, error) {
	return service.UserRepo.GetSuggestionsForUser(id)
}

func (service *UserService) CheckIfFollowingConnectionExist(id1, id2 int64) (bool, error) {
	return service.UserRepo.CheckIfFollowingConnectionExist(id1, id2)
}

func (service *UserService) CheckIfExist(id int64) (bool, error) {
	user, err := service.UserRepo.FindById(id)
	if err != nil {
		return false, err
	}

	if user.ID == 0 {
		return false, nil
	} else {
		return true, nil
	}
}
