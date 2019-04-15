package characters

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
)

// Manager is the global manager instance for Character objects
type Manager struct {
	Characters []*Character `json:"characters"`
	dataFile string
	mux sync.Mutex
}

// Init creates a new character Manager instance
func Init(dataPath string) *Manager {
	m := &Manager{
		dataFile: fmt.Sprintf("%s/characters.json", dataPath),
	}

	m.loadCharacters()
	registerCommands()

	return m
}

func (m *Manager) loadCharacters() {
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