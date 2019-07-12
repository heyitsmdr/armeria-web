package armeria

import (
	"sync"
)

// Coords store positional information relative to an Area.
type Coords struct {
	sync.RWMutex
	X int `json:"x"`
	Y int `json:"y"`
	Z int `json:"z"`
	I int `json:"-"`
}

// Location stores where something is within the world.
type Location struct {
	sync.RWMutex
	UnsafeAreaUUID string  `json:"area"`
	Coords         *Coords `json:"coords"`
}

// NewLocation creates and returns a new Location at instance 0.
func NewLocation(areaUuid string, x int, y int, z int) *Location {
	return &Location{
		UnsafeAreaUUID: areaUuid,
		Coords: &Coords{
			X: x,
			Y: y,
			Z: z,
			I: 0,
		},
	}
}

// SetAreaUUID sets the Area UUID for the Location.
func (l *Location) SetAreaUUID(uuid string) {
	l.Lock()
	defer l.Unlock()

	l.UnsafeAreaUUID = uuid
}

// Area returns the Area object referenced by the Location.
func (l *Location) Area() *Area {
	l.RLock()
	defer l.RUnlock()

	for _, a := range Armeria.worldManager.Areas() {
		if a.Id() == l.UnsafeAreaUUID {
			return a
		}
	}

	return nil
}

// Room returns the Room object referenced by the Location.
func (l *Location) Room() *Room {
	a := l.Area()
	if a == nil {
		return nil
	}

	return a.RoomAt(l.Coords)
}

func (c *Coords) Get() {

}

// Set sets the x, y, z and i values of the Coords.
func (c *Coords) Set(x int, y int, z int, i int) {
	c.Lock()
	defer c.Unlock()

	c.X = x
	c.Y = y
	c.Z = z
	c.I = i
}
