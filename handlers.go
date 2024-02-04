package main

import (
	"errors"
	"fmt"
	"net/http"

	"gorm.io/gorm"

	"github.com/labstack/echo/v4"
)

func (u *UserData) BeforeCreate(tx *gorm.DB) (err error) {
	if u.User.Name == "" {
		return errors.New("name is required")
	}

	if u.User.Email == "" {
		return errors.New("email is required")
	}
	return
}

func (u *UserData) BeforeUpdate(tx *gorm.DB) (err error) {
	if u.User.Name == "" {
		return errors.New("name is required")
	}

	if u.User.Email == "" {
		return errors.New("email is required")
	}
	return
}

func (u *UserData) AfterSave(tx *gorm.DB) (err error) {
	tx.Find(&users)
	return
}

func (u *UserData) AfterDelete(tx *gorm.DB) (err error) {
	tx.Find(&users)
	return
}

func saveUser(c echo.Context) error {
	u := new(User)
	if err := c.Bind(u); err != nil {
		return err
	}

	for _, user := range users {
		if user.User.Name == u.Name {
			return c.JSON(http.StatusConflict, "user already exists")
		}
	}

	c.Get("db").(*gorm.DB).Create(&UserData{User: *u})
	return c.JSON(http.StatusCreated, u)
}

func deleteUser(c echo.Context) error {
	id := c.Param("id")

	for i, user := range users {
		if id == fmt.Sprintf("%d", user.ID) {
			c.Get("db").(*gorm.DB).Delete(&users[i], id)
			return c.JSON(http.StatusOK, id)
		}
	}

	return c.JSON(http.StatusNotFound, "user not found")
}

func updateUser(c echo.Context) error {
	id := c.Param("id")
	u := new(User)
	if err := c.Bind(u); err != nil {
		return err
	}

	for i, user := range users {
		if id == fmt.Sprintf("%d", user.ID) {
			c.Get("db").(*gorm.DB).Model(&users[i]).Updates(UserData{User: *u})
			return c.JSON(http.StatusOK, users[i])
		}
	}

	return c.JSON(http.StatusNotFound, "user not found")
}

func getUsers(c echo.Context) error {
	return c.JSON(http.StatusOK, users)
}

func getUser(c echo.Context) error {
	id := c.Param("id")

	for _, u := range users {
		if id == fmt.Sprintf("%d", u.ID) {
			c.Get("db").(*gorm.DB).First(&u, id)
			return c.JSON(http.StatusOK, u)
		}
	}

	return c.JSON(http.StatusNotFound, "user not found")
}

func healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, "status: ok")
}
