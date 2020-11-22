package armeria

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"go.uber.org/zap"
)

type Ticker struct {
	sync.RWMutex
	Name         string
	Handler      func()
	Interval     time.Duration
	RunAtBoot    bool
	LastStart    time.Time
	LastDuration time.Duration
	Iterations   int
}

type TickManager struct {
	Tickers []*Ticker
}

// NewTickManager creates a new TickManager.
func NewTickManager() *TickManager {
	m := &TickManager{
		Tickers: []*Ticker{
			{
				Name:      "WipeDanglingInstances",
				Handler:   WipeDanglingInstances,
				Interval:  1 * time.Hour,
				RunAtBoot: true,
			},
			{
				Name:     "PeriodicGameSave",
				Handler:  PeriodicGameSave,
				Interval: 2 * time.Minute,
			},
		},
	}

	m.Start()

	return m
}

// Start starts the tickers and immediately runs anything designated to run at boot.
func (m *TickManager) Start() {
	for _, ticker := range m.Tickers {
		if ticker.RunAtBoot {
			ticker.Run()
		}

		c := time.Tick(ticker.Interval)
		go func(t *Ticker) {
			for range c {
				t.Run()
			}
		}(ticker)
	}

	Armeria.log.Info("tickers started",
		zap.Int("count", len(m.Tickers)),
	)
}

func (t *Ticker) Run() {
	t.Lock()
	defer t.Unlock()
	t.LastStart = time.Now()
	t.Handler()
	t.LastDuration = time.Since(t.LastStart)
	t.Iterations = t.Iterations + 1
}

// LastRanString returns the date string for which the ticker last ran.
func (t *Ticker) LastRanString() string {
	t.RLock()
	defer t.RUnlock()
	return fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
		t.LastStart.Year(), t.LastStart.Month(), t.LastStart.Day(),
		t.LastStart.Hour(), t.LastStart.Minute(), t.LastStart.Second())
}

// LastDurationString returns the duration string for which the ticker last ran.
func (t *Ticker) LastDurationString() string {
	t.RLock()
	defer t.RUnlock()
	return t.LastDuration.String()
}

// IterationsString returns the number of iterations, as a string, for which the ticker ran.
func (t *Ticker) IterationsString() string {
	t.RLock()
	defer t.RUnlock()
	return strconv.Itoa(t.Iterations)
}

// WipeDanglingInstances searches for mob and item instances in the database that that have no valid container
// parent. The instances are wiped from the registry and on the next database save, will be removed from the database.
func WipeDanglingInstances() {
	var mobsAndItems []interface{}

	mi := Armeria.registry.GetAllFromType(RegistryTypeMobInstance)
	mobsAndItems = append(mobsAndItems, mi...)
	ii := Armeria.registry.GetAllFromType(RegistryTypeItemInstance)
	mobsAndItems = append(mobsAndItems, ii...)

	for _, o := range mobsAndItems {
		obj := o.(ContainerObject)
		container := Armeria.registry.GetObjectContainer(obj.ID())
		if container == nil {
			Armeria.log.Info(
				"found dangling object instance",
				zap.String("uuid", obj.ID()),
			)

			if obj.Type() == ContainerObjectTypeMob {
				obj.(*MobInstance).Parent.DeleteInstance(obj.(*MobInstance))
			} else if obj.Type() == ContainerObjectTypeItem {
				obj.(*ItemInstance).Parent.DeleteInstance(obj.(*ItemInstance))
			}

			Armeria.log.Info(
				"dangling instance deleted",
				zap.String("uuid", obj.ID()),
			)
		}
	}
}

// PeriodicGameSave flushes the game data to disk.
func PeriodicGameSave() {
	Armeria.Save()
}
