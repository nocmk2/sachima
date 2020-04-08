package server

// User demo
type User struct {
	UserName  string `form:"user" json:"username" binding:"required"`
	Password  string `form:"password" json:"password" binding:"required"`
	FirstName string
	LastName  string
	Email     string `json:"email"`
}
