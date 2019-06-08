package armeria

import (
	"armeria/internal/pkg/misc"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"go.uber.org/zap"
)

type Character struct {
	UnsafeName           string            `json:"name"`
	UnsafePassword       string            `json:"password"`
	UnsafeLocation       *Location         `json:"location"`
	UnsafeAttributes     map[string]string `json:"attributes"`
	UnsafeTempAttributes map[string]string
	player               *Player
	mux                  sync.Mutex
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

// ValidCharacterAttributes returns an array of valid attributes that can be permanently set.
func ValidCharacterAttributes() []string {
	return []string{
		"picture",
		"role",
		"title",
	}
}

// CharacterAttributeDefault returns the default value for a particular attribute.
func CharacterAttributeDefault(name string) string {
	switch name {

	}

	return ""
}

// Type returns the object type, since Character implements the Object interface.
func (c *Character) Type() int {
	return ObjectTypeCharacter
}

// UnsafeName returns the raw character name.
func (c *Character) Name() string {
	c.mux.Lock()
	defer c.mux.Unlock()
	return c.UnsafeName
}

// FormattedName returns the formatted character name.
func (c *Character) FormattedName() string {
	c.mux.Lock()
	defer c.mux.Unlock()
	return fmt.Sprintf("[b]%s[/b]", c.UnsafeName)
}

// FormattedNameWithTitle returns the formatted character name including the character's title (if set).
func (c *Character) FormattedNameWithTitle() string {
	c.mux.Lock()
	defer c.mux.Unlock()
	title := c.UnsafeAttributes["title"]
	if title != "" {
		return fmt.Sprintf("[b]%s[/b] (%s)", c.UnsafeName, title)
	}
	return fmt.Sprintf("[b]%s[/b]", c.UnsafeName)
}

// Password returns the character's password.
func (c *Character) Password() string {
	c.mux.Lock()
	defer c.mux.Unlock()
	return c.UnsafePassword
}

// SaltedPasswordHash returns the character's salted password as an md5 hash.
func (c *Character) SaltedPasswordHash(salt string) string {
	c.mux.Lock()
	defer c.mux.Unlock()
	b := []byte(c.UnsafePassword + salt)
	return fmt.Sprintf("%x", md5.Sum(b))
}

// Player returns the player that is playing the character.
func (c *Character) Player() *Player {
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

// UnsafeLocation returns the character's location.
func (c *Character) Location() *Location {
	c.mux.Lock()
	defer c.mux.Unlock()
	return c.UnsafeLocation
}

// LocationData returns the character's location as a JSON-dump.
func (c *Character) LocationData() string {
	c.mux.Lock()
	defer c.mux.Unlock()

	locationJson, err := json.Marshal(c.UnsafeLocation.Coords)
	if err != nil {
		log.Fatalf("[character] failed to marshal location data: %s", err)
	}

	return string(locationJson)
}

// SetLocation sets the character's location.
func (c *Character) SetLocation(l *Location) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.UnsafeLocation = l
}

// Room returns the room that the character is in.
func (c *Character) Room() *Room {
	return Armeria.worldManager.RoomFromLocation(c.UnsafeLocation)
}

// Area returns the area that the character is in.
func (c *Character) Area() *Area {
	return Armeria.worldManager.AreaFromLocation(c.UnsafeLocation)
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
	room := Armeria.worldManager.RoomFromLocation(c.UnsafeLocation)
	area := Armeria.worldManager.AreaFromLocation(c.UnsafeLocation)

	// Add character to room
	if room == nil {
		log.Fatalf("[character] character %s logged in to an invalid room", c.Name())
	}
	room.AddObjectToRoom(c)

	// Use command: /look
	Armeria.commandManager.ProcessCommand(c.Player(), "look")

	// Show message to others in the same room
	roomChars := room.Characters(c)
	for _, char := range roomChars {
		pc := char.Player()
		pc.clientActions.ShowText(
			fmt.Sprintf("%s connected and appeared here with you.", c.Name()),
		)
	}

	area.CharacterEntered(c, true)
	room.CharacterEntered(c, true)

	Armeria.log.Info("character entered the game",
		zap.String("character", c.Name()),
	)
}

// LoggedOut handles everything that needs to happen when a character leaves the game.
func (c *Character) LoggedOut() {
	room := Armeria.worldManager.RoomFromLocation(c.UnsafeLocation)
	area := Armeria.worldManager.AreaFromLocation(c.UnsafeLocation)

	// Remove character from room
	if room == nil {
		log.Fatalf("[character] character %s logged out in an invalid room", c.Name())
	}
	room.RemoveObjectFromRoom(c)

	// Show message to others in the same room
	roomChars := room.Characters(nil)
	for _, char := range roomChars {
		pc := char.Player()
		pc.clientActions.ShowText(
			fmt.Sprintf("%s disconnected and is no longer here with you.", c.Name()),
		)
	}

	area.CharacterLeft(c, true)
	room.CharacterLeft(c, true)

	// Clear temp attributes
	for key, _ := range c.UnsafeTempAttributes {
		delete(c.UnsafeTempAttributes, key)
	}

	Armeria.log.Info("character left the game",
		zap.String("character", c.Name()),
	)
}

// GetTempAttribute retrieves a previously-saved temp attribute.
func (c *Character) GetTempAttribute(name string) string {
	c.mux.Lock()
	defer c.mux.Unlock()
	return c.UnsafeTempAttributes[name]
}

// TempAttribute sets a temporary attribute, which is cleared on log out. Additionally, these
// attributes are not validated.
func (c *Character) TempAttribute(name string, value string) {
	c.mux.Lock()
	defer c.mux.Unlock()

	if c.UnsafeTempAttributes == nil {
		c.UnsafeTempAttributes = make(map[string]string)
	}

	c.UnsafeTempAttributes[name] = value
}

// SetAttribute sets a permanent attribute and only valid attributes can be set.
func (c *Character) SetAttribute(name string, value string) {
	c.mux.Lock()
	defer c.mux.Unlock()

	if c.UnsafeAttributes == nil {
		c.UnsafeAttributes = make(map[string]string)
	}

	if !misc.Contains(ValidCharacterAttributes(), name) {
		Armeria.log.Fatal("attempted to set invalid attribute",
			zap.String("attribute", name),
			zap.String("value", value),
		)
	}

	c.UnsafeAttributes[name] = value
}

// Attribute returns a permanent attribute.
func (c *Character) Attribute(name string) string {
	c.mux.Lock()
	defer c.mux.Unlock()

	if len(c.UnsafeAttributes[name]) == 0 {
		return CharacterAttributeDefault(name)
	}

	return c.UnsafeAttributes[name]
}

// MoveAllowed will check if moving to a particular location is valid/allowed.
func (c *Character) MoveAllowed(to *Location) (bool, string) {
	newRoom := Armeria.worldManager.RoomFromLocation(to)
	if newRoom == nil {
		return false, "You cannot move that way."
	}

	return true, ""
}

// Move will move the character to a new location (no move checks are performed).
func (c *Character) Move(to *Location, msgToChar string, msgToOld string, msgToNew string) {
	oldRoom := c.Room()
	newRoom := Armeria.worldManager.RoomFromLocation(to)

	oldRoom.RemoveObjectFromRoom(c)
	newRoom.AddObjectToRoom(c)

	for _, char := range oldRoom.Characters(nil) {
		char.Player().clientActions.ShowText(msgToOld)
	}

	for _, char := range newRoom.Characters(c) {
		char.Player().clientActions.ShowText(msgToNew)
	}

	c.Player().clientActions.ShowText(msgToChar)

	c.SetLocation(to)

	// Trigger character entered / left events on the new and old rooms, respectively
	newRoom.CharacterEntered(c, false)
	oldRoom.CharacterLeft(c, false)
}

// EditorData returns the JSON used for the object editor.
func (c *Character) GetEditorData() *ObjectEditorData {
	var props []*ObjectEditorDataProperty
	for _, attrName := range ValidCharacterAttributes() {
		propType := "editable"
		if attrName == "picture" {
			propType = "picture"
		}

		props = append(props, &ObjectEditorDataProperty{
			PropType: propType,
			Name:     attrName,
			Value:    c.Attribute(attrName),
		})
	}

	return &ObjectEditorData{
		Name:       c.Name(),
		ObjectType: "character",
		Properties: props,
	}
}
