package armeria

import (
	"armeria/internal/pkg/misc"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/google/uuid"

	"go.uber.org/zap"

	lua "github.com/yuin/gopher-lua"
)

// Room is a physical room that exists within an Area.
type Room struct {
	sync.RWMutex
	UUID             string            `json:"uuid"`
	UnsafeAttributes map[string]string `json:"attributes"`
	UnsafeHere       *ObjectContainer  `json:"here"`
	Coords           *Coords           `json:"coords"`
	ParentArea       *Area             `json:"-"`
}

// AdjacentRooms holds all of the Room objects that are adjacent to the current room.
type AdjacentRooms struct {
	North *Room
	South *Room
	East  *Room
	West  *Room
	Up    *Room
	Down  *Room
}

// Id returns the uuid of the room.
func (r *Room) Id() string {
	r.RLock()
	defer r.RUnlock()
	return r.UUID
}

// Init is called when the Room is created or loaded from disk.
func (r *Room) Init(a *Area) {
	// initialize uuid
	if r.UUID == "" {
		r.UUID = uuid.New().String()
	}
	// initialize UnsafeHere on rooms that don't have it defined
	if r.UnsafeHere == nil {
		r.UnsafeHere = NewObjectContainer(0)
	}
	// attach area
	r.ParentArea = a
	// attach self as container's parent
	r.UnsafeHere.AttachParent(r, ContainerParentTypeRoom)
	// sync container
	r.UnsafeHere.Sync()
	// register room with registry
	Armeria.registry.Register(r, r.UUID, RegistryTypeRoom)
}

// Deinit is called when the Room is deleted.
func (r *Room) Deinit() {
	Armeria.registry.Unregister(r.Id())
}

// SetAttribute sets a persistent attribute for the Room.
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

// Attribute retrieves a persistent attribute from the Room.
func (r *Room) Attribute(name string) string {
	r.RLock()
	defer r.RUnlock()

	if len(r.UnsafeAttributes[name]) == 0 {
		return RoomAttributeDefault(name)
	}

	return r.UnsafeAttributes[name]
}

// Here returns all the objects in the room via the ObjectContainer.
func (r *Room) Here() *ObjectContainer {
	r.RLock()
	defer r.RUnlock()

	return r.UnsafeHere
}

// RoomTargetData returns the JSON used for rendering the room objects on the client.
func (r *Room) RoomTargetData() string {
	r.RLock()
	defer r.RUnlock()

	var roomObjects []map[string]interface{}

	for _, obj := range r.Here().All() {
		o := obj.(ContainerObject)
		roomObjects = append(roomObjects, map[string]interface{}{
			"uuid":    o.Id(),
			"name":    o.Name(),
			"type":    o.Type(),
			"sort":    ObjectSortOrder(o.Type()),
			"picture": o.Attribute(AttributePicture),
			"rarity":  o.Attribute(AttributeRarity),
			"title":   o.Attribute(AttributeTitle),
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

	tc := fmt.Sprintf("%d,%d,%d", r.Coords.UnsafeX, r.Coords.UnsafeY, r.Coords.UnsafeZ)

	return &ObjectEditorData{
		Name:       r.Attribute(AttributeTitle),
		ObjectType: "room",
		Properties: props,
		TextCoords: tc,
	}
}

// CharacterEntered is called when the character is moved to the room (or logged in).
func (r *Room) CharacterEntered(c *Character, causedByLogin bool) {
	ca := c.Player().client
	ca.SyncMapLocation()
	ca.SyncRoomTitle()

	for _, char := range r.Here().Characters(true, nil) {
		char.Player().client.SyncRoomObjects()
	}

	for _, mi := range r.Here().Mobs() {
		go CallMobFunc(
			c,
			mi,
			"character_entered",
			lua.LString(c.Name()),
		)
	}
}

// CharacterLeft is called when the character left the room (or logged out).
func (r *Room) CharacterLeft(c *Character, causedByLogout bool) {
	for _, char := range r.Here().Characters(true, nil) {
		char.Player().client.SyncRoomObjects()
	}

	for _, mi := range r.Here().Mobs() {
		go CallMobFunc(
			c,
			mi,
			"character_left",
			lua.LString(c.Name()),
		)
	}
}

// AdjacentRooms returns the Room objects that are adjacent to the current room.
func (r *Room) AdjacentRooms() *AdjacentRooms {
	return &AdjacentRooms{
		North: Armeria.worldManager.RoomInDirection(r, NorthDirection),
		South: Armeria.worldManager.RoomInDirection(r, SouthDirection),
		East:  Armeria.worldManager.RoomInDirection(r, EastDirection),
		West:  Armeria.worldManager.RoomInDirection(r, WestDirection),
		Up:    Armeria.worldManager.RoomInDirection(r, UpDirection),
		Down:  Armeria.worldManager.RoomInDirection(r, DownDirection),
	}
}
