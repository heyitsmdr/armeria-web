package armeria

import (
	"armeria/internal/pkg/misc"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// A Character is the player's logged in character.
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
	UnsafeMobConvo       *Conversation     `json:"-"`
	player               *Player
}

// PronounType is used to determine the correct pronoun (he/she etc.)
type PronounType int

// Character constants
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

// ValidSettings returns all valid settings for a unsafeCharacter.
func ValidSettings() []string {
	return []string{
		SettingBrief,
	}
}

// SettingDesc is used to retrieve the description of a unsafeCharacter setting.
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
	// register the unsafeCharacter with registry
	Armeria.registry.Register(c, c.ID(), RegistryTypeCharacter)
}

// ID returns the uuid of the unsafeCharacter.
func (c *Character) ID() string {
	return c.UUID
}

// Type returns the object type, since Character implements the ContainerObject interface.
func (c *Character) Type() ContainerObjectType {
	return ContainerObjectTypeCharacter
}

// Name returns the raw unsafeCharacter name.
func (c *Character) Name() string {
	c.RLock()
	defer c.RUnlock()
	return c.UnsafeName
}

// FormattedName returns the formatted unsafeCharacter name.
func (c *Character) FormattedName() string {
	c.RLock()
	defer c.RUnlock()
	return TextStyle(c.UnsafeName, WithBold())
}

