package datastore

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/TeaMeow/KitSvc/model"
	// The mysql driver for gorm.
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type datastore struct {
	*gorm.DB
}

// setupDatabase initialize the database tables.
func setupDatabase(db *gorm.DB) {
	db.AutoMigrate(&model.User{}, &model.Permission{})
}

// cleanDatabase tear downs the database tables.
func cleanDatabase(db *gorm.DB) {
	db.DropTable(&model.User{}, &model.Permission{})
}

// Open opens a new database connection and returns a store.
func Open(user string, password string, host string, name string, charset string, parseTime bool, loc string) *datastore {
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=%s",
		user, password, host, name, charset, parseTime, loc,
	))
	if err != nil {
		logrus.Errorln(err)
		logrus.Fatalln("Database connection failed.")
	}

	cleanDatabase(db)
	setupDatabase(db)

	return &datastore{db}
}
