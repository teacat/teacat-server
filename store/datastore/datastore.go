package datastore

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/TeaMeow/KitSvc/model"
	"github.com/TeaMeow/KitSvc/module/logger"
	"github.com/TeaMeow/KitSvc/store"
	// MySQL driver.
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	// SQLite driver.
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type datastore struct {
	*gorm.DB
}

// setupDatabase initialize the database tables.
func setupDatabase(db *gorm.DB) {
	db.AutoMigrate(&model.User{})
}

// cleanDatabase tear downs the database tables.
func cleanDatabase(db *gorm.DB) {
	db.DropTable(&model.User{})
}

// resetDatabase resets the database tables.
func resetDatabase(db *gorm.DB) {
	cleanDatabase(db)
	setupDatabase(db)
}

// From returns a Store using an existing database connection.
func From(db *gorm.DB) store.Store {
	return &datastore{db}
}

// New implement the Store interface with the database connection.
func New(driver, config string) store.Store {
	return From(
		open(driver, config),
	)
}

// open a new database connection and returns a store.
func open(driver, config string) (db *gorm.DB) {
	db, err := gorm.Open(driver, config)
	if err != nil {
		logger.FatalFields("Database connection failed.", logrus.Fields{
			"err": err,
		})
	}
	resetDatabase(db)
	return
}

func openTest() *gorm.DB {
	driver := "sqlite3"
	config := ":memory:"
	if os.Getenv("KITSVC_DATABASE_DRIVER") != "" {
		driver = os.Getenv("KITSVC_DATABASE_DRIVER")
		config = os.Getenv("KITSVC_DATABASE_CONFIG")
	}
	return open(driver, config)
}
