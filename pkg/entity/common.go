// Package entity defines all the entities used in the application.
package entity

// Pagination Get Users.
type Pagination struct {
	Page  int `validate:"gte=0,default=1"`
	Limit int `validate:"gte=0,default=10"`
}
