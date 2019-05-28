package armeria

import (
	"log"
	"os"
	"os/exec"
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
	publicPath       string
	dataPath         string
	objectImagesPath string
	scriptsPath      string
}

var (
	Armeria *GameState
)

func Init(production bool, publicPath string, dataPath string, scriptsPath string, httpPort int) {
	Armeria = &GameState{
		production:       production,
		publicPath:       publicPath,
		dataPath:         dataPath,
		objectImagesPath: dataPath + "/object-images",
		scriptsPath:      scriptsPath,
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("error initializing zap logger: %s", err)
	}
	Armeria.log = logger

	Armeria.playerManager = NewPlayerManager()
	Armeria.commandManager = NewCommandManager()
	Armeria.characterManager = NewCharacterManager()
	Armeria.worldManager = NewWorldManager()

	RegisterGameCommands()
	InitWeb(httpPort)
}

func (gs *GameState) Save() {
	gs.characterManager.SaveCharacters()
	gs.worldManager.SaveWorld()
}

func (gs *GameState) Reload(callingPlayer *Player, component string) {
	steps := make(chan string, 2)

	callingPlayer.clientActions.ShowText("Please wait while the requested components are updated and built..")

	steps <- "start_update_script"
	go func() {
		for stepName := range steps {
			if stepName == "send_warning" {
				for _, c := range gs.characterManager.GetCharacters() {
					c.GetPlayer().clientActions.ShowText("The game server is about to go down for a restart in 5 seconds.")
				}
				time.Sleep(5 * time.Second)
				steps <- "save_world"
			} else if stepName == "save_world" {
				gs.Save()
				steps <- "terminate_server"
			} else if stepName == "start_update_script" {
				output, err := exec.Command(gs.scriptsPath+"/update.sh", component).CombinedOutput()
				if err != nil {
					callingPlayer.clientActions.ShowText(
						"An error occurred when attempting to update. Check the logs for more info.",
					)
					log.Printf("[state] an error occurred when trying to execute update.sh: %s", err)
					close(steps)
				} else {
					callingPlayer.clientActions.ShowText(string(output))

					if component == "client" {
						close(steps)
						callingPlayer.clientActions.ShowText("The client has been updated and rebuilt. Refresh!")
					} else {
						steps <- "send_warning"
					}
				}
			} else if stepName == "terminate_server" {
				close(steps)
				cmd := exec.Command(gs.scriptsPath + "/restart.sh")
				err := cmd.Start()
				if err != nil {
					callingPlayer.clientActions.ShowText(
						"An error occurred when attempting to restart. Check the logs for more info.",
					)
					log.Printf("[state] an error occurred when trying to execute restart.sh: %s", err)
					return
				}
				os.Exit(0)
			}
		}
	}()
}
