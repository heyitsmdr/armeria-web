package main

import (
	"armeria/internal/pkg/armeria"
	"flag"
)

func main() {
	configPath := flag.String("config", "./config/development.yml", "path to the config file")

	flag.Parse()

	armeria.Init(*configPath)
}
