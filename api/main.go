package main

import "github.com/labstack/echo/v4"

func main() {
	c := NewContainer()

	if err := c.Invoke(func(e *echo.Echo) {
		e.Logger.Fatal(e.Start(":8080"))
	}); err != nil {
		panic(err)
	}
}
