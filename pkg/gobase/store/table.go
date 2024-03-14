package store

import (
	"context"

	"github.com/Mr-LvGJ/jota/id"

	v1 "github.com/Mr-LvGJ/gobase/pkg/gobase/model/v1"
)

var _ TableStore = &tables{}

type tables struct {
	*datastore
}

func newTables(ds *datastore) *tables {
	return &tables{
		datastore: ds,
	}
}

func (t *tables) List(ctx context.Context, author string) (*v1.TableList, error) {
	var ret = &v1.TableList{}
	tableDal := t.Q.Table
	ts, err := tableDal.WithContext(ctx).Where(tableDal.Author.Eq(author)).Find()
	if err != nil {
		return nil, err
	}
	ret.Items = ts

	//if err := t.db.Where("author = ?", author).Find(&ret.Items).Error; err != nil {
	//	return nil, err
	//}
	return ret, nil
}

func (t *tables) Create(ctx context.Context, table *v1.Table) error {
	table.Id = id.Next()
	return t.db.Model(&v1.Table{}).Create(table).Error
}
