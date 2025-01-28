package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/yasv98/movies-api/cmd/runner"
)

var configPath = flag.String("configPath", "config/config.yaml", "path to config file")

func main() {
	flag.Parse()
	if err := runner.Run(context.Background(), *configPath); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
