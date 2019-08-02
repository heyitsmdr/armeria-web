package armeria

import (
	"errors"
	"sync"
)

// ObjectContainer is a container of game objects that can be persisted to disk. These can be used for
// things in a room, a character's inventory, a chest, etc.
type ObjectContainer struct {
	sync.RWMutex
	UnsafeItems      []*ObjectContainerDefinition `json:"objects"`
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
		UnsafeItems:   make([]*ObjectContainerDefinition, 0),
		UnsafeMaxSize: maxSize,
	}
}

func (oc *ObjectContainer) AttachParent(p interface{}, t ContainerParentType) {
	oc.Lock()
	defer oc.Unlock()

	oc.UnsafeParent = p
	oc.UnsafeParentType = t
}
