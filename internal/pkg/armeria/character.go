package armeria

import (
	"armeria/internal/pkg/misc"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type Character struct {
	sync.RWMutex
	UUID                 string            `json:"uuid"`
	UnsafeName           string            `json:"name"`
	UnsafePassword       string            `json:"password"`
	UnsafeAttributes     map[string]string `json:"attributes"`
	UnsafeSettings       map[string]string `json:"settings"`
	UnsafeInventory      *ObjectContainer  `json:"inventory"`
	UnsafeTempAttributes map[string]string `json:"-"`
	UnsafeLastSeen       time.Time         `json:"lastSeen"`
	player               *Player
}

type PronounType int

const (
	ColorRoomTitle int = iota
	ColorSay
	ColorMovement
	ColorMovementAlt
	ColorError
	ColorRoomDirs
	ColorWhisper
	ColorSuccess
	ColorCmdHelp
	ColorChannelGeneral
	ColorChannelCore
	ColorChannelBuilders

	SettingBrief string = "brief"

	PronounSubjective PronounType = iota
	PronounPossessiveAdjective
	PronounPossessiveAbsolute
	PronounObjective
)

// ValidSettings returns all valid settings for a character.
func ValidSettings() []string {
	return []string{
		SettingBrief,
	}
}

// SettingDesc is used to retrieve the description of a character setting.
func SettingDesc(name string) string {
	switch name {
	case SettingBrief:
		return "Toggle short room descriptions when moving."
	}

	return ""
}

// SettingDefault is used as a fallback for setting values.
func SettingDefault(name string) string {
	switch name {
	case SettingBrief:
		return "false"
	}

	return ""
}

// Init is called when the Character is created or loaded from disk.
func (c *Character) Init() {
	// initialize UnsafeInventory on characters that don't have it defined
	if c.UnsafeInventory == nil {
		c.UnsafeInventory = NewObjectContainer(35)
	}
	// attach self as container's parent
	c.UnsafeInventory.AttachParent(c, ContainerParentTypeCharacter)
	// sync container
	c.UnsafeInventory.Sync()
	// register the character with registry
	Armeria.registry.Register(c, c.ID(), RegistryTypeCharacter)
}

// ID returns the uuid of the character.
func (c *Character) ID() string {
	return c.UUID
}

// Type returns the object type, since Character implements the ContainerObject interface.
func (c *Character) Type() ContainerObjectType {
	return ContainerObjectTypeCharacter
}

// Name returns the raw character name.
func (c *Character) Name() string {
	c.RLock()
	defer c.RUnlock()
	return c.UnsafeName
}

// FormattedName returns the formatted character name.
func (c *Character) FormattedName() string {
	c.RLock()
	defer c.RUnlock()
	return TextStyle(c.UnsafeName, TextStyleBold)
}

// FormattedNameWithTitle returns the formatted character name including the character's title (if set).
func (c *Character) FormattedNameWithTitle() string {
	c.RLock()
	defer c.RUnlock()

	title := c.UnsafeAttributes["title"]
	if title != "" {
		return fmt.Sprintf("%s (%s)", TextStyle(c.UnsafeName, TextStyleBold), title)
	}

	return TextStyle(c.UnsafeName, TextStyleBold)
}

// CheckPassword returns a bool indicating whether the password is correct or not.
func (c *Character) CheckPassword(pw string) bool {
	c.RLock()
	defer c.RUnlock()

	byteHash := []byte(c.UnsafePassword)
	err := bcrypt.CompareHashAndPassword(byteHash, []byte(pw))
	if err != nil {
		return false
	}

	return true
}

// SetPassword hashes and sets a new password for the Character.
func (c *Character) SetPassword(pw string) {
	c.Lock()
	defer c.Unlock()

	hash, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.MinCost)
	if err != nil {
		Armeria.log.Fatal("error generating password hash",
			zap.Error(err),
		)
	}

	c.UnsafePassword = string(hash)
}

