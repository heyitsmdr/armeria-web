package armeria

type ClientActions struct {
	player *Player
}

func NewClientActions(p *Player) *ClientActions {
	return &ClientActions{
		player: p,
	}
}

// ShowText displays text on the player's main text window
func (ca *ClientActions) ShowText(text string) {
	ca.player.CallClientAction("showText", "\n"+text)
}

// ShowRawText displays raw text on the player's main text window
func (ca *ClientActions) ShowRawText(text string) {
	ca.player.CallClientAction("showText", text)
}

// RenderMap displays the current area on the minimap
func (ca *ClientActions) RenderMap() {
	minimap := ca.player.GetCharacter().GetArea().GetMinimapData()
	ca.player.CallClientAction("setMapData", minimap)
}

// SyncMapLocation sets the character location on the minimap
func (ca *ClientActions) SyncMapLocation() {
	loc := ca.player.GetCharacter().GetLocationData()
	ca.player.CallClientAction("setCharacterLocation", loc)
}

// SyncRoomObjects sets the current room objects on the client
func (ca *ClientActions) SyncRoomObjects() {
	obj := ca.player.GetCharacter().GetRoom().GetObjectData()
	ca.player.CallClientAction("setRoomObjects", obj)
}

// Disconnect requests that the client disconnects from the server
func (ca *ClientActions) Disconnect() {
	ca.player.CallClientAction("disconnect", nil)
}
