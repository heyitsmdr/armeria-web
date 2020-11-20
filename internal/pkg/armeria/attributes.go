package armeria

import (
	"strconv"
	"strings"
)

const (
	AttributeChannels    string = "channels"
	AttributeColor       string = "color"
	AttributeDescription string = "description"
	AttributeDown        string = "down"
	AttributeEast        string = "east"
	AttributeGender      string = "gender"
	AttributeMoney       string = "money"
	AttributeMusic       string = "music"
	AttributeNorth       string = "north"
	AttributePermissions string = "permissions"
	AttributePicture     string = "picture"
	AttributeRarity      string = "rarity"
	AttributeScript      string = "script"
	AttributeSouth       string = "south"
	AttributeTitle       string = "title"
	AttributeType        string = "type"
	AttributeUp          string = "up"
	AttributeWest        string = "west"

	TempAttributeEditorOpen string = "editorOpen"
	TempAttributeGhost      string = "ghost"
	TempAttributeReplyTo    string = "replyTo"
)

// AttributeEditorType returns the object editor "type" for each attribute
func AttributeEditorType(ot ObjectType, attr string) string {
	switch attr {
	case AttributePicture:
		return "picture"
	case AttributeScript:
		return "script"
	case AttributeMusic:
		return "enum:track-one|track-two"
	case AttributeRarity:
		return "enum:common|uncommon"
	case AttributeGender:
		return "enum:male|female|thing"
	case AttributeColor:
		return "color"
	case AttributeType:
		switch ot {
		case ObjectTypeItem:
			return "enum:generic|mob-spawner"
		default:
			return "editable"
		}
	}

	return "editable"
}

// ValidAreaAttributes returns an array of valid attributes that can be permanently set.
func ValidAreaAttributes() []string {
	return []string{
		AttributeMusic,
	}
}

// ValidCharacterAttributes returns an array of valid attributes that can be permanently set.
func ValidCharacterAttributes() []string {
	return []string{
		AttributePicture,
		AttributeTitle,
		AttributePermissions,
		AttributeChannels,
		AttributeGender,
		AttributeMoney,
	}
}

// ValidItemAttributes returns an array of valid attributes that can be permanently set.
func ValidItemAttributes() []string {
	return []string{
		AttributePicture,
		AttributeType,
		AttributeRarity,
		AttributeDescription,
	}
}

// ValidItemInstanceAttributes returns an array of attributes that can be overriden from the parent.
func ValidItemInstanceAttributes() []string {
	return []string{
		AttributeRarity,
		AttributeDescription,
	}
}

// ValidMobInstanceAttributes returns an array of attributes that can be overriden from the parent.
func ValidMobInstanceAttributes() []string {
	return []string{
		AttributeGender,
		AttributeTitle,
	}
}

// ValidMobAttributes returns an array of valid attributes that can be permanently set.
func ValidMobAttributes() []string {
	return []string{
		AttributePicture,
		AttributeScript,
		AttributeGender,
		AttributeTitle,
	}
}

// ValidRoomAttributes returns an array of valid attributes that can be permanently set.
func ValidRoomAttributes() []string {
	return []string{
		AttributeTitle,
		AttributeDescription,
		AttributeColor,
		AttributeType,
		AttributeNorth,
		AttributeEast,
		AttributeSouth,
		AttributeWest,
		AttributeUp,
		AttributeDown,
	}
}

// AreaAttributeDefault returns the default value for a particular attribute.
func AreaAttributeDefault(name string) string {
	switch name {

	}

	return ""
}

// CharacterAttributeDefault returns the default value for a particular attribute.
func CharacterAttributeDefault(name string) string {
	switch name {
	case AttributeGender:
		return "male"
	case AttributeMoney:
		return "0"
	}

	return ""
}

// ItemAttributeDefault returns the default value for a particular attribute.
func ItemAttributeDefault(name string) string {
	switch name {
	case AttributeRarity:
		return "common"
	}

	return ""
}

// MobAttributeDefault returns the default value for a particular attribute.
func MobAttributeDefault(name string) string {
	switch name {
	case AttributeGender:
		return "thing"
	}

	return ""
}

// RoomAttributeDefault returns the default value for a particular attribute.
func RoomAttributeDefault(name string) string {
	switch name {
	case AttributeTitle:
		return "Empty Room"
	case AttributeDescription:
		return "You are in a newly created empty room. Make it a good one!"
	case AttributeColor:
		return "190,190,190"
	}

	return ""
}

// ValidateItemAttribute returns a bool indicating whether a particular value is allowed
// for a particular attribute.
func ValidateItemAttribute(name string, value string) (bool, string) {
	return true, ""
}

// ValidateMobAttribute returns a bool indicating whether a particular value is allowed
// for a particular attribute.
func ValidateMobAttribute(name, value string) (bool, string) {
	switch name {
	case AttributeScript:
		return false, "script cannot be set explicitly"
	case AttributeGender:
		vlc := strings.ToLower(value)
		if vlc != "male" && vlc != "female" && vlc != "thing" {
			return false, "gender can only be male, female, or thing"
		}
	}

	return true, ""
}

// ValidateCharacterAttribute returns a bool indicating whether a particular value is allowed
// for a particular attribute.
func ValidateCharacterAttribute(name, value string) (bool, string) {
	switch name {
	case AttributeGender:
		vlc := strings.ToLower(value)
		if vlc != "male" && vlc != "female" {
			return false, "gender can only be male or female"
		}
	case AttributeMoney:
		if _, err := strconv.ParseFloat(value, 64); err != nil {
			return false, "money can only be a numeric value"
		}
	}

	return true, ""
}
