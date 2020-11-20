package armeria

import (
	"armeria/internal/pkg/misc"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/google/uuid"

	"go.uber.org/zap"
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

// ID returns the uuid of the room.
func (r *Room) ID() string {
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
	Armeria.registry.Unregister(r.ID())
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

// RoomTargetJSON returns the JSON used for rendering the room objects on the client.
func (r *Room) RoomTargetJSON() string {
	r.RLock()
	defer r.RUnlock()

	var roomObjects []map[string]interface{}

	for _, obj := range r.Here().All() {
		o := obj.(ContainerObject)

		if o.Type() == ContainerObjectTypeCharacter && o.(*Character).Player() == nil {
			continue
		}

		rarityColor := ""
		if o.Type() == ContainerObjectTypeItem {
			rarityColor = o.(*ItemInstance).RarityColor()
		}

		roomObjects = append(roomObjects, map[string]interface{}{
			"uuid":    o.ID(),
			"name":    o.Name(),
			"type":    o.Type(),
			"sort":    ObjectSortOrder(o.Type()),
			"picture": o.Attribute(AttributePicture),
			"color":   rarityColor,
			"title":   o.Attribute(AttributeTitle),
		})
	}

	roomObjectJSON, err := json.Marshal(roomObjects)
	if err != nil {
		Armeria.log.Fatal("failed to marshal room object data",
			zap.String("room", r.UUID),
			zap.Error(err),
		)
	}

	return string(roomObjectJSON)
}

// EditorData returns the JSON used for the object editor.
func (r *Room) EditorData() *ObjectEditorData {
	var props []*ObjectEditorDataProperty
	for _, attrName := range ValidRoomAttributes() {
		props = append(props, &ObjectEditorDataProperty{
			PropType: AttributeEditorType(ObjectTypeMob, attrName),
			Name:     attrName,
			Value:    r.Attribute(attrName),
		})
	}

	tc := fmt.Sprintf("%d,%d,%d", r.Coords.UnsafeX, r.Coords.UnsafeY, r.Coords.UnsafeZ)

	return &ObjectEditorData{
		UUID:       r.ID(),
		Name:       r.Attribute(AttributeTitle),
		ObjectType: "room",
		Properties: props,
		TextCoords: tc,
	}
}

// CharacterEntered is called when the unsafeCharacter is moved to the room (or logged in).
func (r *Room) CharacterEntered(c *Character, causedByLogin bool) {
	ca := c.Player().client
	ca.SyncMapLocation()
	ca.SyncRoomTitle()

	for _, char := range r.Here().Characters(true) {
		char.Player().client.SyncRoomObjects()
	}

	for _, mi := range r.Here().Mobs() {
		go CallMobFunc(
			c,
			mi,
			"character_entered",
		)
	}
}

// CharacterLeft is called when the unsafeCharacter left the room (or logged out).
func (r *Room) CharacterLeft(c *Character, causedByLogout bool) {
	for _, char := range r.Here().Characters(true, c) {
		char.Player().client.SyncRoomObjects()
	}

	for _, mi := range r.Here().Mobs() {
		go CallMobFunc(
			c,
			mi,
			"character_left",
		)
	}
}

// AdjacentRooms returns the Room objects that are adjacent to the current room.
func (r *Room) AdjacentRooms() *AdjacentRooms {
	return &AdjacentRooms{
		North: r.ConnectedRoom(NorthDirection),
		South: r.ConnectedRoom(SouthDirection),
		East:  r.ConnectedRoom(EastDirection),
		West:  r.ConnectedRoom(WestDirection),
		Up:    r.ConnectedRoom(UpDirection),
		Down:  r.ConnectedRoom(DownDirection),
	}
}

// ConnectedRoom returns the adjacent explicit or implicit Room object.
func (r *Room) ConnectedRoom(direction string) *Room {
	// check for an explicit exit first
	explicitDirection := r.Attribute(direction)
	edSections := strings.Split(explicitDirection, ",")
	if len(explicitDirection) > 0 {
		if len(edSections) == 3 {
			return r.ParentArea.RoomAt(NewCoordsFromString(explicitDirection))
		} else if len(edSections) == 4 {
			a := Armeria.worldManager.AreaByName(edSections[0])
			if a != nil {
				return a.RoomAt(NewCoordsFromString(strings.Join(edSections[1:], ",")))
			}
		} else if explicitDirection[0:1] == "!" {
			return nil
		}
	}

	// check for an implicit exit
	offsets := misc.DirectionOffsets(direction)
	if offsets == nil {
		return nil
	}

	x := r.Coords.X() + offsets["x"]
	y := r.Coords.Y() + offsets["y"]
	z := r.Coords.Z() + offsets["z"]

	loc := NewCoords(x, y, z, 0)

	return r.ParentArea.RoomAt(loc)
}

// LocationString returns the location of the room within the game world as a string.
func (r *Room) LocationString() string {
	a := r.ParentArea
	return fmt.Sprintf("%s,%d,%d,%d", a.Name(), r.Coords.X(), r.Coords.Y(), r.Coords.Z())
}

// DistanceBetween returns the Coords-based distance between two rooms.
func (r *Room) DistanceBetween(rm *Room) *Coords {
	co := &Coords{
		UnsafeX: rm.Coords.X() - r.Coords.X(),
		UnsafeY: rm.Coords.Y() - r.Coords.Y(),
		UnsafeZ: rm.Coords.Z() - r.Coords.Z(),
	}
	return co
}
