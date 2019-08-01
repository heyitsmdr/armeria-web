package armeria

// Registry is the global registry for anything in-game that has a UUID. This registry can act as a quick
// lookup and retrieval for any game object based on it's UUID. Note that when retrieving an object from
// the registry, you will need to cast it to the appropriate type.
type Registry struct {
	entries map[string]interface{}
	types   map[string]RegistryType
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
)

// Register registers a unique object instance with the global registry.
func (r *Registry) Register(o interface{}, uuid string, rt RegistryType) {
	r.entries[uuid] = o
	r.types[uuid] = rt
}

// Unregister removes a unique object instance from the global registry.
func (r *Registry) Unregister(uuid string) {
	delete(r.entries, uuid)
	delete(r.types, uuid)
}

// Get returns a unique object instance from the global registry.
func (r *Registry) Get(uuid string) (interface{}, RegistryType) {
	if r.entries[uuid] == nil {
		return nil, RegistryTypeUnknown
	}

	return r.entries[uuid], r.types[uuid]
}
