package store

import (
	"github.com/TeaMeow/KitSvc/service/model"
	"github.com/jinzhu/gorm"
)

// Store stores the database connection and the operations to process with the database.
type Store struct {
	DB *gorm.DB
}

// Upstream builds the database tables when the service is started.
func (s Store) Upstream() {
	s.DB.AutoMigrate(model.String{})
}

// Downstream tears down the database tables when the service is going offline.
func (s Store) Downstream() {
	s.DB.DropTable(model.String{})
}
