package armeria

import (
	"armeria/internal/pkg/misc"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/google/uuid"

	"go.uber.org/zap"
)

type WorldManager struct {
	sync.RWMutex
	dataFile    string
	UnsafeWorld []*Area `json:"world"`
}

func NewWorldManager() *WorldManager {
	m := &WorldManager{
		dataFile: fmt.Sprintf("%s/world.json", Armeria.dataPath),
	}

	m.LoadWorld()

	return m
}

func (m *WorldManager) LoadWorld() {
	m.Lock()
	defer m.Unlock()

	worldFile, err := os.Open(m.dataFile)
	defer worldFile.Close()

	if err != nil {
		Armeria.log.Fatal("failed to load data file",
			zap.String("file", m.dataFile),
			zap.Error(err),
		)
	}

	jsonParser := json.NewDecoder(worldFile)

	err = jsonParser.Decode(m)
	if err != nil {
		Armeria.log.Fatal("failed to decode data file",
			zap.String("file", m.dataFile),
			zap.Error(err),
		)
	}

	Armeria.log.Info("areas loaded",
		zap.Int("count", len(m.UnsafeWorld)),
	)
}

func (m *WorldManager) SaveWorld() {
	m.RLock()
	defer m.RUnlock()

	worldFile, err := os.Create(m.dataFile)
	defer worldFile.Close()

	raw, err := json.Marshal(m)
	if err != nil {
		Armeria.log.Fatal("failed to marshal data",
			zap.Error(err),
		)
	}

	bytes, err := worldFile.Write(raw)
	if err != nil {
		Armeria.log.Fatal("failed to write data file",
			zap.String("file", m.dataFile),
			zap.Error(err),
		)
	}

	_ = worldFile.Sync()

	Armeria.log.Info("wrote data to file",
		zap.String("file", m.dataFile),
		zap.Int("bytes", bytes),
	)
}

func (m *WorldManager) CreateRoom(c *Coords) *Room {
	return &Room{
		UnsafeCoords:     c,
		UnsafeAttributes: map[string]string{},
	}
}

func (m *WorldManager) CreateArea(name string) *Area {
	m.Lock()
	defer m.Unlock()

	a := &Area{
		UUID:             uuid.New().String(),
		UnsafeName:       name,
		UnsafeAttributes: make(map[string]string),
	}

	r := m.CreateRoom(&Coords{0, 0, 0, 0})
	a.AddRoom(r)
	m.UnsafeWorld = append(m.UnsafeWorld, a)

	return a
}

func (m *WorldManager) RoomInDirection(a *Area, r *Room, direction string) *Room {
	o := misc.DirectionOffsets(direction)
	if o == nil {
		Armeria.log.Fatal("invalid direction provided",
			zap.String("direction", direction),
		)
	}

	loc := &Location{
		AreaUUID: a.Id(),
		Coords: &Coords{
			X: r.Coords().X + o["x"],
			Y: r.Coords().Y + o["y"],
			Z: r.Coords().Z + o["z"],
		},
	}

	return loc.Room()
}

func (m *WorldManager) AreaByName(name string) *Area {
	m.RLock()
	defer m.RUnlock()

	for _, a := range m.UnsafeWorld {
		if strings.ToLower(a.Name()) == strings.ToLower(name) {
			return a
		}
	}

	return nil
}

func (m *WorldManager) Areas() []*Area {
	m.RLock()
	defer m.RUnlock()

	return m.UnsafeWorld
}
