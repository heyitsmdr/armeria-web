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

// AttributeList returns the valid attributes for a given ObjectType.
func AttributeList(ot ObjectType) []string {
	switch ot {
	case ObjectTypeCharacter:
		return []string{
			AttributePicture,
			AttributeTitle,
			AttributePermissions,
			AttributeChannels,
			AttributeGender,
			AttributeMoney,
		}
	case ObjectTypeArea:
		return []string{
			AttributeMusic,
		}
	case ObjectTypeRoom:
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
	case ObjectTypeItem:
		return []string{
			AttributePicture,
			AttributeType,
			AttributeRarity,
			AttributeDescription,
		}
	case ObjectTypeItemInstance:
		return []string{
			AttributeRarity,
			AttributeDescription,
		}
	case ObjectTypeMob:
		return []string{
			AttributePicture,
			AttributeScript,
			AttributeGender,
			AttributeTitle,
		}
	case ObjectTypeMobInstance:
		return []string{
			AttributeTitle,
		}
	}

	return []string{}
}

// AttributeEditorType returns the object editor "type" string of an attribute for a given ObjectType.
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

// AttributeDefault returns the default value of an attribute for a given ObjectType.
func AttributeDefault(ot ObjectType, attr string) string {
	switch attr {
	case AttributeGender:
		switch ot {
		case ObjectTypeCharacter:
			return "male"
		case ObjectTypeMob:
			return "thing"
		}
	case AttributeMoney:
		return "0"
	case AttributeRarity:
		return "common"
	case AttributeType:
		return "generic"
	case AttributeTitle:
		switch ot {
		case ObjectTypeRoom:
			return "Empty Room"
		}
	case AttributeDescription:
		switch ot {
		case ObjectTypeRoom:
			return "You are in a newly created empty room. Make it a good one!"
		}
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
