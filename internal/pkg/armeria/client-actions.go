package armeria

import (
	"encoding/json"
	"fmt"
	"strings"

	"go.uber.org/zap"
)

type ClientActions struct {
	parent *Player
}

type ObjectEditorData struct {
	Name       string                      `json:"name"`
	ObjectType string                      `json:"objectType"`
	Properties []*ObjectEditorDataProperty `json:"properties"`
	AccessKey  string                      `json:"accessKey"`
	TextCoords string                      `json:"textCoords"`
}

type ObjectEditorDataProperty struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	PropType string `json:"propType"`
}

func NewClientActions(p *Player) ClientActions {
	return ClientActions{
		parent: p,
	}
}

// ShowColorizedText displays color-formatted text if there is a character attached to
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

// RenderMap displays the current area on the minimap.
func (ca *ClientActions) RenderMap() {
	minimap := ca.parent.Character().Room().ParentArea.MinimapJSON()
	ca.parent.CallClientAction("setMapData", minimap)
}

// SyncMapLocation sets the character location on the minimap.
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
		ca.parent.CallClientAction("setRoomTitle", r.Attribute("title"))
	}
}

// SyncInventory renders the inventory on the client.
func (ca *ClientActions) SyncInventory() {
	inv := ca.parent.Character().InventoryJSON()
	ca.parent.CallClientAction("setInventory", inv)
}

// ShowObjectEditor displays the object editor on the client.
func (ca *ClientActions) ShowObjectEditor(editorData *ObjectEditorData) {
	// add access key
	c := ca.parent.Character()
	editorData.AccessKey = c.Name() + "/" + c.PasswordHash()
	j, err := json.Marshal(editorData)
	if err != nil {
		Armeria.log.Fatal("failed to marshal data",
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
