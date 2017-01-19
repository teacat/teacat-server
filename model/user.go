package model

type User struct {
	ID       int
	Username string `json:"username" gorm:"not null" binding:"required"`
	Password string `json:"password" gorm:"not null" binding:"required"`
}
