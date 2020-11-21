package armeria

import (
	"armeria/internal/pkg/misc"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"sync"

	"go.uber.org/zap"

	"github.com/google/uuid"
)

type Mob struct {
	sync.RWMutex
	UnsafeName        string            `json:"name"`
	UnsafeAttributes  map[string]string `json:"attributes"`
	UnsafeInstances   []*MobInstance    `json:"instances"`
	UnsafeScript      string            `json:"-"`
	UnsafeScriptFuncs []string          `json:"-"`
}

// Init is called when the Mob is created or loaded from disk.
func (m *Mob) Init() {
	m.CacheScript()
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
		return AttributeDefault(ObjectTypeMob, name)
	}

	return m.UnsafeAttributes[name]
}

// SetAttribute sets a permanent attribute and only valid attributes can be set.
func (m *Mob) SetAttribute(name string, value string) {
	m.Lock()
	defer m.Unlock()

	if !misc.Contains(AttributeList(ObjectTypeMob), name) {
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
	for _, attrName := range AttributeList(ObjectTypeMob) {
		props = append(props, &ObjectEditorDataProperty{
			PropType: AttributeEditorType(ObjectTypeMob, attrName),
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
func (m *Mob) CreateInstance() *MobInstance {
	m.Lock()
	defer m.Unlock()

	mi := &MobInstance{
		UUID:             uuid.New().String(),
		UnsafeAttributes: make(map[string]string),
		UnsafeInventory:  NewObjectContainer(0),
		Parent:           m,
	}

	m.UnsafeInstances = append(m.UnsafeInstances, mi)

	mi.Init()

	Armeria.log.Info("instance created",
		zap.String("uuid", mi.ID()),
		zap.String("name", m.UnsafeName),
	)

	return mi
}

// DeleteInstance removes the MobInstance from memory.
func (m *Mob) DeleteInstance(mi *MobInstance) bool {
	m.Lock()
	defer m.Unlock()

	mi.Deinit()

	for i, inst := range m.UnsafeInstances {
		if inst.ID() == mi.ID() {
			m.UnsafeInstances[i] = m.UnsafeInstances[len(m.UnsafeInstances)-1]
			m.UnsafeInstances = m.UnsafeInstances[:len(m.UnsafeInstances)-1]
			Armeria.log.Info("instance deleted",
				zap.String("uuid", mi.ID()),
				zap.String("name", m.UnsafeName),
			)
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

// scriptFile returns the full path to the associated Lua script file. This DOES NOT request a lock and IS NOT
// thread safe.
func (m *Mob) scriptFile() string {
	return fmt.Sprintf(
		"%s/scripts/mob-%s.lua",
		Armeria.dataPath,
		strings.ToLower(strings.ReplaceAll(m.UnsafeName, " ", "-")),
	)
}

// ScriptFile returns the full path to the associated Lua script file.
func (m *Mob) ScriptFile() string {
	m.RLock()
	defer m.RUnlock()
	return m.scriptFile()
}

// CacheScript reads the script contents and caches the contents and individual functions.
func (m *Mob) CacheScript() {
	m.Lock()
	defer m.Unlock()

	if _, err := os.Stat(m.scriptFile()); err != nil {
		return
	}

	b, err := ioutil.ReadFile(m.scriptFile())
	if err != nil {
		return
	}

	m.UnsafeScriptFuncs = []string{}

	re := regexp.MustCompile("function ([a-zA-Z_]+)")
	matches := re.FindAllSubmatch(b, -1)
	if matches != nil {
		for _, match := range matches {
			m.UnsafeScriptFuncs = append(m.UnsafeScriptFuncs, string(match[1]))
		}
	}

	m.UnsafeScript = string(b)
}

// Script returns the cached script contents.
func (m *Mob) Script() string {
	m.RLock()
	defer m.RUnlock()
	return m.UnsafeScript
}

// ScriptFuncs returns the cached functions within the Mob's script file.
func (m *Mob) ScriptFuncs() []string {
	m.RLock()
	defer m.RUnlock()
	return m.UnsafeScriptFuncs
}
