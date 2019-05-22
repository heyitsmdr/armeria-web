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

// DirectionOffsets returns the X/Y/Z offsets for each direction as a map
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
func SliceRemove(s []interface{}, i int) []interface{} {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
