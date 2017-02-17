package datastore

import (
	"testing"

	"github.com/TeaMeow/KitSvc/model"
	"github.com/TeaMeow/KitSvc/store"
	"github.com/stretchr/testify/assert"
)

var (
	db store.Store
)

func TestMain(t *testing.T) {
	db = openTest()
	//db.DB.Close()
	//defer db.Close()
}

func TestCreateUser(t *testing.T) {
	assert := assert.New(t)

	u := model.User{
		Username: "YamiOdymel",
		Password: "testtest",
	}
	err := db.CreateUser(&u)

	assert.True(err != nil)
	assert.True(u.ID != 0)
}

func TestGetUser(t *testing.T) {

}

func TestGetLastUser(t *testing.T) {

}

func TestGetUserAfter(t *testing.T) {

}

func TestDeleteUser(t *testing.T) {

}

func TestUpdateUser(t *testing.T) {

}
