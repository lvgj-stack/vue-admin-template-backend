package v1

import (
	"github.com/go-playground/validator/v10"

	"github.com/Mr-LvGJ/jota/models"
)

type User struct {
	models.BaseModel

	Username string `json:"username"           gorm:"column:username;not null" binding:"required" validate:"min=1,max=32"`
	Password string `json:"password,omitempty" gorm:"column:password;not null" binding:"required" validate:"min=5,max=128"`
	NickName string `json:"nickname"           gorm:"column:nickname"          binding:"required" validate:"required,min=1,max=30"`
	Email    string `json:"email"              gorm:"column:email"             binding:"required" validate:"min=1,max=100"`
	Avatar   string `json:"avatar"             gorm:"column:avatar"`
}

func (u *User) TableName() string {
	return "user"
}

func (u *User) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}

type UserList struct {
	Items []*User `json:"items"`
}
