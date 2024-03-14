package store

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/Mr-LvGJ/gobase/pkg/common/errno"
	v1 "github.com/Mr-LvGJ/gobase/pkg/gobase/model/v1"
)

type users struct {
	*datastore
}

func newUsers(ds *datastore) *users {
	return &users{datastore: ds}
}

func (u *users) Create(ctx context.Context, user *v1.User) error {
	return u.db.Create(&user).Error
}

func (u *users) Update(ctx context.Context, user *v1.User) error {
	return u.db.Save(user).Error
}

func (u *users) Delete(ctx context.Context, username string) error {
	err := u.db.Where("username = ?", username).Delete(&v1.User{}).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errno.ErrUserNotFound
	}
	return nil
}

func (u *users) Get(ctx context.Context, username string) (*v1.User, error) {
	user := &v1.User{}
	err := u.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return user, nil
}

func (u *users) List(ctx context.Context) (*v1.UserList, error) {
	ret := &v1.UserList{}
	//ol := gormutil.Unpointer(options.Limit, options.Offset)
	//selector, _ := fields.ParseSelector(options.FieldSelector)
	//username, _ := selector.RequiresExactMatch("username")
	//d := u.db.Where("username like ?", "%"+username+"%").
	//	Offset(ol.Offset).
	//	Limit(ol.Limit).
	//	Order("id desc").
	//	Find(&ret.Items)
	return ret, nil
}
