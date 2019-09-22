package armeria

import (
	"armeria/internal/pkg/github"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

// GameState stores the manager singletons and any other global state.
type GameState struct {
	log              *zap.Logger
	production       bool
	playerManager    *PlayerManager
	commandManager   *CommandManager
	characterManager *CharacterManager
	worldManager     *WorldManager
	mobManager       *MobManager
	itemManager      *ItemManager
	convoManager     *ConversationManager
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

	Armeria.registry = NewRegistry()
	Armeria.commandManager = NewCommandManager()
	Armeria.playerManager = NewPlayerManager()
	Armeria.characterManager = NewCharacterManager()
	Armeria.worldManager = NewWorldManager()
	Armeria.mobManager = NewMobManager()
	Armeria.itemManager = NewItemManager()
	Armeria.channels = NewChannels()
	Armeria.convoManager = NewConversationManager()

	Armeria.github = github.New()

	Armeria.setupGracefulExit()
	Armeria.setupPeriodicSaves()
	Armeria.setupAncillaryTasks()

	Armeria.startTime = time.Now()

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

func (gs *GameState) setupAncillaryTasks() {
	searchForDanglingInstances()

	c := time.Tick(1 * time.Hour)
	go func() {
		for range c {
			searchForDanglingInstances()
		}
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

// Save writes the in-memory data to disk.
func (gs *GameState) Save() {
	gs.characterManager.SaveCharacters()
	gs.worldManager.SaveWorld()
	gs.mobManager.SaveMobs()
	gs.itemManager.SaveItems()
}
