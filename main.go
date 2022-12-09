package main

import (
	"log"

	"openconnect-wrapper/internal/app"
)

func main() {
	if err := app.Main(); err != nil {
		log.Fatal(err)
	}
}
