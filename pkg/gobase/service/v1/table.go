package v1

import (
	"context"

	v1 "github.com/Mr-LvGJ/gobase/pkg/gobase/model/v1"
	"github.com/Mr-LvGJ/gobase/pkg/gobase/store"
)

var _ TableSrv = &tableService{}

type tableService struct {
	store store.Factory
}

func (t *tableService) Create(ctx context.Context, table *v1.Table) error {
	return t.store.Tables().Create(ctx, table)
}

func (t *tableService) List(ctx context.Context, author string) (*v1.TableList, error) {
	return t.store.Tables().List(ctx, author)
}

func newTables(srv *service) *tableService {
	return &tableService{store: srv.store}
}
