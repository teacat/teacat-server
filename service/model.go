package main

import (
	"strings"

	"github.com/jinzhu/gorm"
)

// ToUpper converts the string to uppercase.
func (m Model) ToUpper(s string) (string, error) {
	if s == "" {
		return "", Err{
			Message: ErrEmpty,
		}
	}

	return strings.ToUpper(s), nil
}

// Count counts the length of the string.
func (m Model) Count(s string) int {
	return len(s)
}

// The functions, structs down below are the core methods,
// you shouldn't edit them until you know what you're doing,
// or you understand how KitSvc works.
//
// Or if you are brave enough ;)

// Model represents the model layer of the service.
type Model struct {
	DB *gorm.DB
}

// createModel creates the model of the service with the database connection.
func createModel(db *gorm.DB) Model {
	return Model{db}
}
