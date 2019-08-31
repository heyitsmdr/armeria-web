package armeria

import (
	"armeria/internal/pkg/misc"
	"errors"
	"fmt"
	"sync"
)

// ItemInstance is an instance of an Item.
type ItemInstance struct {
	sync.RWMutex
	UUID             string            `json:"uuid"`
	UnsafeAttributes map[string]string `json:"attributes"`
	Parent           *Item             `json:"-"`
}

// Init is called when the ItemInstance is created or loaded from disk.
func (ii *ItemInstance) Init() {
	Armeria.registry.Register(ii, ii.ID(), RegistryTypeItemInstance)
}

// Deinit is called when the ItemInstance is deleted.
func (ii *ItemInstance) Deinit() {
	Armeria.registry.Unregister(ii.ID())
}

// ID returns the UUID of the instance.
func (ii *ItemInstance) ID() string {
	return ii.UUID
}

// Type returns the object type, since Item implements the ContainerObject interface.
func (ii *ItemInstance) Type() ContainerObjectType {
	return ContainerObjectTypeItem
}

// Name returns the raw Item name.
func (ii *ItemInstance) Name() string {
	return ii.Parent.Name()
}

// FormattedName returns the formatted Item name.
func (ii *ItemInstance) FormattedName() string {
	return TextStyle(
		ii.RarityColor()+"["+ii.Parent.Name()+"]",
		TextStyleBold,
		TextStyleColor,
	)
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

// InstanceAttribute returns an attribute on the ItemInstance, with no fallback to the parent Item.
func (ii *ItemInstance) InstanceAttribute(name string) string {
	ii.RLock()
	defer ii.RUnlock()

	return ii.UnsafeAttributes[name]
}

// Character returns the Character that has the ItemInstance.
func (ii *ItemInstance) Character() *Character {
	oc := Armeria.registry.GetObjectContainer(ii.ID())
	if oc == nil {
		return nil
	}
	return oc.ParentCharacter()
}

// Room returns the ItemInstance's Room based on the object container it is within.
func (ii *ItemInstance) Room() *Room {
	oc := Armeria.registry.GetObjectContainer(ii.ID())
	if oc == nil {
		return nil
	}
	return oc.ParentRoom()
}

// RarityColor returns the HTML color code that represents the rarity of the item.
func (ii *ItemInstance) RarityColor() string {
	switch ii.Attribute(AttributeRarity) {
	case "0":
		return "#ffffff"
	default:
		return "#ffffff"
	}
}

// RarityName returns the name of human-readable rarity name as a string.
func (ii *ItemInstance) RarityName() string {
	switch ii.Attribute(AttributeRarity) {
	case "0":
		return "Common"
	default:
		return "Common"
	}
}

// EditorData returns the JSON used for the object editor.
func (ii *ItemInstance) EditorData() *ObjectEditorData {
	props := []*ObjectEditorDataProperty{
		{PropType: "parent", Name: "parent", Value: ii.Name()},
	}

	for _, attrName := range ValidItemInstanceAttributes() {
		props = append(props, &ObjectEditorDataProperty{
			PropType:    "editable",
			Name:        attrName,
			Value:       ii.InstanceAttribute(attrName),
			ParentValue: ii.Parent.Attribute(attrName),
		})
	}

	return &ObjectEditorData{
		UUID:       ii.ID(),
		Name:       ii.Name(),
		ObjectType: "specific-item",
		Properties: props,
	}
}

// TooltipHTML generates the HTML string to be sent to the game client for rendering the associated tooltip.
func (ii *ItemInstance) TooltipHTML() string {
	return fmt.Sprintf(
		`
			<div class="name">%s</div>
			<div class="type">%s</div>
		`,
		ii.Name(),
		ii.RarityName(),
	)
}
