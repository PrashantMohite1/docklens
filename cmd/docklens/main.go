package main

import (
	"log"

	"github.com/PrashantMohite1/docklens/internal/cli"
	"github.com/PrashantMohite1/docklens/internal/logger"
)

func main() {
	logger.Init()
	if err := cli.Execute(); err != nil {
		log.Fatal(err)
	}
}
