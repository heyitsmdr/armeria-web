package armeria

import (
	"armeria/internal/pkg/misc"
	"sync"

	"github.com/google/uuid"

	"go.uber.org/zap"
)

type Item struct {
	sync.RWMutex
	UnsafeName       string            `json:"name"`
	UnsafeAttributes map[string]string `json:"attributes"`
	UnsafeInstances  []*ItemInstance   `json:"instances"`
}

// Init is called when the Item is created or loaded from disk.
func (i *Item) Init() {}

// Name returns the name of the Item.
func (i *Item) Name() string {
	i.RLock()
	defer i.RUnlock()

	return i.UnsafeName
}

// Instances returns all of the ItemInstance instances.
func (i *Item) Instances() []*ItemInstance {
	i.RLock()
	defer i.RUnlock()

	return i.UnsafeInstances
}

// CreateInstance creates a new ItemInstance and adds it in-memory.
func (i *Item) CreateInstance() *ItemInstance {
	i.Lock()
	defer i.Unlock()

	ii := &ItemInstance{
		UUID:             uuid.New().String(),
		UnsafeAttributes: make(map[string]string),
		Parent:           i,
	}

	i.UnsafeInstances = append(i.UnsafeInstances, ii)

	ii.Init()

	return ii
}

// DeleteInstance removes the ItemInstance from memory.
func (i *Item) DeleteInstance(ii *ItemInstance) bool {
	i.Lock()
	defer i.Unlock()

	ii.Deinit()

	for idx, inst := range i.UnsafeInstances {
		if inst.ID() == ii.ID() {
			i.UnsafeInstances[idx] = i.UnsafeInstances[len(i.UnsafeInstances)-1]
			i.UnsafeInstances = i.UnsafeInstances[:len(i.UnsafeInstances)-1]
			return true
		}
	}

	return false
}

// Attribute returns a permanent attribute.
func (i *Item) Attribute(name string) string {
	i.RLock()
	defer i.RUnlock()

	if len(i.UnsafeAttributes[name]) == 0 {
		return AttributeDefault(ObjectTypeItem, name)
	}

	return i.UnsafeAttributes[name]
}

// SetAttribute sets a permanent attribute and only valid attributes can be set.
func (i *Item) SetAttribute(name string, value string) {
	i.Lock()
	defer i.Unlock()

	if !misc.Contains(AttributeList(ObjectTypeItem), name) {
		Armeria.log.Fatal("attempted to set invalid attribute",
			zap.String("attribute", name),
			zap.String("value", value),
		)
	}

	i.UnsafeAttributes[name] = value
}

// EditorData returns the JSON used for the object editor.
func (i *Item) EditorData() *ObjectEditorData {
	var props []*ObjectEditorDataProperty
	for _, attrName := range AttributeList(ObjectTypeItem) {
		props = append(props, &ObjectEditorDataProperty{
			PropType: AttributeEditorType(ObjectTypeItem, attrName),
			Name:     attrName,
			Value:    i.Attribute(attrName),
		})
	}

	return &ObjectEditorData{
		Name:       i.UnsafeName,
		ObjectType: "item",
		Properties: props,
	}
}
