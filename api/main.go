package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func main() {
	c := NewContainer()

	if err := c.Invoke(func(db *gorm.DB) error {
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		return sqlDB.Ping()
	}); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := c.Invoke(func(e *echo.Echo) {
		e.Logger.Fatal(e.Start(":8080"))
	}); err != nil {
		panic(err)
	}
}
