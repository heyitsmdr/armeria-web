package armeria

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
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

// NewCoordsFromString creates and returns a new Coords from a single string.
func NewCoordsFromString(c string) *Coords {
	sections := strings.Split(c, ",")
	if len(sections) != 3 {
		return nil
	}

	var x, y, z int
	if i, err := strconv.Atoi(sections[0]); err == nil {
		x = i
	}
	if i, err := strconv.Atoi(sections[1]); err == nil {
		y = i
	}
	if i, err := strconv.Atoi(sections[2]); err == nil {
		z = i
	}

	return NewCoords(x, y, z, 0)
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

// SetFrom sets the coordinates from another Coords struct.
func (c *Coords) SetFrom(co *Coords) {
	c.Set(co.X(), co.Y(), co.Z(), co.I())
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

// String returns the coordinates as a string.
func (c *Coords) String() string {
	return fmt.Sprintf("%d,%d,%d", c.UnsafeX, c.UnsafeY, c.UnsafeZ)
}
