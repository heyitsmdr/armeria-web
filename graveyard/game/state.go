package game

import (
	"armeria/internal/pkg/characters"
	schemaCharacters "armeria/internal/pkg/characters/schema"
	"armeria/internal/pkg/commands"
	schemaCommands "armeria/internal/pkg/commands/schema"
	"armeria/internal/pkg/players"
	schemaPlayers "armeria/internal/pkg/players/schema"
	"armeria/internal/pkg/web"
)

type gameState struct {
	playerManager  schemaPlayers.IPlayerManager
	commandManager schemaCommands.ICommandManager
	characterManager schemaCharacters.ICharacterManager
}

func Init(publicPath string, dataPath string) {
	gs := &gameState{}

	gs.playerManager = players.Init(gs)
	gs.commandManager = commands.Init(gs)
	gs.characterManager = characters.Init(gs, dataPath)

	// Initialize the web server last since it will start accepting player connections
	web.Init(gs, publicPath)
}

// PlayerManager returns the player manager singleton
func (gs gameState) PlayerManager() schemaPlayers.IPlayerManager {
	return gs.playerManager
}

// CommandManager returns the command manager singleton
func (gs gameState) CommandManager() schemaCommands.ICommandManager {
	return gs.commandManager
}

// CharacterManager returns the character manager singleton
func (gs gameState) CharacterManager() schemaCharacters.ICharacterManager {
	return gs.characterManager
}