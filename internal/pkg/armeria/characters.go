package armeria

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

type CharacterManager struct {
	gameState  *GameState
	dataFile   string
	Characters []*Character `json:"characters"`
}

func NewCharacterManager(state *GameState, dataPath string) *CharacterManager {
	m := &CharacterManager{
		gameState: state,
		dataFile:  fmt.Sprintf("%s/characters.json", dataPath),
	}

	m.loadCharacters()

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

	for _, c := range m.Characters {
		c.Init(m.gameState)
	}

	log.Printf("[characters] loaded %d characters from file", len(m.Characters))
}

func (m *CharacterManager) GetCharacterByName(name string) (*Character, error) {
	for _, c := range m.Characters {
		if strings.ToLower(c.GetName()) == strings.ToLower(name) {
			return c, nil
		}
	}

	return nil, errors.New("character not found")
}
