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
	sync.RWMutex
	dataFile    string
	UnsafeItems []*Item `json:"items"`
}

func NewItemManager() *ItemManager {
	m := &ItemManager{
		dataFile: fmt.Sprintf("%s/items.json", Armeria.dataPath),
	}

	m.LoadItems()
	m.AttachParents()
	m.AddItemInstancesToRooms()

	return m
}

// LoadItems loads the items from disk into memory.
func (m *ItemManager) LoadItems() {
	m.Lock()
	defer m.Unlock()

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
	m.RLock()
	defer m.RUnlock()

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

// AttachParents attaches a pointer to ItemInstance that references the parent Item.
func (m *ItemManager) AttachParents() {
	m.RLock()
	defer m.RUnlock()

	for _, i := range m.UnsafeItems {
		for _, ii := range i.Instances() {
			ii.Parent = i
		}
	}
}

// AddItemInstancesToRooms adds ItemInstance objects that are in Rooms to their
// respective Room objects.
func (m *ItemManager) AddItemInstancesToRooms() {
	m.RLock()
	defer m.RUnlock()

	for _, i := range m.UnsafeItems {
		for _, ii := range i.Instances() {
			if ii.LocationType() != ItemLocationRoom {
				continue
			}

			r := ii.Location.Room()
			if r == nil {
				Armeria.log.Fatal("item instance in invalid room",
					zap.String("item", ii.Name()),
					zap.String("uuid", ii.Id()),
					zap.String("location", fmt.Sprintf("%v", ii.Location)),
				)
				return
			}
			r.AddObjectToRoom(ii)
		}
	}
}

// ItemByName returns the matching Item, by name.
func (m *ItemManager) ItemByName(name string) *Item {
	m.RLock()
	defer m.RUnlock()

	for _, i := range m.UnsafeItems {
		if strings.ToLower(i.Name()) == strings.ToLower(name) {
			return i
		}
	}

	return nil
}

// ItemInstanceById returns the matching ItemInstance, by uuid.
func (m *ItemManager) ItemInstanceById(uuid string) *ItemInstance {
	m.RLock()
	defer m.RUnlock()

	for _, i := range m.UnsafeItems {
		for _, ii := range i.Instances() {
			if ii.Id() == uuid {
				return ii
			}
		}
	}

	return nil
}

// Items returns all of the in-memory Items.
func (m *ItemManager) Items() []*Item {
	m.RLock()
	defer m.RUnlock()

	return m.UnsafeItems
}

// CreateItem creates a new Item instance, but doesn't add it to memory.
func (m *ItemManager) CreateItem(name string) *Item {
	return &Item{
		UnsafeName:       name,
		UnsafeAttributes: make(map[string]string),
	}
}

// AddItem adds a new Item reference to memory.
func (m *ItemManager) AddItem(i *Item) {
	m.Lock()
	defer m.Unlock()

	m.UnsafeItems = append(m.UnsafeItems, i)
}
