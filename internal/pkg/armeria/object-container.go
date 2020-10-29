package armeria

import (
	"armeria/internal/pkg/misc"
	"errors"
	"strings"
	"sync"
)

// ObjectContainer is a container of game objects that can be persisted to disk. These can be used for
// things in a room, a unsafeCharacter's inventory, a chest, etc.
type ObjectContainer struct {
	sync.RWMutex
	UnsafeObjects    []*ObjectContainerDefinition `json:"objects"`
	UnsafeMaxSize    int                          `json:"maxSize"` // 0 = unlimited
	UnsafeParent     interface{}                  `json:"-"`
	UnsafeParentType ContainerParentType          `json:"-"`
}

// ObjectContainerDefinition contains a definition for an object within a container.
type ObjectContainerDefinition struct {
	UUID string `json:"uuid"`
	Slot int    `json:"slot"`
}

// ContainerObject is an interface that describes an object that can go within an ObjectContainer.
type ContainerObject interface {
	ID() string
	Type() ContainerObjectType
	Name() string
	FormattedName() string
	Attribute(name string) string
	SetAttribute(name string, value string) error
}

// ContainerObjectType is an int representing the object type.
type ContainerObjectType int

// Constants representing the various object types.
const (
	ContainerObjectTypeCharacter ContainerObjectType = iota
	ContainerObjectTypeMob
	ContainerObjectTypeItem
)

var (
	// ErrContainerNoRoom is an error for when the container is bounded and full.
	ErrContainerNoRoom = errors.New("no space in container")
	// ErrContainerDuplicate is an error for when the container already contains a specific uuid.
	ErrContainerDuplicate = errors.New("object already in container")
)

// ContainerParentType is an int representing the container's parent type.
type ContainerParentType int

// Constants representing the various parent types.
const (
	ContainerParentTypeRoom ContainerParentType = iota
	ContainerParentTypeCharacter
	ContainerParentTypeMobInstance
)

// NewObjectContainer will return a new object container with the specified max size.
func NewObjectContainer(maxSize int) *ObjectContainer {
	return &ObjectContainer{
		UnsafeObjects: make([]*ObjectContainerDefinition, 0),
		UnsafeMaxSize: maxSize,
	}
}

// ObjectSortOrder returns the sort order for each type of ContainerObject. This will affect
// how it will appear in the client's room list. Sorting is in descending order.
func ObjectSortOrder(ot ContainerObjectType) int {
	switch ot {
	case ContainerObjectTypeMob:
		return 75
	case ContainerObjectTypeCharacter:
		return 50
	case ContainerObjectTypeItem:
		return 25
	}

	return 0
}

// AttachParent sets the parent object that the object container belongs to.
func (oc *ObjectContainer) AttachParent(p interface{}, t ContainerParentType) {
	oc.Lock()
	defer oc.Unlock()

	oc.UnsafeParent = p
	oc.UnsafeParentType = t
}

// ParentRoom returns the parent Room if the object has the appropriate parent type.
func (oc *ObjectContainer) ParentRoom() *Room {
	oc.RLock()
	defer oc.RUnlock()

	if oc.UnsafeParentType != ContainerParentTypeRoom {
		return nil
	}

	return oc.UnsafeParent.(*Room)
}

// ParentCharacter returns the parent Character if the object has the appropriate parent type.
func (oc *ObjectContainer) ParentCharacter() *Character {
	oc.RLock()
	defer oc.RUnlock()

	if oc.UnsafeParentType != ContainerParentTypeCharacter {
		return nil
	}

	return oc.UnsafeParent.(*Character)
}

// ParentMobInstance returns the parent MobInstance if the object has the appropriate parent type.
func (oc *ObjectContainer) ParentMobInstance() *MobInstance {
	oc.RLock()
	defer oc.RUnlock()

	if oc.UnsafeParentType != ContainerParentTypeMobInstance {
		return nil
	}

	return oc.UnsafeParent.(*MobInstance)
}

// ParentType returns the ContainerParentType that owns this object container.
func (oc *ObjectContainer) ParentType() ContainerParentType {
	oc.RLock()
	defer oc.RUnlock()

	return oc.UnsafeParentType
}

