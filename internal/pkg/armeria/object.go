package armeria

// Object is the interface that describes an in-game object, which is
// implemented by Character, MobInstance, and ItemInstance.
type Object interface {
	Id() string
	Type() ObjectType
	Name() string
	FormattedName() string
	Attribute(name string) string
	SetAttribute(name string, value string) error
}

type ObjectType int

const (
	ObjectTypeCharacter ObjectType = iota
	ObjectTypeMob
	ObjectTypeItem
)

// ObjectSortOrder returns the sort order for each type of Object. This will affect
// how it will appear in the client's room list. Sorting is in descending order.
func ObjectSortOrder(ot ObjectType) int {
	switch ot {
	case ObjectTypeMob:
		return 75
	case ObjectTypeCharacter:
		return 50
	case ObjectTypeItem:
		return 25
	}

	return 0
}
