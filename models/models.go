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
	Id          uint64 `json:"id" db:"id"`
	Name        string `json:"name" db:"name" binding:"required"`
	Surname     string `json:"surname" db:"surname" binding:"required"`
	Patronymic  string `json:"patronymic" db:"patronymic"`
	IsuNumber   string `json:"number" db:"isu_number" binding:"required"`
	AddedBy     string `json:"added-by" db:"added_by"`
	Title       string `json:"title" db:"title" binding:"required"`
	Description string `json:"description" db:"description"`

	Time string `json:"time" db:"reg_date"`
}
