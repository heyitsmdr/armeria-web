package armeria

import (
	"armeria/internal/pkg/misc"
	"encoding/json"
	"log"
	"sync"
)

type Room struct {
	Attributes map[string]string `json:"attributes"`
	Coords     *Coords           `json:"coords"`
	objects    []Object
	mux        sync.Mutex
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

func GetValidRoomAttributes() []string {
	return []string{
		"title",
		"description",
		"color",
	}
}

func GetRoomAttributeDefault(name string) string {
	switch name {
	case "title":
		return "Empty Room"
	case "description":
		return "You are in a newly created empty room. Make it a good one!"
	case "color":
		return "255,255,255"
	}

	return ""
}

func (r *Room) GetCoords() *Coords {
	r.mux.Lock()
	defer r.mux.Unlock()
	return r.Coords
}

func (r *Room) SetAttribute(name string, value string) {
	r.mux.Lock()
	defer r.mux.Unlock()

	if r.Attributes == nil {
		r.Attributes = make(map[string]string)
	}

	if !misc.Contains(GetValidRoomAttributes(), name) {
		log.Fatalf("[area] attempted set-attribute on a room using an invalid attribute: %s", name)
	}

	r.Attributes[name] = value
}

func (r *Room) GetAttribute(name string) string {
	r.mux.Lock()
	defer r.mux.Unlock()

	if len(r.Attributes[name]) == 0 {
		return GetRoomAttributeDefault(name)
	}

	return r.Attributes[name]
}

func (r *Room) GetObjects() []Object {
	r.mux.Lock()
	defer r.mux.Unlock()
	return r.objects
}

// GetCharacters returns online characters within the room
func (r *Room) GetCharacters(except *Character) []*Character {
	r.mux.Lock()
	defer r.mux.Unlock()

	var returnChars []*Character

	for _, o := range r.objects {
		if o.GetType() == ObjectTypeCharacter {
			if except == nil || o.GetName() != except.GetName() {
				char := o.(*Character)
				if char.GetPlayer() != nil {
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
		if o.GetType() == obj.GetType() && o.GetName() == obj.GetName() {
			r.objects[i] = r.objects[len(r.objects)-1]
			r.objects = r.objects[:len(r.objects)-1]
			return true
		}
	}

	return false
}

// GetObjectData returns the JSON used for rendering the room objects on the client
func (r *Room) GetObjectData() string {
	r.mux.Lock()
	defer r.mux.Unlock()

	var roomObjects []map[string]interface{}

	for _, o := range r.objects {
		roomObjects = append(roomObjects, map[string]interface{}{
			"name": o.GetName(),
			"type": o.GetType(),
		})
	}

	roomObjectJson, err := json.Marshal(roomObjects)
	if err != nil {
		log.Fatalf("[area] failed to marshal room object data: %s", err)
	}

	return string(roomObjectJson)
}

// GetEditorData returns the JSON used for the object editor
func (r *Room) GetEditorData() *ObjectEditorData {
	var props []*ObjectEditorDataProperty
	for _, attrName := range GetValidRoomAttributes() {
		props = append(props, &ObjectEditorDataProperty{
			PropType: "editable",
			Name:     attrName,
			Value:    r.GetAttribute(attrName),
		})
	}

	return &ObjectEditorData{
		Name:       r.GetAttribute("title"),
		ObjectType: "room",
		Properties: props,
	}
}

// OnCharacterEntered is called when the character is moved to the room (or logged in)
func (r *Room) OnCharacterEntered(c *Character, causedByLogin bool) {
	c.GetPlayer().clientActions.SyncMapLocation()

	for _, char := range r.GetCharacters(nil) {
		char.GetPlayer().clientActions.SyncRoomObjects()
	}
}

// OnCharacterLeft is called when the character left the room (or logged out)
func (r *Room) OnCharacterLeft(c *Character, causedByLogout bool) {
	for _, char := range r.GetCharacters(nil) {
		char.GetPlayer().clientActions.SyncRoomObjects()
	}
}
