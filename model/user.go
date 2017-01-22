package model

import validator "gopkg.in/go-playground/validator.v9"

type User struct {
	ID       int
	Username string `json:"username" gorm:"not null" binding:"required" validate:"min=1,max=32"`
	Password string `json:"password" gorm:"not null" binding:"required" validate:"min=8,max=128"`
}

func (u *User) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}
