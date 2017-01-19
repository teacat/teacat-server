package datastore

import "github.com/TeaMeow/KitSvc/model"

//
func (db *datastore) CreateUser(u *model.User) error {
	return db.Create(&u).Error
}

//
func (db *datastore) GetUser(username string) (*model.User, error) {
	u := &model.User{}
	d := db.Where(&model.User{Username: username}).First(&u)
	return u, d.Error
}

//
func (db *datastore) DeleteUser(id int) error {
	return db.Delete(&model.User{ID: id}).Error
}

//
func (db *datastore) UpdateUser(u *model.User) error {
	return db.Save(&u).Error
}
