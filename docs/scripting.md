# Mob Scripting

## Index

* Function: [say](#saytext)
* Event: [character_entered](#character_enteredcharacter_name)
* Event: [character_left](#character_leftcharacter_name)

## Functions

### say(text)

The mob will say `text` in the same room it's in. All other characters in the room will see this.

## Events

### character_entered(character_name)

Triggered when a character enters the room. The `character_name` parameter will be the name of the Character entering the room.

### character_left(character_name)

Triggered when a character leaves the room. The `character_name` parameter will be the name of the Character leaving the room.