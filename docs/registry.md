# Global Registry

The global registry, `registry.go`, stores references to the in-memory objects based on their UUID. This registry
is not persisted to disk and is re-constructed at startup time. When objects are created or deleted while Armeria
is running, they will be registered or unregistered as needed.

The registry supports registration of the following types:

* `RegistryTypeCharacter`
* `RegistryTypeItemInstance`
* `RegistryTypeMobInstance`
* `RegistryTypeArea`
* `RegistryTypeRoom`

Everything within the registry are things that exist in-game in some shape or form.

## Container Objects

The registry also contains a mapping of a `ContainerObject`'s UUID mapped to the `ObjectContainer` it is within. This is
used to quickly lookup which container a `ContainerObject` is within without having to iterate over all `ObjectContainer`
instances. 