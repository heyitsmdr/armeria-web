package armeria

import (
	"armeria/internal/pkg/misc"
	"errors"
	"sync"
)

type MobInstance struct {
	sync.RWMutex
	UUID              string            `json:"uuid"`
	UnsafeAttributes  map[string]string `json:"attributes"`
	UnsafeInventory   *ObjectContainer  `json:"inventory"`
	UnsafeItemLedgers []*Ledger         `json:"-"`
	Parent            *Mob              `json:"-"`
}

// Init is called when the MobInstance is created or loaded from disk.
func (mi *MobInstance) Init() {
	// register mob instance with registry
	Armeria.registry.Register(mi, mi.ID(), RegistryTypeMobInstance)
	// attach self as container's parent
	mi.UnsafeInventory.AttachParent(mi, ContainerParentTypeMobInstance)
	// sync container
	mi.UnsafeInventory.Sync()
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
	return TextStyle(mi.Parent.Name(), WithBold())
}

// SetAttribute sets a permanent attribute on the MobInstance.
func (mi *MobInstance) SetAttribute(name string, value string) error {
	mi.Lock()
	defer mi.Unlock()

	if mi.UnsafeAttributes == nil {
		mi.UnsafeAttributes = make(map[string]string)
	}

	if !misc.Contains(ValidMobAttributes(), name) {
		return errors.New("attribute name is invalid")
	}

	mi.UnsafeAttributes[name] = value
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

	for _, attrName := range ValidMobInstanceAttributes() {
		props = append(props, &ObjectEditorDataProperty{
			PropType:    AttributeEditorType(attrName),
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
