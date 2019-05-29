package armeria

// Object is the interface that describes an in-game object
type Object interface {
	GetType() int
	GetName() string
	GetFName() string
	GetAttribute(name string) string
	SetAttribute(name string, value string)
}

const (
	ObjectTypeCharacter int = 0
	ObjectTypeMob       int = 1
)
