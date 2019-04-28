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
	ColorRoomTitle int = 0
	ColorSay       int = 1
	ColorMovement  int = 2
)

const RoleAdmin int = 0

// Init initializes the character when loaded from disk
func (c *Character) Init(state *GameState) {
	c.gameState = state
}

// GetType returns the object type, since Character uses the Object interface
func (c *Character) GetType() int {
	return ObjectTypeCharacter
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

// GetPassword returns the character's password
func (c *Character) GetPassword() string {
	c.mux.Lock()
	defer c.mux.Unlock()
	return c.Password
}

// GetPlayer returns the player that is playing the character
func (c *Character) GetPlayer() *Player {
	c.mux.Lock()
	defer c.mux.Unlock()
	return c.player
}

// SetPlayer sets the player that is playing the character
func (c *Character) SetPlayer(p *Player) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.player = p
}

// GetLocation returns the character's location
func (c *Character) GetLocation() *Location {
	c.mux.Lock()
	defer c.mux.Unlock()
	return c.Location
}

// SetLocation sets the character's location
func (c *Character) SetLocation(l *Location) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.Location = l
}

// GetRoom returns the room that the character is in
func (c *Character) GetRoom() *Room {
	return c.gameState.worldManager.GetRoomFromLocation(c.Location)
}

// GetRole returns the character's permission role
func (c *Character) GetRole() int {
	c.mux.Lock()
	defer c.mux.Unlock()
	return c.Role
}

// Colorize will color text according to the character's color settings
func (c *Character) Colorize(text string, color int) string {
	c.mux.Lock()
	defer c.mux.Unlock()
	switch color {
	case ColorRoomTitle:
		return fmt.Sprintf("<span style='color:#6e94ff;font-weight:600'>%s</span>", text)
	case ColorSay:
		return fmt.Sprintf("<span style='color:#ffeb3b'>%s</span>", text)
	case ColorMovement:
		return fmt.Sprintf("<span style='color:#00bcd4'>%s</span>", text)
	default:
		return text
	}
}

// LoggedIn handles everything that needs to happen when a character enters the game
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

// LoggedOut handles everything that needs to happen when a character leaves the game
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

// MoveAllowed will check if moving to a particular location is valid/allowed
func (c *Character) MoveAllowed(to *Location) (bool, string) {
	newRoom := c.gameState.worldManager.GetRoomFromLocation(to)
	if newRoom == nil {
		return false, "You cannot move that way."
	}

	return true, ""
}

// Move will move the character to a new location (no move checks are performed)
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
