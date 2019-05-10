package main

import (
	"armeria/internal/pkg/armeria"
	"flag"
)

func main() {
	publicPath := flag.String("public", "", "public directory of client")
	dataPath := flag.String("data", "", "data directory")
	scriptsPath := flag.String("scripts", "", "scripts directory")

	flag.Parse()

	armeria.Init(*publicPath, *dataPath, *scriptsPath)
}
