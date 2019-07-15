package main

import (
	"armeria/internal/pkg/armeria"
	"flag"
)

func main() {
	publicPath := flag.String("public", "./client/dist", "public directory of client")
	dataPath := flag.String("data", "./data", "data directory")
	httpPort := flag.Int("port", 8081, "http listen port")
	prodFlag := flag.Bool("prod", false, "sets production flag")

	flag.Parse()

	armeria.Init(*prodFlag, *publicPath, *dataPath, *httpPort)
}
