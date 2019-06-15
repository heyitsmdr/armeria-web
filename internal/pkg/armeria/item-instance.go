package armeria

import (
	"armeria/internal/pkg/misc"
	"fmt"
	"sync"

	"go.uber.org/zap"
)

type ItemInstance struct {
	UUID               string            `json:"uuid"`
	UnsafeParent       string            `json:"parent"`
	UnsafeLocationType int               `json:"location_type"`
	UnsafeLocation     *Location         `json:"location"`
	UnsafeCharacter    string            `json:"character"`
	UnsafeAttributes   map[string]string `json:"attributes"`
	mux                sync.Mutex
}

const (
	ItemLocationRoom      int = 0
	ItemLocationCharacter int = 1
)

// Id returns the UUID of the instance.
func (ii *ItemInstance) Id() string {
	return ii.UUID
}

// Parent returns the Item parent.
func (ii *ItemInstance) Parent() *Item {
	return Armeria.itemManager.ItemByName(ii.UnsafeParent)
}

// Type returns the object type, since Item implements the Object interface.
func (ii *ItemInstance) Type() int {
	return ObjectTypeItem
}

// UnsafeName returns the raw Item name.
func (ii *ItemInstance) Name() string {
	return ii.UnsafeParent
}

// FormattedName returns the formatted Item name.
func (ii *ItemInstance) FormattedName() string {
	return fmt.Sprintf("[b]%s[/b]", ii.UnsafeParent)
}

// SetAttribute sets a permanent attribute on the ItemInstance.
func (ii *ItemInstance) SetAttribute(name string, value string) {
	ii.mux.Lock()
	defer ii.mux.Unlock()

	if ii.UnsafeAttributes == nil {
		ii.UnsafeAttributes = make(map[string]string)
	}

	if !misc.Contains(ValidMobAttributes(), name) {
		Armeria.log.Fatal("attempted to set invalid attribute",
			zap.String("attribute", name),
			zap.String("value", value),
		)
	}

	ii.UnsafeAttributes[name] = value
}

// Attribute returns an attribute on the ItemInstance, and falls back to the parent Item.
func (ii *ItemInstance) Attribute(name string) string {
	ii.mux.Lock()
	defer ii.mux.Unlock()

	if len(ii.UnsafeAttributes[name]) == 0 {
		return ii.Parent().Attribute(name)
	}

	return ii.UnsafeAttributes[name]
}

// LocationType returns the location type (room or character).
func (ii *ItemInstance) LocationType() int {
	ii.mux.Lock()
	defer ii.mux.Unlock()
	return ii.UnsafeLocationType
}

// Location returns the location of the ItemInstance.
func (ii *ItemInstance) Location() *Location {
	ii.mux.Lock()
	defer ii.mux.Unlock()
	return ii.UnsafeLocation
}

// SetLocation sets the location of the ItemInstance.
func (ii *ItemInstance) SetLocation(l *Location) {
	ii.mux.Lock()
	defer ii.mux.Unlock()
	ii.UnsafeLocationType = ItemLocationRoom
	ii.UnsafeLocation = l
}

// Character returns the Character that has the ItemInstance.
func (ii *ItemInstance) Character() *Character {
	ii.mux.Lock()
	defer ii.mux.Unlock()
	return Armeria.characterManager.CharacterByName(ii.UnsafeCharacter)
}

// SetCharacter sets the character that has the ItemInstance.
func (ii *ItemInstance) SetCharacter(c *Character) {
	ii.mux.Lock()
	defer ii.mux.Unlock()
	ii.UnsafeLocationType = ItemLocationCharacter
	ii.UnsafeCharacter = c.Name()
}
