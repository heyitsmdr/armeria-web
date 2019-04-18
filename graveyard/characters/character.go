package characters

import "log"

type Character struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (c *Character) ShowName() {
	log.Printf("Name of character is: %s", c.Name)
}
