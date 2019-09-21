# Mob Scripting

## Index

* Function: [c_attr](#c_attrcharacter_name-attribute-temp)
* Function: [c_set_attr](#c_set_attrcharacter_name-attribute-value-temp)
* Function: [i_name](#i_nameuuid)
* Function: [inv_give](#inv_givecharacter_uuid-item_uuid)
* Function: [say](#saytext)
* Function: [sleep](#sleepduration)

* Event: [character_entered](#character_enteredcharacter_name)
* Event: [character_left](#character_leftcharacter_name)
* Event: [character_said](#character_saidcharacter_name-text)
* Event: [received_item](#received_itemcharacter_name-uuid)

)
## Functions

### c_attr(character_name, attribute, temp)

**Arguments**:
* `character_name` `(string)`: name of the character
* `attribute` `(string)`: attribute name to retrieve
* `temp` `(bool)`: whether to retrieve a persistent or temporary attribute

**Returns**
* A `string` containing the attribute value or an `int` of `-1` when the character was not found. An invalid attribute
will return an empty string

Returns the value of a character's persistent or temporary attribute. A temporary attribute only exists for the
duration of the character's session.

### c_set_attr(character_name, attribute, value, temp)

**Arguments**:
* `character_name` `(string)`: name of the character
* `attribute` `(string)`: attribute name to set
* `value` `(string)`: attribute value to set
* `temp` `(bool)`: whether to set a persistent or temporary attribute

**Returns**
* An `int` set to either `0` for success, `-1` for an invalid character, or `-2` for an invalid persistent attribute

Sets the value of a character's persistent or temporary attribute. A temporary attribute only exists for the
duration of the character's session.

### i_name(uuid)

**Arguments**:
* `uuid` `(string)`: item uuid

**Returns**
* A `string` containing the formatted name of the item or an `int` of `0` if the uuid was either not found or not referencing an item.

Returns the formatted name of an item based on the item UUID.

### inv_give(character_uuid, item_uuid)

**Arguments**:
* `character_uuid` `(string)`: character uuid to receive item
* `item_uuid` `(string)`: item uuid of item in mob inventory

Gives an item from the mob's inventory to a character.

### say(text)

**Arguments**:
* `text` `(string)`: the text to say

The mob will say `text` in the same room it's in. All other characters in the room will see this.

### sleep(duration)

**Arguments**:
* `duration` `(string)`: duration to sleep for (ie: `30s`)

Delays the mob script for a particular duration.

## Events

### character_entered(character_name)

**Parameters**:
* `character_name` `(string)`: name of the character

Triggered when a character enters the room. The `character_name` parameter will be the name of the Character
entering the room.

### character_left(character_name)

**Parameters**:
* `character_name` `(string)`: name of the character
    
Triggered when a character leaves the room. The `character_name` parameter will be the name of the Character
leaving the room.

### character_said(character_name, text)

**Parameters**:
* `character_name` `(string)`: name of the character
* `text` `(string)`: text said by the character
    
Triggered when a character says something in the room.

### received_item(character_uuid, item_uuid)

**Parameters**:
* `character_uuid` `(string)`: character uuid
* `item_uuid` `(string)`: object uuid of received item
    
Triggered when a character gives an item to a mob.