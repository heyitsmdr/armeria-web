package armeria

import (
	"armeria/internal/pkg/misc"
	"fmt"
	"sync"

	"go.uber.org/zap"
)

type MobInstance struct {
	UnsafeId         string            `json:"id"`
	UnsafeParent     string            `json:"parent"`
	UnsafeLocation   *Location         `json:"location"`
	UnsafeAttributes map[string]string `json:"attributes"`
	mux              sync.Mutex
}

// Parent returns the Mob parent.
func (mi *MobInstance) Parent() *Mob {
	return Armeria.mobManager.MobByName(mi.UnsafeParent)
}

// Location returns the location of the mob.
func (mi *MobInstance) Location() *Location {
	mi.mux.Lock()
	defer mi.mux.Unlock()
	return mi.UnsafeLocation
}

// Room returns the Room of the mob.
func (mi *MobInstance) Room() *Room {
	return Armeria.worldManager.RoomFromLocation(mi.Location())
}

// Type returns the object type, since Mob implements the Object interface.
func (mi *MobInstance) Type() int {
	return ObjectTypeMob
}

// UnsafeName returns the raw mob name.
func (mi *MobInstance) Name() string {
	mi.mux.Lock()
	defer mi.mux.Unlock()
	return mi.UnsafeParent
}

// FormattedName returns the formatted mob name.
func (mi *MobInstance) FormattedName() string {
	mi.mux.Lock()
	defer mi.mux.Unlock()
	return fmt.Sprintf("[b]%s[/b]", mi.UnsafeParent)
}

// SetAttribute sets a permanent attribute on the mob instance.
func (mi *MobInstance) SetAttribute(name string, value string) {
	mi.mux.Lock()
	defer mi.mux.Unlock()

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

// Attribute returns an attribute on the mob instance, and falls back to the parent Mob.
func (mi *MobInstance) Attribute(name string) string {
	mi.mux.Lock()
	defer mi.mux.Unlock()

	if len(mi.UnsafeAttributes[name]) == 0 {
		return mi.Parent().Attribute(name)
	}

	return mi.UnsafeAttributes[name]
}
