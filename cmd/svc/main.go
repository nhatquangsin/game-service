package main

import (
	"log"
	"os"

	// Automatically set GOMAXPROCS to match Linux container CPU quota.
	_ "go.uber.org/automaxprocs"
)

func main() {
	app := newApp()
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
