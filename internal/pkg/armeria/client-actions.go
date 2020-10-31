package armeria

import (
	"armeria/internal/pkg/sfx"
	"encoding/json"
	"fmt"
	"strings"

	"go.uber.org/zap"
)

// ClientActions is a struct that references a Player and can send data to the client to trigger client-based
// actions within the Vuex store.
type ClientActions struct {
	parent *Player
}

// ObjectEditorData is a struct that contains the json fields for the Object Editor on the client.
type ObjectEditorData struct {
	UUID       string                      `json:"uuid"`
	Name       string                      `json:"name"`
	ObjectType string                      `json:"objectType"`
	Properties []*ObjectEditorDataProperty `json:"properties"`
	AccessKey  string                      `json:"accessKey"`
	TextCoords string                      `json:"textCoords"`
	IsChild    bool                        `json:"isChild"`
}

// ObjectEditorDataProperty is a struct that contains the json fields for each individual property within the
// Object Editor.
type ObjectEditorDataProperty struct {
	Name        string `json:"name"`
	Value       string `json:"value"`
	ParentValue string `json:"parentValue"`
	PropType    string `json:"propType"`
}

// NewClientActions returns a new instance of the ClientActions struct.
func NewClientActions(p *Player) ClientActions {
	return ClientActions{
		parent: p,
	}
}

// ShowColorizedText displays color-formatted text if there is a unsafeCharacter attached to
// the parent instance.
func (ca *ClientActions) ShowColorizedText(text string, color int) {
	var t string
	c := ca.parent.Character()
	if c != nil {
		t = c.Colorize(text, color)
	} else {
		t = text
	}

	ca.ShowText(t)
}

// ShowText displays text on the parent's main text window.
func (ca *ClientActions) ShowText(text string) {
	ca.parent.CallClientAction("showText", "\n"+text)
}

// ShowRawText displays raw text on the parent's main text window.
func (ca *ClientActions) ShowRawText(text string) {
	ca.parent.CallClientAction("showText", text)
}

// SyncMap displays the current area on the minimap.
func (ca *ClientActions) SyncMap() {
	minimap := ca.parent.Character().Room().ParentArea.MinimapJSON()
	ca.parent.CallClientAction("setMapData", minimap)
}

// SyncMapLocation sets the unsafeCharacter location on the minimap.
func (ca *ClientActions) SyncMapLocation() {
	loc := ca.parent.Character().Room().Coords.JSON()
	ca.parent.CallClientAction("setCharacterLocation", loc)
}

// SyncRoomObjects sets the current room objects on the client.
func (ca *ClientActions) SyncRoomObjects() {
	obj := ca.parent.Character().Room().RoomTargetJSON()
	ca.parent.CallClientAction("setRoomObjects", obj)
}

// SyncRoomTitle sets the current room title on the client.
func (ca *ClientActions) SyncRoomTitle() {
	r := ca.parent.Character().Room()
	if ca.parent.Character().HasPermission("CAN_BUILD") {
		c := r.Coords
		ca.parent.CallClientAction("setRoomTitle",
			fmt.Sprintf("%s (%d,%d,%d)", r.Attribute("title"), c.X(), c.Y(), c.Z()),
		)
	} else {
		ca.parent.CallClientAction("setRoomTitle", r.Attribute(AttributeTitle))
	}
}

// SyncInventory renders the inventory on the client.
func (ca *ClientActions) SyncInventory() {
	inv := ca.parent.Character().InventoryJSON()
	ca.parent.CallClientAction("setInventory", inv)
}

// SyncPermissions sets the character permissions on the client (to allow/disallow certain client actions / UI tweaks).
func (ca *ClientActions) SyncPermissions() {
	ca.parent.CallClientAction("setPermissions", ca.parent.Character().Attribute(AttributePermissions))
}

// SyncMoney sets the character's money on the client.
func (ca *ClientActions) SyncMoney() {
	ca.parent.CallClientAction("setMoney", ca.parent.Character().Attribute(AttributeMoney))
}

// SyncPlayerInfo sets the character/player information on the client.
func (ca *ClientActions) SyncPlayerInfo() {
	ca.parent.CallClientAction("setPlayerInfo", ca.parent.Character().Player().PlayerInfoJSON())
}

// SyncCommands sends all of the valid commands to the client (used for auto-complete).
func (ca *ClientActions) SyncCommands() {
	ca.parent.CallClientAction("setCommandDictionary",
		Armeria.commandManager.CharacterCommandDictionaryJSON(ca.parent.Character().Player()),
	)
}

// ShowObjectEditor displays the object editor on the client.
func (ca *ClientActions) ShowObjectEditor(editorData *ObjectEditorData) {
	// add access key
	c := ca.parent.Character()
	editorData.AccessKey = c.Name() + "/" + c.PasswordHash()
	j, err := json.Marshal(editorData)
	if err != nil {
		Armeria.log.Fatal("failed to marshal data for client action: ShowObjectEditor",
			zap.Error(err),
		)
	}

	ca.parent.CallClientAction("setObjectEditorData", string(j))
}

// Disconnect requests that the client disconnects from the server.
func (ca *ClientActions) Disconnect() {
	ca.parent.CallClientAction("disconnect", nil)
}

// ToggleAutologin sets (or disables) auto-login on the client.
func (ca *ClientActions) ToggleAutologin() {
	ca.parent.CallClientAction(
		"toggleAutoLogin",
		strings.ToLower(ca.parent.Character().Name())+":"+ca.parent.Character().PasswordHash(),
	)
}

// SetItemTooltipHTML sets the item's tooltip HTML on the client and stores it in the client-side cache.
func (ca *ClientActions) SetItemTooltipHTML(ii *ItemInstance) {
	ca.parent.CallClientAction("setItemTooltipHTML", ii.TooltipContentJSON())
}

// SetItemTooltipHTMLRaw sets an item's tooltip HTML on the client to some arbitrary value.
func (ca *ClientActions) SetItemTooltipHTMLRaw(uuid, content string) {
	tt := map[string]string{
		"uuid":   uuid,
		"html":   content,
		"rarity": "ffffff",
	}

	ttJSON, err := json.Marshal(tt)
	if err != nil {
		Armeria.log.Fatal("failed to marshal data for client action: SetItemTooltipHTMLRaw",
			zap.Error(err),
		)
	}

	ca.parent.CallClientAction("setItemTooltipHTML", string(ttJSON))
}

// PlaySFX plays a sound effect on the client.
func (ca *ClientActions) PlaySFX(id sfx.ClientSoundEffect, volume float32) {
	data := map[string]interface{}{
		"id":     string(id),
		"volume": volume,
	}
	dataJSON, err := json.Marshal(data)
	if err != nil {
		Armeria.log.Fatal("failed to marshal data for client action: PlaySFX",
			zap.Error(err),
		)
	}
	ca.parent.CallClientAction("playSFX", string(dataJSON))
}
