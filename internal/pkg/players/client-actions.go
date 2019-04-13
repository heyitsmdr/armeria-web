package players

type clientAction struct {
	player *Player
}

func newClientAction(p *Player) *clientAction {
	return &clientAction{
		player: p,
	}
}

func (ca *clientAction) ShowText(text string) {
	ca.player.CallClientAction("showText", text)
}