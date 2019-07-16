# Mob Scripting

## Index

* Function [c_attr](#c_attrcharacter_name-attribute-temp)
* Function: [say](#saytext)

* Event: [character_entered](#character_enteredcharacter_name)
* Event: [character_left](#character_leftcharacter_name)
* Event: [character_said](#character_saidcharacter_name-text
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

### c_set_attr(character_name, attribute, temp)

**Arguments**:
* `character_name` `(string)`: name of the character
* `attribute` `(string)`: attribute name to set
* `value` `(string)`: attribute value to set
* `temp` `(bool)`: whether to set a persistent or temporary attribute

**Returns**
* A `int` set to either `0` for success, `-1` for an invalid character, or `-2` for an invalid persistent attribute

Sets the value of a character's persistent or temporary attribute. A temporary attribute only exists for the
duration of the character's session.

### say(text)

**Arguments**:
* `text` `(string)`: the text to say

The mob will say `text` in the same room it's in. All other characters in the room will see this.

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