package arcadia

import "fmt"

// Arcadia holds the global state of the game
type Arcadia struct {
	SocketManager *SocketManager
	HTTPManager   *HTTPManager
}

// NewArcadia creates a new Arcadia server instance
func NewArcadia() *Arcadia {
	return &Arcadia{
		SocketManager: NewSocketManager(),
		HTTPManager:   NewHTTPManager(),
	}
}

// Test is a temporary test function
func (a *Arcadia) Test() {
	fmt.Println("Testing")
	a.HTTPManager.ReadyToServe()
}
