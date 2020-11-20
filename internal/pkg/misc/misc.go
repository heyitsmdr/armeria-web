package misc

import (
	"strings"

	"github.com/google/uuid"

	"github.com/leekchan/accounting"
)

var Money = accounting.Accounting{Symbol: "$", Precision: 2}

// Contains tells whether a contains x.
func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

// Find returns the smallest index i at which x == a[i],
// or len(a) if there is no such index.
func Find(a []string, x string) int {
	for i, n := range a {
		if x == n {
			return i
		}
	}
	return len(a)
}

// DirectionOffsets returns the X/Y/Z offsets for each direction as a map.
func DirectionOffsets(dir string) map[string]int {
	offsets := map[string]map[string]int{
		"north": {"x": 0, "y": 1, "z": 0},
		"south": {"x": 0, "y": -1, "z": 0},
		"east":  {"x": 1, "y": 0, "z": 0},
		"west":  {"x": -1, "y": 0, "z": 0},
		"up":    {"x": 0, "y": 0, "z": 1},
		"down":  {"x": 0, "y": 0, "z": -1},
	}

	return offsets[dir]
}

// OppositeDirection returns the opposite direction string.
func OppositeDirection(dir string) string {
	switch dir {
	case "north":
		return "south"
	case "south":
		return "north"
	case "east":
		return "west"
	case "west":
		return "east"
	case "up":
		return "down"
	case "down":
		return "up"
	default:
		return ""
	}
}

// ParseArguments parses a string and returns an array of arguments.
func ParseArguments(args []string) []string {
	var parsed []string

	var recording bool
	var recorded string
	for _, a := range args {
		start := a[0:1]
		end := a[len(a)-1:]
		if start == "\"" && end == "\"" {
			parsed = append(parsed, a[1:len(a)-1])
		} else if start == "\"" {
			recording = true
			recorded = recorded + a[1:]
		} else if end == "\"" {
			recording = false
			parsed = append(parsed, recorded+" "+a[:len(a)-1])
			recorded = ""
		} else if recording {
			recorded = recorded + " " + a
		} else {
			parsed = append(parsed, a)
		}
	}
	return parsed
}

// IsStringBool returns true if the string is a boolean value.
func IsStringBool(s string) bool {
	lc := strings.ToLower(s)
	if lc == "true" || lc == "false" {
		return true
	}
	return false
}

// ToggleStringBool toggles a string like a bool.
func ToggleStringBool(s string) string {
	lc := strings.ToLower(s)
	if lc == "true" {
		return "false"
	} else if lc == "false" {
		return "true"
	}
	return s
}

// BoolToWords returns a string depending on whether or not a bool is true.
func BoolToWords(b bool, true string, false string) string {
	if b {
		return true
	}

	return false
}

// IsUUID returns a bool indicating whether the input string is a parsable UUID.
func IsUUID(in string) bool {
	if _, err := uuid.Parse(in); err != nil {
		return false
	}

	return true
}
