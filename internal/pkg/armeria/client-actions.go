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
func (ca *ClientActions) RenderMap(minimapData string) {
	ca.player.CallClientAction("setMapData", minimapData)
}

// SetLocation moves the character on the minimap
func (ca *ClientActions) SetLocation(locationData string) {
	ca.player.CallClientAction("setCharacterLocation", locationData)
}
