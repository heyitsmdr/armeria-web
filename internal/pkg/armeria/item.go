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

const (
	ItemTypeGeneric    string = "generic"
	ItemTypeMobSpawner        = "mob-spawner"
	ItemTypeTrashCan          = "trash-can"
	ItemTypeBreadcrumb        = "mob-breadcrumb"
	ItemTypeBankCard          = "bank-card"

	ItemRarityCommon   string = "common"
	ItemRarityUncommon        = "uncommon"
)

// ItemTypes return the possible item types.
func ItemTypes() []string {
	return []string{
		ItemTypeGeneric,
		ItemTypeMobSpawner,
		ItemTypeBreadcrumb,
		ItemTypeTrashCan,
		ItemTypeBankCard,
	}
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

	Armeria.log.Info("instance created",
		zap.String("uuid", ii.ID()),
		zap.String("name", i.UnsafeName),
	)

	return ii
}

// DeleteInstance uninitializes the ItemInstance, unregisters it from the registrar, and
// removes it from memory.
func (i *Item) DeleteInstance(ii *ItemInstance) bool {
	i.Lock()
	defer i.Unlock()

	ii.Deinit()

	for idx, inst := range i.UnsafeInstances {
		if inst.ID() == ii.ID() {
			i.UnsafeInstances[idx] = i.UnsafeInstances[len(i.UnsafeInstances)-1]
			i.UnsafeInstances = i.UnsafeInstances[:len(i.UnsafeInstances)-1]
			Armeria.log.Info("instance deleted",
				zap.String("uuid", ii.ID()),
				zap.String("name", i.UnsafeName),
			)
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
