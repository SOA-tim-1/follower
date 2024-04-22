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

func (handler *UserHandler) GetUserById(writer http.ResponseWriter, req *http.Request) {
	idStr := mux.Vars(req)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		// Handle parsing error
		http.Error(writer, "Invalid ID", http.StatusBadRequest)
		return
	}

	user, err := handler.UserService.FindById(id)
	writer.Header().Set("Content-Type", "application/json")
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(user)
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

	err = handler.UserService.CreateFollowConnection(followDto.FollowingId, followDto.FollowedId)
	if err != nil {
		println("Error while creating a new follow connection")
		writer.WriteHeader(http.StatusExpectationFailed)
		return
	}

	// Set response headers and write the JSON response
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
}

func (handler *UserHandler) DeleteFollowConnection(writer http.ResponseWriter, req *http.Request) {
	var followDto dtos.FollowDto

	values := req.URL.Query()
	followedId := values.Get("followedId")
	followingId := values.Get("followingId")

	// Convert the string parameters to integers
	followedIdInt, err := strconv.ParseInt(followedId, 10, 64)
	if err != nil {
		http.Error(writer, "Invalid followedId", http.StatusBadRequest)
		return
	}

	followingIdInt, err := strconv.ParseInt(followingId, 10, 64)
	if err != nil {
		http.Error(writer, "Invalid followingId", http.StatusBadRequest)
		return
	}

	// Map the URL parameters to the FollowDto
	followDto = dtos.FollowDto{
		FollowedId:  followedIdInt,
		FollowingId: followingIdInt,
	}

	err = handler.UserService.DeleteFollowConnection(followDto.FollowingId, followDto.FollowedId)
	if err != nil {
		println("Error while delete follow connection")
		writer.WriteHeader(http.StatusExpectationFailed)
		return
	}

	// Set response headers and write the JSON response
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
}

func (handler *UserHandler) CheckIfFirstFollowSecond(writer http.ResponseWriter, req *http.Request) {
	var followDto dtos.FollowDto

	values := req.URL.Query()
	followedId := values.Get("followedId")
	followingId := values.Get("followingId")

	// Convert the string parameters to integers
	followedIdInt, err := strconv.ParseInt(followedId, 10, 64)
	if err != nil {
		http.Error(writer, "Invalid followedId", http.StatusBadRequest)
		return
	}

	followingIdInt, err := strconv.ParseInt(followingId, 10, 64)
	if err != nil {
		http.Error(writer, "Invalid followingId", http.StatusBadRequest)
		return
	}

	// Map the URL parameters to the FollowDto
	followDto = dtos.FollowDto{
		FollowedId:  followedIdInt,
		FollowingId: followingIdInt,
	}

	isFollowing, err := handler.UserService.CheckIfFollowingConnectionExist(followDto.FollowingId, followDto.FollowedId)
	if err != nil {
		println("Error while checking follow")
		writer.WriteHeader(http.StatusExpectationFailed)
		return
	}

	jsonData, err := json.Marshal(isFollowing)
	if err != nil {
		http.Error(writer, "Error encoding JSON response", http.StatusInternalServerError)
		return
	}

	// Write the response
	writer.WriteHeader(http.StatusOK)
	writer.Write(jsonData)
}
