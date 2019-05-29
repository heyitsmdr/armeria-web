package armeria

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	"go.uber.org/zap"
)

type MobManager struct {
	dataFile string
	Mobs     []*Mob `json:"mobs"`
	mux      sync.Mutex
}

func NewMobManager() *MobManager {
	m := &MobManager{
		dataFile: fmt.Sprintf("%s/mobs.json", Armeria.dataPath),
	}

	m.LoadMobs()
	m.AddMobInstancesToRooms()

	return m
}

// LoadMobs loads the mobs from disk into memory.
func (m *MobManager) LoadMobs() {
	mobsFile, err := os.Open(m.dataFile)
	defer mobsFile.Close()

	if err != nil {
		Armeria.log.Fatal("failed to load data file",
			zap.String("file", m.dataFile),
			zap.Error(err),
		)
	}

	jsonParser := json.NewDecoder(mobsFile)

	err = jsonParser.Decode(m)
	if err != nil {
		Armeria.log.Fatal("failed to decode data file",
			zap.String("file", m.dataFile),
			zap.Error(err),
		)
	}

	Armeria.log.Info("mobs loaded",
		zap.Int("count", len(m.Mobs)),
	)
}

// SaveMobs writes the in-memory mobs to disk.
func (m *MobManager) SaveMobs() {
	mobsFile, err := os.Create(m.dataFile)
	defer mobsFile.Close()

	raw, err := json.Marshal(m)
	if err != nil {
		Armeria.log.Fatal("failed to marshal data",
			zap.Error(err),
		)
	}

	bytes, err := mobsFile.Write(raw)
	if err != nil {
		Armeria.log.Fatal("failed to write data file",
			zap.String("file", m.dataFile),
			zap.Error(err),
		)
	}

	mobsFile.Sync()

	Armeria.log.Info("wrote data to file",
		zap.String("file", m.dataFile),
		zap.Int("bytes", bytes),
	)
}

func (m *MobManager) AddMobInstancesToRooms() {
	m.mux.Lock()
	defer m.mux.Unlock()

	for _, m := range m.Mobs {
		for _, mi := range m.Instances {
			r := Armeria.worldManager.GetRoomFromLocation(mi.GetLocation())
			if r == nil {
				Armeria.log.Fatal("mob instance in invalid room",
					zap.String("mob", mi.GetName()),
					zap.String("location", fmt.Sprintf("%v", mi.GetLocation())),
				)
				return
			}
			r.AddObjectToRoom(mi)
		}
	}
}

// GetMobByName returns the matching Mob.
func (m *MobManager) GetMobByName(name string) *Mob {
	m.mux.Lock()
	defer m.mux.Unlock()

	for _, m := range m.Mobs {
		if strings.ToLower(m.Name) == strings.ToLower(name) {
			return m
		}
	}

	return nil
}

func (m *MobManager) CreateMob(mob *Mob) {
	m.mux.Lock()
	defer m.mux.Unlock()

	m.Mobs = append(m.Mobs, mob)
}
