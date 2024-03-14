package store

import (
	"context"

	v1 "github.com/Mr-LvGJ/gobase/pkg/gobase/model/v1"
)

type Factory interface {
	Users() UserStore
	Tables() TableStore
	Close() error
}

type UserStore interface {
	Create(ctx context.Context, user *v1.User) error
	Update(ctx context.Context, user *v1.User) error
	Delete(ctx context.Context, username string) error
	Get(ctx context.Context, username string) (*v1.User, error)
	List(ctx context.Context) (*v1.UserList, error)
}

type TableStore interface {
	List(ctx context.Context, author string) (*v1.TableList, error)
	Create(ctx context.Context, table *v1.Table) error
}
