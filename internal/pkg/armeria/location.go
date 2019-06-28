package armeria

type Coords struct {
	X int `json:"x"`
	Y int `json:"y"`
	Z int `json:"z"`
	I int `json:"-"`
}

type Location struct {
	AreaUUID string  `json:"area"`
	Coords   *Coords `json:"coords"`
}

func (l *Location) Area() *Area {
	for _, a := range Armeria.worldManager.Areas() {
		if a.Id() == l.AreaUUID {
			return a
		}
	}

	return nil
}

func (l *Location) Room() *Room {
	a := l.Area()
	if a == nil {
		return nil
	}

	return a.RoomAt(l.Coords)
}
