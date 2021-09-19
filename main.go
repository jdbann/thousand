package main

import (
	"log"
	"os"

	"emailaddress.horse/thousand/cmd"
)

func main() {
	cliApp := cmd.BuildCLIApp()
	if err := cliApp.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
