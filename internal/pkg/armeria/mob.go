package armeria

import (
	"armeria/internal/pkg/misc"
	"fmt"
	"strings"
	"sync"

	"go.uber.org/zap"

	"github.com/google/uuid"
)

type Mob struct {
	sync.RWMutex
	UnsafeName       string            `json:"name"`
	UnsafeAttributes map[string]string `json:"attributes"`
	UnsafeInstances  []*MobInstance    `json:"instances"`
}

// ValidMobAttributes returns an array of valid attributes that can be permanently set.
func ValidMobAttributes() []string {
	return []string{
		"picture",
		"script",
	}
}

// MobAttributeDefault returns the default value for a particular attribute.
func MobAttributeDefault(name string) string {
	switch name {

	}

	return ""
}

// ValidateMobAttribute returns a bool indicating whether a particular value is allowed
// for a particular attribute.
func ValidateMobAttribute(name string, value string) (bool, string) {
	switch name {
	case "script":
		return false, "script cannot be set explicitly"
	}

	return true, ""
}

// Name returns the name of the Mob.
func (m *Mob) Name() string {
	m.RLock()
	defer m.RUnlock()
	return m.UnsafeName
}

// Attribute returns a permanent attribute.
func (m *Mob) Attribute(name string) string {
	m.RLock()
	defer m.RUnlock()

	if len(m.UnsafeAttributes[name]) == 0 {
		return MobAttributeDefault(name)
	}

	return m.UnsafeAttributes[name]
}

// SetAttribute sets a permanent attribute and only valid attributes can be set.
func (m *Mob) SetAttribute(name string, value string) {
	m.Lock()
	defer m.Unlock()

	if !misc.Contains(ValidMobAttributes(), name) {
		Armeria.log.Fatal("attempted to set invalid attribute",
			zap.String("attribute", name),
			zap.String("value", value),
		)
	}

	m.UnsafeAttributes[name] = value
}

// EditorData returns the JSON used for the object editor.
func (m *Mob) EditorData() *ObjectEditorData {
	var props []*ObjectEditorDataProperty
	for _, attrName := range ValidMobAttributes() {
		propType := "editable"
		if attrName == "picture" {
			propType = "picture"
		} else if attrName == "script" {
			propType = "script"
		}

		props = append(props, &ObjectEditorDataProperty{
			PropType: propType,
			Name:     attrName,
			Value:    m.Attribute(attrName),
		})
	}

	return &ObjectEditorData{
		Name:       m.UnsafeName,
		ObjectType: "mob",
		Properties: props,
	}
}

// CreateInstance creates a new MobInstance, adds it to the Mob
// and returns the MobInstance.
func (m *Mob) CreateInstance(l *Location) *MobInstance {
	m.Lock()
	defer m.Unlock()

	mi := &MobInstance{
		UUID:             uuid.New().String(),
		UnsafeParent:     m.UnsafeName,
		Location:         CopyLocation(l),
		UnsafeAttributes: make(map[string]string),
	}

	m.UnsafeInstances = append(m.UnsafeInstances, mi)

	return mi
}

// DeleteInstance removes the MobInstance from memory.
func (m *Mob) DeleteInstance(mi *MobInstance) bool {
	m.Lock()
	defer m.Unlock()

	for i, inst := range m.UnsafeInstances {
		if inst.Id() == mi.Id() {
			m.UnsafeInstances[i] = m.UnsafeInstances[len(m.UnsafeInstances)-1]
			m.UnsafeInstances = m.UnsafeInstances[:len(m.UnsafeInstances)-1]
			return true
		}
	}

	return false
}

// InstanceByUUID returns a MobInstance by the instance identifier.
func (m *Mob) InstanceByUUID(uuid string) *MobInstance {
	m.RLock()
	defer m.RUnlock()

	for _, mi := range m.UnsafeInstances {
		if mi.UUID == uuid {
			return mi
		}
	}

	return nil
}

// Instances returns all of the mob instances in memory.
func (m *Mob) Instances() []*MobInstance {
	m.RLock()
	defer m.RUnlock()
	return m.UnsafeInstances
}

// ScriptFile returns the full path to the associated Lua script file.
func (m *Mob) ScriptFile() string {
	m.RLock()
	defer m.RUnlock()
	return fmt.Sprintf(
		"%s/scripts/mob-%s.lua",
		Armeria.dataPath,
		strings.ToLower(strings.ReplaceAll(m.UnsafeName, " ", "-")),
	)
}
