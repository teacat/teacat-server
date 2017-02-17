package datastore

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/TeaMeow/KitSvc/model"
	"github.com/TeaMeow/KitSvc/module/logger"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
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

// From returns a Store using an existing database connection.
//func From(db *gorm.DB) store.Store {
//	return &datastore{db}
//}

// Open opens a new database connection and returns a store.
func Open(driver, config string) *datastore {
	db, err := gorm.Open(driver, config)
	if err != nil {
		logger.FatalFields("Database connection failed.", logrus.Fields{
			"err": err,
		})
	}

	cleanDatabase(db)
	setupDatabase(db)

	return &datastore{db}
}

func openTest() *datastore {
	driver := "sqlite3"
	config := ":memory:"
	if os.Getenv("KITSVC_DATABASE_DRIVER") != "" {
		driver = os.Getenv("KITSVC_DATABASE_DRIVER")
		config = os.Getenv("KITSVC_DATABASE_CONFIG")
	}
	return Open(driver, config)
}
