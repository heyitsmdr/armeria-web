package main

import (
	"armeria/internal/pkg/armeria"
	"flag"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "no Path", "path to the config file")
	flag.StringVar(&configPath, "c", "no Path", "path to the config file  (shorthand)")

	flag.Parse()

	armeria.Init(configPath)
}
