package armeria

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"
)

type CharacterManager struct {
	dataFile   string
	Characters []*Character `json:"characters"`
}

func NewCharacterManager() *CharacterManager {
	m := &CharacterManager{
		dataFile: fmt.Sprintf("%s/characters.json", Armeria.dataPath),
	}

	m.loadCharacters()

	return m
}

func (m *CharacterManager) loadCharacters() {
	charactersFile, err := os.Open(m.dataFile)
	defer charactersFile.Close()

	if err != nil {
		Armeria.log.Fatal("failed to load data file",
			zap.String("file", m.dataFile),
			zap.Error(err),
		)
	}

	jsonParser := json.NewDecoder(charactersFile)

	err = jsonParser.Decode(m)
	if err != nil {
		Armeria.log.Fatal("failed to decode data file",
			zap.String("file", m.dataFile),
			zap.Error(err),
		)
	}

	Armeria.log.Info("characters loaded",
		zap.Int("count", len(m.Characters)),
	)
}

func (m *CharacterManager) SaveCharacters() {
	charactersFile, err := os.Create(m.dataFile)
	defer charactersFile.Close()

	raw, err := json.Marshal(m)
	if err != nil {
		Armeria.log.Fatal("failed to marshal data",
			zap.Error(err),
		)
	}

	bytes, err := charactersFile.Write(raw)
	if err != nil {
		Armeria.log.Fatal("failed to write data file",
			zap.String("file", m.dataFile),
			zap.Error(err),
		)
	}

	charactersFile.Sync()

	Armeria.log.Info("wrote data to file",
		zap.String("file", m.dataFile),
		zap.Int("bytes", bytes),
	)
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
