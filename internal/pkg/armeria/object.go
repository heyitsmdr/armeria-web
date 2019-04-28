package armeria

type Object interface {
	GetType() int
	GetName() string
	GetFName() string
}

const ObjectTypeCharacter int = 0
