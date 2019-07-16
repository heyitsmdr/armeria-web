# Mob Scripting

## Index

* Function: [mob_say](#character_enteredcharacter_name)
* Event: [character_entered](#character_enteredcharacter_name)
* Event: [character_left](#character_leftcharacter_name)

## Functions

### say(text)

Have the mob say `text` in the same room. All other characters in the room will see this.

## Events

### character_entered(character_name)

Triggered when a Character enters the room. The `character_name` parameter will be the name of the Character entering the room.

### character_left(character_name)

Triggered when a Character leaves the room. The `character_name` parameter will be the name of the Character leaving the room.