func (oc *ObjectContainer) MaxSize() int {
	oc.RLock()
	defer oc.RUnlock()

	return oc.UnsafeMaxSize
}

// Contains returns a bool indicating whether the object container contains something with the
// specified uuid.
func (oc *ObjectContainer) Contains(uuid string) bool {
	o, _, _ := oc.Get(uuid)
	return o != nil
}

// Get retrieves an object from the object container, based on the uuid.
func (oc *ObjectContainer) Get(uuid string) (interface{}, *ObjectContainerDefinition, RegistryType) {
	oc.RLock()
	defer oc.RUnlock()

	for _, ocd := range oc.UnsafeObjects {
		if ocd.UUID == uuid {
			o, ot := Armeria.registry.Get(ocd.UUID)
			return o, ocd, ot
		}
	}

	return nil, nil, RegistryTypeUnknown
}

// GetByName retrieves an object from the object container, based on the name of the object.
func (oc *ObjectContainer) GetByName(name string) (interface{}, *ObjectContainerDefinition, RegistryType) {
	oc.RLock()
	defer oc.RUnlock()

	for _, ocd := range oc.UnsafeObjects {
		o, ot := Armeria.registry.Get(ocd.UUID)
		if strings.ToLower(o.(ContainerObject).Name()) == strings.ToLower(name) {
			return o, ocd, ot
		}
	}

	return nil, nil, RegistryTypeUnknown
}

// GetByUUIDOrName attempts to retrieve an object by it's uuid, and then by it's name.
func (oc *ObjectContainer) GetByUUIDOrName(UUIDOrName string) (interface{}, *ObjectContainerDefinition, RegistryType) {
	if obj, ocd, rt := oc.Get(UUIDOrName); rt != RegistryTypeUnknown {
		return obj, ocd, rt
	}

	return oc.GetByName(UUIDOrName)
}

// Slot returns the slot that the uuid is within. If the uuid does not exist, slot 0 will be returned,
// which could result in a false positive. Check existance of the uuid before using this function.
func (oc *ObjectContainer) Slot(uuid string) int {
	o, ocd, _ := oc.Get(uuid)
	if o == nil {
		return 0
	}

	return ocd.Slot
}

// SetSlot explicitly sets an item slot without checking if another item already exists in that slot. Use this
// function carefully and as-needed (ie: swapping items).
func (oc *ObjectContainer) SetSlot(uuid string, slot int) {
	_, ocd, _ := oc.Get(uuid)
	ocd.Slot = slot
}

// AtSlot retrieves an object from a specific slot, or nil if the container is unbounded.
func (oc *ObjectContainer) AtSlot(slot int) (interface{}, *ObjectContainerDefinition, RegistryType) {
	oc.RLock()
	defer oc.RUnlock()

	if oc.UnsafeMaxSize == 0 {
		return nil, nil, RegistryTypeUnknown
	}

	for _, ocd := range oc.UnsafeObjects {
		if ocd.Slot == slot {
			o, ot := Armeria.registry.Get(ocd.UUID)
			return o, ocd, ot
		}
	}

	return nil, nil, RegistryTypeUnknown
}

// Count returns the number of objects within the container.
func (oc *ObjectContainer) Count() int {
	oc.RLock()
	defer oc.RUnlock()

	return len(oc.UnsafeObjects)
}

// All returns all the objects within the container.
func (oc *ObjectContainer) All() []interface{} {
	oc.RLock()
	defer oc.RUnlock()

	var everything []interface{}
	for _, ocd := range oc.UnsafeObjects {
		o, _ := Armeria.registry.Get(ocd.UUID)
		everything = append(everything, o)
	}

	return everything
}

// Sync will make sure all objects within the container are registered to the global registry. Note that this
// will NOT remove objects from the global registry that have since been removed from the container. The Remove()
// function on the ObjectContainer will handle that.
func (oc *ObjectContainer) Sync() {
	oc.RLock()
	defer oc.RUnlock()

	for _, ocd := range oc.UnsafeObjects {
		Armeria.registry.RegisterContainerObject(ocd.UUID, oc)
	}
}

