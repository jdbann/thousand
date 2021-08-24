package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"emailaddress.horse/thousand/app"
)

func main() {
	app := app.NewApp()

	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	addr := fmt.Sprintf(":%s", port)
	if err := http.ListenAndServe(addr, app.Engine()); err != nil {
		log.Fatal(err)
	}
}
