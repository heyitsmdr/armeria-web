package main

import (
	"armeria/pkg/armeria"
	"armeria/pkg/web"
	"fmt"
)

func main() {
	armeria.NewServer()
	fmt.Printf("Welcome to %s", armeria.GameState.Name)
	web.Init()
}
