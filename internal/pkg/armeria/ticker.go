package armeria

import (
	"armeria/internal/pkg/misc"
	"armeria/internal/pkg/sfx"
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
			{
				Name:     "MobSpawner",
				Handler:  MobSpawner,
				Interval: 1 * time.Minute,
			},
			{
				Name:     "MobMovement",
				Handler:  MobMovement,
				Interval: 5 * time.Second,
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

// MobSpawner handles the spawning of mobs into the game world from mob spawners.
func MobSpawner() {
	mobSpawnerItems := Armeria.itemManager.ItemsByAttribute(AttributeType, ItemTypeMobSpawner)
	for _, spawner := range mobSpawnerItems {
		for _, inst := range spawner.Instances() {
			// Find the mob.
			mobStr := inst.Attribute(AttributeSpawnMob)
			mob := Armeria.mobManager.MobByName(mobStr)
			if mob == nil {
				// Let builders know.
				Armeria.channels[ChannelBuilders].Broadcast(
					nil,
					fmt.Sprintf(
						"%s cannot spawn mob '%s' as it does not match any existing mobs.",
						inst.FormattedName(),
						mobStr,
					),
				)
				continue
			}
			// Check the limit. If we reached it, move on.
			mobLimit := inst.AttributeInt(AttributeSpawnLimit)
			existingSpawns := mob.InstancesFromSpawner(inst)
			if len(existingSpawns) >= mobLimit {
				continue
			}
			// Check that the mob spawner is in a room (and not on a character, etc).
			if inst.Room() == nil {
				continue
			}
			// Spawn the mob.
			mobInst := mob.CreateInstance()
			mobInst.SetMobSpawnerUUID(inst.ID())
			_ = inst.Room().Here().Add(mobInst.ID())
			// Refresh the room.
			spawnSFX := mob.Attribute(AttributeSpawnSFX)
			for _, c := range inst.Room().Here().Characters(true) {
				c.Player().client.ShowText(
					fmt.Sprintf("With a flash of light, a %s appeared out of nowhere!", mobInst.FormattedName()),
				)
				c.Player().client.SyncRoomObjects()
				if len(spawnSFX) > 0 {
					c.Player().client.PlaySFX(sfx.ClientSoundEffect(spawnSFX))
				}
			}
		}
	}
}

// MobMovement handles the movement of mobs around the game world.
func MobMovement() {
	for _, m := range Armeria.mobManager.Mobs() {
		for _, mi := range m.Instances() {
			if len(mi.Attribute(AttributeFollowCrumb)) == 0 {
				continue
			}

			// Increment the ticks and determine if we should attempt mob movement.
			mi.IncMoveTicks()
			if mi.MoveTicks() < mi.AttributeInt(AttributeFollowSpeed) {
				continue
			}
			// Reset the tick counter and attempt movement.
			mi.ResetMoveTicks()
			// Find a new random direction, following the crumb.
			possibleRooms := mi.Room().AdjacentRoomsWithItem(mi.Attribute(AttributeFollowCrumb))
			dirStr, newRoom := possibleRooms.Random()
			if newRoom == nil {
				// Let builders know.
				Armeria.channels[ChannelBuilders].Broadcast(
					nil,
					fmt.Sprintf(
						"Mob %s tried to follow breadcrumb '%s' but cannot find an adjacent room with the breadcrumb.",
						mi.FormattedName(),
						mi.Attribute(AttributeFollowCrumb),
					),
				)
				continue
			}
			// Move the mob.
			oldRoom := mi.Room()
			oldRoom.Here().Remove(mi.ID())
			newRoom.Here().Add(mi.ID())
			mobNameString := fmt.Sprintf("A %s", mi.FormattedName())
			if mi.Attribute(AttributeGender) != "thing" {
				mobNameString = mi.FormattedName()
			}
			for _, c := range oldRoom.Here().Characters(true) {
				c.Player().client.ShowText(
					TextStyle(
						fmt.Sprintf("%s travels %s.", mobNameString, misc.MoveToStringFromDir("to the", dirStr)),
						WithUserColor(c, ColorMovement),
					),
				)
				c.Player().client.SyncRoomObjects()
			}
			for _, c := range newRoom.Here().Characters(true) {
				c.Player().client.ShowText(
					TextStyle(
						fmt.Sprintf(
							"%s entered from %s.",
							mobNameString,
							misc.MoveToStringFromDir("the", misc.OppositeDirection(dirStr)),
						),
						WithUserColor(c, ColorMovement),
					),
				)
				c.Player().client.SyncRoomObjects()
			}
		}
	}
}
