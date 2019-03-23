package main

import (
	"armeria/internal/pkg/armeria"
	"armeria/internal/pkg/web"
	"fmt"
)

func main() {
	armeria.NewServer()
	fmt.Printf("Welcome to %s", armeria.GameState.Name)
	web.Init()
}
