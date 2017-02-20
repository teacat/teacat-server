package client

import "github.com/TeaMeow/KitSvc/model"

// Client is used to communicate with the service.
type Client interface {
	// PostUser creates a new user account.
	PostUser(*model.User) (*model.User, error)

	// GetUser gets an user by the user identifier.
	GetUser(string) (*model.User, error)

	// PutUser updates an user account information.
	PutUser(int, *model.User) (*model.User, error)

	// DeleteUser deletes the user by the user identifier.
	DeleteUser(int) error

	// PostToken generates the authentication token
	// if the password was matched with the specified account.
	PostToken(*model.User) (*model.Token, error)
}
