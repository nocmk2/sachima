package server

import "github.com/jinzhu/gorm"

// User demo
type User struct {
	UserName  string `form:"user" json:"username" binding:"required"`
	Password  string `form:"password" json:"password" binding:"required"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	gorm.Model
}
