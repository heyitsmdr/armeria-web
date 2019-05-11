package armeria

import (
	"log"
	"os/exec"
)

type GameState struct {
	playerManager    *PlayerManager
	commandManager   *CommandManager
	characterManager *CharacterManager
	worldManager     *WorldManager
	publicPath       string
	dataPath         string
	scriptsPath      string
}

func Init(publicPath string, dataPath string, scriptsPath string) {
	state := &GameState{}

	state.publicPath = publicPath
	state.dataPath = dataPath
	state.scriptsPath = scriptsPath

	state.playerManager = NewPlayerManager(state)
	state.commandManager = NewCommandManager(state)
	state.characterManager = NewCharacterManager(state)
	state.worldManager = NewWorldManager(state)

	RegisterGameCommands(state)
	InitWeb(state)
}

func (gs *GameState) Save() {
	gs.characterManager.SaveCharacters()
	gs.worldManager.SaveWorld()
}

func (gs *GameState) Reload(callingPlayer *Player, component string) {
	steps := make(chan string, 2)

	steps <- "send_warning"

	go func() {
		for stepName := range steps {
			if stepName == "send_warning" {
				for _, c := range gs.characterManager.GetCharacters() {
					c.GetPlayer().clientActions.ShowText("The game server is about to go down for a restart.")
				}
				steps <- "save_world"
			} else if stepName == "save_world" {
				gs.Save()
				steps <- "start_update_script"
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
					steps <- "terminate_server"
				}
			} else if stepName == "terminate_server" {
				close(steps)
				//os.Exit(0)
			}
		}
	}()
}
