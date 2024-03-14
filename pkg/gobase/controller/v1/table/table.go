package table

import (
	"github.com/gin-gonic/gin"

	"github.com/Mr-LvGJ/gobase/pkg/common/constant"
	"github.com/Mr-LvGJ/gobase/pkg/common/core"
	v1 "github.com/Mr-LvGJ/gobase/pkg/gobase/model/v1"
	srvv1 "github.com/Mr-LvGJ/gobase/pkg/gobase/service/v1"
	"github.com/Mr-LvGJ/gobase/pkg/gobase/store"
)

type TableController struct {
	srv srvv1.Service
}

func NewTableController(store store.Factory) *TableController {
	return &TableController{
		srv: srvv1.NewService(store),
	}
}

func (t *TableController) List(c *gin.Context) {
	u := c.Request.Context().Value(constant.XUsernameKey).(string)
	tableList, err := t.srv.Tables().List(c, u)
	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}
	core.WriteResponse(c, nil, tableList)
}

func (t *TableController) Create(c *gin.Context) {
	tb := v1.Table{
		Author: c.Value(constant.XUsernameKey).(string),
	}
	if err := c.ShouldBindJSON(&tb); err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	if err := t.srv.Tables().Create(c, &tb); err != nil {
		core.WriteResponse(c, err, nil)
		return
	}
	core.WriteResponse(c, nil, nil)
}
