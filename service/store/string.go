package store

import "github.com/TeaMeow/KitSvc/service/model"

func (s Store) CreateString(input string, output string) {
	s.DB.Create(&model.String{
		Input:  input,
		Output: output,
	})
}

func (s Store) GetLastString() model.String {
	var str model.String

	s.DB.Last(&str)

	return str
}
