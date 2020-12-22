package armeria

import (
	"armeria/internal/pkg/misc"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
)

type MobInstance struct {
	sync.RWMutex
	UUID                 string            `json:"uuid"`
	UnsafeAttributes     map[string]string `json:"attributes"`
	UnsafeInventory      *ObjectContainer  `json:"inventory"`
	UnsafeItemLedgers    []*Ledger         `json:"-"`
	Parent               *Mob              `json:"-"`
	UnsafeMobSpawnerUUID string            `json:"spawnerUUID"`
	UnsafeMoveTicks      int               `json:"moveTicks"`
	UnsafeConvoText      map[string]string `json:"-"`
}

// Init is called when the MobInstance is created or loaded from disk.
func (mi *MobInstance) Init() {
	// Register mob instance with registry.
	Armeria.registry.Register(mi, mi.ID(), RegistryTypeMobInstance)
	// Attach self as container's parent.
	mi.UnsafeInventory.AttachParent(mi, ContainerParentTypeMobInstance)
	// Sync container.
	mi.UnsafeInventory.Sync()
	// Initialize some properties.
	mi.UnsafeConvoText = make(map[string]string)
}

// Deinit is called when the MobInstance is deleted.
func (mi *MobInstance) Deinit() {
	Armeria.registry.Unregister(mi.ID())
}

// ID returns the UUID of the instance.
func (mi *MobInstance) ID() string {
	return mi.UUID
}

// Type returns the object type, since Mob implements the ContainerObject interface.
func (mi *MobInstance) Type() ContainerObjectType {
	return ContainerObjectTypeMob
}

// UnsafeName returns the raw Mob name.
func (mi *MobInstance) Name() string {
	return mi.Parent.Name()
}

// FormattedName returns the formatted Mob name.
func (mi *MobInstance) FormattedName() string {
	return TextStyle(
		mi.Parent.Name(),
		WithContextMenu(
			mi.Name(),
			"mob",
			"d48a3e",
			[]string{
				fmt.Sprintf("Look @|/look %s", mi.ID()),
				fmt.Sprintf("Interact @|/interact %s", mi.ID()),
				fmt.Sprintf("Jump @|/tp %s||CAN_BUILD", mi.Room().LocationString()),
				fmt.Sprintf("Edit @|/mob iedit %s||CAN_BUILD", mi.ID()),
				fmt.Sprintf("Edit-Parent @|/mob edit %s||CAN_BUILD", mi.Name()),
			},
		),
		WithBold(),
		WithColor("d48a3e"),
	)
}

// MoveTicks returns the number of mob movement ticks that have passed.
func (mi *MobInstance) MoveTicks() int {
	mi.RLock()
	defer mi.RUnlock()
	return mi.UnsafeMoveTicks
}

// IncMoveTicks increments the number of mob movement ticks since the last move.
func (mi *MobInstance) IncMoveTicks() {
	mi.Lock()
	defer mi.Unlock()
	mi.UnsafeMoveTicks = mi.UnsafeMoveTicks + 1
}

// ResetMoveTicks resets the number of mob movement ticks.
func (mi *MobInstance) ResetMoveTicks() {
	mi.Lock()
	defer mi.Unlock()
	mi.UnsafeMoveTicks = 0
}

// SetConvoText sets the display text for a particular conversation option. Used for caching.
func (mi *MobInstance) SetConvoText(optionId, displayText string) {
	mi.Lock()
	defer mi.Unlock()
	mi.UnsafeConvoText[optionId] = displayText
}

// ConvoText retrieves the cached display text for a particular conversation option. No entry returns a blank string.
func (mi *MobInstance) ConvoText(optionId string) string {
	mi.RLock()
	defer mi.RUnlock()
	if displayText, found := mi.UnsafeConvoText[optionId]; found {
		return displayText
	}
	return ""
}

// MobSpawnerUUID returns the UUID of the associated mob spawner, if any.
func (mi *MobInstance) MobSpawnerUUID() string {
	mi.RLock()
	defer mi.RUnlock()

	return mi.UnsafeMobSpawnerUUID
}

// SetMobSpawnerUUID sets the UUID of the mob spawner.
func (mi *MobInstance) SetMobSpawnerUUID(uuid string) {
	mi.Lock()
	defer mi.Unlock()

	mi.UnsafeMobSpawnerUUID = uuid
}

