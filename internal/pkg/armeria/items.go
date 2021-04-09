package armeria

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	"go.uber.org/zap"
)

// ItemManager holds all items found in the server data.
type ItemManager struct {
	sync.RWMutex
	dataFile    string
	UnsafeItems []*Item `json:"items"`
}

// NewItemManager creates a new ItemManager.
func NewItemManager() *ItemManager {
	m := &ItemManager{
		dataFile: fmt.Sprintf("%s/items.json", Armeria.dataPath),
	}

	m.LoadItems()
	m.AttachParents()

	return m
}

// LoadItems loads the items from disk into memory.
func (m *ItemManager) LoadItems() {
	m.Lock()
	defer m.Unlock()

	err := json.Unmarshal(Armeria.storageManager.ReadFile("items.json"), m)
	if err != nil {
		Armeria.log.Fatal("failed to unmarshal data file",
			zap.String("file", m.dataFile),
			zap.Error(err),
		)
	}

	for _, i := range m.UnsafeItems {
		i.Init()

		for _, ii := range i.Instances() {
			ii.Init()
		}
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

// ItemInstanceByID returns the matching ItemInstance, by uuid.
func (m *ItemManager) ItemInstanceByID(uuid string) *ItemInstance {
	m.RLock()
	defer m.RUnlock()

	for _, i := range m.UnsafeItems {
		for _, ii := range i.Instances() {
			if ii.ID() == uuid {
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

// ItemsByAttribute returns all of the in-memory Items that have a particular attribute + value defined.
func (m *ItemManager) ItemsByAttribute(a, t string) []*Item {
	m.RLock()
	defer m.RUnlock()

	matches := make([]*Item, 0)
	for _, i := range m.UnsafeItems {
		if i.Attribute(a) == t {
			matches = append(matches, i)
		}
	}

	return matches
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

// RemoveItem removes an existing Item reference from memory.
func (m *ItemManager) RemoveItem(item *Item) {
	m.Lock()
	defer m.Unlock()

	var found bool
	for idx, inst := range m.UnsafeItems {
		if inst.Name() == item.Name() {
			m.UnsafeItems[idx] = m.UnsafeItems[len(m.UnsafeItems)-1]
			m.UnsafeItems = m.UnsafeItems[:len(m.UnsafeItems)-1]
			found = true
			break
		}
	}

	if !found {
		return
	}

	// Delete the picture file.
	picture := item.Attribute(AttributePicture)
	if len(picture) > 0 {
		DeleteObjectPictureFromDisk(picture)
	}
}
