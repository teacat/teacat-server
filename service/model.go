package main

import (
	"strings"

	"github.com/jinzhu/gorm"
)

type Model struct {
	DB *gorm.DB
}

func createModel(db *gorm.DB) Model {
	return Model{db}
}

func (m Model) ToUpper(s string) (string, error) {
	if s == "" {
		return "", Err{
			Message: ErrEmpty,
		}
	}

	return strings.ToUpper(s), nil
}

func (m Model) Count(s string) int {
	return len(s)
}
