package armeria

import "sync"

type Item struct {
	UnsafeName       string            `json:"name"`
	UnsafeAttributes map[string]string `json:"attributes"`
	//UnsafeInstances  []*ItemInstance    `json:"instances"`
	mux sync.Mutex
}

func (i *Item) Name() string {
	i.mux.Lock()
	defer i.mux.Unlock()
	return i.UnsafeName
}
