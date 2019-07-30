package armeria

import (
	"errors"
	"sync"
)

type ObjectContainer struct {
	sync.RWMutex
	UnsafeItems   []*ObjectContainerDefinition `json:"objects"`
	UnsafeMaxSize int                          `json:"maxSize"`
}

type ObjectContainerDefinition struct {
	Id     string `json:"id"`
	Slot   int    `json:"slot"`
	Object Object `json:"-"`
}

var (
	ErrNoRoom = errors.New("no space in container")
)

func NewObjectContainer(maxSize int) *ObjectContainer {
	return &ObjectContainer{
		UnsafeItems:   make([]*ObjectContainerDefinition, 0),
		UnsafeMaxSize: maxSize,
	}
}
