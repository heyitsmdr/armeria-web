package armeria

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	"go.uber.org/zap"
)

type ItemManager struct {
	dataFile    string
	UnsafeItems []*Item `json:"items"`
	mux         sync.Mutex
}

func NewItemManager() *ItemManager {
	m := &ItemManager{
		dataFile: fmt.Sprintf("%s/items.json", Armeria.dataPath),
	}

	m.LoadItems()

	return m
}

// LoadItems loads the items from disk into memory.
func (m *ItemManager) LoadItems() {
	itemsFile, err := os.Open(m.dataFile)
	defer itemsFile.Close()

	if err != nil {
		Armeria.log.Fatal("failed to load data file",
			zap.String("file", m.dataFile),
			zap.Error(err),
		)
	}

	jsonParser := json.NewDecoder(itemsFile)

	err = jsonParser.Decode(m)
	if err != nil {
		Armeria.log.Fatal("failed to decode data file",
			zap.String("file", m.dataFile),
			zap.Error(err),
		)
	}

	Armeria.log.Info("items loaded",
		zap.Int("count", len(m.UnsafeItems)),
	)
}

// SaveItems writes the in-memory items to disk.
func (m *ItemManager) SaveItems() {
	itemsFile, err := os.Create(m.dataFile)
	defer itemsFile.Close()

	raw, err := json.Marshal(m)
	if err != nil {
		Armeria.log.Fatal("failed to marshal data",
			zap.Error(err),
		)
	}

	bytes, err := itemsFile.Write(raw)
	if err != nil {
		Armeria.log.Fatal("failed to write data file",
			zap.String("file", m.dataFile),
			zap.Error(err),
		)
	}

	_ = itemsFile.Sync()

	Armeria.log.Info("wrote data to file",
		zap.String("file", m.dataFile),
		zap.Int("bytes", bytes),
	)
}

// ItemByName returns the matching Item.
func (m *ItemManager) ItemByName(name string) *Item {
	m.mux.Lock()
	defer m.mux.Unlock()

	for _, i := range m.UnsafeItems {
		if strings.ToLower(i.Name()) == strings.ToLower(name) {
			return i
		}
	}

	return nil
}

// Items returns all of the in-memory Items.
func (m *ItemManager) Items() []*Item {
	return m.UnsafeItems
}

// AddItem adds a new Item reference to memory.
func (m *ItemManager) AddItem(i *Item) {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.UnsafeItems = append(m.UnsafeItems, i)
}
