package armeria

import "sync"

type Ledger struct {
	sync.RWMutex
	UnsafeName string `json:"name"`
}

func (l *Ledger) Name() string {
	l.RLock()
	defer l.RUnlock()

	return l.UnsafeName
}
