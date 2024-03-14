package user

import (
	"github.com/gin-gonic/gin"

	"github.com/Mr-LvGJ/jota/log"

	"github.com/Mr-LvGJ/gobase/pkg/common/constant"

	"github.com/Mr-LvGJ/gobase/pkg/common/auth"
	"github.com/Mr-LvGJ/gobase/pkg/common/core"
	"github.com/Mr-LvGJ/gobase/pkg/common/errno"
	"github.com/Mr-LvGJ/gobase/pkg/common/token"
	v1 "github.com/Mr-LvGJ/gobase/pkg/gobase/model/v1"
	srvv1 "github.com/Mr-LvGJ/gobase/pkg/gobase/service/v1"
	"github.com/Mr-LvGJ/gobase/pkg/gobase/store"
)

type UserController struct {
	srv srvv1.Service
}

func NewUserController(store store.Factory) *UserController {
	return &UserController{
		srv: srvv1.NewService(store),
	}
}

func (u *UserController) Get(c *gin.Context) {
	log.Info(c, "get user function called.")
	user, err := u.srv.Users().Get(c, c.Param("name"))
	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}
	core.WriteResponse(c, nil, user)

}

func (u *UserController) Info(c *gin.Context) {
	username := c.Request.Context().Value(constant.XUsernameKey).(string)
	user, err := u.srv.Users().Get(c, username)
	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}
	core.WriteResponse(c, nil, user)
}

func (u *UserController) Login(c *gin.Context) {
	var r LoginRequest

	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(c, errno.ErrBind, nil)

		return
	}
	log.Info(c.Request.Context(), "LoginRequest", "request", r, "username", r.Username)

	user, err := u.srv.Users().Get(c, r.Username)
	if err != nil {
		core.WriteResponse(c, errno.ErrUserNotFound, nil)
		return
	}

	if err = auth.Compare(user.Password, r.Password); err != nil {
		core.WriteResponse(c, errno.ErrPasswordIncorrect, nil)
		return
	}

	t, err := token.Sign(r.Username)
	if err != nil {
		core.WriteResponse(c, errno.ErrToken, nil)
		return
	}

	core.WriteResponse(c, nil, LoginResponse{Token: t})

}

func (u *UserController) Logout(c *gin.Context) {
	username := c.Request.Context().Value(constant.XUsernameKey).(string)
	log.Info(c.Request.Context(), "success logout", "username", username)
	core.WriteResponse(c, nil, username)
}

func (u *UserController) Create(c *gin.Context) {
	log.Info(c.Request.Context(), "user create func called", c.Request)
	var r v1.User

	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(c, errno.ErrBind, nil)
		return
	}

	if err := r.Validate(); err != nil {
		core.WriteResponse(c, errno.ErrValidation, nil)
		return
	}
	var err error
	r.Password, err = auth.Encrypt(r.Password)
	if err != nil {
		core.WriteResponse(c, errno.ErrEncrypt, nil)
		return
	}

	if err = u.srv.Users().Create(c, &r); err != nil {
		core.WriteResponse(c, err, nil)
		return
	}
	core.WriteResponse(c, errno.OK, r)
}

func (u *UserController) List(c *gin.Context) {
	log.Info(c.Request.Context(), "list user func called.")
	//if err := c.ShouldBindQuery(&r); err != nil {
	//	core.WriteResponse(c, errno.ErrBind, nil)
	//	return
	//}
	//
	//users, err := u.srv.Users().List(c, r)
	//if err != nil {
	//	core.WriteResponse(c, err, nil)
	//	return
	//}

	core.WriteResponse(c, nil, nil)
}

func (u *UserController) Update(c *gin.Context) {
	log.Info(c.Request.Context(), "update user info func called.", c.Request)
	var r UpdateRequest

	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(c, errno.ErrBind, nil)
		return
	}

	user, err := u.srv.Users().Get(c, c.Param("name"))
	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	if r.Nickname != nil {
		user.NickName = *r.Nickname
	}

	if r.Email != nil {
		user.Email = *r.Email
	}

	if err = u.srv.Users().Update(c, user); err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, user)

}

func (u *UserController) Delete(c *gin.Context) {
	log.Info(c.Request.Context(), "deleted func called.", c.Request)

	if err := u.srv.Users().Delete(c, c.Param("name")); err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, errno.OK)
}
