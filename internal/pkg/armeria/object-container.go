package armeria

import (
	"errors"
	"sync"

	"github.com/google/uuid"
)

// ObjectContainer is a container of game objects that can be persisted to disk. These can be used for
// things in a room, a character's inventory, a chest, etc.
type ObjectContainer struct {
	sync.RWMutex
	UUID             string                       `json:"uuid"`
	UnsafeObjects    []*ObjectContainerDefinition `json:"objects"`
	UnsafeMaxSize    int                          `json:"maxSize"` // 0 = unlimited
	UnsafeParent     interface{}
	UnsafeParentType ContainerParentType
}

type ObjectContainerDefinition struct {
	UUID string `json:"uuid"`
	Slot int    `json:"slot"`
}

var (
	ErrContainerNoRoom    = errors.New("no space in container")
	ErrContainerDuplicate = errors.New("object already in container")
)

type ContainerParentType int

const (
	ContainerParentTypeRoom ContainerParentType = iota
)

func NewObjectContainer(maxSize int) *ObjectContainer {
	return &ObjectContainer{
		UUID:          uuid.New().String(),
		UnsafeObjects: make([]*ObjectContainerDefinition, 0),
		UnsafeMaxSize: maxSize,
	}
}

// Id returns the uuid of the object container.
func (oc *ObjectContainer) Id() string {
	return oc.UUID
}

// AttachParent sets the parent object that the object container belongs to.
func (oc *ObjectContainer) AttachParent(p interface{}, t ContainerParentType) {
	oc.Lock()
	defer oc.Unlock()

	oc.UnsafeParent = p
	oc.UnsafeParentType = t
}

// ParentRoom returns the parent Room if the object has an appropriate parent type.
func (oc *ObjectContainer) ParentRoom() *Room {
	oc.RLock()
	defer oc.RUnlock()

	if oc.UnsafeParentType != ContainerParentTypeRoom {
		return nil
	}

	return oc.UnsafeParent.(*Room)
}

// ParentType returns the ContainerParentType that owns this object container.
func (oc *ObjectContainer) ParentType() ContainerParentType {
	oc.RLock()
	defer oc.RUnlock()

	return oc.UnsafeParentType
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
func (oc *ObjectContainer) Characters(except *Character) []*Character {
	oc.RLock()
	defer oc.RUnlock()

	var chars []*Character
	for _, ocd := range oc.UnsafeObjects {
		c, ot := Armeria.registry.Get(ocd.UUID)
		if ot == RegistryTypeCharacter {
			if except == nil || c.(*Character).Id() != except.Id() {
				chars = append(chars, c.(*Character))
			}
		}
	}

	return chars
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

// Add attempts to add an object to the container. This can fail if the object already exists within the container
// or if the container is already at the maximum size.
func (oc *ObjectContainer) Add(uuid string) error {
	oc.Lock()
	defer oc.Unlock()

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

	oc.UnsafeObjects = append(oc.UnsafeObjects, ocd)

	return nil
}
