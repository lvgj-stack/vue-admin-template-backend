package v1

import (
	"time"

	"github.com/go-playground/validator/v10"

	"github.com/Mr-LvGJ/jota/models"

	"github.com/Mr-LvGJ/gobase/pkg/common/constant"
	"github.com/Mr-LvGJ/gobase/pkg/gobase/model/v1/entity"
)

type Table struct {
	models.BaseModel

	Title       string               `json:"title"        gorm:"column:title"`
	Status      constant.TableStatus `json:"status"       gorm:"column:status;type:varchar(31);not null"                                                     binding:"required"`
	Author      string               `json:"author"       gorm:"column:author;type:varchar(255);not null"                                                    binding:"required" validate:"min=3,max=32"`
	DisplayTime *time.Time           `json:"display_time" gorm:"column:display_time;type:datetime(6);not null;index:created_at;default:current_timestamp(6)"`
	PageViews   int                  `json:"pageviews"    gorm:"column:page_views;type:int;not null;default:0"`
}

func (t *Table) TableName() string {
	return "table"
}

func (t *Table) Validate() error {
	validate := validator.New()
	return validate.Struct(t)
}

type TableList struct {
	Items []*entity.Table `json:"items"`
}
