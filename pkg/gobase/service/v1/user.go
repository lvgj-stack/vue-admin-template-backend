package v1

import (
	"context"
	"regexp"

	"github.com/Mr-LvGJ/jota/id"

	"github.com/Mr-LvGJ/gobase/pkg/common/errno"
	v1 "github.com/Mr-LvGJ/gobase/pkg/gobase/model/v1"
	"github.com/Mr-LvGJ/gobase/pkg/gobase/store"
)

var userIdGetter = mustNewIdGenerator()

type userService struct {
	store store.Factory
}

func mustNewIdGenerator() *id.Generator {
	generator, err := id.NewGenerator()
	if err != nil {
		panic(err)
	}
	return generator
}

func newUsers(srv *service) *userService {
	return &userService{store: srv.store}
}

func (u *userService) Create(ctx context.Context, user *v1.User) error {
	user.Id = userIdGetter.Next().WithPrefix("u-")
	if err := u.store.Users().Create(ctx, user); err != nil {
		if match, _ := regexp.MatchString("Duplicate entry '.* for key 'username'", err.Error()); match {
			return errno.ErrUserAlreadyExist
		}
		return err
	}
	return nil
}

func (u *userService) Update(ctx context.Context, user *v1.User) error {
	return u.store.Users().Update(ctx, user)
}

func (u *userService) Delete(ctx context.Context, username string) error {
	return u.store.Users().Delete(ctx, username)
}

func (u *userService) Get(ctx context.Context, username string) (*v1.User, error) {
	user, err := u.store.Users().Get(ctx, username)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userService) List(ctx context.Context) (*v1.UserList, error) {
	return &v1.UserList{}, nil
}
