package main

import (
	"armeria/internal/pkg/armeria"
	"flag"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "./config/development.yml", "path to the config file")
	flag.StringVar(&configPath, "c", "./config/development.yml", "path to the config file  (shorthand)")

	flag.Parse()

	armeria.Init(configPath)
}
