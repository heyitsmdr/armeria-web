package armeria

import (
	"armeria/internal/pkg/sfx"
	"armeria/internal/pkg/validate"
	"strings"
)

const (
	AttributeChannels    string = "channels"
	AttributeColor       string = "color"
	AttributeDescription string = "description"
	AttributeDown        string = "down"
	AttributeEast        string = "east"
	AttributeFollowCrumb string = "followCrumb"
	AttributeFollowSpeed string = "followSpeed"
	AttributeGender      string = "gender"
	AttributeHoldable    string = "holdable"
	AttributeMoney       string = "money"
	AttributeMusic       string = "music"
	AttributeNorth       string = "north"
	AttributePermissions string = "permissions"
	AttributePicture     string = "picture"
	AttributeRarity      string = "rarity"
	AttributeScript      string = "script"
	AttributeSpawnLimit  string = "spawnLimit"
	AttributeSpawnMob    string = "spawnMob"
	AttributeSpawnSFX    string = "spawnSFX"
	AttributeSouth       string = "south"
	AttributeTitle       string = "title"
	AttributeType        string = "type"
	AttributeUp          string = "up"
	AttributeVisible     string = "visible"
	AttributeWest        string = "west"

	TempAttributeEditorOpen string = "editorOpen"
	TempAttributeGhost      string = "ghost"
	TempAttributeReplyTo    string = "replyTo"
)

// AttributeCasing returns the correct casing for a given object type and attribute.
func AttributeCasing(attr string) string {
	allAttrs := make([]string, 0)
	allAttrs = append(allAttrs, AttributeList(ObjectTypeCharacter)...)
	allAttrs = append(allAttrs, AttributeList(ObjectTypeArea)...)
	allAttrs = append(allAttrs, AttributeList(ObjectTypeRoom)...)
	allAttrs = append(allAttrs, AttributeList(ObjectTypeItem)...)
	allAttrs = append(allAttrs, AttributeList(ObjectTypeItemInstance)...)
	allAttrs = append(allAttrs, AttributeList(ObjectTypeMob)...)
	allAttrs = append(allAttrs, AttributeList(ObjectTypeMobInstance)...)

	for _, a := range allAttrs {
		if strings.ToLower(a) == strings.ToLower(attr) {
			return a
		}
	}

	return ""
}

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
			AttributeHoldable,
			AttributeVisible,
			AttributeSpawnMob,
			AttributeSpawnLimit,
		}
	case ObjectTypeItemInstance:
		return []string{
			AttributeRarity,
			AttributeDescription,
			AttributeHoldable,
			AttributeVisible,
			AttributeSpawnLimit,
		}
	case ObjectTypeMob:
		return []string{
			AttributePicture,
			AttributeScript,
			AttributeGender,
			AttributeTitle,
			AttributeSpawnSFX,
			AttributeFollowCrumb,
			AttributeFollowSpeed,
		}
	case ObjectTypeMobInstance:
		return []string{
			AttributeTitle,
		}
	}

	return []string{}
}

// AttributeEditorType returns the object editor "type" string of an attribute for a given ObjectType. Case sensitive.
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
		switch ot {
		case ObjectTypeCharacter:
			return "enum:male|female"
		case ObjectTypeMob:
			return "enum:male|female|thing"
		}
	case AttributeColor:
		return "color"
	case AttributeType:
		switch ot {
		case ObjectTypeItem:
			return "enum:" + strings.Join(ItemTypes(), "|")
		case ObjectTypeRoom:
			return "enum:generic|track|bank|armor|sword|home|wand"
		default:
			return "editable"
		}
	case AttributeHoldable:
		return "enum:true|false"
	case AttributeVisible:
		return "enum:true|false"
	case AttributeSpawnSFX:
		return "enum:" + strings.Join(sfx.List(), "|")
	}

	return "editable"
}

// AttributeDefault returns the default value of an attribute for a given ObjectType. Case sensitive.
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
	case AttributeHoldable:
		return "true"
	case AttributeVisible:
		return "true"
	case AttributeSpawnLimit:
		return "0"
	case AttributeFollowSpeed:
		return "12"
	}

	return ""
}

// AttributeValidate returns the validation result of an attribute value for a given ObjectType. Case sensitive.
func AttributeValidate(ot ObjectType, attr, val string) validate.ValidationResult {
	var validatorString string
	switch ot {
	case ObjectTypeMob:
		switch attr {
		case AttributeScript:
			validatorString = "empty"
			break
		case AttributeGender:
			validatorString = "in:thing,male,female"
			break
		case AttributeSpawnSFX:
			validatorString = "in:" + strings.Join(sfx.List(), ",")
			break
		case AttributeFollowSpeed:
			validatorString = "num|min:1|max:60"
			break
		}
	case ObjectTypeCharacter:
		switch attr {
		case AttributeGender:
			validatorString = "in:male,female"
			break
		case AttributeMoney:
			validatorString = "num|min:0"
			break
		}
	case ObjectTypeItem:
		switch attr {
		case AttributeType:
			validatorString = "in:" + strings.Join(ItemTypes(), ",")
			break
		case AttributeRarity:
			validatorString = "in:common,uncommon"
			break
		case AttributeHoldable:
			validatorString = "bool"
			break
		case AttributeVisible:
			validatorString = "bool"
			break
		case AttributeSpawnLimit:
			validatorString = "num|min:0|max:100"
			break
		}
	case ObjectTypeRoom:
		switch attr {
		case AttributeType:
			validatorString = "in:generic,track,bank,armor,sword,home,wand"
			break
		}
	}

	return validate.Check(val, validatorString)
}
