package handler

import (
	"context"
	"follower/dtos"
	follower "follower/proto"
	"follower/service"

	"google.golang.org/grpc"
)

type UserHandlergRPC struct {
	UserService service.IUserService
	follower.UnimplementedUserServiceServer
}

func (handler *UserHandler) WriteUserRpc(ctx context.Context, in *follower.WriteUserRequest) (*follower.Empty, error) {

	userDto := in.GetUserDto()
	dto := dtos.UserDto{
		ID: userDto.GetId(),
	}

	error := handler.UserService.WriteUser(&dto)
	if error != nil {
		println("Error while creating a new user")
		return new(follower.Empty), nil
	}
	return nil, error
}

func (handler *UserHandler) FindByIdRpc(ctx context.Context, in *follower.FindByIdRequest) (*follower.FindByIdResponse, error) {

	user, err := handler.UserService.FindById(in.GetId())
	if err != nil {
		return nil, err
	}

	userDto := &follower.UserDto{
		Id: user.ID,
	}

	response := &follower.FindByIdResponse{
		UserDto: userDto,
	}

	return response, nil
}

func (handler *UserHandler) CreateFollowConnectionRpc(ctx context.Context, in *follower.CreateFollowConnectionRequest) (*follower.Empty, error) {

	err := handler.UserService.CreateFollowConnection(in.GetFirstId(), in.GetSecondId())
	if err != nil {
		return nil, err
	}

	return new(follower.Empty), nil
}

func (handler *UserHandler) GetFollowsRpc(ctx context.Context, in *follower.GetFollowsRequest, opts ...grpc.CallOption) (*follower.FollowsResponse, error) {

	follows, err := handler.UserService.GetFollows(in.GetId())
	if err != nil {
		return nil, err
	}

	if len(*follows) == 0 {
		return &follower.FollowsResponse{}, nil
	}

	response := &follower.FollowsResponse{
		Follows: *follows,
	}

	return response, nil
}

func (handler *UserHandler) GetFollowersRpc(ctx context.Context, in *follower.GetFollowersRequest) (*follower.FollowersResponse, error) {

	followers, err := handler.UserService.GetFollowers(in.GetId())
	if err != nil {
		return nil, err
	}

	if len(*followers) == 0 {
		return &follower.FollowersResponse{}, nil
	}

	response := &follower.FollowersResponse{
		Followers: *followers,
	}

	return response, nil
}

func (handler *UserHandler) GetSuggestionsForUserRpc(ctx context.Context, in *follower.GetSuggestionsRequest) (*follower.SuggestionsResponse, error) {

	suggestions, err := handler.UserService.GetSuggestionsForUser(in.GetId())
	if err != nil {
		return nil, err
	}

	if len(*suggestions) == 0 {
		return &follower.SuggestionsResponse{}, nil
	}

	response := &follower.SuggestionsResponse{
		Suggestions: *suggestions,
	}

	return response, nil
}

func (handler *UserHandler) CheckIfFollowingConnectionExistRpc(ctx context.Context, in *follower.CheckIfFollowingConnectionExistRequest) (*follower.CheckResponse, error) {

	isFollowing, err := handler.UserService.CheckIfFollowingConnectionExist(in.GetId1(), in.GetId2())
	if err != nil {
		return nil, err
	}

	response := &follower.CheckResponse{
		Exists: isFollowing,
	}

	return response, nil
}

func (handler *UserHandler) DeleteFollowConnectionRpc(ctx context.Context, in *follower.DeleteFollowConnectionRequest) (*follower.Empty, error) {

	err := handler.UserService.DeleteFollowConnection(in.GetId1(), in.GetId2())
	if err != nil {
		return nil, err
	}
	return new(follower.Empty), nil
}
