package armeria

import (
	"encoding/json"
	"sync"

	"go.uber.org/zap"
)

// Coords store positional information relative to an ParentArea.
type Coords struct {
	sync.RWMutex
	UnsafeX int `json:"x"`
	UnsafeY int `json:"y"`
	UnsafeZ int `json:"z"`
	UnsafeI int `json:"-"`
}

// NewCoords creates and returns a new Coords.
func NewCoords(x int, y int, z int, i int) *Coords {
	return &Coords{
		UnsafeX: x,
		UnsafeY: y,
		UnsafeZ: z,
		UnsafeI: i,
	}
}

// CopyCoords copies the contents of a Coords pointer and returns a fresh Coords pointer.
func CopyCoords(c *Coords) *Coords {
	return &Coords{
		UnsafeX: c.X(),
		UnsafeY: c.Y(),
		UnsafeZ: c.Z(),
		UnsafeI: c.I(),
	}
}

// X returns the x-coordinate.
func (c *Coords) X() int {
	c.RLock()
	defer c.RUnlock()

	return c.UnsafeX
}

// Y returns the y-coordinate.
func (c *Coords) Y() int {
	c.RLock()
	defer c.RUnlock()

	return c.UnsafeY
}

// Z returns the z-coordinate.
func (c *Coords) Z() int {
	c.RLock()
	defer c.RUnlock()

	return c.UnsafeZ
}

// I returns the i-coordinate.
func (c *Coords) I() int {
	c.RLock()
	defer c.RUnlock()

	return c.UnsafeI
}

// XYZ returns an integer array for the x, y, and z coordinates.
func (c *Coords) XYZ() []int {
	c.RLock()
	defer c.RUnlock()

	return []int{c.UnsafeX, c.UnsafeY, c.UnsafeZ}
}

// XYZI returns an integer array for the x, y, z, and i coordinates.
func (c *Coords) XYZI() []int {
	c.RLock()
	defer c.RUnlock()

	return []int{c.UnsafeX, c.UnsafeY, c.UnsafeZ, c.UnsafeI}
}

// Set sets the x, y, z and i values of the Coords.
func (c *Coords) Set(x int, y int, z int, i int) {
	c.Lock()
	defer c.Unlock()

	c.UnsafeX = x
	c.UnsafeY = y
	c.UnsafeZ = z
	c.UnsafeI = i
}

// JSON returns the JSON-encoded string of the coordinates.
func (c *Coords) JSON() string {
	c.RLock()
	defer c.RUnlock()

	j, err := json.Marshal(c)
	if err != nil {
		Armeria.log.Fatal("failed to marshal coords",
			zap.Error(err),
		)
	}

	return string(j)
}
