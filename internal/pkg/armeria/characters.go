package armeria

import (
	"encoding/json"
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

func NewCharacterManager(state *GameState) *CharacterManager {
	m := &CharacterManager{
		gameState: state,
		dataFile:  fmt.Sprintf("%s/characters.json", state.dataPath),
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

	log.Printf("[characters] loaded %d characters from file", len(m.Characters))
}

func (m *CharacterManager) SaveCharacters() {
	charactersFile, err := os.Create(m.dataFile)
	defer charactersFile.Close()

	raw, err := json.Marshal(m)
	if err != nil {
		log.Fatalf("[world] failed to marshal characters: %s", err)
	}

	bytes, err := charactersFile.Write(raw)
	if err != nil {
		log.Fatalf("[world] failed to write to characters file: %s", err)
	}

	charactersFile.Sync()

	log.Printf("[world] wrote %d bytes to characters file", bytes)
}

// GetCharacterByName returns the matching Character.
func (m *CharacterManager) GetCharacterByName(name string) *Character {
	for _, c := range m.Characters {
		if strings.ToLower(c.GetName()) == strings.ToLower(name) {
			return c
		}
	}

	return nil
}

// GetCharacters returns the characters logged in to the game.
func (m *CharacterManager) GetCharacters() []*Character {
	var chars []*Character
	for _, c := range m.Characters {
		if c.GetPlayer() != nil {
			chars = append(chars, c)
		}
	}
	return chars
}
