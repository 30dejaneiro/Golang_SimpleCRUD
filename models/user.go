package models

type User struct {
	Username string `json:"username" gorm:"primaryKey" validate:"required"`
	Pass     string `json:"password" validate:"required"`
	IsAdmin  bool   `json:"isAdmin"`
}
