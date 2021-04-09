package armeria

import (
	"armeria/internal/pkg/cloud"
	"armeria/internal/pkg/github"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"go.uber.org/zap"
)

// GameState stores the manager singletons and any other global state.
type GameState struct {
	log              *zap.Logger
	production       bool
	storageManager   *cloud.StorageManager
	playerManager    *PlayerManager
	commandManager   *CommandManager
	characterManager *CharacterManager
	worldManager     *WorldManager
	mobManager       *MobManager
	itemManager      *ItemManager
	convoManager     *ConversationManager
	ledgerManager    *LedgerManager
	tickManager      *TickManager
	registry         *Registry
	channels         map[string]*Channel
	publicPath       string
	dataPath         string
	objectImagesPath string
	startTime        time.Time
	github           *github.ArmeriaRepo
}

var (
	// Armeria contains the global state for the game server.
	Armeria *GameState
)

// Init loads the Armeria game server and starts serving requests.
func Init(configFilePath string, serveTraffic bool) {
	c := parseConfigFile(configFilePath)

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

	if !serveTraffic {
		return
	}

	verifySchemaVersion()

	Armeria.storageManager = cloud.NewStorageManager(c.GCSBucket, c.GCSServiceAccount)
	Armeria.registry = NewRegistry()
	Armeria.commandManager = NewCommandManager()
	Armeria.playerManager = NewPlayerManager()
	Armeria.characterManager = NewCharacterManager()
	Armeria.worldManager = NewWorldManager()
	Armeria.mobManager = NewMobManager()
	Armeria.itemManager = NewItemManager()
	Armeria.channels = NewChannels()
	Armeria.convoManager = NewConversationManager()
	Armeria.ledgerManager = NewLedgerManager()
	Armeria.tickManager = NewTickManager()

	Armeria.github = github.New()

	Armeria.setupGracefulExit()

	Armeria.startTime = time.Now()

	RegisterGameCommands()

	port := c.HTTPPort
	// For Heroku, we must listen on a specific port.
	if len(os.Getenv("PORT")) > 0 {
		port, err = strconv.Atoi(os.Getenv("PORT"))
		if err != nil {
			log.Fatalf("error parsing PORT environment variable: %s", err)
		}
	}
	InitWeb(port)
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

// Save writes the in-memory data to disk.
func (gs *GameState) Save() {
	gs.characterManager.SaveCharacters()
	gs.worldManager.SaveWorld()
	gs.mobManager.SaveMobs()
	gs.itemManager.SaveItems()
	gs.ledgerManager.SaveLedgers()
	gs.storageManager.CloseClient()
}
