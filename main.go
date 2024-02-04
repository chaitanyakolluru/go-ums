package main

import (
	"fmt"

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
