package armeria

import "sync"

// Registry is an in-memory global registry for anything in-game that has a UUID. This registry can act as a
// quick lookup and retrieval for any game object based on it's UUID. Note that when retrieving an object from
// the registry, you will need to cast it to the appropriate type.
type Registry struct {
	sync.RWMutex
	entries          map[string]interface{}
	types            map[string]RegistryType
	containerEntries map[string]*ObjectContainer
}

func NewRegistry() *Registry {
	return &Registry{
		entries: make(map[string]interface{}),
		types:   make(map[string]RegistryType),
	}
}

type RegistryType int

const (
	RegistryTypeUnknown RegistryType = iota
	RegistryTypeCharacter
	RegistryTypeItemInstance
	RegistryTypeMobInstance
	RegistryTypeArea
	RegistryTypeObjectContainer
)

// Register registers a unique object instance with the global registry.
func (r *Registry) Register(o interface{}, uuid string, rt RegistryType) {
	r.Lock()
	defer r.Unlock()

	r.entries[uuid] = o
	r.types[uuid] = rt
}

// Unregister removes a unique object instance from the global registry.
func (r *Registry) Unregister(uuid string) {
	r.Lock()
	defer r.Unlock()

	delete(r.entries, uuid)
	delete(r.types, uuid)
}

// RegisterContainerObject registers a unique object instance with the global registry.
func (r *Registry) RegisterContainerObject(ouuid string, oc *ObjectContainer) {
	r.Lock()
	defer r.Unlock()

	r.containerEntries[ouuid] = oc
}

// RegisterContainerObject registers a unique object instance with the global registry.
func (r *Registry) UnregisterContainerObject(ouuid string) {
	r.Lock()
	defer r.Unlock()

	delete(r.containerEntries, ouuid)
}

// Get returns a unique object instance from the global registry.
func (r *Registry) Get(uuid string) (interface{}, RegistryType) {
	r.RLock()
	defer r.RUnlock()

	if r.entries[uuid] == nil {
		return nil, RegistryTypeUnknown
	}

	return r.entries[uuid], r.types[uuid]
}

// GetAllFromType returns all objects matching a specific type.
func (r *Registry) GetAllFromType(rt RegistryType) []interface{} {
	r.RLock()
	defer r.RUnlock()

	var o []interface{}

	for uuid, t := range r.types {
		if t == rt {
			o = append(o, r.entries[uuid])
		}
	}

	return o
}

// GetObjectContainer returns the ObjectContainer that an Object is within.
func (r *Registry) GetObjectContainer(ouuid string) *ObjectContainer {
	r.RLock()
	defer r.RUnlock()

	if r.containerEntries[ouuid] == nil {
		return nil
	}

	return r.containerEntries[ouuid]
}
