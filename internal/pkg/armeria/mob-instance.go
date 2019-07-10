package armeria

import (
	"armeria/internal/pkg/misc"
	"fmt"
	"sync"

	"go.uber.org/zap"
)

type MobInstance struct {
	sync.RWMutex
	UUID             string            `json:"uuid"`
	UnsafeParent     string            `json:"parent"`
	UnsafeLocation   *Location         `json:"location"`
	UnsafeAttributes map[string]string `json:"attributes"`
}

// Id returns the UUID of the instance.
func (mi *MobInstance) Id() string {
	return mi.UUID
}

// Parent returns the Mob parent.
func (mi *MobInstance) Parent() *Mob {
	mi.RLock()
	defer mi.RUnlock()
	return Armeria.mobManager.MobByName(mi.UnsafeParent)
}

// Location returns the location of the MobInstance.
func (mi *MobInstance) Location() *Location {
	mi.RLock()
	defer mi.RUnlock()
	return mi.UnsafeLocation
}

// Room returns the Room of the mob.
func (mi *MobInstance) Room() *Room {
	mi.RLock()
	defer mi.RUnlock()
	return mi.UnsafeLocation.Room()
}

// Type returns the object type, since Mob implements the Object interface.
func (mi *MobInstance) Type() int {
	return ObjectTypeMob
}

// UnsafeName returns the raw Mob name.
func (mi *MobInstance) Name() string {
	mi.RLock()
	defer mi.RUnlock()
	return mi.UnsafeParent
}

// FormattedName returns the formatted Mob name.
func (mi *MobInstance) FormattedName() string {
	mi.RLock()
	defer mi.RUnlock()
	return fmt.Sprintf("[b]%s[/b]", mi.UnsafeParent)
}

// SetAttribute sets a permanent attribute on the MobInstance.
func (mi *MobInstance) SetAttribute(name string, value string) {
	mi.Lock()
	defer mi.Unlock()

	if mi.UnsafeAttributes == nil {
		mi.UnsafeAttributes = make(map[string]string)
	}

	if !misc.Contains(ValidMobAttributes(), name) {
		Armeria.log.Fatal("attempted to set invalid attribute",
			zap.String("attribute", name),
			zap.String("value", value),
		)
	}

	mi.UnsafeAttributes[name] = value
}

// Attribute returns an attribute on the MobInstance, and falls back to the parent Mob.
func (mi *MobInstance) Attribute(name string) string {
	mi.RLock()
	defer mi.RUnlock()

	if len(mi.UnsafeAttributes[name]) == 0 {
		return mi.Parent().Attribute(name)
	}

	return mi.UnsafeAttributes[name]
}
