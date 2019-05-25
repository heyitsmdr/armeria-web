package armeria

import (
	"armeria/internal/pkg/misc"
	"encoding/json"
	"fmt"
	"log"
	"sync"
)

type Character struct {
	Name           string            `json:"name"`
	Password       string            `json:"password"`
	Location       *Location         `json:"location"`
	Attributes     map[string]string `json:"attributes"`
	TempAttributes map[string]string
	player         *Player
	mux            sync.Mutex
}

const (
	ColorRoomTitle int = 0
	ColorSay       int = 1
	ColorMovement  int = 2
	ColorError     int = 3
	ColorRoomDirs  int = 4
	ColorWhisper   int = 5
	ColorSuccess   int = 6
)

// GetValidCharacterAttributes returns an array of valid attributes that can be permanently set.
func GetValidCharacterAttributes() []string {
	return []string{
		"role",
	}
}

// GetCharacterAttributeDefault returns the default value for a particular attribute.
func GetCharacterAttributeDefault(name string) string {
	switch name {
	case "role":
		return ""
	}

	return ""
}

// GetType returns the object type, since Character uses the Object interface.
func (c *Character) GetType() int {
	return ObjectTypeCharacter
}

// GetName returns the raw character name.
func (c *Character) GetName() string {
	c.mux.Lock()
	defer c.mux.Unlock()
	return c.Name
}

// GetFName returns the formatted character name.
func (c *Character) GetFName() string {
	c.mux.Lock()
	defer c.mux.Unlock()
	return fmt.Sprintf("[b]%s[/b]", c.Name)
}

// GetPassword returns the character's password.
func (c *Character) GetPassword() string {
	c.mux.Lock()
	defer c.mux.Unlock()
	return c.Password
}

// GetPlayer returns the player that is playing the character.
func (c *Character) GetPlayer() *Player {
	c.mux.Lock()
	defer c.mux.Unlock()
	return c.player
}

// SetPlayer sets the player that is playing the character.
func (c *Character) SetPlayer(p *Player) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.player = p
}

// GetLocation returns the character's location.
func (c *Character) GetLocation() *Location {
	c.mux.Lock()
	defer c.mux.Unlock()
	return c.Location
}

// GetLocationData returns the character's location as a JSON-dump.
func (c *Character) GetLocationData() string {
	c.mux.Lock()
	defer c.mux.Unlock()

	locationJson, err := json.Marshal(c.Location.Coords)
	if err != nil {
		log.Fatalf("[character] failed to marshal location data: %s", err)
	}

	return string(locationJson)
}

// SetLocation sets the character's location.
func (c *Character) SetLocation(l *Location) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.Location = l
}

// GetRoom returns the room that the character is in.
func (c *Character) GetRoom() *Room {
	return Armeria.worldManager.GetRoomFromLocation(c.Location)
}

// GetArea returns the area that the character is in.
func (c *Character) GetArea() *Area {
	return Armeria.worldManager.GetAreaFromLocation(c.Location)
}

// Colorize will color text according to the character's color settings.
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
	case ColorError:
		return fmt.Sprintf("<span style='color:#e91e63'>%s</span>", text)
	case ColorRoomDirs:
		return fmt.Sprintf("<span style='color:#4c9af3'>%s</span>", text)
	case ColorWhisper:
		return fmt.Sprintf("<span style='color:#b730f7'>%s</span>", text)
	case ColorSuccess:
		return fmt.Sprintf("<span style='color:#8ee22b'>%s</span>", text)
	default:
		return text
	}
}

// LoggedIn handles everything that needs to happen when a character enters the game.
func (c *Character) LoggedIn() {
	room := Armeria.worldManager.GetRoomFromLocation(c.Location)
	area := Armeria.worldManager.GetAreaFromLocation(c.Location)

	// Add character to room
	if room == nil {
		log.Fatalf("[character] character %s logged in to an invalid room", c.GetName())
	}
	room.AddObjectToRoom(c)

	// Use command: /look
	Armeria.commandManager.ProcessCommand(c.GetPlayer(), "look")

	// Show message to others in the same room
	roomChars := room.GetCharacters(c)
	for _, char := range roomChars {
		pc := char.GetPlayer()
		pc.clientActions.ShowText(
			fmt.Sprintf("%s connected and appeared here with you.", c.GetName()),
		)
	}

	area.OnCharacterEntered(c, true)
	room.OnCharacterEntered(c, true)
}

