package main

import (
	"log"

	"github.com/mladenovic-13/bank-api/app"
)

// @title Bank API
// @version 1.0
// @contact.name API Support
// @contact.url http://www.mladenovic13.com
// @contact.email mladenovic13.dev@gmail.com
// @BasePath /api/v1
func main() {
	err := app.Run()

	if err != nil {
		log.Println(err.Error())
		panic("Internal server error")
	}
}
