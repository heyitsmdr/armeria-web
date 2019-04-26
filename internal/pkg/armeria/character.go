package armeria

import (
	"fmt"
	"log"
	"sync"
)

type Character struct {
	gameState *GameState
	Name      string    `json:"name"`
	Password  string    `json:"password"`
	Location  *Location `json:"location"`
	Role      int       `json:"role"`
	player    *Player
	mux       sync.Mutex
}

const (
	COLOR_ROOM_TITLE int = 0
	COLOR_SAY        int = 1
	COLOR_MOVEMENT   int = 2
)

const ROLE_ADMIN int = 0

func (c *Character) Init(state *GameState) {
	c.gameState = state
}

func (c *Character) GetType() int {
	return OBJECT_TYPE_CHARACTER
}

// GetName returns the raw character name
func (c *Character) GetName() string {
	c.mux.Lock()
	defer c.mux.Unlock()
	return c.Name
}

// GetFName returns the formatted character name
func (c *Character) GetFName() string {
	c.mux.Lock()
	defer c.mux.Unlock()
	return fmt.Sprintf("[b]%s[/b]", c.Name)
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

func (c *Character) GetRoom() *Room {
	return c.gameState.worldManager.GetRoomFromLocation(c.Location)
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
	case COLOR_MOVEMENT:
		return fmt.Sprintf("<span style='color:#00bcd4'>%s</span>", text)
	default:
		return text
	}
}

func (c *Character) LoggedIn() {
	// Add character to room
	room := c.gameState.worldManager.GetRoomFromLocation(c.Location)
	if room == nil {
		log.Fatalf("[character] character %s logged in to an invalid room", c.GetName())
	}
	room.AddObjectToRoom(c)

	// Use command: /look
	c.gameState.commandManager.ProcessCommand(c.GetPlayer(), "look")

	// Show message to others in the same room
	roomChars := room.GetCharacters(c)
	for _, char := range roomChars {
		pc := char.GetPlayer()
		pc.clientActions.ShowText(
			fmt.Sprintf("%s connected to the game, and is now here with you.", c.GetName()),
		)
	}
}

func (c *Character) LoggedOut() {
	// Remove character from room
	room := c.gameState.worldManager.GetRoomFromLocation(c.Location)
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

func (c *Character) MoveAllowed(to *Location) (bool, string) {
	newRoom := c.gameState.worldManager.GetRoomFromLocation(to)
	if newRoom == nil {
		return false, "You cannot move that way."
	}

	return true, ""
}

func (c *Character) Move(to *Location, msgToChar string, msgToOld string, msgToNew string) {
	oldRoom := c.GetRoom()
	newRoom := c.gameState.worldManager.GetRoomFromLocation(to)

	oldRoom.RemoveObjectFromRoom(c)
	newRoom.AddObjectToRoom(c)

	for _, char := range oldRoom.GetCharacters(nil) {
		char.GetPlayer().clientActions.ShowText(msgToOld)
	}

	for _, char := range newRoom.GetCharacters(c) {
		char.GetPlayer().clientActions.ShowText(msgToNew)
	}

	c.GetPlayer().clientActions.ShowText(msgToChar)

	c.SetLocation(to)
}
