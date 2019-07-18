package armeria

import (
	"armeria/internal/pkg/misc"
	"encoding/json"
	"log"
	"sync"

	"go.uber.org/zap"

	lua "github.com/yuin/gopher-lua"
)

type Room struct {
	sync.RWMutex
	UnsafeAttributes map[string]string `json:"attributes"`
	Coords           *Coords           `json:"coords"`
	objects          []Object
}

func (r *Room) SetAttribute(name string, value string) {
	r.Lock()
	defer r.Unlock()

	if r.UnsafeAttributes == nil {
		r.UnsafeAttributes = make(map[string]string)
	}

	if !misc.Contains(ValidRoomAttributes(), name) {
		log.Fatalf("[area] attempted set-attribute on a room using an invalid attribute: %s", name)
	}

	r.UnsafeAttributes[name] = value
}

func (r *Room) Attribute(name string) string {
	r.RLock()
	defer r.RUnlock()

	if len(r.UnsafeAttributes[name]) == 0 {
		return RoomAttributeDefault(name)
	}

	return r.UnsafeAttributes[name]
}

func (r *Room) Objects() []Object {
	r.RLock()
	defer r.RUnlock()

	return r.objects
}

// OnlineCharacters returns online characters within the room.
func (r *Room) Characters(except *Character) []*Character {
	r.RLock()
	defer r.RUnlock()

	var returnChars []*Character

	for _, o := range r.objects {
		if o.Type() == ObjectTypeCharacter {
			if except == nil || o.Name() != except.Name() {
				char := o.(*Character)
				if char.Player() != nil {
					returnChars = append(returnChars, char)
				}
			}
		}
	}

	return returnChars
}

// AddObjectToRoom adds an Object to the Room.
func (r *Room) AddObjectToRoom(obj Object) {
	r.Lock()
	defer r.Unlock()

	r.objects = append(r.objects, obj)
}

// RemoveObjectFromRoom attempts to remove the Object from the Room, and returns
// a bool indicating whether it was successful or not.
func (r *Room) RemoveObjectFromRoom(obj Object) bool {
	r.Lock()
	defer r.Unlock()

	for i, o := range r.objects {
		if o.Id() == obj.Id() {
			r.objects[i] = r.objects[len(r.objects)-1]
			r.objects = r.objects[:len(r.objects)-1]
			return true
		}
	}

	return false
}

// ObjectData returns the JSON used for rendering the room objects on the client.
func (r *Room) ObjectData() string {
	r.RLock()
	defer r.RUnlock()

	var roomObjects []map[string]interface{}

	for _, o := range r.objects {
		roomObjects = append(roomObjects, map[string]interface{}{
			"uuid":    o.Id(),
			"name":    o.Name(),
			"type":    o.Type(),
			"sort":    ObjectSortOrder(o.Type()),
			"picture": o.Attribute(AttributePicture),
			"rarity":  o.Attribute(AttributeRarity),
		})
	}

	roomObjectJson, err := json.Marshal(roomObjects)
	if err != nil {
		Armeria.log.Fatal("failed to marshal room object data",
			zap.Error(err),
		)
	}

	return string(roomObjectJson)
}

// EditorData returns the JSON used for the object editor.
func (r *Room) EditorData() *ObjectEditorData {
	var props []*ObjectEditorDataProperty
	for _, attrName := range ValidRoomAttributes() {
		props = append(props, &ObjectEditorDataProperty{
			PropType: "editable",
			Name:     attrName,
			Value:    r.Attribute(attrName),
		})
	}

	return &ObjectEditorData{
		Name:       r.Attribute(AttributeTitle),
		ObjectType: "room",
		Properties: props,
	}
}

// CharacterEntered is called when the character is moved to the room (or logged in).
func (r *Room) CharacterEntered(c *Character, causedByLogin bool) {
	ca := c.Player().client
	ca.SyncMapLocation()
	ca.SyncRoomTitle()

	for _, char := range r.Characters(nil) {
		char.Player().client.SyncRoomObjects()
	}

	for _, o := range r.Objects() {
		if o.Type() == ObjectTypeMob {
			go CallMobFunc(
				c,
				o.(*MobInstance),
				"character_entered",
				lua.LString(c.Name()),
			)
		}
	}
}

// CharacterLeft is called when the character left the room (or logged out).
func (r *Room) CharacterLeft(c *Character, causedByLogout bool) {
	for _, char := range r.Characters(nil) {
		char.Player().client.SyncRoomObjects()
	}

	for _, o := range r.Objects() {
		if o.Type() == ObjectTypeMob {
			go CallMobFunc(
				c,
				o.(*MobInstance),
				"character_left",
				lua.LString(c.Name()),
			)
		}
	}
}
