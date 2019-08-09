package armeria

import (
	"armeria/internal/pkg/misc"
	"errors"
	"fmt"
	"sync"
)

type ItemInstance struct {
	sync.RWMutex
	UUID              string            `json:"uuid"`
	UnsafeCharacterId string            `json:"character"`
	UnsafeAttributes  map[string]string `json:"attributes"`
	Parent            *Item             `json:"-"`
}

type ItemLocationType int

const (
	ItemLocationRoom ItemLocationType = iota
	ItemLocationCharacter
)

// Init is called when the ItemInstance is created or loaded from disk.
func (ii *ItemInstance) Init() {
	Armeria.registry.Register(ii, ii.Id(), RegistryTypeItemInstance)
}

// Deinit is called when the ItemInstance is deleted.
func (ii *ItemInstance) Deinit() {
	Armeria.registry.Unregister(ii.Id())
}

// Id returns the UUID of the instance.
func (ii *ItemInstance) Id() string {
	return ii.UUID
}

// Type returns the object type, since Item implements the ContainerObject interface.
func (ii *ItemInstance) Type() ContainerObjectType {
	return ContainerObjectTypeItem
}

// UnsafeName returns the raw Item name.
func (ii *ItemInstance) Name() string {
	return ii.Parent.Name()
}

// FormattedName returns the formatted Item name.
func (ii *ItemInstance) FormattedName() string {
	return fmt.Sprintf("[b]%s[/b]", ii.Parent.Name())
}

// SetAttribute sets a permanent attribute on the ItemInstance.
func (ii *ItemInstance) SetAttribute(name string, value string) error {
	ii.Lock()
	defer ii.Unlock()

	if ii.UnsafeAttributes == nil {
		ii.UnsafeAttributes = make(map[string]string)
	}

	if !misc.Contains(ValidItemAttributes(), name) {
		return errors.New("attribute name is invalid")
	}

	ii.UnsafeAttributes[name] = value
	return nil
}

// Attribute returns an attribute on the ItemInstance, and falls back to the parent Item.
func (ii *ItemInstance) Attribute(name string) string {
	ii.RLock()
	defer ii.RUnlock()

	if len(ii.UnsafeAttributes[name]) == 0 {
		return ii.Parent.Attribute(name)
	}

	return ii.UnsafeAttributes[name]
}

// Character returns the Character that has the ItemInstance.
func (ii *ItemInstance) Character() *Character {
	ii.RLock()
	defer ii.RUnlock()

	return Armeria.characterManager.CharacterById(ii.UnsafeCharacterId)
}

// SetCharacter sets the character that has the ItemInstance.
func (ii *ItemInstance) SetCharacter(c *Character) {
	ii.Lock()
	defer ii.Unlock()

	ii.UnsafeCharacterId = c.Id()
}

// ItemInstance returns the ItemInstance's Room based on the object container it is within.
func (ii *ItemInstance) Room() *Room {
	oc := Armeria.registry.GetObjectContainer(ii.Id())
	if oc == nil {
		return nil
	}
	return oc.ParentRoom()
}
