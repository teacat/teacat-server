package model

type String struct {
	ID     int
	Input  string `gorm:"not null"`
	Output string `gorm:"not null"`
}
