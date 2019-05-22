package armeria

import (
	"armeria/internal/pkg/misc"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type WorldManager struct {
	gameState *GameState
	dataFile  string
	World     []*Area `json:"world"`
}

func NewWorldManager(state *GameState) *WorldManager {
	m := &WorldManager{
		gameState: state,
		dataFile:  fmt.Sprintf("%s/world.json", state.dataPath),
	}

	m.LoadWorld()

	return m
}

func (m *WorldManager) LoadWorld() {
	worldFile, err := os.Open(m.dataFile)
	defer worldFile.Close()

	if err != nil {
		log.Fatalf("[world] failed to load from %s: %s", m.dataFile, err)
	}

	jsonParser := json.NewDecoder(worldFile)

	err = jsonParser.Decode(m)
	if err != nil {
		log.Fatalf("[world] failed to decode world file: %s", err)
	}

	log.Printf("[world] loaded %d areas from world file", len(m.World))
}

func (m *WorldManager) SaveWorld() {
	worldFile, err := os.Create(m.dataFile)
	defer worldFile.Close()

	raw, err := json.Marshal(m)
	if err != nil {
		log.Fatalf("[world] failed to marshal world: %s", err)
	}

	bytes, err := worldFile.Write(raw)
	if err != nil {
		log.Fatalf("[world] failed to write to world file: %s", err)
	}

	worldFile.Sync()

	log.Printf("[world] wrote %d bytes to world file", bytes)
}

func (m *WorldManager) GetAreaFromLocation(l *Location) *Area {
	for _, a := range m.World {
		if a.GetName() == l.AreaName {
			return a
		}
	}

	return nil
}

func (m *WorldManager) GetRoomFromLocation(l *Location) *Room {
	a := m.GetAreaFromLocation(l)
	if a == nil {
		return nil
	}

	return a.GetRoom(l.Coords)
}

func (m *WorldManager) GetRoomInDirection(a *Area, r *Room, direction string) *Room {
	o := misc.DirectionOffsets(direction)
	if o == nil {
		log.Fatal("[room] invalid direction")
	}

	loc := &Location{
		AreaName: a.Name,
		Coords: &Coords{
			X: r.Coords.X + o["x"],
			Y: r.Coords.Y + o["y"],
			Z: r.Coords.Z + o["z"],
		},
	}

	return m.GetRoomFromLocation(loc)
}
