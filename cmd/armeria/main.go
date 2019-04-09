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
	web.Init(*publicPath)
}
