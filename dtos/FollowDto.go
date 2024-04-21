package dtos

type FollowDto struct {
	UserId     int64 `json:"userId"`
	FollowedId int64 `json:"followedId"`
}
