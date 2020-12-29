package armeria

import (
	"armeria/internal/pkg/misc"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"go.uber.org/zap"
)

// Force verify that ItemInstance implements ContainerObject.
var _ ContainerObject = (*ItemInstance)(nil)

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
		fmt.Sprintf("[%s]", ii.Parent.Name()),
		WithItemTooltip(ii.ID()),
		WithContextMenu(
			ii.Name(),
			"item",
			ii.RarityColor(),
			[]string{
				fmt.Sprintf("Look @|/look %s", ii.ID()),
			},
		),
		WithBold(),
		WithColor(ii.RarityColor()),
	)
}

// SetAttribute sets a permanent attribute on the ItemInstance.
func (ii *ItemInstance) SetAttribute(name string, value string) error {
	ii.Lock()
	defer ii.Unlock()

	if ii.UnsafeAttributes == nil {
		ii.UnsafeAttributes = make(map[string]string)
	}

	if !misc.Contains(AttributeList(ObjectTypeItemInstance), name) {
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

// AttributeBool returns an attribute on the ItemInstance as a bool.
func (ii *ItemInstance) AttributeBool(name string) bool {
	v := ii.Attribute(name)
	if v == "true" {
		return true
	}

	return false
}

// AttributeInt returns an attribute on the ItemInstance as an int.
func (ii *ItemInstance) AttributeInt(name string) int {
	v := ii.Attribute(name)
	i, err := strconv.Atoi(v)
	if err != nil {
		return 0
	}

	return i
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

// Room returns the Room that has the ItemInstance.
func (ii *ItemInstance) Room() *Room {
	oc := Armeria.registry.GetObjectContainer(ii.ID())
	if oc == nil {
		return nil
	}
	return oc.ParentRoom()
}

// MobInstance returns the MobInstance that has the ItemInstance.
func (ii *ItemInstance) MobInstance() *MobInstance {
	oc := Armeria.registry.GetObjectContainer(ii.ID())
	if oc == nil {
		return nil
	}
	return oc.ParentMobInstance()
}

// RarityColor returns the HTML color code that represents the rarity of the item.
func (ii *ItemInstance) RarityColor() string {
	switch ii.Attribute(AttributeRarity) {
	case ItemRarityCommon:
		return "ffffff"
	case ItemRarityUncommon:
		return "00ff00"
	default:
		return "ffffff"
	}
}

// RarityName returns the capitalized name of the rarity type.
func (ii *ItemInstance) RarityName() string {
	switch ii.Attribute(AttributeRarity) {
	case ItemRarityCommon:
		return "Common"
	case ItemRarityUncommon:
		return "Uncommon"
	default:
		return "Common"
	}
}

// EditorData returns the JSON used for the object editor.
func (ii *ItemInstance) EditorData() *ObjectEditorData {
	props := []*ObjectEditorDataProperty{
		{PropType: "parent", Name: "parent", Value: ii.Name()},
	}

	for _, attrName := range AttributeList(ObjectTypeItemInstance) {
		props = append(props, &ObjectEditorDataProperty{
			PropType:    AttributeEditorType(ObjectTypeItemInstance, attrName),
			Name:        attrName,
			Value:       ii.InstanceAttribute(attrName),
			ParentValue: ii.Parent.Attribute(attrName),
		})
	}

	return &ObjectEditorData{
		UUID:       ii.ID(),
		Name:       ii.Name(),
		ObjectType: "specific-item",
		IsChild:    true,
		Properties: props,
	}
}

// TooltipContentJSON generates the HTML string to be sent to the game client in JSON format.
func (ii *ItemInstance) TooltipContentJSON() string {
	qualitiesSlice := make([]string, 0)

	if ii.Attribute(AttributeHoldable) == "false" {
		qualitiesSlice = append(qualitiesSlice, "Not Holdable")
	}

	if ii.Attribute(AttributeVisible) == "false" {
		qualitiesSlice = append(qualitiesSlice, "Not Visible")
	}

	tt := map[string]string{
		"uuid": ii.ID(),
		"html": fmt.Sprintf(
			`
			<div class="name" style="color:%s">%s</div>
			<div class="type">%s</div>
			<div class="qualities">%s</div>
			`,
			ii.RarityColor(),
			ii.Name(),
			ii.RarityName(),
			strings.Join(qualitiesSlice, "<br />"),
		),
		"rarity":  ii.RarityColor(),
		"picture": ii.Attribute(AttributePicture),
	}

	ttJSON, err := json.Marshal(tt)
	if err != nil {
		Armeria.log.Fatal("failed to marshal item tooltip content",
			zap.String("uuid", ii.ID()),
			zap.Error(err),
		)
	}

	return string(ttJSON)
}

// Delete removes the item instance from the game. It should be manually removed from containers
// first before calling this function!
func (ii *ItemInstance) Delete() {
	ii.Parent.DeleteInstance(ii)
}
