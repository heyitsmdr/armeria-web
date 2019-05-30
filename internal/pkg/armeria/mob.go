package armeria

import (
	"armeria/internal/pkg/misc"
	"fmt"
	"strconv"
	"sync"
	"time"

	"go.uber.org/zap"
)

type Mob struct {
	Name       string            `json:"name"`
	Attributes map[string]string `json:"attributes"`
	Instances  []*MobInstance    `json:"instances"`
	mux        sync.Mutex
}

// GetValidMobAttributes returns an array of valid attributes that can be permanently set.
func GetValidMobAttributes() []string {
	return []string{
		"picture",
		"script",
	}
}

// GetMobAttributeDefault returns the default value for a particular attribute.
func GetMobAttributeDefault(name string) string {
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

// SetAttribute sets a permanent attribute and only valid attributes can be set.
func (m *Mob) SetAttribute(name string, value string) {
	m.mux.Lock()
	defer m.mux.Unlock()

	if m.Attributes == nil {
		m.Attributes = make(map[string]string)
	}

	if !misc.Contains(GetValidMobAttributes(), name) {
		Armeria.log.Fatal("attempted to set invalid attribute",
			zap.String("attribute", name),
			zap.String("value", value),
		)
	}

	m.Attributes[name] = value
}

// GetAttribute returns a permanent attribute.
func (m *Mob) GetAttribute(name string) string {
	m.mux.Lock()
	defer m.mux.Unlock()

	if len(m.Attributes[name]) == 0 {
		return GetMobAttributeDefault(name)
	}

	return m.Attributes[name]
}

// GetEditorData returns the JSON used for the object editor.
func (m *Mob) GetEditorData() *ObjectEditorData {
	var props []*ObjectEditorDataProperty
	for _, attrName := range GetValidMobAttributes() {
		propType := "editable"
		if attrName == "picture" {
			propType = "picture"
		} else if attrName == "script" {
			propType = "script"
		}

		props = append(props, &ObjectEditorDataProperty{
			PropType: propType,
			Name:     attrName,
			Value:    m.GetAttribute(attrName),
		})
	}

	return &ObjectEditorData{
		Name:       m.Name,
		ObjectType: "mob",
		Properties: props,
	}
}

func (m *Mob) CreateInstance(loc *Location) *MobInstance {
	m.mux.Lock()
	defer m.mux.Unlock()

	mi := &MobInstance{
		Id:       strconv.FormatInt(time.Now().UnixNano(), 10),
		Parent:   m.Name,
		Location: loc,
	}

	m.Instances = append(m.Instances, mi)

	return mi
}

func (m *Mob) GetInstanceById(id string) *MobInstance {
	m.mux.Lock()
	defer m.mux.Unlock()

	for _, mi := range m.Instances {
		if mi.Id == id {
			return mi
		}
	}

	return nil
}

func (m *Mob) GetScriptFile() string {
	return fmt.Sprintf("%s/scripts/mob-%s.lua", Armeria.dataPath, m.Name)
}
