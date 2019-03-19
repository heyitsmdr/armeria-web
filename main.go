package main

import (
	"arcadia/server"
	"fmt"
)

func main() {
	fmt.Println("Welcome to Arcadia.")
	a := arcadia.NewArcadia()
	a.Test()
}
