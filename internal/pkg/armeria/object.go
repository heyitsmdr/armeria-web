package armeria

// Object is the interface that describes an in-game object, which is
// implemented by Character, MobInstance, and ItemInstance.
type Object interface {
	Id() string
	Type() int
	Name() string
	FormattedName() string
	Attribute(name string) string
	SetAttribute(name string, value string)
}

const (
	ObjectTypeCharacter int = 0
	ObjectTypeMob       int = 1
	ObjectTypeItem      int = 2
)
