package main

import "github.com/jinzhu/gorm"

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
