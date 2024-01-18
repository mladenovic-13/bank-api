package main

import (
	"fmt"

	_ "github.com/go-chi/chi/v5"
	_ "github.com/go-chi/cors"
	_ "github.com/google/uuid"
	"github.com/mladenovic-13/bank-api/app"
)

func main() {
	err := app.SetupAndRunApp()

	if err != nil {
		fmt.Printf("error %+v", err)
		panic("internal server error")
	}
}
