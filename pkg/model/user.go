package model

import (
	"gorm.io/gorm"
	"errors"
)

type UserData struct {
	gorm.Model
	User `json:"user"`
}

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

var Users []UserData

func validateParams(u *UserData) (err error) {
	if u.User.Name == "" {
		return errors.New("name is required")
	}

	if u.User.Email == "" {
		return errors.New("email is required")
	}
	return
}

func (u *UserData) BeforeCreate(tx *gorm.DB) (err error) {
	return validateParams(u)
}

func (u *UserData) BeforeUpdate(tx *gorm.DB) (err error) {
	return validateParams(u)
}

func (u *UserData) AfterSave(tx *gorm.DB) (err error) {
	tx.Find(&Users)
	return
}

func (u *UserData) AfterDelete(tx *gorm.DB) (err error) {
	tx.Find(&Users)
	return
}

