package armeria

type GameState struct {
	playerManager    *PlayerManager
	commandManager   *CommandManager
	characterManager *CharacterManager
}

func Init(publicPath string, dataPath string) {
	state := &GameState{}

	state.playerManager = NewPlayerManager(state)
	state.commandManager = NewCommandManager(state)
	state.characterManager = NewCharacterManager(state, dataPath)

	InitWeb(state, publicPath)
}

// PlayerManager returns the player manager instance, so other packages can access it
func (gs *GameState) PlayerManager() *PlayerManager {
	return gs.playerManager
}
