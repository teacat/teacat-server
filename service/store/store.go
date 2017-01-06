package store

import (
	"github.com/TeaMeow/KitSvc/service/model"
	"github.com/jinzhu/gorm"
)

type Store struct {
	DB *gorm.DB
}

func (s Store) Upstream() {
	s.DB.AutoMigrate(model.String{})
}

func (s Store) Downstream() {
	s.DB.DropTable(model.String{})
}
