package main

import (
	"fmt"
	"net/http"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/labstack/echo/v4"
)

type UserData struct {
	gorm.Model
	User `json:"user"`
}

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

var users []UserData

func saveUser(c echo.Context) error {
	u := new(User)
	if err := c.Bind(u); err != nil {
		return err
	}

	c.Get("db").(*gorm.DB).Find(&users)
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

	c.Get("db").(*gorm.DB).Find(&users)
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

	c.Get("db").(*gorm.DB).Find(&users)
	for i, user := range users {
		if id == fmt.Sprintf("%d", user.ID) {
			c.Get("db").(*gorm.DB).Model(&users[i]).Updates(UserData{User: *u})
			return c.JSON(http.StatusOK, users[i])
		}
	}

	return c.JSON(http.StatusNotFound, "user not found")
}

func getUsers(c echo.Context) error {
	c.Get("db").(*gorm.DB).Find(&users)
	return c.JSON(http.StatusOK, users)
}

func getUser(c echo.Context) error {
	id := c.Param("id")

	c.Get("db").(*gorm.DB).Find(&users)
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

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}

	db.AutoMigrate(&UserData{})
	e := echo.New()

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("db", db)
			return next(c)
		}
	})

	e.GET("/healthz", healthCheck)
	e.POST("/users", saveUser)
	e.GET("/users", getUsers)
	e.GET("/users/:id", getUser)
	e.DELETE("/users/:id", deleteUser)
	e.PUT("/users/:id", updateUser)

	e.Logger.Fatal(e.Start(":8080"))
}
