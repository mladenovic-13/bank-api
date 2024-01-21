package main

import (
	"fmt"

	"github.com/mladenovic-13/bank-api/app"
)

func main() {
	err := app.Run()

	if err != nil {
		fmt.Printf("%+v\n", err)
		panic("internal server error")
	}
}
