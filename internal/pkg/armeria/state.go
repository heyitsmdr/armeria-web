package armeria

import (
	"log"

	"go.uber.org/zap"
)

type GameState struct {
	log              *zap.Logger
	production       bool
	playerManager    *PlayerManager
	commandManager   *CommandManager
	characterManager *CharacterManager
	worldManager     *WorldManager
	mobManager       *MobManager
	itemManager      *ItemManager
	publicPath       string
	dataPath         string
	objectImagesPath string
}

var (
	Armeria *GameState
)

func Init(production bool, publicPath string, dataPath string, httpPort int) {
	Armeria = &GameState{
		production:       production,
		publicPath:       publicPath,
		dataPath:         dataPath,
		objectImagesPath: dataPath + "/object-images",
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("error initializing zap logger: %s", err)
	}
	Armeria.log = logger

	Armeria.commandManager = NewCommandManager()
	Armeria.playerManager = NewPlayerManager()
	Armeria.characterManager = NewCharacterManager()
	Armeria.worldManager = NewWorldManager()
	Armeria.mobManager = NewMobManager()
	Armeria.itemManager = NewItemManager()

	RegisterGameCommands()
	InitWeb(httpPort)
}

func (gs *GameState) Save() {
	gs.characterManager.SaveCharacters()
	gs.worldManager.SaveWorld()
	gs.mobManager.SaveMobs()
	gs.itemManager.SaveItems()
}