// SetAttribute sets a permanent attribute on the MobInstance.
func (mi *MobInstance) SetAttribute(name string, value string) error {
	mi.Lock()
	defer mi.Unlock()

	if mi.UnsafeAttributes == nil {
		mi.UnsafeAttributes = make(map[string]string)
	}

	if !misc.Contains(AttributeList(ObjectTypeMobInstance), name) {
		return errors.New("attribute name is invalid")
	}

	mi.UnsafeAttributes[strings.ToLower(name)] = value
	return nil
}

// Attribute returns an attribute on the MobInstance, and falls back to the parent Mob.
func (mi *MobInstance) Attribute(name string) string {
	mi.RLock()
	defer mi.RUnlock()

	if len(mi.UnsafeAttributes[name]) == 0 {
		return mi.Parent.Attribute(name)
	}

	return mi.UnsafeAttributes[name]
}

// AttributeBool returns an attribute on the MobInstance as a bool.
func (mi *MobInstance) AttributeBool(name string) bool {
	v := mi.Attribute(name)
	if v == "true" {
		return true
	}

	return false
}

// AttributeInt returns an attribute on the MobInstance as an int.
func (mi *MobInstance) AttributeInt(name string) int {
	v := mi.Attribute(name)
	i, err := strconv.Atoi(v)
	if err != nil {
		return 0
	}

	return i
}

// InstanceAttribute returns an attribute on the MobInstance, with no fallback to the parent Mob.
func (mi *MobInstance) InstanceAttribute(name string) string {
	mi.RLock()
	defer mi.RUnlock()

	return mi.UnsafeAttributes[name]
}

// AddItemLedger adds a "known" item ledger to the mob for use with buying and selling items.
func (mi *MobInstance) AddItemLedger(ledger *Ledger) {
	mi.Lock()
	defer mi.Unlock()

	mi.UnsafeItemLedgers = append(mi.UnsafeItemLedgers, ledger)
}

// ItemLedgers returns the "known" item ledgers for use with buying and selling items.
func (mi *MobInstance) ItemLedgers() []*Ledger {
	mi.RLock()
	defer mi.RUnlock()

	return mi.UnsafeItemLedgers
}

// MobInstance returns the MobInstance's Room based on the object container it is within.
func (mi *MobInstance) Room() *Room {
	oc := Armeria.registry.GetObjectContainer(mi.ID())
	if oc == nil {
		return nil
	}
	return oc.ParentRoom()
}

// Inventory returns the unsafeCharacter's inventory.
func (mi *MobInstance) Inventory() *ObjectContainer {
	mi.RLock()
	defer mi.RUnlock()

	return mi.UnsafeInventory
}

// EditorData returns the JSON used for the object editor.
func (mi *MobInstance) EditorData() *ObjectEditorData {
	props := []*ObjectEditorDataProperty{
		{PropType: "parent", Name: "parent", Value: mi.Name()},
	}

	for _, attrName := range AttributeList(ObjectTypeMobInstance) {
		props = append(props, &ObjectEditorDataProperty{
			PropType:    AttributeEditorType(ObjectTypeMobInstance, attrName),
			Name:        attrName,
			Value:       mi.InstanceAttribute(attrName),
			ParentValue: mi.Parent.Attribute(attrName),
		})
	}

	return &ObjectEditorData{
		UUID:       mi.ID(),
		Name:       mi.Name(),
		ObjectType: "specific-mob",
		IsChild:    true,
		Properties: props,
	}
}

func (mi *MobInstance) Pronoun(pt PronounType) string {
	gender := mi.Attribute(AttributeGender)
	if gender == "male" {
		if pt == PronounSubjective {
			return "he"
		} else if pt == PronounPossessiveAbsolute {
			return "his"
		} else if pt == PronounPossessiveAdjective {
			return "his"
		} else if pt == PronounObjective {
			return "him"
		}
	} else if gender == "female" {
		if pt == PronounSubjective {
			return "she"
		} else if pt == PronounPossessiveAbsolute {
			return "hers"
		} else if pt == PronounPossessiveAdjective {
			return "her"
		} else if pt == PronounObjective {
			return "her"
		}
	} else if gender == "thing" {
		if pt == PronounSubjective {
			return "it"
		} else if pt == PronounPossessiveAbsolute {
			return "its"
		} else if pt == PronounPossessiveAdjective {
			return "its"
		} else if pt == PronounObjective {
			return "it"
		}
	}

	return ""
}

// Delete removes the mob instance from the game. It should be manually removed from containers
// first before calling this function!
func (mi *MobInstance) Delete() {
	mi.Parent.DeleteInstance(mi)
}
