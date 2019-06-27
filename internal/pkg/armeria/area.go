package armeria

import (
	"armeria/internal/pkg/misc"
	"encoding/json"
	"log"
	"sync"

	"go.uber.org/zap"
)

type Area struct {
	UUID             string            `json:"uuid"`
	UnsafeName       string            `json:"name"`
	UnsafeRooms      []*Room           `json:"rooms"`
	UnsafeAttributes map[string]string `json:"attributes"`
	mux              sync.Mutex
}

type AdjacentRooms struct {
	North *Room
	South *Room
	East  *Room
	West  *Room
	Up    *Room
	Down  *Room
}

const (
	NorthDirection = "north"
	SouthDirection = "south"
	EastDirection  = "east"
	WestDirection  = "west"
	UpDirection    = "up"
	DownDirection  = "down"
)

// ValidAreaAttributes returns an array of valid attributes that can be permanently set.
func ValidAreaAttributes() []string {
	return []string{}
}

// AreaAttributeDefault returns the default value for a particular attribute.
func AreaAttributeDefault(name string) string {
	switch name {

	}

	return ""
}

// Id returns the UUID of the Area..
func (a *Area) Id() string {
	return a.UUID
}

// UnsafeName returns the name of the area.
func (a *Area) Name() string {
	a.mux.Lock()
	defer a.mux.Unlock()
	return a.UnsafeName
}

// RoomAt returns the Room at a particular UnsafeCoords within the same area.
func (a *Area) RoomAt(c *Coords) *Room {
	a.mux.Lock()
	defer a.mux.Unlock()

	for _, r := range a.UnsafeRooms {
		rc := r.Coords()
		if rc.X == c.X && rc.Y == c.Y && rc.Z == c.Z && rc.I == c.I {
			return r
		}
	}

	return nil
}

// MinimapData returns the JSON used for minimap rendering on the client.
func (a *Area) MinimapData() string {
	a.mux.Lock()
	defer a.mux.Unlock()

	var rooms []map[string]interface{}
	for _, r := range a.UnsafeRooms {
		rooms = append(rooms, map[string]interface{}{
			"title": r.Attribute("title"),
			"color": r.Attribute("color"),
			"type":  r.Attribute("type"),
			"x":     r.Coords().X,
			"y":     r.Coords().Y,
			"z":     r.Coords().Z,
		})
	}

	minimap := map[string]interface{}{
		"name":  a.UnsafeName,
		"rooms": rooms,
	}

	mapJson, err := json.Marshal(minimap)
	if err != nil {
		log.Fatalf("[area] failed to marshal minimap data: %s", err)
	}

	return string(mapJson)

}

// EditorData returns the JSON used for the object editor.
func (a *Area) EditorData() *ObjectEditorData {
	var props []*ObjectEditorDataProperty
	for _, attrName := range ValidAreaAttributes() {
		props = append(props, &ObjectEditorDataProperty{
			PropType: "editable",
			Name:     attrName,
			Value:    a.Attribute(attrName),
		})
	}

	return &ObjectEditorData{
		Name:       a.Name(),
		ObjectType: "area",
		Properties: props,
	}
}

// SetAttribute sets a permanent attribute and only valid attributes can be set.
func (a *Area) SetAttribute(name string, value string) {
	a.mux.Lock()
	defer a.mux.Unlock()

	if !misc.Contains(ValidAreaAttributes(), name) {
		Armeria.log.Fatal("attempted to set invalid attribute",
			zap.String("attribute", name),
			zap.String("value", value),
		)
	}

	a.UnsafeAttributes[name] = value
}

// Attribute returns a permanent attribute.
func (a *Area) Attribute(name string) string {
	a.mux.Lock()
	defer a.mux.Unlock()

	if len(a.UnsafeAttributes[name]) == 0 {
		return AreaAttributeDefault(name)
	}

	return a.UnsafeAttributes[name]
}

// CharacterEntered is called when the character is moved into the area (or logged in).
func (a *Area) CharacterEntered(c *Character, causedByLogin bool) {
	c.Player().clientActions.RenderMap()
}

// CharacterLeft is called when the character left the area (or logged out).
func (a *Area) CharacterLeft(c *Character, causedByLogout bool) {

}

// AddRoom adds a Room to the area.
func (a *Area) AddRoom(r *Room) {
	a.mux.Lock()
	defer a.mux.Unlock()
	a.UnsafeRooms = append(a.UnsafeRooms, r)
}

func (a *Area) RemoveRoom(r *Room) {

}

// UnsafeCharacters returns online characters within the area, with an optional Character exception.
func (a *Area) Characters(except *Character) []*Character {
	a.mux.Lock()
	defer a.mux.Unlock()

	var returnChars []*Character

	for _, r := range a.UnsafeRooms {
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
	}

	return returnChars
}

// AdjacentRooms returns the Room objects that are adjacent to the current room.
func (a *Area) AdjacentRooms(r *Room) *AdjacentRooms {
	return &AdjacentRooms{
		North: Armeria.worldManager.RoomInDirection(a, r, NorthDirection),
		South: Armeria.worldManager.RoomInDirection(a, r, SouthDirection),
		East:  Armeria.worldManager.RoomInDirection(a, r, EastDirection),
		West:  Armeria.worldManager.RoomInDirection(a, r, WestDirection),
		Up:    Armeria.worldManager.RoomInDirection(a, r, UpDirection),
		Down:  Armeria.worldManager.RoomInDirection(a, r, DownDirection),
	}
}
