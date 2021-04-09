package armeria

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	"go.uber.org/zap"
)

type LedgerManager struct {
	sync.RWMutex
	dataFile      string
	UnsafeLedgers []*Ledger `json:"ledgers"`
}

// NewLedgerManager creates a new LedgerManager.
func NewLedgerManager() *LedgerManager {
	m := &LedgerManager{
		dataFile: fmt.Sprintf("%s/ledgers.json", Armeria.dataPath),
	}

	m.LoadLedgers()

	return m
}

// LoadLedgers loads the ledgers from disk into memory.
func (m *LedgerManager) LoadLedgers() {
	m.Lock()
	defer m.Unlock()

	err := json.Unmarshal(Armeria.storageManager.ReadFile("ledgers.json"), m)
	if err != nil {
		Armeria.log.Fatal("failed to unmarshal data file",
			zap.String("file", m.dataFile),
			zap.Error(err),
		)
	}

	Armeria.log.Info("ledgers loaded",
		zap.Int("count", len(m.UnsafeLedgers)),
	)
}

// SaveLedgers writes the in-memory ledgers to disk.
func (m *LedgerManager) SaveLedgers() {
	m.RLock()
	defer m.RUnlock()

	ledgersFile, err := os.Create(m.dataFile)
	defer ledgersFile.Close()

	raw, err := json.Marshal(m)
	if err != nil {
		Armeria.log.Fatal("failed to marshal data",
			zap.Error(err),
		)
	}

	bytes, err := ledgersFile.Write(raw)
	if err != nil {
		Armeria.log.Fatal("failed to write data file",
			zap.String("file", m.dataFile),
			zap.Error(err),
		)
	}

	_ = ledgersFile.Sync()

	Armeria.log.Info("wrote data to file",
		zap.String("file", m.dataFile),
		zap.Int("bytes", bytes),
	)
}

// Ledgers returns all of the in-memory Ledgers.
func (m *LedgerManager) Ledgers() []*Ledger {
	m.RLock()
	defer m.RUnlock()

	return m.UnsafeLedgers
}

// LedgerByName returns the matching Ledger, by name.
func (m *LedgerManager) LedgerByName(name string) *Ledger {
	m.RLock()
	defer m.RUnlock()

	for _, l := range m.UnsafeLedgers {
		if strings.ToLower(l.Name()) == strings.ToLower(name) {
			return l
		}
	}

	return nil
}

// CreateLedger creates a new Ledger instance, but doesn't add it to memory.
func (m *LedgerManager) CreateLedger(name string) *Ledger {
	return &Ledger{
		UnsafeName: name,
	}
}

// AddLedger adds a new Item reference to memory.
func (m *LedgerManager) AddLedger(l *Ledger) {
	m.Lock()
	defer m.Unlock()

	m.UnsafeLedgers = append(m.UnsafeLedgers, l)
}
