package armeria

import (
	"encoding/json"
	"log"
	"sync"
)

type Area struct {
	Name  string  `json:"name"`
	Rooms []*Room `json:"rooms"`
	mux   sync.Mutex
}

type Room struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Coords      *Coords `json:"coords"`
	objects     []Object
	mux         sync.Mutex
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

func (a *Area) GetName() string {
	a.mux.Lock()
	defer a.mux.Unlock()
	return a.Name
}

func (a *Area) GetRoom(c *Coords) *Room {
	a.mux.Lock()
	defer a.mux.Unlock()

	for _, r := range a.Rooms {
		rc := r.GetCoords()
		if rc.X == c.X && rc.Y == c.Y && rc.Z == c.Z && rc.I == c.I {
			return r
		}
	}

	return nil
}

// GetMinimapData returns the JSON used for minimap rendering on the client
func (a *Area) GetMinimapData() string {
	a.mux.Lock()
	defer a.mux.Unlock()

	var rooms []map[string]interface{}
	for _, r := range a.Rooms {
		rooms = append(rooms, map[string]interface{}{
			"title": r.GetTitle(),
			"x":     r.GetCoords().X,
			"y":     r.GetCoords().Y,
			"z":     r.GetCoords().Z,
		})
	}

	minimap := map[string]interface{}{
		"name":  a.Name,
		"rooms": rooms,
	}

	mapJson, err := json.Marshal(minimap)
	if err != nil {
		log.Fatalf("[area] failed to marshal minimap data: %s", err)
	}

	return string(mapJson)

}

// OnCharacterEntered is called when the character is moved into the area (or logged in)
func (a *Area) OnCharacterEntered(c *Character, causedByLogin bool) {
	c.GetPlayer().clientActions.RenderMap()
}

// OnCharacterLeft is called when the character left the area (or logged out)
func (a *Area) OnCharacterLeft(c *Character, causedByLogout bool) {

}

func (r *Room) GetCoords() *Coords {
	r.mux.Lock()
	defer r.mux.Unlock()
	return r.Coords
}

func (r *Room) SetTitle(title string) {
	r.mux.Lock()
	defer r.mux.Unlock()
	r.Title = title
}

func (r *Room) GetTitle() string {
	r.mux.Lock()
	defer r.mux.Unlock()
	return r.Title
}

func (r *Room) GetDescription() string {
	r.mux.Lock()
	defer r.mux.Unlock()
	return r.Description
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
	r.mux.Lock()
	defer r.mux.Unlock()

	editorData := &ObjectEditorData{
		Name:       r.Title,
		ObjectType: "room",
		Properties: []*ObjectEditorDataProperty{
			{
				PropType: "editable",
				Name:     "title",
				Value:    r.Title,
			},
			{
				PropType: "editable",
				Name:     "description",
				Value:    r.Description,
			},
		},
	}

	return editorData
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
