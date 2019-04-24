package armeria

import (
	"fmt"
	"log"
	"sync"
)

type Character struct {
	Name     string    `json:"name"`
	Password string    `json:"password"`
	Location *Location `json:"location"`
	Role     int       `json:"role"`
	player   *Player
	mux      sync.Mutex
}

const COLOR_ROOM_TITLE int = 0
const COLOR_SAY int = 1

const ROLE_ADMIN int = 0

func (c *Character) GetType() int {
	return OBJECT_TYPE_CHARACTER
}

func (c *Character) GetName() string {
	c.mux.Lock()
	defer c.mux.Unlock()
	return c.Name
}

func (c *Character) GetPassword() string {
	c.mux.Lock()
	defer c.mux.Unlock()
	return c.Password
}

func (c *Character) GetPlayer() *Player {
	c.mux.Lock()
	defer c.mux.Unlock()
	return c.player
}

func (c *Character) SetPlayer(p *Player) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.player = p
}

func (c *Character) GetLocation() *Location {
	c.mux.Lock()
	defer c.mux.Unlock()
	return c.Location
}

func (c *Character) SetLocation(l *Location) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.Location = l
}

func (c *Character) GetRole() int {
	c.mux.Lock()
	defer c.mux.Unlock()
	return c.Role
}

func (c *Character) Colorize(text string, color int) string {
	c.mux.Lock()
	defer c.mux.Unlock()
	switch color {
	case COLOR_ROOM_TITLE:
		return fmt.Sprintf("<span style='color:#6e94ff;font-weight:600'>%s</span>", text)
	case COLOR_SAY:
		return fmt.Sprintf("<span style='color:#ffeb3b'>%s</span>", text)
	default:
		return text
	}
}

func (c *Character) LoggedIn(state *GameState) {
	// Add character to room
	room := state.worldManager.GetRoomFromLocation(c.Location)
	if room == nil {
		log.Fatalf("[character] character %s logged in to an invalid room", c.GetName())
	}
	room.AddObjectToRoom(c)

	// Use command: /look
	state.commandManager.ProcessCommand(c.GetPlayer(), "/look")

	// Show message to others in the same room
	roomChars := room.GetCharacters(c)
	for _, char := range roomChars {
		pc := char.GetPlayer()
		pc.clientActions.ShowText(
			fmt.Sprintf("%s connected to the game, and is now here with you.", c.GetName()),
		)
	}
}

func (c *Character) LoggedOut(state *GameState) {
	// Remove character from room
	room := state.worldManager.GetRoomFromLocation(c.Location)
	if room == nil {
		log.Fatalf("[character] character %s logged out in an invalid room", c.GetName())
	}
	room.RemoveObjectFromRoom(c)

	// Show message to others in the same room
	roomChars := room.GetCharacters(nil)
	for _, char := range roomChars {
		pc := char.GetPlayer()
		pc.clientActions.ShowText(
			fmt.Sprintf("%s disconnected, and is no longer here with you.", c.GetName()),
		)
	}
}

func (c *Character) MoveAllowed(state *GameState, to *Location) (bool, string) {
	newRoom := state.worldManager.GetRoomFromLocation(to)
	if newRoom == nil {
		return false, "You cannot move that way."
	}

	return true, ""
}

func (c *Character) Move(to *Location, msgToOld string, msgToNew string) {
	c.mux.Lock()
	defer c.mux.Unlock()

	if msgToOld != "" {

	}

	c.Location = to
}
