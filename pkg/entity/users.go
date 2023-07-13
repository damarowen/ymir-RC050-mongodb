// Package entity defines all the entities used in the application.
package entity

import (
	"time"
)

// User represents a user in the collection.
type User struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string    `bson:"name,omitempty" json:"name,omitempty" validate:"required,min=3,max=100"`
	Email     string    `bson:"email,omitempty" json:"email,omitempty" validate:"required,email"`
	Age       int       `bson:"age,omitempty" json:"age,omitempty" validate:"required"`
	CreatedAt time.Time `bson:"created_at,omitempty" json:"created_at,omitempty"`
}

// RequestGetUsers represents a parameter to get user with pagination in the collection.
type RequestGetUsers struct {
	Pagination `json:"pagination"`
}

// ResponseGetUsers represents a parameter to get user with pagination in the collection.
type ResponseGetUsers struct {
	Users      []User `json:"users"`
	Pagination `json:"pagination"`
}
