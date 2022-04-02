package models

type User struct {
	Id       int    `json:"-"   db:"id"`
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

type User2 struct {
	Id       uint64 `form:"id"   db:"id"`
	Name     string `form:"name"`
	Username string `form:"username"`
	Password string `form:"password"`
	Email    string `form:"email"`
}

type Student struct {
}
