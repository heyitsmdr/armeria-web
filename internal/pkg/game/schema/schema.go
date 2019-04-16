package schema

import (
	schemaCharacters "armeria/internal/pkg/characters/schema"
	schemaCommands "armeria/internal/pkg/commands/schema"
	schemaPlayers "armeria/internal/pkg/players/schema"
)

type IGameState interface {
	PlayerManager() schemaPlayers.IPlayerManager
	CommandManager() schemaCommands.ICommandManager
	CharacterManager() schemaCharacters.ICharacterManager
}