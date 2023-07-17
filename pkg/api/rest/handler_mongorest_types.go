// Package rest handles the port operations.
package rest

import "github.com/kubuskotak/ymir-test/pkg/entity"

// GetListUsersRequest is a struct that embeds Pagination fields
// for getting users request.
type GetListUsersRequest struct {
	entity.Pagination `json:"pagination"`
}

// ResponseMessage is a struct for response
// that holds a message.
type ResponseMessage struct {
	Message string
}

// GetListUsersResponse is a struct for response
// that holds a slice of User objects.
type GetListUsersResponse struct {
	Data []entity.User
}

// GetUserResponse is a struct for response
// that return User objects.
type GetUserResponse struct {
	entity.User
}

// GetRequestParam is a struct for request
// that holds a UserId from param.
type GetRequestParam struct {
	UserID string
}

// UpsertUserRequest is a struct for request
// that holds a User object.
type UpsertUserRequest struct {
	GetRequestParam
	entity.User
}
