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
	ErrNoRoom = errors.New("no space in container")
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

// Parent returns the parent object that owns this object container.
func (oc *ObjectContainer) Parent() interface{} {
	oc.RLock()
	defer oc.RUnlock()

	return oc.UnsafeParent
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

// Get retrieves an object from the object container.
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
