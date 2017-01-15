package datastore

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/TeaMeow/KitSvc/model"
	"github.com/jinzhu/gorm"
)

type datastore struct {
	*gorm.DB
}

// setupDatabase initialize the database tables.
func setupDatabase(db *gorm.DB) {
	db.AutoMigrate(model.User{})
}

// cleanDatabase tear downs the database tables.
func cleanDatabase(db *gorm.DB) {
	db.DropTable(model.User{})
}

// Open opens a new database connection and returns a store.
func Open() *datastore {
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=%s&loc=%s",
		os.Getenv("KITSVC_DATABASE_USER"),
		os.Getenv("KITSVC_DATABASE_PASSWORD"),
		os.Getenv("KITSVC_DATABASE_HOST"),
		os.Getenv("KITSVC_DATABASE_NAME"),
		os.Getenv("KITSVC_DATABASE_CHARSET"),
		os.Getenv("KITSVC_DATABASE_PARSE_TIME"),
		os.Getenv("KITSVC_DATABASE_LOC"),
	))
	if err != nil {
		logrus.Errorln(err)
		logrus.Fatalln("Database connection failed.")
	}

	setupDatabase(db)
	cleanDatabase(db)

	return &datastore{db}
}
