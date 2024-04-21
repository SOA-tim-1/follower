package handler

import (
	"encoding/json"
	"follower/dtos"
	"follower/service"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	UserService service.IUserService
}

func (handler *UserHandler) WriteUser(writer http.ResponseWriter, req *http.Request) {
	var userDto dtos.UserDto
	err := json.NewDecoder(req.Body).Decode(&userDto)
	if err != nil {
		println("Error while parsing json")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	error := handler.UserService.WriteUser(&userDto)
	if error != nil {
		println("Error while creating a new user")
		writer.WriteHeader(http.StatusExpectationFailed)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	return
}

func (handler *UserHandler) GetFollows(writer http.ResponseWriter, req *http.Request) {
	idStr := mux.Vars(req)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		// Handle parsing error
		http.Error(writer, "Invalid ID", http.StatusBadRequest)
		return
	}

	follows, err := handler.UserService.GetFollows(id)
	writer.Header().Set("Content-Type", "application/json")
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	if len(*follows) == 0 {
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte("[]")) // Write empty array
		return
	}

	// Encode tours into JSON
	jsonData, err := json.Marshal(follows)
	if err != nil {
		http.Error(writer, "Error encoding JSON response", http.StatusInternalServerError)
		return
	}

	// Write the response
	writer.WriteHeader(http.StatusOK)
	writer.Write(jsonData)
}

func (handler *UserHandler) GetFollowers(writer http.ResponseWriter, req *http.Request) {
	idStr := mux.Vars(req)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		// Handle parsing error
		http.Error(writer, "Invalid ID", http.StatusBadRequest)
		return
	}

	follows, err := handler.UserService.GetFollowers(id)
	writer.Header().Set("Content-Type", "application/json")
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	if len(*follows) == 0 {
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte("[]")) // Write empty array
		return
	}

	// Encode tours into JSON
	jsonData, err := json.Marshal(follows)
	if err != nil {
		http.Error(writer, "Error encoding JSON response", http.StatusInternalServerError)
		return
	}

	// Write the response
	writer.WriteHeader(http.StatusOK)
	writer.Write(jsonData)
}

func (handler *UserHandler) GetSuggestionsForUser(writer http.ResponseWriter, req *http.Request) {
	idStr := mux.Vars(req)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		// Handle parsing error
		http.Error(writer, "Invalid ID", http.StatusBadRequest)
		return
	}

	follows, err := handler.UserService.GetSuggestionsForUser(id)
	writer.Header().Set("Content-Type", "application/json")
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	if len(*follows) == 0 {
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte("[]")) // Write empty array
		return
	}

	// Encode tours into JSON
	jsonData, err := json.Marshal(follows)
	if err != nil {
		http.Error(writer, "Error encoding JSON response", http.StatusInternalServerError)
		return
	}

	// Write the response
	writer.WriteHeader(http.StatusOK)
	writer.Write(jsonData)
}

func (handler *UserHandler) CreateFollowConnection(writer http.ResponseWriter, req *http.Request) {
	var followDto dtos.FollowDto
	err := json.NewDecoder(req.Body).Decode(&followDto)
	if err != nil {
		println("Error while parsing json")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	err = handler.UserService.CreateFollowConnection(followDto.UserId, followDto.FollowedId)
	if err != nil {
		println("Error while creating a new tour")
		writer.WriteHeader(http.StatusExpectationFailed)
		return
	}

	// Set response headers and write the JSON response
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
}
