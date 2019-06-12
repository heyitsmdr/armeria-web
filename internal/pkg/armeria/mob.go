package armeria

import (
	"armeria/internal/pkg/misc"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
)

type Mob struct {
	UnsafeName       string            `json:"name"`
	UnsafeAttributes map[string]string `json:"attributes"`
	UnsafeInstances  []*MobInstance    `json:"instances"`
	mux              sync.Mutex
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
	m.mux.Lock()
	defer m.mux.Unlock()
	return m.UnsafeName
}

// SetAttribute sets a permanent attribute and only valid attributes can be set.
func (m *Mob) SetAttribute(name string, value string) {
	m.mux.Lock()
	defer m.mux.Unlock()

	if m.UnsafeAttributes == nil {
		m.UnsafeAttributes = make(map[string]string)
	}

	if !misc.Contains(ValidMobAttributes(), name) {
		Armeria.log.Fatal("attempted to set invalid attribute",
			zap.String("attribute", name),
			zap.String("value", value),
		)
	}

	m.UnsafeAttributes[name] = value
}

// Attribute returns a permanent attribute.
func (m *Mob) Attribute(name string) string {
	m.mux.Lock()
	defer m.mux.Unlock()

	if len(m.UnsafeAttributes[name]) == 0 {
		return MobAttributeDefault(name)
	}

	return m.UnsafeAttributes[name]
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

// CreateInstance creates a new MobInstance, adds it to the Mob and returns the MobInstance.
func (m *Mob) CreateInstance(loc *Location) *MobInstance {
	m.mux.Lock()
	defer m.mux.Unlock()

	mi := &MobInstance{
		UnsafeId:       strconv.FormatInt(time.Now().UnixNano(), 10),
		UnsafeParent:   m.UnsafeName,
		UnsafeLocation: loc,
	}

	m.UnsafeInstances = append(m.UnsafeInstances, mi)

	return mi
}

// InstanceById returns a MobInstance by the instance identifier.
func (m *Mob) InstanceById(id string) *MobInstance {
	m.mux.Lock()
	defer m.mux.Unlock()

	for _, mi := range m.UnsafeInstances {
		if mi.UnsafeId == id {
			return mi
		}
	}

	return nil
}

// ScriptFile returns the full path to the associated Lua script file.
func (m *Mob) ScriptFile() string {
	return fmt.Sprintf(
		"%s/scripts/mob-%s.lua",
		Armeria.dataPath,
		strings.ReplaceAll(m.UnsafeName, " ", "-"),
	)
}
