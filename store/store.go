package store

import (
	"github.com/TeaMeow/KitSvc/model"
	"github.com/gin-gonic/gin"
)

// Store wraps the functions that interactive with the database,
// just like the Model in MVC architecture.
type Store interface {
	CreateUser(*model.User) error
	GetUser(string) (*model.User, error)
	DeleteUser(int) error
	UpdateUser(*model.User) error
}

// CreateUser creates a new user account.
func CreateUser(c *gin.Context, user *model.User) error {
	return FromContext(c).CreateUser(user)
}

// GetUser gets an user by the user identifier.
func GetUser(c *gin.Context, username string) (*model.User, error) {
	return FromContext(c).GetUser(username)
}

// DeleteUser deletes the user by the user identifier.
func DeleteUser(c *gin.Context, id int) error {
	return FromContext(c).DeleteUser(id)
}

// UpdateUser updates an user account information.
func UpdateUser(c *gin.Context, user *model.User) error {
	return FromContext(c).UpdateUser(user)
}
