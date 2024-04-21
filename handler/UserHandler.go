package handler

import (
	"encoding/json"
	"fmt"
	"follower/dtos"
	"follower/service"
	"net/http"
)

type UserHandler struct {
	UserService service.IUserService
}

func (handler *UserHandler) Create(writer http.ResponseWriter, req *http.Request) {
	var userDto dtos.UserDto
	err := json.NewDecoder(req.Body).Decode(&userDto)
	if err != nil {
		println("Error while parsing json")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	createdUserDto, err := handler.UserService.Create(&userDto)
	if err != nil {
		println("Error while creating a new tour")
		writer.WriteHeader(http.StatusExpectationFailed)
		return
	}

	fmt.Println(createdUserDto)

	// Marshal the createdTourDto object into JSON
	tourJSON, err := json.Marshal(createdUserDto)
	if err != nil {
		println("Error while encoding tourDto to JSON")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set response headers and write the JSON response
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	_, err = writer.Write(tourJSON)
	if err != nil {
		println("Error while writing JSON response")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}
