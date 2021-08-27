package main

import (
	"fmt"
	"log"
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
	if err := app.Start(addr); err != nil {
		log.Fatal(err)
	}
}
