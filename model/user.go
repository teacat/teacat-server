package model

import (
	"github.com/TeaMeow/KitSvc/shared/auth"
	validator "gopkg.in/go-playground/validator.v9"
)

// User represents a registered user.
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username" gorm:"not null" binding:"required" validate:"min=1,max=32"`
	Password string `json:"password" gorm:"not null" binding:"required" validate:"min=8,max=128"`
}

// Token represents a JSON web token.
type Token struct {
	Token string `json:"token"`
}

// Compare with the plain text password. Returns true if it's the same as the encrypted one (in the `User` struct).
func (u *User) Compare(pwd string) (err error) {
	err = auth.Compare(u.Password, pwd)
	return
}

// Encrypt the user password.
func (u *User) Encrypt() (err error) {
	u.Password, err = auth.Encrypt(u.Password)
	return
}

// Validate the fields.
func (u *User) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}
