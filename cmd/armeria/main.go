package main

import (
	"armeria/internal/pkg/armeria"
	"flag"
)

func main() {
	publicPath := flag.String("public", "./client/dist", "public directory of client")
	dataPath := flag.String("data", "./data", "data directory")
	scriptsPath := flag.String("scripts", "./scripts", "scripts directory")
	httpPort := flag.Int("port", 8081, "http listen port")

	flag.Parse()

	armeria.Init(*publicPath, *dataPath, *scriptsPath, *httpPort)
}
