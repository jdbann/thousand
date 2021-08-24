package main

import (
	"fmt"
	"log"
	"net/http"

	"emailaddress.horse/thousand/app"
)

func main() {
	app := app.NewApp()
	fmt.Println("Listening on http://localhost:4000")
	if err := http.ListenAndServe(":4000", app.Engine()); err != nil {
		log.Fatal(err)
	}
}
