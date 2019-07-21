package armeria

import "fmt"

const (
	TextStyleMonospace int = 0
)

const (
	TextStatement int = 0
	TextQuestion  int = 1
	TextExclaim   int = 2
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

// TextPunctuation will automatically punctuate a string and return the punctuation type.
func TextPunctuation(text string) (string, int) {
	lastChar := text[len(text)-1:]

	if lastChar == "." {
		return text, TextStatement
	} else if lastChar == "?" {
		return text, TextQuestion
	} else if lastChar == "!" {
		return text, TextExclaim
	} else {
		return text + ".", TextStatement
	}
}
