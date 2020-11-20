package validate

import (
	"fmt"
	"strconv"
	"strings"
)

type ValidationResult struct {
	Result bool
	Checks map[string]bool
	Errors map[string]string
}

const (
	// Bool checks if the input string is a bool value (ie: true or false).
	Bool string = "bool"
	// Min checks if the input string is >= the min value (inclusive).
	Min = "min"
	// Max checks if the input string is <= the max value (inclusive).
	Max = "max"
	// Min checks if the input string is in the provided set.
	In = "in"
	// Num checks if the input string is a number.
	Num = "num"
	// Empty checks if the input string is empty.
	Empty = "empty"
)

func Check(str, validatorString string) ValidationResult {
	result := ValidationResult{Result: true}
	result.Checks = make(map[string]bool)
	result.Errors = make(map[string]string)

	validations := strings.Split(validatorString, "|")
	for _, validation := range validations {
		sections := strings.Split(validation, ":")
		switch sections[0] {
		case Bool:
			validate(&result, Bool, checkBool(str))
		case Min:
			validate(&result, Min, checkMin(str, sections[1]))
		case Max:
			validate(&result, Max, checkMax(str, sections[1]))
		case In:
			validate(&result, In, checkIn(str, sections[1]))
		case Num:
			validate(&result, Num, checkNum(str))
		case Empty:
			validate(&result, Empty, checkEmpty(str))
		}
	}

	for _, pass := range result.Checks {
		if !pass {
			result.Result = false
			break
		}
	}

	return result
}

func validate(result *ValidationResult, cn, r string) {
	passed := len(r) == 0
	result.Checks[cn] = passed
	result.Errors[cn] = r
}

func checkBool(str string) string {
	switch strings.ToLower(str) {
	case "true":
		return ""
	case "false":
		return ""
	}
	return "not a boolean value"
}

func checkMin(str, min string) string {
	i, err := strconv.Atoi(str)
	if err != nil {
		return fmt.Sprintf("less than %s", min)
	}

	m, err := strconv.Atoi(min)
	if err != nil {
		return "min value not an int"
	}

	if i < m {
		return fmt.Sprintf("less than %d", m)
	}

	return ""
}

func checkMax(str, max string) string {
	i, err := strconv.Atoi(str)
	if err != nil {
		return fmt.Sprintf("greater than %s", max)
	}

	m, err := strconv.Atoi(max)
	if err != nil {
		return "max value not an int"
	}

	if i > m {
		return fmt.Sprintf("greater than %d", m)
	}

	return ""
}

func checkIn(str, in string) string {
	inSlice := strings.Split(in, ",")

	found := false
	for _, s := range inSlice {
		if str == s {
			found = true
			break
		}
	}

	if !found {
		return fmt.Sprintf("not in set %v", inSlice)
	}

	return ""
}

func checkNum(str string) string {
	_, err := strconv.Atoi(str)
	if err != nil {
		return "not an int"
	}

	return ""
}

func checkEmpty(str string) string {
	if len(str) > 0 {
		return "not empty"
	}

	return ""
}

// OnlyErrors returns a slice of only the error strings. Useful for sending back to the client when needed.
func (vr ValidationResult) OnlyErrors() []string {
	errors := make([]string, 0)
	for name, _ := range vr.Checks {
		if !vr.Checks[name] {
			errors = append(errors, vr.Errors[name])
		}
	}

	return errors
}

// String returns the error(s) as a one-line comma-separated string. This can be called
// implicitly by simply concatenating the ValidationResult to a string.
func (vr ValidationResult) String() string {
	return strings.Join(vr.OnlyErrors(), ", ")
}