// FormattedNameWithTitle returns the formatted unsafeCharacter name including the unsafeCharacter's title (if set).
func (c *Character) FormattedNameWithTitle() string {
	c.RLock()
	defer c.RUnlock()

	title := c.UnsafeAttributes["title"]
	if title != "" {
		return fmt.Sprintf("%s (%s)", TextStyle(c.UnsafeName, WithBold()), title)
	}

	return TextStyle(c.UnsafeName, WithBold())
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

// PasswordHash returns the unsafeCharacter's already-encrypted password as an md5 hash.
func (c *Character) PasswordHash() string {
	c.RLock()
	defer c.RUnlock()

	b := []byte(c.UnsafePassword)
	return fmt.Sprintf("%x", md5.Sum(b))
}

// Inventory returns the unsafeCharacter's inventory.
func (c *Character) Inventory() *ObjectContainer {
	c.RLock()
	defer c.RUnlock()

	return c.UnsafeInventory
}

// Player returns the parent that is playing the unsafeCharacter.
func (c *Character) Player() *Player {
	c.RLock()
	defer c.RUnlock()

	return c.player
}

// SetPlayer sets the parent that is playing the unsafeCharacter.
func (c *Character) SetPlayer(p *Player) {
	c.Lock()
	defer c.Unlock()

	c.player = p
}

// MobConvo returns the active mob conversation for the unsafeCharacter.
func (c *Character) MobConvo() *Conversation {
	c.RLock()
	defer c.RUnlock()

	return c.UnsafeMobConvo
}

// SetMobConvo sets the active mob conversation with the unsafeCharacter.
func (c *Character) SetMobConvo(convo *Conversation) {
	c.Lock()
	defer c.Unlock()

	c.UnsafeMobConvo = convo
}

// Room returns the unsafeCharacter's Room based on the object container it is within.
func (c *Character) Room() *Room {
	oc := Armeria.registry.GetObjectContainer(c.ID())
	if oc == nil {
		return nil
	}
	return oc.ParentRoom()
}

// Colorize will color text according to the unsafeCharacter's color settings.
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

// LastSeen returns the Time the unsafeCharacter last successfully logged into the game.
func (c *Character) LastSeen() time.Time {
	c.RLock()
	defer c.RUnlock()
	return c.UnsafeLastSeen
}

// SetLastSeen sets the time the unsafeCharacter last logged into the game.
func (c *Character) SetLastSeen(seen time.Time) {
	c.Lock()
	defer c.Unlock()
	c.UnsafeLastSeen = seen
}

// LoggedIn handles everything that needs to happen when a unsafeCharacter enters the game.
func (c *Character) LoggedIn() {
	room := c.Room()
	area := c.Room().ParentArea

	// Add unsafeCharacter to room
	if room == nil || area == nil {
		Armeria.log.Fatal("unsafeCharacter logged into an invalid area/room",
			zap.String("unsafeCharacter", c.Name()),
		)
		return
	}

	// Show server / unsafeCharacter info
	c.Player().client.ShowText(
		fmt.Sprintf(
			"The server has been running for %s.\n"+
				"You last logged in at %s (server time). ",
			TextStyle(time.Since(Armeria.startTime), WithBold()),
			TextStyle(c.LastSeen().Format("Mon Jan 2 2006 15:04:05 MST"), WithBold()),
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
	c.Player().client.SyncMoney()

	Armeria.log.Info("character entered the game",
		zap.String("character", c.Name()),
	)
}

// LoggedOut handles everything that needs to happen when a unsafeCharacter leaves the game.
func (c *Character) LoggedOut() {
	room := c.Room()
	area := c.Room().ParentArea

	// Remove unsafeCharacter from room
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
	for key := range c.UnsafeTempAttributes {
		delete(c.UnsafeTempAttributes, key)
	}

	// Stop any on-going mob conversations
	if c.MobConvo() != nil {
		c.MobConvo().Cancel()
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

// Money returns the character's money as a float.
func (c *Character) Money() float64 {
	money := c.Attribute(AttributeMoney)
	f, err := strconv.ParseFloat(money, 64)
	if err != nil {
		Armeria.log.Fatal("unable to convert money to float64",
			zap.Error(err),
		)
	}
	return f
}

// RemoveMoney attempts to remove money from the character and returns True if they can afford it.
func (c *Character) RemoveMoney(amount float64) bool {
	money := c.Money()
	if amount > money {
		return false
	}

	_ = c.SetAttribute(AttributeMoney, fmt.Sprintf("%.2f", money-amount))

	return true
}

// AddMoney adds money to the character.
func (c *Character) AddMoney(amount float64) {
	_ = c.SetAttribute(AttributeMoney, fmt.Sprintf("%.2f", c.Money()+amount))
}

// SetSetting sets a unsafeCharacter setting and only valid settings can be set.
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
		return false, CommonInvalidDirection
	}

	if len(c.TempAttribute(TempAttributeGhost)) > 0 {
		return true, ""
	}

	if r.Attribute("type") == "track" {
		return false, "You cannot walk onto the train tracks!"
	}

	return true, ""
}

// Move will move the unsafeCharacter to a new location (no move checks are performed).
func (c *Character) Move(to *Room, msgToChar string, msgToOld string, msgToNew string) {
	oldRoom := c.Room()

	oldRoom.Here().Remove(c.ID())
	if err := to.Here().Add(c.ID()); err != nil {
		Armeria.log.Fatal("error adding unsafeCharacter to destination room")
	}

	for _, char := range oldRoom.Here().Characters(true) {
		char.Player().client.ShowText(msgToOld)
	}

	for _, char := range to.Here().Characters(true, c) {
		char.Player().client.ShowText(msgToNew)
	}

	c.Player().client.ShowText(msgToChar)

	oldArea := oldRoom.ParentArea
	newArea := to.ParentArea
	if oldArea.ID() != newArea.ID() {
		oldArea.CharacterLeft(c, false)
		newArea.CharacterEntered(c, false)
	}

	oldRoom.CharacterEntered(c, false)
	to.CharacterEntered(c, false)

	// Stop any on-going mob conversations
	if c.MobConvo() != nil {
		c.MobConvo().Cancel()
	}
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

// Channels returns the Channel objects for the channels this unsafeCharacter is within.
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

// Online is used to see if the character is online.
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
		})
	}

	inventoryJSON, err := json.Marshal(inventory)
	if err != nil {
		Armeria.log.Fatal("failed to marshal inventory data",
			zap.String("unsafeCharacter", c.UUID),
			zap.Error(err),
		)
	}

	return string(inventoryJSON)
}

// Pronoun is used to determine the appropriate pronoun for the character.
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
