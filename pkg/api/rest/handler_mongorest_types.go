// Package rest handles the port operations.
package rest

import "github.com/kubuskotak/ymir-test/pkg/entity"

// GetUsersRequest is a struct that embeds Pagination fields
// for getting users request.
type GetUsersRequest struct {
	entity.Pagination `json:"pagination"`
}

// GetUsersResponse is a struct for response
// that holds a slice of User objects.
type GetUsersResponse struct {
	Data []entity.User
}

// SaveUserResponse is a struct for response
// that holds a User object.
type SaveUserResponse struct {
	entity.User
}

// SaveUserRequest is a struct for request
// that holds a User object.
type SaveUserRequest struct {
	entity.User
}