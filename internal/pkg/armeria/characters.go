package armeria

import (
	"armeria/internal/pkg/cloud"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

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

	err := json.Unmarshal(Armeria.storageManager.ReadFile(cloud.CharactersFile), m)
	if err != nil {
		Armeria.log.Fatal("failed to unmarshal data file",
			zap.Error(err),
		)
	}

	for _, c := range m.UnsafeCharacters {
		c.Init()
	}

	Armeria.log.Info("characters loaded",
		zap.Int("count", len(m.UnsafeCharacters)),
	)
}

func (m *CharacterManager) SaveCharacters() {
	m.RLock()
	defer m.RUnlock()

	raw, err := json.Marshal(m)
	if err != nil {
		Armeria.log.Fatal("failed to marshal data",
			zap.Error(err),
		)
	}

	bytes := Armeria.storageManager.WriteFile(cloud.CharactersFile, "application/json", raw)

	Armeria.log.Info("wrote data to cloud",
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
		if c.ID() == uuid {
			return c
		}
	}

	return nil
}

// CreateCharacter creates a new Character, adds it to memory, initializes it and returns the Character.
func (m *CharacterManager) CreateCharacter(name, password string) *Character {
	m.Lock()
	defer m.Unlock()
	c := &Character{
		UUID:                 uuid.New().String(),
		UnsafeName:           name,
		UnsafeAttributes:     make(map[string]string),
		UnsafeSettings:       make(map[string]string),
		UnsafeTempAttributes: make(map[string]string),
		UnsafeLastSeen:       time.Time{},
	}

	c.SetPassword(password)
	_ = c.SetAttribute(AttributeChannels, "General")

	m.UnsafeCharacters = append(m.UnsafeCharacters, c)

	c.Init()

	Armeria.log.Info("character created",
		zap.String("name", name),
	)

	return c
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
