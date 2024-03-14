package v1

import "github.com/Mr-LvGJ/gobase/pkg/gobase/store"

var _ Service = &service{}

type service struct {
	store store.Factory
}

func (s *service) Tables() TableSrv {
	return newTables(s)
}

func (s *service) Users() UserSrv {
	return newUsers(s)
}

func NewService(store store.Factory) Service {
	return &service{
		store: store,
	}
}
