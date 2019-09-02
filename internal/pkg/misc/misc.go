package misc

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

// SliceRemove removes the element at the specified index; this does not
// preserve ordering.
func SliceRemove(s interface{}, i int) interface{} {
	sa := s.([]*interface{})
	sa[i] = sa[len(sa)-1]
	return sa[:len(sa)-1]
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

func ToggleStringBool(s string) string {
	if s == "true" {
		return "false"
	}
	if s == "false" {
		return "true"
	}
	return s
}
