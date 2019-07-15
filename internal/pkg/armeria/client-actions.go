package armeria

import (
	"encoding/json"
	"fmt"

	"go.uber.org/zap"
)

type ClientActions struct {
	player *Player
}

type ObjectEditorData struct {
	Name       string                      `json:"name"`
	ObjectType string                      `json:"objectType"`
	Properties []*ObjectEditorDataProperty `json:"properties"`
	AccessKey  string                      `json:"accessKey"`
}

type ObjectEditorDataProperty struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	PropType string `json:"propType"`
}

func NewClientActions(p *Player) *ClientActions {
	return &ClientActions{
		player: p,
	}
}

// ShowColorizedText displays color-formatted text if there is a character attached to
// the player instance.
func (ca *ClientActions) ShowColorizedText(text string, color int) {
	c := ca.player.Character()
	if c != nil {
		ca.ShowText(c.Colorize(text, color))
		return
	}
	ca.ShowText(text)
}

// ShowText displays text on the player's main text window.
func (ca *ClientActions) ShowText(text string) {
	ca.player.CallClientAction("showText", "\n"+text)
}

// ShowRawText displays raw text on the player's main text window.
func (ca *ClientActions) ShowRawText(text string) {
	ca.player.CallClientAction("showText", text)
}

// RenderMap displays the current area on the minimap.
func (ca *ClientActions) RenderMap() {
	minimap := ca.player.Character().Area().MinimapJSON()
	ca.player.CallClientAction("setMapData", minimap)
}

// SyncMapLocation sets the character location on the minimap.
func (ca *ClientActions) SyncMapLocation() {
	loc := ca.player.Character().Location().Coords.JSON()
	ca.player.CallClientAction("setCharacterLocation", loc)
}

// SyncRoomObjects sets the current room objects on the client.
func (ca *ClientActions) SyncRoomObjects() {
	obj := ca.player.Character().Room().ObjectData()
	ca.player.CallClientAction("setRoomObjects", obj)
}

// SyncRoomTitle sets the current room title on the client.
func (ca *ClientActions) SyncRoomTitle() {
	r := ca.player.Character().Room()
	if ca.player.Character().HasPermission("CAN_BUILD") {
		c := r.Coords
		ca.player.CallClientAction("setRoomTitle",
			fmt.Sprintf("%s (%d,%d,%d)", r.Attribute("title"), c.X(), c.Y(), c.Z()),
		)
	} else {
		ca.player.CallClientAction("setRoomTitle", r.Attribute("title"))
	}
}

// ShowObjectEditor displays the object editor on the client.
func (ca *ClientActions) ShowObjectEditor(editorData *ObjectEditorData) {
	// add access key
	c := ca.player.Character()
	editorData.AccessKey = c.Name() + "/" + c.PasswordHash()
	j, err := json.Marshal(editorData)
	if err != nil {
		Armeria.log.Fatal("failed to marshal data",
			zap.Error(err),
		)
	}

	ca.player.CallClientAction("setObjectEditorData", string(j))
}

// Disconnect requests that the client disconnects from the server.
func (ca *ClientActions) Disconnect() {
	ca.player.CallClientAction("disconnect", nil)
}
