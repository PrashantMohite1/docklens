package main

import (
	"log"

	"github.com/PrashantMohite1/docklens/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		log.Fatal(err)
	}
}
