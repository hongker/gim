package dto

type ChatroomCreateRequest struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
type ChatroomCreateResponse struct{}

type ChatroomUpdateRequest struct{}
type ChatroomUpdateResponse struct{}

type ChatroomJoinRequest struct {
	Id string `json:"id"`
}
type ChatroomJoinResponse struct{}

type ChatroomLeaveRequest struct{}
type ChatroomLeaveResponse struct{}

type ChatroomDismissRequest struct{}
type ChatroomDismissResponse struct{}
