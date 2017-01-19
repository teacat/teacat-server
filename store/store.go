package store

import (
	"github.com/TeaMeow/KitSvc/model"
	"github.com/gin-gonic/gin"
)

type Store interface {
	// CreateUser creates a new user account.
	CreateUser(*model.User) error

	// GetUser gets an user by the user id.
	GetUser(string) (*model.User, error)

	// DeleteUser deletes the user by the user id.
	DeleteUser(int) error

	// UpdateUser updates an user account.
	UpdateUser(*model.User) error

	//
	Can(*model.Permission) bool
}

// CreateUser creates a new user account.
func CreateUser(c *gin.Context, user *model.User) error {
	return FromContext(c).CreateUser(user)
}

// GetUser gets an user by the user id.
func GetUser(c *gin.Context, username string) (*model.User, error) {
	return FromContext(c).GetUser(username)
}

// DeleteUser deletes the user by the user id.
func DeleteUser(c *gin.Context, id int) error {
	return FromContext(c).DeleteUser(id)
}

// UpdateUser updates an user account.
func UpdateUser(c *gin.Context, user *model.User) error {
	return FromContext(c).UpdateUser(user)
}

func Can(c *gin.Context, p *model.Permission) bool {
	return FromContext(c).Can(p)
}
