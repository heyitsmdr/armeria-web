package armeria

import (
	"armeria/internal/pkg/misc"
	"encoding/json"
	"log"
	"sync"

	"go.uber.org/zap"
)

// Area is a container for rooms.
type Area struct {
	sync.RWMutex
	UUID             string            `json:"uuid"`
	UnsafeName       string            `json:"name"`
	UnsafeRooms      []*Room           `json:"rooms"`
	UnsafeAttributes map[string]string `json:"attributes"`
}

// Direction strings.
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
	Armeria.registry.Register(a, a.ID(), RegistryTypeArea)
}

// Deinit is called when the Area is deleted.
func (a *Area) Deinit() {
	Armeria.registry.Unregister(a.ID())
}

// ID returns the UUID of the Area.
func (a *Area) ID() string {
	return a.UUID
}

// Name returns the name of the area.
func (a *Area) Name() string {
	a.RLock()
	defer a.RUnlock()

	return a.UnsafeName
}

// Rooms returns all the rooms within the area.
func (a *Area) Rooms() []*Room {
	a.RLock()
	defer a.RUnlock()

	return a.UnsafeRooms
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
		var north, south, east, west, up, down string
		if cr := r.ConnectedRoom(NorthDirection); cr != nil {
			north = cr.LocationString()
		}
		if cr := r.ConnectedRoom(SouthDirection); cr != nil {
			south = cr.LocationString()
		}
		if cr := r.ConnectedRoom(EastDirection); cr != nil {
			east = cr.LocationString()
		}
		if cr := r.ConnectedRoom(WestDirection); cr != nil {
			west = cr.LocationString()
		}
		if cr := r.ConnectedRoom(UpDirection); cr != nil {
			up = cr.LocationString()
		}
		if cr := r.ConnectedRoom(DownDirection); cr != nil {
			down = cr.LocationString()
		}
		rooms = append(rooms, map[string]interface{}{
			"title": r.Attribute("title"),
			"color": r.Attribute("color"),
			"type":  r.Attribute("type"),
			"x":     r.Coords.X(),
			"y":     r.Coords.Y(),
			"z":     r.Coords.Z(),
			"north": north,
			"south": south,
			"east":  east,
			"west":  west,
			"up":    up,
			"down":  down,
		})
	}

	minimap := map[string]interface{}{
		"name":  a.UnsafeName,
		"rooms": rooms,
	}

	mapJSON, err := json.Marshal(minimap)
	if err != nil {
		log.Fatalf("[area] failed to marshal minimap data: %s", err)
	}

	return string(mapJSON)

}

// EditorData returns the JSON used for the object editor.
func (a *Area) EditorData() *ObjectEditorData {
	var props []*ObjectEditorDataProperty
	for _, attrName := range AttributeList(ObjectTypeArea) {
		props = append(props, &ObjectEditorDataProperty{
			PropType: AttributeEditorType(ObjectTypeArea, attrName),
			Name:     attrName,
			Value:    a.Attribute(attrName),
		})
	}

	return &ObjectEditorData{
		UUID:       a.ID(),
		Name:       a.Name(),
		ObjectType: "area",
		Properties: props,
	}
}

// SetAttribute sets a permanent attribute and only valid attributes can be set.
func (a *Area) SetAttribute(name string, value string) {
	a.Lock()
	defer a.Unlock()

	if !misc.Contains(AttributeList(ObjectTypeArea), name) {
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
		return AttributeDefault(ObjectTypeArea, name)
	}

	return a.UnsafeAttributes[name]
}

// CharacterEntered is called when the unsafeCharacter is moved into the area (or logged in).
func (a *Area) CharacterEntered(c *Character, causedByLogin bool) {
	c.Player().client.SyncMap()
}

// CharacterLeft is called when the unsafeCharacter left the area (or logged out).
func (a *Area) CharacterLeft(c *Character, causedByLogout bool) {

}

// AddRoom adds a Room to the area.
func (a *Area) AddRoom(r *Room) {
	a.Lock()
	defer a.Unlock()

	r.Init(a)

	a.UnsafeRooms = append(a.UnsafeRooms, r)
}

// RemoveRoom removes a room.
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
func (a *Area) Characters(exceptions ...*Character) []*Character {
	a.RLock()
	defer a.RUnlock()

	var c []*Character
	for _, r := range a.UnsafeRooms {
		c = append(c, r.Here().Characters(true, exceptions...)...)
	}

	return c
}
