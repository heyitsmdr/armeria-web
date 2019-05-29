package armeria

import (
	"encoding/json"

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
	c := ca.player.GetCharacter()
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
	minimap := ca.player.GetCharacter().GetArea().GetMinimapData()
	ca.player.CallClientAction("setMapData", minimap)
}

// SyncMapLocation sets the character location on the minimap.
func (ca *ClientActions) SyncMapLocation() {
	loc := ca.player.GetCharacter().GetLocationData()
	ca.player.CallClientAction("setCharacterLocation", loc)
}

// SyncRoomObjects sets the current room objects on the client.
func (ca *ClientActions) SyncRoomObjects() {
	obj := ca.player.GetCharacter().GetRoom().GetObjectData()
	ca.player.CallClientAction("setRoomObjects", obj)
}

// ShowObjectEditor displays the object editor on the client.
func (ca *ClientActions) ShowObjectEditor(editorData *ObjectEditorData) {
	// add access key
	c := ca.player.GetCharacter()
	editorData.AccessKey = c.GetName() + "/" + c.GetSaltedPasswordHash("ARM0bj3ct3d1t0rERIA")
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
