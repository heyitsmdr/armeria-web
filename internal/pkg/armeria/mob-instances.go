package armeria

import (
	"armeria/internal/pkg/misc"
	"fmt"
	"sync"

	"go.uber.org/zap"
)

type MobInstance struct {
	Parent     string            `json:"parent"`
	Location   *Location         `json:"location"`
	Attributes map[string]string `json:"attributes"`
	mux        sync.Mutex
}

// GetParent returns the Mob parent.
func (mi *MobInstance) GetParent() *Mob {
	return Armeria.mobManager.GetMobByName(mi.Parent)
}

// GetLocation returns the location of the mob.
func (mi *MobInstance) GetLocation() *Location {
	mi.mux.Lock()
	defer mi.mux.Unlock()
	return mi.Location
}

// GetType returns the object type, since Mob implements the Object interface.
func (mi *MobInstance) GetType() int {
	return ObjectTypeMob
}

// GetName returns the raw mob name.
func (mi *MobInstance) GetName() string {
	mi.mux.Lock()
	defer mi.mux.Unlock()
	return mi.Parent
}

// GetFName returns the formatted mob name.
func (mi *MobInstance) GetFName() string {
	mi.mux.Lock()
	defer mi.mux.Unlock()
	return fmt.Sprintf("[b]%s[/b]", mi.Parent)
}

// SetAttribute sets a permanent attribute on the mob instance.
func (mi *MobInstance) SetAttribute(name string, value string) {
	mi.mux.Lock()
	defer mi.mux.Unlock()

	if mi.Attributes == nil {
		mi.Attributes = make(map[string]string)
	}

	if !misc.Contains(GetValidMobAttributes(), name) {
		Armeria.log.Fatal("attempted to set invalid attribute",
			zap.String("attribute", name),
			zap.String("value", value),
		)
	}

	mi.Attributes[name] = value
}

// GetAttribute returns an attribute on the mob instance, and falls back to the parent Mob.
func (mi *MobInstance) GetAttribute(name string) string {
	mi.mux.Lock()
	defer mi.mux.Unlock()

	if len(mi.Attributes[name]) == 0 {
		return mi.GetParent().GetAttribute(name)
	}

	return mi.Attributes[name]
}
