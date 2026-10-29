package main

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
)

func main() {
	db, error := ConnectDatabase()

	if error != nil {
		fmt.Println(error, db)
	}

	//initialize fiber
	app := fiber.New()

	//route activate
	studentRoute(app, db)

	//start server
	if error := app.Listen(":3000"); error != nil {
		fmt.Println(error)
	}
}
