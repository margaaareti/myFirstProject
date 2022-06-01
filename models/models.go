package models

import (
	"github.com/pkg/errors"
	"time"
)

type User struct {
	Id       int    `json:"-"   db:"id"`
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

type User2 struct {
	Id         uint64 `form:"id"   db:"id"`
	Name       string `form:"name"`
	Surname    string `form:"surname"`
	Patronymic string `form:"patronymic"`
	Username   string `form:"username"`
	Password   string `form:"password"`
	Email      string `form:"email"`
}

type Student struct {
	Id          uint64    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name" binding:"required"`
	Surname     string    `json:"surname" db:"surname" binding:"required"`
	Patronymic  string    `json:"patronymic" db:"patronymic"`
	IsuNumber   string    `json:"number" db:"isu_number" binding:"required"`
	AddedBy     string    `json:"added-by" db:"added_by"`
	Title       string    `json:"title" db:"title" binding:"required"`
	Description string    `json:"description" db:"description"`
	Time        time.Time `json:"time" db:"reg_date"`
}

type UpdateStudentInput struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	//Done        *bool   `json:"done"`
}

func (i UpdateStudentInput) Validate() error {
	if i.Title == nil && i.Description == nil /*&& i.Done == nil*/ {
		return errors.New("update structures has no value")
	}
	return nil
}
