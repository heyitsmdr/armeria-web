package armeria

import "fmt"

const (
	TextStyleMonospace int = 0
)

// TextStyle will style text according to one or more styling options.
func TextStyle(text string, opts ...int) string {
	t := text

	for _, o := range opts {
		if o == TextStyleMonospace {
			t = fmt.Sprintf("<span class='monospace'>%s</span>", t)
		}
	}

	return t
}
