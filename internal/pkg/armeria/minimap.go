package armeria

import (
	"encoding/json"
	"log"
)

type MinimapArea struct {
	Name  string        `json:"name"`
	Rooms []MinimapRoom `json:"rooms"`
}

type MinimapRoom struct {
	Title string `json:"title"`
	X     int    `json:"x"`
	Y     int    `json:"y"`
	Z     int    `json:"z"`
}

// GetMinimap returns the JSON used for minimap rendering for a particular area
func (a *Area) GetMinimap() string {
	a.mux.Lock()
	defer a.mux.Unlock()

	var rooms []MinimapRoom

	for _, r := range a.Rooms {
		rooms = append(rooms, MinimapRoom{
			Title: r.GetTitle(),
			X:     r.GetCoords().X,
			Y:     r.GetCoords().Y,
			Z:     r.GetCoords().Z,
		})
	}

	mapArea := &MinimapArea{
		Name:  a.Name,
		Rooms: rooms,
	}

	mapJson, err := json.Marshal(mapArea)
	if err != nil {
		log.Fatalf("[minimap] failed to marshal minimap data: %s", err)
	}

	return string(mapJson)

}
