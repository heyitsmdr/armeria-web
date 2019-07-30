package armeria

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

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

func Init(fileLocation string) {
	c := parseConfigFile(fileLocation)
	Armeria = &GameState{
		production:       c.Production,
		publicPath:       c.PublicPath,
		dataPath:         c.DataPath,
		objectImagesPath: c.DataPath + "/object-images",
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

	Armeria.setupGracefulExit()
	Armeria.setupPeriodicSaves()

	RegisterGameCommands()
	InitWeb(c.HTTPPort)
}

func (gs *GameState) setupGracefulExit() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM)
	go func() {
		<-sigs
		gs.Save()
		os.Exit(0)
	}()
}

func (gs *GameState) setupPeriodicSaves() {
	c := time.Tick(2 * time.Minute)
	go func() {
		for range c {
			gs.Save()
		}
	}()
}

func (gs *GameState) Save() {
	gs.characterManager.SaveCharacters()
	gs.worldManager.SaveWorld()
	gs.mobManager.SaveMobs()
	gs.itemManager.SaveItems()
}
