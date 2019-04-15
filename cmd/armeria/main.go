package main

import (
	"armeria/internal/pkg/game"
	"flag"
)

func main() {
	publicPath := flag.String("public", "", "public directory of client")
	dataPath := flag.String("data", "", "data directory")
	flag.Parse()

	game.Init(*publicPath, *dataPath)
}