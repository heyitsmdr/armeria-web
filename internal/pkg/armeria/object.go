package armeria

type Object interface {
	GetType() int
	GetName() string
	GetFName() string
}

const OBJECT_TYPE_CHARACTER int = 0
