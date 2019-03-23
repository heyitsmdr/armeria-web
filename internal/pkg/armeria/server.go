package armeria

var GameState *gameState

type gameState struct {
	Name string
}

// NewServer creates a new instance of the game server
func NewServer() {
	GameState = &gameState{
		Name: "Armeriaaa",
	}
}
