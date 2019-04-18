package characters

import (
	schemaGame "armeria/internal/pkg/game/schema"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type manager struct {
	gameState  schemaGame.IGameState
	dataFile   string
	Characters []*Character `json:"characters"`
}

// Init creates a new character Manager instance
func Init(gs schemaGame.IGameState, dataPath string) *manager {
	m := &manager{
		gameState: gs,
		dataFile:  fmt.Sprintf("%s/characters.json", dataPath),
	}

	m.loadCharacters()
	m.registerCommands()

	return m
}

func (m *manager) loadCharacters() {
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

func (m *manager) GetCharacterByName(name string) *Character {
	for _, c := range m.Characters {
		if c.Name == name {
			return c
		}
	}

	return nil
}
