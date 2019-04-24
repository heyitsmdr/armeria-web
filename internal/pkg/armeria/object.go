package armeria

type Object interface {
	GetType() int
	GetName() string
}

const (
	OBJECT_TYPE_CHARACTER int = 0
)
