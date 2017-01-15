package model

// User
type User struct {
	ID       int
	Username string `gorm:"not null"`
	Password string `gorm:"not null"`
}
