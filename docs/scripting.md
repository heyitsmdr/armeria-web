# Mob Scripting

## Index

### Functions

- [c_attr](#c_attrcharacter_name-attribute-temp)
- [c_set_attr](#c_set_attrcharacter_name-attribute-value-temp)
- [i_name](#i_nameuuid)
- [inv_give](#inv_givecharacter_uuid-item_uuid)
- [say](#saytext)
- [sleep](#sleepduration)
- [start_convo](#start_convo)
- [end_convo](#end_convo)
- [room_text](#room_texttext)

### Events

- [character_entered](#character_enteredcharacter_name)
- [character_left](#character_leftcharacter_name)
- [character_said](#character_saidcharacter_name-text)
- [received_item](#received_itemcharacter_name-uuid)
- [conversation_tick](#conversation_ticktick_count)

### Global Variables

- **invoker_uuid**: the uuid of the character who invoked the event
- **invoker_name**: the name of the character who invoked the event
- **mob_uuid**: the uuid of the current mob
- **mob_name**: the name of the current mob

## Functions

### c_attr(uuid, attribute, temp)

**Arguments**

- `uuid (string)`: uuid of the character to retrieve
- `attribute (string)`: attribute name to retrieve
- `temp (bool)`: whether to retrieve a persistent or temporary attribute

**Returns**

- A `string` containing the attribute value or an `int` of `-1` when the character was not found. An
  invalid attribute name will return an empty string.

Returns the value of a character's persistent or temporary attribute. A temporary attribute only
exists for the duration of the character's session.

### c_set_attr(uuid, attribute, value, temp)

**Arguments**

- `uuid (string)`: uuid of the character to alter
- `attribute (string)`: attribute name to alter
- `value (string)`: new attribute value
- `temp (bool)`: whether to set a persistent or temporary attribute

**Returns**

- An `int` set to either `0` for success, `-1` for an invalid character, or `-2` for an invalid
  persistent attribute.

Sets the value of a character's persistent or temporary attribute. A temporary attribute only exists
for the duration of the character's session.

### i_name(uuid)

**Arguments**

- `uuid (string)`: item uuid

**Returns**

- A `string` containing the formatted name of the item or an `int` of `0` if the uuid was either not
  found or not referencing an item.

Returns the formatted name of an item based on the item UUID.

### inv_give(uuid, item_uuid)

**Arguments**

- `uuid (string)`: character uuid to give the item to
- `item_uuid (string)`: item uuid of item in mob inventory

Gives an item from the mob's inventory to a character.

### say(text)

**Arguments**:

- `text (string)`: the text to say

The mob will say `text` in the same room it's in. All other characters in the room will see this.

### sleep(duration)

**Arguments**:

- `duration (string)`: duration to sleep for (ie: `30s`)

Delays the mob script for a particular duration. Standard Golang durations are valid.

### start_convo()

Starts a conversation with a character causing the `conversation_tick` event to fire every second.
Note that if a character moves or logs off, the conversation will be automatically ended.

### end_convo()

Ends a conversation with a character stopping the `conversation_tick` event from firing.

### room_text(text)

**Arguments**:

- `text (string)`: text to send to the room

Sends arbitrary text to the current room. Useful for conversations. Everyone in the room will see
this text.

## Events

### character_entered()

Triggered when a character enters the room.

### character_left()

Triggered when a character leaves the room.

### character_said(text)

**Parameters**

- `text (string)`: text said by the character

Triggered when a character says something in the room.

### received_item(item_uuid)

**Parameters**:

- `item_uuid (string)`: object uuid of received item

Triggered when a character gives an item to a mob. The item is automatically added to the mob's
inventory.

### conversation_tick(tick_count)

**Parameters**:

- `tick_count (int)`: current tick count after conversation started

Triggered every second after a conversation with a character is started. The `tick_count` will be
set to the number of ticks (seconds) that have passed since the start of the convo allowing you to
time out events that may occur during a conversation.
