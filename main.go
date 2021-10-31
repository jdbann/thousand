package main

import (
	"log"
	"os"

	"emailaddress.horse/thousand/cmd"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	cliApp := cmd.BuildCLIApp()
	if err := cliApp.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
