package main

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/chaitanyakolluru/go-ums-backend/pkg/controller"
	"github.com/chaitanyakolluru/go-ums-backend/pkg/model"
	"github.com/labstack/echo/v4"
)

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}

	db.AutoMigrate(&model.UserData{})
	e := echo.New()

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("db", db)
			c.Get("db").(*gorm.DB).Find(&[]model.UserData{})
			return next(c)
		}
	})

	e.GET("/healthz", controller.HealthCheck)
	e.POST("/users", controller.SaveUser)
	e.GET("/users", controller.GetUsers)
	e.GET("/users/:id", controller.GetUser)
	e.DELETE("/users/:id", controller.DeleteUser)
	e.PUT("/users/:id", controller.UpdateUser)

	e.Logger.Fatal(e.Start(":8080"))
}
