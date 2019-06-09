package armeria

import (
	"armeria/internal/pkg/misc"
	"encoding/json"
	"fmt"
	"os"

	"go.uber.org/zap"
)

type WorldManager struct {
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

func (m *WorldManager) AreaFromLocation(l *Location) *Area {
	for _, a := range m.UnsafeWorld {
		if a.Name() == l.AreaName {
			return a
		}
	}

	return nil
}

func (m *WorldManager) RoomFromLocation(l *Location) *Room {
	a := m.AreaFromLocation(l)
	if a == nil {
		return nil
	}

	return a.RoomAt(l.Coords)
}

func (m *WorldManager) RoomInDirection(a *Area, r *Room, direction string) *Room {
	o := misc.DirectionOffsets(direction)
	if o == nil {
		Armeria.log.Fatal("invalid direction provided",
			zap.String("direction", direction),
		)
	}

	loc := &Location{
		AreaName: a.UnsafeName,
		Coords: &Coords{
			X: r.UnafeCoords.X + o["x"],
			Y: r.UnafeCoords.Y + o["y"],
			Z: r.UnafeCoords.Z + o["z"],
		},
	}

	return m.RoomFromLocation(loc)
}
