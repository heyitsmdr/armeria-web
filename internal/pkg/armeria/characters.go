package armeria

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type CharacterManager struct {
	gameState  *GameState
	dataFile   string
	Characters []*Character `json:"characters"`
}

type Character struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func NewCharacterManager(state *GameState, dataPath string) *CharacterManager {
	m := &CharacterManager{
		gameState: state,
		dataFile:  fmt.Sprintf("%s/characters.json", dataPath),
	}

	m.loadCharacters()
	m.registerCommands()

	return m
}

func (m *CharacterManager) loadCharacters() {
	charactersFile, err := os.Open(m.dataFile)
	defer charactersFile.Close()

	if err != nil {
		log.Fatalf("[characters] failed to load from %s: %s", m.dataFile, err)
	}

	jsonParser := json.NewDecoder(charactersFile)

	err = jsonParser.Decode(m)
	if err != nil {
		log.Fatalf("[characters] failed to decode file: %s", err)
	}

	log.Printf("[characters] loaded %d characters from file", len(m.Characters))
}

func (m *CharacterManager) GetCharacterByName(name string) *Character {
	for _, c := range m.Characters {
		if c.Name == name {
			return c
		}
	}

	return nil
}

func (m *CharacterManager) registerCommands() {
	m.gameState.commandManager.RegisterCommand(&Command{
		Name:       "login",
		SyntaxHelp: "/login [character] [password]",
		Handler:    handleLoginCommand,
	})
}

func handleLoginCommand(r *CommandRequest) {
	if len(r.Args) != 2 {
		r.Player.clientActions.ShowText(fmt.Sprintf("[b]Syntax:[/b] %s", r.Command.SyntaxHelp))
		return
	}

	//character := r.Args[0]
	//password := r.Args[1]

	r.Player.clientActions.ShowText("Trying to login!")
}
