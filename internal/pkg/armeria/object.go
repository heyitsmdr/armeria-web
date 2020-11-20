package armeria

type ObjectType string

const (
	ObjectTypeMob          ObjectType = "mob"
	ObjectTypeMobInstance  ObjectType = "mob-instance"
	ObjectTypeItem         ObjectType = "item"
	ObjectTypeItemInstance ObjectType = "item-instance"
	ObjectTypeCharacter    ObjectType = "character"
	ObjectTypeRoom         ObjectType = "room"
	ObjectTypeArea         ObjectType = "area"
)
