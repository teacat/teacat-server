package datastore

import "github.com/TeaMeow/KitSvc/model"

func (db *datastore) CreateUser(user *model.User) error {
	return db.Create(user).Error
}

func (db *datastore) GetUser(id int) (*model.User, error) {
	var user *model.User
	d := db.Where(&model.User{ID: id}).First(&user)
	return user, d.Error
}

func (db *datastore) DeleteUser(id int) error {
	return db.Delete(&model.User{ID: id}).Error
}

func (db *datastore) UpdateUser(user *model.User) error {
	return db.Save(&user).Error
}
