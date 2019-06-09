package armeria

import (
	"armeria/internal/pkg/misc"
	"encoding/json"
	"log"
	"sync"

	lua "github.com/yuin/gopher-lua"
)

type Room struct {
	UnsafeAttributes map[string]string `json:"attributes"`
	UnafeCoords      *Coords           `json:"coords"`
	objects          []Object
	mux              sync.Mutex
}

type Coords struct {
	X int `json:"x"`
	Y int `json:"y"`
	Z int `json:"z"`
	I int `json:"-"`
}

type Location struct {
	AreaName string  `json:"area"`
	Coords   *Coords `json:"coords"`
}

func ValidRoomAttributes() []string {
	return []string{
		"title",
		"description",
		"color",
	}
}

func RoomAttributeDefault(name string) string {
	switch name {
	case "title":
		return "Empty Room"
	case "description":
		return "You are in a newly created empty room. Make it a good one!"
	case "color":
		return "190,190,190"
	}

	return ""
}

func (r *Room) Coords() *Coords {
	r.mux.Lock()
	defer r.mux.Unlock()
	return r.UnafeCoords
}

func (r *Room) SetAttribute(name string, value string) {
	r.mux.Lock()
	defer r.mux.Unlock()

	if r.UnsafeAttributes == nil {
		r.UnsafeAttributes = make(map[string]string)
	}

	if !misc.Contains(ValidRoomAttributes(), name) {
		log.Fatalf("[area] attempted set-attribute on a room using an invalid attribute: %s", name)
	}

	r.UnsafeAttributes[name] = value
}

func (r *Room) Attribute(name string) string {
	r.mux.Lock()
	defer r.mux.Unlock()

	if len(r.UnsafeAttributes[name]) == 0 {
		return RoomAttributeDefault(name)
	}

	return r.UnsafeAttributes[name]
}

func (r *Room) Objects() []Object {
	r.mux.Lock()
	defer r.mux.Unlock()
	return r.objects
}

// Characters returns online characters within the room.
func (r *Room) Characters(except *Character) []*Character {
	r.mux.Lock()
	defer r.mux.Unlock()

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

func (r *Room) AddObjectToRoom(obj Object) {
	r.mux.Lock()
	defer r.mux.Unlock()
	r.objects = append(r.objects, obj)
}

func (r *Room) RemoveObjectFromRoom(obj Object) bool {
	r.mux.Lock()
	defer r.mux.Unlock()

	for i, o := range r.objects {
		if o.Type() == obj.Type() && o.Name() == obj.Name() {
			r.objects[i] = r.objects[len(r.objects)-1]
			r.objects = r.objects[:len(r.objects)-1]
			return true
		}
	}

	return false
}

// ObjectData returns the JSON used for rendering the room objects on the client.
func (r *Room) ObjectData() string {
	r.mux.Lock()
	defer r.mux.Unlock()

	var roomObjects []map[string]interface{}

	for _, o := range r.objects {
		roomObjects = append(roomObjects, map[string]interface{}{
			"name":    o.Name(),
			"type":    o.Type(),
			"picture": o.Attribute("picture"),
		})
	}

	roomObjectJson, err := json.Marshal(roomObjects)
	if err != nil {
		log.Fatalf("[area] failed to marshal room object data: %s", err)
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
		Name:       r.Attribute("title"),
		ObjectType: "room",
		Properties: props,
	}
}

// CharacterEntered is called when the character is moved to the room (or logged in).
func (r *Room) CharacterEntered(c *Character, causedByLogin bool) {
	ca := c.Player().clientActions
	ca.SyncMapLocation()
	ca.SyncRoomTitle()

	for _, char := range r.Characters(nil) {
		char.Player().clientActions.SyncRoomObjects()
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
		char.Player().clientActions.SyncRoomObjects()
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
