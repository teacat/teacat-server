package store

import "github.com/TeaMeow/KitSvc/model"

type Store interface {
	// CreateUser creates a new user account.
	CreateUser(*model.User) error

	// GetUser gets an user by the user id.
	GetUser(int) (*model.User, error)

	// DeleteUser deletes the user by the user id.
	DeleteUser(int) error

	// UpdateUser updates an user account.
	UpdateUser(*model.User) error
}