// LoggedOut handles everything that needs to happen when a character leaves the game.
func (c *Character) LoggedOut() {
	room := Armeria.worldManager.GetRoomFromLocation(c.Location)
	area := Armeria.worldManager.GetAreaFromLocation(c.Location)

	// Remove character from room
	if room == nil {
		log.Fatalf("[character] character %s logged out in an invalid room", c.GetName())
	}
	room.RemoveObjectFromRoom(c)

	// Show message to others in the same room
	roomChars := room.GetCharacters(nil)
	for _, char := range roomChars {
		pc := char.GetPlayer()
		pc.clientActions.ShowText(
			fmt.Sprintf("%s disconnected and is no longer here with you.", c.GetName()),
		)
	}

	area.OnCharacterLeft(c, true)
	room.OnCharacterLeft(c, true)

	// Clear temp attributes
	for key, _ := range c.TempAttributes {
		delete(c.TempAttributes, key)
	}
}

// GetTempAttribute retrieves a previously-saved temp attribute.
func (c *Character) GetTempAttribute(name string) string {
	c.mux.Lock()
	defer c.mux.Unlock()
	return c.TempAttributes[name]
}

// SetTempAttribute sets a temporary attribute, which is cleared on log out. Additionally, these
// attributes are not validated.
func (c *Character) SetTempAttribute(name string, value string) {
	c.mux.Lock()
	defer c.mux.Unlock()

	if c.TempAttributes == nil {
		c.TempAttributes = make(map[string]string)
	}

	c.TempAttributes[name] = value
}

// SetAttribute sets a permanent character attribute and only valid attributes can be set.
func (c *Character) SetAttribute(name string, value string) {
	c.mux.Lock()
	defer c.mux.Unlock()

	if c.Attributes == nil {
		c.Attributes = make(map[string]string)
	}

	if !misc.Contains(GetValidCharacterAttributes(), name) {
		log.Fatalf("[character] attempted set-attribute on a character using an invalid attribute: %s", name)
	}

	c.Attributes[name] = value
}

func (c *Character) GetAttribute(name string) string {
	c.mux.Lock()
	defer c.mux.Unlock()

	if len(c.Attributes[name]) == 0 {
		return GetCharacterAttributeDefault(name)
	}

	return c.Attributes[name]
}

// MoveAllowed will check if moving to a particular location is valid/allowed.
func (c *Character) MoveAllowed(to *Location) (bool, string) {
	newRoom := Armeria.worldManager.GetRoomFromLocation(to)
	if newRoom == nil {
		return false, "You cannot move that way."
	}

	return true, ""
}

// Move will move the character to a new location (no move checks are performed).
func (c *Character) Move(to *Location, msgToChar string, msgToOld string, msgToNew string) {
	oldRoom := c.GetRoom()
	newRoom := Armeria.worldManager.GetRoomFromLocation(to)

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

	// Trigger character entered / left events on the new and old rooms, respectively
	newRoom.OnCharacterEntered(c, false)
	oldRoom.OnCharacterLeft(c, false)
}

// GetEditorData returns the JSON used for the object editor.
func (c *Character) GetEditorData() *ObjectEditorData {
	var props []*ObjectEditorDataProperty
	for _, attrName := range GetValidCharacterAttributes() {
		props = append(props, &ObjectEditorDataProperty{
			PropType: "editable",
			Name:     attrName,
			Value:    c.GetAttribute(attrName),
		})
	}

	return &ObjectEditorData{
		Name:       c.GetName(),
		ObjectType: "character",
		Properties: props,
	}
}
