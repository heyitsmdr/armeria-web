package main

import (
	"armeria/internal/pkg/players"
	"armeria/internal/pkg/web"
	"flag"
)

func main() {
	publicPath := flag.String("public", "", "public directory of client")
	flag.Parse()

	players.Init()

	// Initialize the web server last since it will start accepting player connections
	web.Init(*publicPath)
}
