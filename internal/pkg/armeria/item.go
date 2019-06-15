package armeria

import (
	"armeria/internal/pkg/misc"
	"sync"

	"github.com/google/uuid"

	"go.uber.org/zap"
)

type Item struct {
	UnsafeName       string            `json:"name"`
	UnsafeAttributes map[string]string `json:"attributes"`
	UnsafeInstances  []*ItemInstance   `json:"instances"`
	mux              sync.Mutex
}

// ValidItemAttributes returns an array of valid attributes that can be permanently set.
func ValidItemAttributes() []string {
	return []string{
		"picture",
	}
}

// ItemAttributeDefault returns the default value for a particular attribute.
func ItemAttributeDefault(name string) string {
	switch name {

	}

	return ""
}

// ValidateItemAttribute returns a bool indicating whether a particular value is allowed
// for a particular attribute.
func ValidateItemAttribute(name string, value string) (bool, string) {
	switch name {
	case "script":
		return false, "script cannot be set explicitly"
	}

	return true, ""
}

func (i *Item) Name() string {
	i.mux.Lock()
	defer i.mux.Unlock()
	return i.UnsafeName
}

func (i *Item) CreateInstance() *ItemInstance {
	i.mux.Lock()
	defer i.mux.Unlock()

	ii := &ItemInstance{
		UUID:             uuid.New().String(),
		UnsafeParent:     i.UnsafeName,
		UnsafeAttributes: make(map[string]string),
	}

	i.UnsafeInstances = append(i.UnsafeInstances, ii)

	return ii
}

// DeleteInstance removes the ItemInstance from memory.
func (i *Item) DeleteInstance(ii *ItemInstance) bool {
	i.mux.Lock()
	defer i.mux.Unlock()

	for idx, inst := range i.UnsafeInstances {
		if inst.Id() == ii.Id() {
			i.UnsafeInstances[idx] = i.UnsafeInstances[len(i.UnsafeInstances)-1]
			i.UnsafeInstances = i.UnsafeInstances[:len(i.UnsafeInstances)-1]
			return true
		}
	}

	return false
}

// Attribute returns a permanent attribute.
func (i *Item) Attribute(name string) string {
	i.mux.Lock()
	defer i.mux.Unlock()

	if len(i.UnsafeAttributes[name]) == 0 {
		return ItemAttributeDefault(name)
	}

	return i.UnsafeAttributes[name]
}

// SetAttribute sets a permanent attribute and only valid attributes can be set.
func (i *Item) SetAttribute(name string, value string) {
	i.mux.Lock()
	defer i.mux.Unlock()

	if !misc.Contains(ValidItemAttributes(), name) {
		Armeria.log.Fatal("attempted to set invalid attribute",
			zap.String("attribute", name),
			zap.String("value", value),
		)
	}

	i.UnsafeAttributes[name] = value
}
