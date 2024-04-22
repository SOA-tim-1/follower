package dtos

type FollowDto struct {
	FollowedId  int64 `json:"followedId"`
	FollowingId int64 `json:"followingId"`
}
