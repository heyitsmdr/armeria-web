package armeria

import (
	"encoding/json"
	"fmt"
	"os"
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

	ledgersFile, err := os.Open(m.dataFile)
	defer ledgersFile.Close()

	if err != nil {
		Armeria.log.Fatal("failed to load data file",
			zap.String("file", m.dataFile),
			zap.Error(err),
		)
	}

	jsonParser := json.NewDecoder(ledgersFile)

	err = jsonParser.Decode(m)
	if err != nil {
		Armeria.log.Fatal("failed to decode data file",
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
