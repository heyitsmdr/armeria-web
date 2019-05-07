package armeria

// Object is the interface that describes an in-game object
type Object interface {
	GetType() int
	GetName() string
	GetFName() string
}

const ObjectTypeCharacter int = 0
