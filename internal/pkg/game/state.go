package game

import (
	"armeria/internal/pkg/characters"
	"armeria/internal/pkg/commands"
	"armeria/internal/pkg/players"
	"armeria/internal/pkg/web"
)

var GameState *state

type state struct {
	PlayerManager *players.Manager
	CommandManager *commands.Manager
	CharacterManager *characters.Manager
}

func Init(publicPath string, dataPath string) {
	GameState := &state{
		PlayerManager: players.Init(),
		CommandManager: commands.Init(),
		CharacterManager: characters.Init(dataPath),
	}


	// Initialize the web server last since it will start accepting player connections
	web.Init(publicPath)
}