// Package users is implements component logic.
package users

import (
	"context"
	"github.com/kubuskotak/ymir-test/pkg/adapters"
	"github.com/kubuskotak/ymir-test/pkg/entity"
	"github.com/kubuskotak/ymir-test/pkg/usecase"
	"reflect"
)

func init() {
	usecase.Register(usecase.Registration{
		Name: "users",
		Inf:  reflect.TypeOf((*T)(nil)).Elem(),
		New: func() any {
			return &impl{}
		},
	})
}

// T is the interface implemented by all users Component implementations.
type T interface {
	GetAll(ctx context.Context, paging entity.RequestGetUsers) (entity.ResponseGetUsers, error)
	Create(ctx context.Context, user entity.User) (entity.User, error)
	GetById(ctx context.Context, userId string) (entity.User, error)
	DeleteById(ctx context.Context, userId string) error
	UpdateById(ctx context.Context, user entity.User) (entity.User, error)
}

type impl struct {
	adapter *adapters.Adapter
}

// Init initializes the execution of a process involved in a users Component usecase.
func (i *impl) Init(adapter *adapters.Adapter) error {
	i.adapter = adapter
	return nil
}
