package armeria

import (
	"armeria/internal/pkg/cloud"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"go.uber.org/zap"
)

type MobManager struct {
	sync.RWMutex
	dataFile   string
	UnsafeMobs []*Mob `json:"mobs"`
}

func NewMobManager() *MobManager {
	m := &MobManager{
		dataFile: fmt.Sprintf("%s/mobs.json", Armeria.dataPath),
	}

	m.LoadMobs()
	m.AttachParents()

	return m
}

// LoadMobs loads the mobs from disk into memory.
func (m *MobManager) LoadMobs() {
	m.Lock()
	defer m.Unlock()

	err := json.Unmarshal(Armeria.storageManager.ReadFile(cloud.MobsFile), m)
	if err != nil {
		Armeria.log.Fatal("failed to unmarshal data file",
			zap.String("file", m.dataFile),
			zap.Error(err),
		)
	}

	for _, mob := range m.UnsafeMobs {
		mob.Init()

		for _, mi := range mob.Instances() {
			mi.Init()
		}
	}

	Armeria.log.Info("mobs loaded",
		zap.Int("count", len(m.UnsafeMobs)),
	)
}

// SaveMobs writes the in-memory mobs to disk.
func (m *MobManager) SaveMobs() {
	m.RLock()
	defer m.RUnlock()

	mobsFile, err := os.Create(m.dataFile)
	defer mobsFile.Close()

	raw, err := json.Marshal(m)
	if err != nil {
		Armeria.log.Fatal("failed to marshal data",
			zap.Error(err),
		)
	}

	bytes := Armeria.storageManager.WriteFile(cloud.MobsFile, "application/json", raw)

	Armeria.log.Info("wrote data to cloud",
		zap.Int("bytes", bytes),
	)
}

// AttachParents attaches a pointer to MobInstance that references the parent Mob.
func (m *MobManager) AttachParents() {
	m.RLock()
	defer m.RUnlock()

	for _, m := range m.UnsafeMobs {
		for _, mi := range m.Instances() {
			mi.Parent = m
		}
	}
}

// MobByName returns the matching Mob.
func (m *MobManager) MobByName(name string) *Mob {
	m.RLock()
	defer m.RUnlock()

	for _, mob := range m.UnsafeMobs {
		if strings.ToLower(mob.Name()) == strings.ToLower(name) {
			return mob
		}
	}

	return nil
}

// Mobs returns all of the in-memory Mobs.
func (m *MobManager) Mobs() []*Mob {
	m.RLock()
	defer m.RUnlock()

	return m.UnsafeMobs
}

// CreateMob creates a new Mob instance, but doesn't add it to memory.
func (m *MobManager) CreateMob(name string) *Mob {
	mob := &Mob{
		UnsafeName:       name,
		UnsafeAttributes: make(map[string]string),
	}

	// Create the script file.
	content := fmt.Sprintf("-- %s Script", name)
	_ = ioutil.WriteFile(mob.ScriptFile(), []byte(content), 0644)

	return mob
}

// AddMob adds a new Mob reference to memory.
func (m *MobManager) AddMob(mob *Mob) {
	m.Lock()
	defer m.Unlock()

	m.UnsafeMobs = append(m.UnsafeMobs, mob)
}

// RemoveMob removes an existing Mob reference from memory.
func (m *MobManager) RemoveMob(mob *Mob) {
	m.Lock()
	defer m.Unlock()

	var found bool
	for idx, inst := range m.UnsafeMobs {
		if inst.Name() == mob.Name() {
			m.UnsafeMobs[idx] = m.UnsafeMobs[len(m.UnsafeMobs)-1]
			m.UnsafeMobs = m.UnsafeMobs[:len(m.UnsafeMobs)-1]
			found = true
			break
		}
	}

	if !found {
		return
	}

	// Delete the script file.
	_ = os.Remove(mob.ScriptFile())

	// Delete the picture file.
	picture := mob.Attribute(AttributePicture)
	if len(picture) > 0 {
		DeleteObjectPictureFromDisk(picture)
	}
}
