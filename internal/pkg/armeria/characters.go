package armeria

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"
)

type CharacterManager struct {
	dataFile         string
	UnsafeCharacters []*Character `json:"characters"`
}

func NewCharacterManager() *CharacterManager {
	m := &CharacterManager{
		dataFile: fmt.Sprintf("%s/characters.json", Armeria.dataPath),
	}

	m.LoadCharacters()

	return m
}

func (m *CharacterManager) LoadCharacters() {
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
		zap.Int("count", len(m.UnsafeCharacters)),
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

	_ = charactersFile.Sync()

	Armeria.log.Info("wrote data to file",
		zap.String("file", m.dataFile),
		zap.Int("bytes", bytes),
	)
}

// CharacterByName returns the matching Character.
func (m *CharacterManager) CharacterByName(name string) *Character {
	for _, c := range m.UnsafeCharacters {
		if strings.ToLower(c.Name()) == strings.ToLower(name) {
			return c
		}
	}

	return nil
}

// OnlineCharacters returns the characters logged in to the game.
func (m *CharacterManager) OnlineCharacters() []*Character {
	var chars []*Character
	for _, c := range m.UnsafeCharacters {
		if c.Player() != nil {
			chars = append(chars, c)
		}
	}
	return chars
}

// Characters returns all the characters in the database.
func (m *CharacterManager) Characters() []*Character {
	return m.UnsafeCharacters
}
