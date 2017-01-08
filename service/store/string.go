package store

import "github.com/TeaMeow/KitSvc/service/model"

// CreateString creates the string record.
func (s Store) CreateString(input string, output string) {
	s.DB.Create(&model.String{
		Input:  input,
		Output: output,
	})
}

// GetLastString returns the last string record.
func (s Store) GetLastString() model.String {
	var str model.String
	s.DB.Last(&str)

	return str
}
