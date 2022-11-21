package main

import (
	"anyconnect-wrapper/internal/app"
	"log"
)

func main() {
	if err := app.Main(); err != nil {
		log.Fatal(err)
	}
}
