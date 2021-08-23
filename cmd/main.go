package main

import (
	"fmt"
	"log"
	"net/http"

	"emailaddress.horse/thousand/app"
)

func main() {
	app := &app.App{}
	fmt.Println("Listening on http://localhost:4000")
	if err := http.ListenAndServe(":4000", app.Routes()); err != nil {
		log.Fatal(err)
	}
}