// PasswordHash returns the character's already-encrypted password as an md5 hash.
func (c *Character) PasswordHash() string {
	c.RLock()
	defer c.RUnlock()

	b := []byte(c.UnsafePassword)
	return fmt.Sprintf("%x", md5.Sum(b))
}

// Inventory returns the character's inventory.
func (c *Character) Inventory() *ObjectContainer {
	c.RLock()
	defer c.RUnlock()

	return c.UnsafeInventory
}

// Player returns the parent that is playing the character.
func (c *Character) Player() *Player {
	c.RLock()
	defer c.RUnlock()

	return c.player
}

// SetPlayer sets the parent that is playing the character.
func (c *Character) SetPlayer(p *Player) {
	c.Lock()
	defer c.Unlock()

	c.player = p
}

// Room returns the character's Room based on the object container it is within.
func (c *Character) Room() *Room {
	oc := Armeria.registry.GetObjectContainer(c.ID())
	if oc == nil {
		return nil
	}
	return oc.ParentRoom()
}

// Colorize will color text according to the character's color settings.
func (c *Character) Colorize(text string, color int) string {
	c.RLock()
	defer c.RUnlock()

	switch color {
	case ColorRoomTitle:
		return fmt.Sprintf("<span style='color:#6e94ff;font-weight:600'>%s</span>", text)
	case ColorSay:
		return fmt.Sprintf("<span style='color:#ffeb3b'>%s</span>", text)
	case ColorMovement:
		return fmt.Sprintf("<span style='color:#00bcd4'>%s</span>", text)
	case ColorMovementAlt:
		return fmt.Sprintf("<span style='color:#00ffc6'>%s</span>", text)
	case ColorError:
		return fmt.Sprintf("<span style='color:#e91e63'>%s</span>", text)
	case ColorRoomDirs:
		return fmt.Sprintf("<span style='color:#4c9af3'>%s</span>", text)
	case ColorWhisper:
		return fmt.Sprintf("<span style='color:#b730f7'>%s</span>", text)
	case ColorSuccess:
		return fmt.Sprintf("<span style='color:#8ee22b'>%s</span>", text)
	case ColorCmdHelp:
		return fmt.Sprintf("<span style='color:#e9761e'>%s</span>", text)
	case ColorChannelGeneral:
		return fmt.Sprintf("<span style='color:#3bffdc'>%s</span>", text)
	case ColorChannelCore:
		return fmt.Sprintf("<span style='color:#ff5722'>%s</span>", text)
	case ColorChannelBuilders:
		return fmt.Sprintf("<span style='color:#007cff'>%s</span>", text)
	default:
		return text
	}
}

// LastSeen returns the Time the character last successfully logged into the game.
func (c *Character) LastSeen() time.Time {
	c.RLock()
	defer c.RUnlock()
	return c.UnsafeLastSeen
}

// SetLastSeen sets the time the character last logged into the game.
func (c *Character) SetLastSeen(seen time.Time) {
	c.Lock()
	defer c.Unlock()
	c.UnsafeLastSeen = seen
}

// LoggedIn handles everything that needs to happen when a character enters the game.
func (c *Character) LoggedIn() {
	room := c.Room()
	area := c.Room().ParentArea

	// Add character to room
	if room == nil || area == nil {
		Armeria.log.Fatal("character logged into an invalid area/room",
			zap.String("character", c.Name()),
		)
		return
	}

	// Show server / character info
	c.Player().client.ShowText(
		fmt.Sprintf(
			"The server has been running for %s.\n"+
				"You last logged in at %s (server time). ",
			TextStyle(time.Since(Armeria.startTime), TextStyleBold),
			TextStyle(c.LastSeen().Format("Mon Jan 2 2006 15:04:05 MST"), TextStyleBold),
		),
	)

	// Update lastSeen
	c.SetLastSeen(time.Now())

	// Use command: /look
	Armeria.commandManager.ProcessCommand(c.Player(), "look", false)

	// Show message to others in the same room
	for _, char := range room.Here().Characters(true, c) {
		pc := char.Player()
		pc.client.ShowText(
			fmt.Sprintf("%s connected and appeared here with you.", c.Name()),
		)
	}

	area.CharacterEntered(c, true)
	room.CharacterEntered(c, true)

	c.Player().client.SyncInventory()
	c.Player().client.SyncPermissions()
	c.Player().client.SyncPlayerInfo()

	Armeria.log.Info("character entered the game",
		zap.String("character", c.Name()),
	)
}

