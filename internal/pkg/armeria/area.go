package armeria

import (
	"armeria/internal/pkg/misc"
	"encoding/json"
	"log"
	"sync"

	"go.uber.org/zap"
)

type Area struct {
	sync.RWMutex
	UUID             string            `json:"uuid"`
	UnsafeName       string            `json:"name"`
	UnsafeRooms      []*Room           `json:"rooms"`
	UnsafeAttributes map[string]string `json:"attributes"`
}

const (
	NorthDirection = "north"
	SouthDirection = "south"
	EastDirection  = "east"
	WestDirection  = "west"
	UpDirection    = "up"
	DownDirection  = "down"
)

// Init is called when the Area is created or loaded from disk.
func (a *Area) Init() {
	Armeria.registry.Register(a, a.Id(), RegistryTypeArea)
}

// Deinit is called when the Area is deleted.
func (a *Area) Deinit() {
	Armeria.registry.Unregister(a.Id())
}

// ID returns the UUID of the Area.
func (a *Area) Id() string {
	return a.UUID
}

// UnsafeName returns the name of the area.
func (a *Area) Name() string {
	a.RLock()
	defer a.RUnlock()

	return a.UnsafeName
}

// RoomAt returns the Room at a particular Coords within the same area.
func (a *Area) RoomAt(c *Coords) *Room {
	a.RLock()
	defer a.RUnlock()

	for _, r := range a.UnsafeRooms {
		rc := r.Coords
		if rc.X() == c.X() && rc.Y() == c.Y() && rc.Z() == c.Z() && rc.I() == c.I() {
			return r
		}
	}

	return nil
}

// MinimapJSON returns the JSON used for minimap rendering on the client.
func (a *Area) MinimapJSON() string {
	a.RLock()
	defer a.RUnlock()

	var rooms []map[string]interface{}
	for _, r := range a.UnsafeRooms {
		rooms = append(rooms, map[string]interface{}{
			"title": r.Attribute("title"),
			"color": r.Attribute("color"),
			"type":  r.Attribute("type"),
			"x":     r.Coords.X(),
			"y":     r.Coords.Y(),
			"z":     r.Coords.Z(),
			"north": r.Attribute(NorthDirection),
			"east":  r.Attribute(EastDirection),
			"south": r.Attribute(SouthDirection),
			"west":  r.Attribute(WestDirection),
			"up":    r.Attribute(UpDirection),
			"down":  r.Attribute(DownDirection),
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
		UUID:       a.Id(),
		Name:       a.Name(),
		ObjectType: "area",
		Properties: props,
	}
}

// SetAttribute sets a permanent attribute and only valid attributes can be set.
func (a *Area) SetAttribute(name string, value string) {
	a.Lock()
	defer a.Unlock()

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
	a.RLock()
	defer a.RUnlock()

	if len(a.UnsafeAttributes[name]) == 0 {
		return AreaAttributeDefault(name)
	}

	return a.UnsafeAttributes[name]
}

// CharacterEntered is called when the character is moved into the area (or logged in).
func (a *Area) CharacterEntered(c *Character, causedByLogin bool) {
	c.Player().client.SyncMap()
}

// CharacterLeft is called when the character left the area (or logged out).
func (a *Area) CharacterLeft(c *Character, causedByLogout bool) {

}

// AddRoom adds a Room to the area.
func (a *Area) AddRoom(r *Room) {
	a.Lock()
	defer a.Unlock()

	r.Init(a)

	a.UnsafeRooms = append(a.UnsafeRooms, r)
}

func (a *Area) RemoveRoom(r *Room) {
	a.Lock()
	defer a.Unlock()

	r.Deinit()

	for i, rm := range a.UnsafeRooms {
		if rm.ID() == r.ID() {
			a.UnsafeRooms[i] = a.UnsafeRooms[len(a.UnsafeRooms)-1]
			a.UnsafeRooms = a.UnsafeRooms[:len(a.UnsafeRooms)-1]
			break
		}
	}
}

// Characters returns online characters within the area, with an optional Character exception.
func (a *Area) Characters(except *Character) []*Character {
	a.RLock()
	defer a.RUnlock()

	var c []*Character
	for _, r := range a.UnsafeRooms {
		c = append(c, r.Here().Characters(true, except)...)
	}

	return c
}
