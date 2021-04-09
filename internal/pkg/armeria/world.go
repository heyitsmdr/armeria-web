package armeria

import (
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

	err := json.Unmarshal(Armeria.storageManager.ReadFile("world.json"), m)
	if err != nil {
		Armeria.log.Fatal("failed to unmarshal data file",
			zap.String("file", m.dataFile),
			zap.Error(err),
		)
	}

	for _, a := range m.UnsafeWorld {
		a.Init()

		for _, r := range a.UnsafeRooms {
			r.Init(a)
		}
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

func (m *WorldManager) CreateRoom(a *Area, c *Coords) *Room {
	r := &Room{
		UUID:             uuid.New().String(),
		Coords:           CopyCoords(c),
		UnsafeAttributes: map[string]string{},
		UnsafeHere:       NewObjectContainer(0),
		ParentArea:       a,
	}
	a.AddRoom(r)
	return r
}

func (m *WorldManager) CreateArea(name string) *Area {
	m.Lock()
	defer m.Unlock()

	a := &Area{
		UUID:             uuid.New().String(),
		UnsafeName:       name,
		UnsafeAttributes: make(map[string]string),
	}

	a.Init()

	_ = m.CreateRoom(a, NewCoords(0, 0, 0, 0))

	m.UnsafeWorld = append(m.UnsafeWorld, a)

	return a
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