// LoggedOut handles everything that needs to happen when a character leaves the game.
func (c *Character) LoggedOut() {
	room := c.Room()
	area := c.Room().ParentArea

	// Remove character from room
	if room == nil || area == nil {
		Armeria.log.Fatal("character logged out of an invalid area/room",
			zap.String("character", c.Name()),
		)
		return
	}

	// Show message to others in the same room
	for _, char := range room.Here().Characters(true, c) {
		pc := char.Player()
		pc.client.ShowText(
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

// TempAttribute retrieves a previously-saved temp attribute.
func (c *Character) TempAttribute(name string) string {
	c.RLock()
	defer c.RUnlock()

	return c.UnsafeTempAttributes[name]
}

// SetTempAttribute sets a temporary attribute, which is cleared on log out. Additionally, these
// attributes are not validated.
func (c *Character) SetTempAttribute(name string, value string) {
	c.Lock()
	defer c.Unlock()

	if c.UnsafeTempAttributes == nil {
		c.UnsafeTempAttributes = make(map[string]string)
	}

	c.UnsafeTempAttributes[name] = value
}

// SetAttribute sets a permanent attribute and only valid attributes can be set.
func (c *Character) SetAttribute(name string, value string) error {
	c.Lock()
	defer c.Unlock()

	if !misc.Contains(ValidCharacterAttributes(), name) {
		return errors.New("attribute name is invalid")
	}

	c.UnsafeAttributes[name] = value
	return nil
}

// Attribute returns a permanent attribute.
func (c *Character) Attribute(name string) string {
	c.RLock()
	defer c.RUnlock()

	if len(c.UnsafeAttributes[name]) == 0 {
		return CharacterAttributeDefault(name)
	}

	return c.UnsafeAttributes[name]
}

// SetSetting sets a character setting and only valid settings can be set.
func (c *Character) SetSetting(name string, value string) error {
	c.Lock()
	defer c.Unlock()

	if !misc.Contains(ValidSettings(), name) {
		return errors.New("setting name is invalid")
	}
	c.UnsafeSettings[name] = value
	return nil
}

// Setting returns a setting's value.
func (c *Character) Setting(name string) string {
	c.RLock()
	defer c.RUnlock()

	if len(c.UnsafeSettings[name]) == 0 {
		return SettingDefault(name)
	}

	return c.UnsafeSettings[name]
}

// MoveAllowed will check if moving to a particular location is valid/allowed.
func (c *Character) MoveAllowed(r *Room) (bool, string) {
	if r == nil {
		return false, "You cannot move that way."
	}

	if len(c.TempAttribute(TempAttributeGhost)) > 0 {
		return true, ""
	}

	if r.Attribute("type") == "track" {
		return false, "You cannot walk onto the train tracks!"
	}

	return true, ""
}

// Move will move the character to a new location (no move checks are performed).
func (c *Character) Move(to *Room, msgToChar string, msgToOld string, msgToNew string) {
	oldRoom := c.Room()

	oldRoom.Here().Remove(c.ID())
	if err := to.Here().Add(c.ID()); err != nil {
		Armeria.log.Fatal("error adding character to destination room")
	}

	for _, char := range oldRoom.Here().Characters(true, nil) {
		char.Player().client.ShowText(msgToOld)
	}

	for _, char := range to.Here().Characters(true, c) {
		char.Player().client.ShowText(msgToNew)
	}

	c.Player().client.ShowText(msgToChar)

	oldArea := oldRoom.ParentArea
	newArea := to.ParentArea
	if oldArea.Id() != newArea.Id() {
		oldArea.CharacterLeft(c, false)
		newArea.CharacterEntered(c, false)
	}

	oldRoom.CharacterEntered(c, false)
	to.CharacterEntered(c, false)
}

// EditorData returns the JSON used for the object editor.
func (c *Character) EditorData() *ObjectEditorData {
	var props []*ObjectEditorDataProperty
	for _, attrName := range ValidCharacterAttributes() {
		propType := "editable"
		if attrName == AttributePicture {
			propType = "picture"
		}

		props = append(props, &ObjectEditorDataProperty{
			PropType: propType,
			Name:     attrName,
			Value:    c.Attribute(attrName),
		})
	}

	return &ObjectEditorData{
		UUID:       c.ID(),
		Name:       c.Name(),
		ObjectType: "character",
		Properties: props,
	}
}

// HasPermission returns true if the Character has a particular permission.
func (c *Character) HasPermission(p string) bool {
	c.RLock()
	defer c.RUnlock()

	perms := strings.Split(c.UnsafeAttributes[AttributePermissions], " ")
	return misc.Contains(perms, p)
}

// Channels returns the Channel objects for the channels this character is within.
func (c *Character) Channels() []*Channel {
	var channels []*Channel

	for _, channel := range strings.Split(c.Attribute(AttributeChannels), ",") {
		ch := ChannelByName(channel)
		if ch != nil {
			channels = append(channels, ch)
		}
	}

	return channels
}

// InChannel returns true if the Character is in a particular channel.
func (c *Character) InChannel(ch *Channel) bool {
	c.RLock()
	defer c.RUnlock()

	channelsString := c.UnsafeAttributes[AttributeChannels]
	return misc.Contains(strings.Split(strings.ToLower(channelsString), ","), strings.ToLower(ch.Name))
}

func (c *Character) Online() bool {
	return c.Player() != nil
}

// JoinChannel adds a channel to the Character's channel list so that they will receive messages
// on that channel.
func (c *Character) JoinChannel(ch *Channel) {
	chs := strings.Split(c.Attribute(AttributeChannels), ",")
	if len(chs[0]) == 0 {
		chs[0] = ch.Name
	} else {
		chs = append(chs, ch.Name)
	}
	_ = c.SetAttribute(AttributeChannels, strings.Join(chs, ","))
}

// InventoryJSON returns the JSON used for rendering the inventory on the client.
func (c *Character) InventoryJSON() string {
	var inventory []map[string]interface{}

	for _, ii := range c.Inventory().Items() {
		inventory = append(inventory, map[string]interface{}{
			"uuid":    ii.ID(),
			"picture": ii.Attribute(AttributePicture),
			"slot":    c.Inventory().Slot(ii.ID()),
			"color":   ii.RarityColor(),
			"tooltip": ii.TooltipHTML(),
		})
	}

	inventoryJSON, err := json.Marshal(inventory)
	if err != nil {
		Armeria.log.Fatal("failed to marshal inventory data",
			zap.String("character", c.UUID),
			zap.Error(err),
		)
	}

	return string(inventoryJSON)
}

func (c *Character) Pronoun(pt PronounType) string {
	gender := c.Attribute(AttributeGender)
	if gender == "male" {
		if pt == PronounSubjective {
			return "he"
		} else if pt == PronounPossessiveAbsolute {
			return "his"
		} else if pt == PronounPossessiveAdjective {
			return "his"
		} else if pt == PronounObjective {
			return "him"
		}
	} else if gender == "female" {
		if pt == PronounSubjective {
			return "she"
		} else if pt == PronounPossessiveAbsolute {
			return "hers"
		} else if pt == PronounPossessiveAdjective {
			return "her"
		} else if pt == PronounObjective {
			return "her"
		}
	}

	return ""
}
