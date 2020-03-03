package armeria

import (
	"strings"
	"sync"
)

type LedgerEntry struct {
	ItemName  string  `json:"name"`
	BuyPrice  float64 `json:"buy_price"`
	SellPrice float64 `json:"sell_price"`
}
type Ledger struct {
	sync.RWMutex
	UnsafeName    string         `json:"name"`
	UnsafeEntries []*LedgerEntry `json:"entries"`
}

// Name returns the name of the ledger.
func (l *Ledger) Name() string {
	l.RLock()
	defer l.RUnlock()

	return l.UnsafeName
}

// SetName sets a new name for the ledger.
func (l *Ledger) SetName(name string) {
	l.Lock()
	defer l.Unlock()

	l.UnsafeName = name
}

// AddEntry adds an item to the ledger.
func (l *Ledger) AddEntry(le *LedgerEntry) {
	l.Lock()
	defer l.Unlock()

	l.UnsafeEntries = append(l.UnsafeEntries, le)
}

// RemoveEntry removes an item from the ledger.
func (l *Ledger) RemoveEntry(le *LedgerEntry) {
	l.Lock()
	defer l.Unlock()

	for i, entry := range l.UnsafeEntries {
		if entry.ItemName == le.ItemName {
			l.UnsafeEntries[i] = l.UnsafeEntries[len(l.UnsafeEntries)-1]
			l.UnsafeEntries = l.UnsafeEntries[:len(l.UnsafeEntries)-1]
			break
		}
	}
}

// Contains returns a LedgerEntry if any entry within the ledger matches the item name.
func (l *Ledger) Contains(item string) *LedgerEntry {
	l.RLock()
	defer l.RUnlock()

	for _, entry := range l.UnsafeEntries {
		if strings.ToLower(entry.ItemName) == strings.ToLower(item) {
			return entry
		}
	}

	return nil
}

// Entries returns all of the entries within the ledger.
func (l *Ledger) Entries() []*LedgerEntry {
	l.RLock()
	defer l.RUnlock()

	return l.UnsafeEntries
}
