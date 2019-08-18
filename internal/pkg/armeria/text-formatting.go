package armeria

import "fmt"

const (
	TextStyleMonospace int = iota
	TextStyleBold
)

const (
	TextStatement int = iota
	TextQuestion
	TextExclaim
)

// TextStyle will style text according to one or more styling options.
func TextStyle(text string, opts ...int) string {
	t := text

	for _, o := range opts {
		switch o {
		case TextStyleBold:
			t = fmt.Sprintf("<span style='font-weight:600'>%s</span>", t)
		case TextStyleMonospace:
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