// Characters returns all Character objects from the container.
func (oc *ObjectContainer) Characters(onlineOnly bool, exceptions ...*Character) []*Character {
	oc.RLock()
	defer oc.RUnlock()

	var chars []*Character
	var exceptionIds []string

	for _, c := range exceptions {
		exceptionIds = append(exceptionIds, c.ID())
	}

	for _, ocd := range oc.UnsafeObjects {
		c, ot := Armeria.registry.Get(ocd.UUID)
		if ot == RegistryTypeCharacter {
			char := c.(*Character)
			if !onlineOnly || char.Online() {
				if len(exceptionIds) == 0 || !misc.Contains(exceptionIds, char.ID()) {
					chars = append(chars, char)
				}
			}
		}
	}

	return chars
}

// Mobs returns all MobInstance objects from the container.
func (oc *ObjectContainer) Mobs() []*MobInstance {
	oc.RLock()
	defer oc.RUnlock()

	var mobs []*MobInstance
	for _, ocd := range oc.UnsafeObjects {
		m, ot := Armeria.registry.Get(ocd.UUID)
		if ot == RegistryTypeMobInstance {
			mobs = append(mobs, m.(*MobInstance))
		}
	}

	return mobs
}

// Items returns all ItemInstance objects from the container.
func (oc *ObjectContainer) Items() []*ItemInstance {
	oc.RLock()
	defer oc.RUnlock()

	var items []*ItemInstance
	for _, ocd := range oc.UnsafeObjects {
		i, ot := Armeria.registry.Get(ocd.UUID)
		if ot == RegistryTypeItemInstance {
			items = append(items, i.(*ItemInstance))
		}
	}

	return items
}

// NextAvailableSlot returns the next unused slot within the container.
func (oc *ObjectContainer) NextAvailableSlot() (int, error) {
	if oc.UnsafeMaxSize == 0 {
		return 0, nil
	}

	for s := 0; s < oc.UnsafeMaxSize; s++ {
		o, _, _ := oc.AtSlot(s)
		if o == nil {
			return s, nil
		}
	}

	return 0, ErrContainerNoRoom
}

// Remove removes an object from the container.
func (oc *ObjectContainer) Remove(uuid string) {
	oc.Lock()
	defer oc.Unlock()

	for i, ocd := range oc.UnsafeObjects {
		if ocd.UUID == uuid {
			oc.UnsafeObjects[i] = oc.UnsafeObjects[len(oc.UnsafeObjects)-1]
			oc.UnsafeObjects = oc.UnsafeObjects[:len(oc.UnsafeObjects)-1]
		}
	}

	Armeria.registry.UnregisterContainerObject(uuid)
}

// Add attempts to add an object to the container. This can fail if the object already exists within the container
// or if the container is already at the maximum size.
func (oc *ObjectContainer) Add(uuid string) error {
	if oc.Contains(uuid) {
		return ErrContainerDuplicate
	}

	if oc.UnsafeMaxSize > 0 && len(oc.UnsafeObjects) >= oc.UnsafeMaxSize {
		return ErrContainerNoRoom
	}

	slot, err := oc.NextAvailableSlot()
	if err != nil {
		return err
	}

	ocd := &ObjectContainerDefinition{
		UUID: uuid,
		Slot: slot,
	}

	oc.Lock()
	defer oc.Unlock()
	oc.UnsafeObjects = append(oc.UnsafeObjects, ocd)

	Armeria.registry.RegisterContainerObject(uuid, oc)

	return nil
}

// PopulateFromLedger ensures at least one entry from the ledger, with a buy price, exists within the
// object container.
func (oc *ObjectContainer) PopulateFromLedger(ledger *Ledger) {
	for _, entry := range ledger.Entries() {
		if entry.BuyPrice > 0 {
			item := Armeria.itemManager.ItemByName(entry.ItemName)
			if item != nil {
				_, _, rt := oc.GetByName(item.Name())
				if rt == RegistryTypeUnknown {
					ii := item.CreateInstance()
					oc.Add(ii.ID())
				}
			}
		}
	}
}
