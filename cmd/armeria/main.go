package main

import (
	"armeria/internal/pkg/sockets"
	"armeria/internal/pkg/web"
	"flag"
)

func main() {
	publicPath := flag.String("public", "", "public directory of client")
	flag.Parse()

	web.Init(*publicPath)
	sockets.Init()
}
