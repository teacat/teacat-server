package datastore

import (
	"testing"

	"github.com/TeaMeow/KitSvc/model"
	"github.com/TeaMeow/KitSvc/store"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

var (
	db *gorm.DB
	s  store.Store
)

func TestMain(t *testing.T) {
	db = openTest()
	//defer db.Close()
	s = From(db)
}

func TestCreateUser(t *testing.T) {
	assert := assert.New(t)
	resetDatabase(db)
	u := model.User{
		Username: "YamiOdymel",
		Password: "testtest",
	}
	err := s.CreateUser(&u)

	assert.True(err == nil)
	assert.True(u.ID != 0)
}

func TestGetUser(t *testing.T) {
	assert := assert.New(t)
	resetDatabase(db)
	u := model.User{
		Username: "YamiOdymel",
		Password: "testtest",
	}

	s.CreateUser(&u)
	getuser, err := s.GetUser("YamiOdymel")
	assert.True(err == nil)
	assert.Equal(u.ID, getuser.ID)
	assert.Equal(u.Username, getuser.Username)
	// Test database functions only, so the password is plain text.
	assert.Equal(u.Password, getuser.Password)
}

func TestGetLastUser(t *testing.T) {
	assert := assert.New(t)
	resetDatabase(db)
	u := model.User{
		Username: "YamiOdymel",
		Password: "testtest",
	}

	s.CreateUser(&u)
	lastuser, err := s.GetLastUser()
	assert.True(err == nil)
	assert.Equal(u.ID, lastuser.ID)
	assert.Equal(u.Username, lastuser.Username)
	assert.Equal(u.Password, lastuser.Password)
}

func TestGetUserAfter(t *testing.T) {
	assert := assert.New(t)
	resetDatabase(db)
	u := model.User{
		Username: "YamiOdymel",
		Password: "testtest",
	}
	u2 := model.User{
		Username: "Akaria",
		Password: "adminadmin",
	}

	s.CreateUser(&u)
	s.CreateUser(&u2)
	afteruser, err := s.GetUserAfter(u.ID)
	assert.True(err == nil)
	assert.Equal(u2.ID, afteruser.ID)
	assert.Equal(u2.Username, afteruser.Username)
	assert.Equal(u2.Password, afteruser.Password)
}

func TestDeleteUser(t *testing.T) {
	assert := assert.New(t)
	resetDatabase(db)
	u := model.User{
		Username: "YamiOdymel",
		Password: "testtest",
	}

	s.CreateUser(&u)
	err := s.DeleteUser(u.ID)
	assert.True(err == nil)
	getuser, err := s.GetUser("YamiOdymel")
	assert.True(err != nil)
	assert.Equal((&model.User{}), getuser)
}

func TestUpdateUser(t *testing.T) {
	assert := assert.New(t)
	resetDatabase(db)
	u := model.User{
		Username: "YamiOdymel",
		Password: "testtest",
	}
	s.CreateUser(&u)
	u.Username = "Akaria"
	u.Password = "adminadmin"
	err := s.UpdateUser(&u)
	assert.True(err == nil)
	getuser, err := s.GetUser("Akaria")
	assert.True(err == nil)
	assert.Equal(u.ID, getuser.ID)
	assert.Equal(u.Username, getuser.Username)
	assert.Equal(u.Password, getuser.Password)
}
