package armeria

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	"go.uber.org/zap"
)

type CharacterManager struct {
	sync.RWMutex
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
	m.Lock()
	defer m.Unlock()

	charactersFile, err := os.Open(m.dataFile)
	if err != nil {
		Armeria.log.Fatal("failed to load data file",
			zap.String("file", m.dataFile),
			zap.Error(err),
		)
	}
	defer charactersFile.Close()

	jsonParser := json.NewDecoder(charactersFile)

	err = jsonParser.Decode(m)
	if err != nil {
		Armeria.log.Fatal("failed to decode data file",
			zap.String("file", m.dataFile),
			zap.Error(err),
		)
	}

	for _, _ = range m.UnsafeCharacters {
		//if c.UnsafeInventory == nil {
		//	c.UnsafeInventory = NewItemContainer(35)
		//}
	}

	Armeria.log.Info("characters loaded",
		zap.Int("count", len(m.UnsafeCharacters)),
	)
}

func (m *CharacterManager) SaveCharacters() {
	m.RLock()
	defer m.RUnlock()

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

// CharacterByName returns the matching Character, by name.
func (m *CharacterManager) CharacterByName(name string) *Character {
	m.RLock()
	defer m.RUnlock()

	for _, c := range m.UnsafeCharacters {
		if strings.ToLower(c.Name()) == strings.ToLower(name) {
			return c
		}
	}

	return nil
}

// CharacterById returns the matching Character, by uuid.
func (m *CharacterManager) CharacterById(uuid string) *Character {
	m.RLock()
	defer m.RUnlock()

	for _, c := range m.UnsafeCharacters {
		if c.Id() == uuid {
			return c
		}
	}

	return nil
}

// OnlineCharacters returns the characters logged in to the game.
func (m *CharacterManager) OnlineCharacters() []*Character {
	m.RLock()
	defer m.RUnlock()

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
	m.RLock()
	defer m.RUnlock()

	return m.UnsafeCharacters
}
