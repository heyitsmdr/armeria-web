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
			"title": r.GetAttribute("title"),
			"color": r.GetAttribute("color"),
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

// OnCharacterEntered is called when the character is moved into the area (or logged in).
func (a *Area) OnCharacterEntered(c *Character, causedByLogin bool) {
	c.GetPlayer().clientActions.RenderMap()
}

// OnCharacterLeft is called when the character left the area (or logged out).
func (a *Area) OnCharacterLeft(c *Character, causedByLogout bool) {

}

func (a *Area) AddRoom(r *Room) {
	a.mux.Lock()
	defer a.mux.Unlock()
	a.Rooms = append(a.Rooms, r)
}

func (a *Area) RemoveRoom(r *Room) {

}

// GetCharacters returns online characters within the area.
func (a *Area) GetCharacters(except *Character) []*Character {
	a.mux.Lock()
	defer a.mux.Unlock()

	var returnChars []*Character

	for _, r := range a.Rooms {
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
	}

	return returnChars
}

// GetAdjacentRooms returns the Room objects that are adjacent to the current room.
func (a *Area) GetAdjacentRooms(r *Room) *AdjacentRooms {
	return &AdjacentRooms{
		North: Armeria.worldManager.GetRoomInDirection(a, r, NorthDirection),
		South: Armeria.worldManager.GetRoomInDirection(a, r, SouthDirection),
		East:  Armeria.worldManager.GetRoomInDirection(a, r, EastDirection),
		West:  Armeria.worldManager.GetRoomInDirection(a, r, WestDirection),
		Up:    Armeria.worldManager.GetRoomInDirection(a, r, UpDirection),
		Down:  Armeria.worldManager.GetRoomInDirection(a, r, DownDirection),
	}
}
