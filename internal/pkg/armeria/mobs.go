package armeria

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"go.uber.org/zap"
)

type MobManager struct {
	dataFile   string
	UnsafeMobs []*Mob `json:"mobs"`
	mux        sync.Mutex
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
		zap.Int("count", len(m.UnsafeMobs)),
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

	_ = mobsFile.Sync()

	Armeria.log.Info("wrote data to file",
		zap.String("file", m.dataFile),
		zap.Int("bytes", bytes),
	)
}

func (m *MobManager) AddMobInstancesToRooms() {
	m.mux.Lock()
	defer m.mux.Unlock()

	for _, m := range m.UnsafeMobs {
		for _, mi := range m.UnsafeInstances {
			r := Armeria.worldManager.RoomFromLocation(mi.Location())
			if r == nil {
				Armeria.log.Fatal("mob instance in invalid room",
					zap.String("mob", mi.Name()),
					zap.String("uuid", mi.Id()),
					zap.String("location", fmt.Sprintf("%v", mi.Location())),
				)
				return
			}
			r.AddObjectToRoom(mi)
		}
	}
}

// MobByName returns the matching Mob.
func (m *MobManager) MobByName(name string) *Mob {
	m.mux.Lock()
	defer m.mux.Unlock()

	for _, mob := range m.UnsafeMobs {
		if strings.ToLower(mob.Name()) == strings.ToLower(name) {
			return mob
		}
	}

	return nil
}

// Mobs returns all of the in-memory Mobs.
func (m *MobManager) Mobs() []*Mob {
	return m.UnsafeMobs
}

// CreateMob creates a new Mob instance, but doesn't add it to memory.
func (m *MobManager) CreateMob(name string) *Mob {
	mob := &Mob{
		UnsafeName:       name,
		UnsafeAttributes: make(map[string]string),
	}

	// create script file
	content := fmt.Sprintf("-- %s Script", name)
	_ = ioutil.WriteFile(mob.ScriptFile(), []byte(content), 0644)

	return mob
}

// AddMob adds a new Mob reference to memory.
func (m *MobManager) AddMob(mob *Mob) {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.UnsafeMobs = append(m.UnsafeMobs, mob)
}
