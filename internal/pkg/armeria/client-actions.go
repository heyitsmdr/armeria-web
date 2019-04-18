package armeria

type ClientActions struct {
	player *Player
}

func NewClientActions(p *Player) *ClientActions {
	return &ClientActions{
		player: p,
	}
}

// ShowText displays raw text on the player's main text window
func (ca *ClientActions) ShowText(text string) {
	ca.player.CallClientAction("showText", text)
}
