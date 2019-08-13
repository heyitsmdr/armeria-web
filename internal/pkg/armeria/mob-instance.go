package armeria

import (
	"armeria/internal/pkg/misc"
	"errors"
	"fmt"
	"sync"
)

type MobInstance struct {
	sync.RWMutex
	UUID             string            `json:"uuid"`
	UnsafeAttributes map[string]string `json:"attributes"`
	Parent           *Mob              `json:"-"`
}

// Init is called when the MobInstance is created or loaded from disk.
func (mi *MobInstance) Init() {
	Armeria.registry.Register(mi, mi.ID(), RegistryTypeMobInstance)
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
	return fmt.Sprintf("[b]%s[/b]", mi.Parent.Name())
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

// MobInstance returns the MobInstance's Room based on the object container it is within.
func (mi *MobInstance) Room() *Room {
	oc := Armeria.registry.GetObjectContainer(mi.ID())
	if oc == nil {
		return nil
	}
	return oc.ParentRoom()
}
