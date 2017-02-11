package model

import (
	"github.com/TeaMeow/KitSvc/shared/auth"
	validator "gopkg.in/go-playground/validator.v9"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username" gorm:"not null" binding:"required" validate:"min=1,max=32"`
	Password string `json:"password" gorm:"not null" binding:"required" validate:"min=8,max=128"`
}

type Token struct {
	Token string `json:"token"`
}

func (u *User) Compare(pwd string) (err error) {
	err = auth.Compare(u.Password, pwd)
	return
}

func (u *User) Encrypt() (err error) {
	u.Password, err = auth.Encrypt(u.Password)
	return
}

func (u *User) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}
