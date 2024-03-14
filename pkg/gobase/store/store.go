package store

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/Mr-LvGJ/jota/log"
	"github.com/Mr-LvGJ/jota/models"

	"github.com/Mr-LvGJ/gobase/pkg/gobase/model/v1/entity"
	"github.com/Mr-LvGJ/gobase/pkg/gobase/store/dal"

	"github.com/Mr-LvGJ/gobase/pkg/common/setting"
)

var (
	client Factory
	DB     *gorm.DB
)

func Client() Factory {
	return client
}

type datastore struct {
	db *gorm.DB
	Q  *dal.Query
}

func (d *datastore) Tables() TableStore {
	return newTables(d)
}

func (d *datastore) Close() error {
	//TODO implement me
	panic("implement me")
}

func (d *datastore) Users() UserStore {
	return newUsers(d)
}

func dbMigrate() {
	db := DB.Session(&gorm.Session{
		Logger: DB.Logger.LogMode(logger.Warn),
	})
	err := db.AutoMigrate(&entity.User{}, &entity.Table{})
	if err != nil {
		panic(err)
	}
}

func Setup(ctx context.Context) (*gorm.DB, error) {
	log.Info(ctx, "init database", "config", setting.C().Database)
	var err error
	DB, err = models.New(setting.C().Database)
	if err != nil {
		panic(err)
	}
	dal.SetDefault(DB)
	dbMigrate()
	client = &datastore{
		db: DB,
		Q:  dal.Q,
	}
	return DB, nil
}

func GetDB() *gorm.DB {
	return DB
}
