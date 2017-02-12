package datastore

import "github.com/TeaMeow/KitSvc/model"

// CreateUser creates a new user account.
func (db *datastore) CreateUser(u *model.User) error {
	return db.Create(&u).Error
}

// GetUser gets an user by the user identifier.
func (db *datastore) GetUser(username string) (*model.User, error) {
	u := &model.User{}
	d := db.Where(&model.User{Username: username}).First(&u)
	return u, d.Error
}

// DeleteUser deletes the user by the user identifier.
func (db *datastore) DeleteUser(id int) error {
	return db.Delete(&model.User{ID: id}).Error
}

// UpdateUser updates an user account information.
func (db *datastore) UpdateUser(u *model.User) error {
	return db.Save(&u).Error
}
