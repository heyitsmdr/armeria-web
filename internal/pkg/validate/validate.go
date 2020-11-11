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
	Bool string = "bool"
	Min         = "min"
	Max         = "max"
	In          = "in"
	Num         = "num"
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
		return "input not convertible to an int"
	}

	m, err := strconv.Atoi(min)
	if err != nil {
		return "min value not convertible to an int"
	}

	if i < m {
		return fmt.Sprintf("lower than %d", m)
	}

	return ""
}

func checkMax(str, max string) string {
	i, err := strconv.Atoi(str)
	if err != nil {
		return "input not convertible to an int"
	}

	m, err := strconv.Atoi(max)
	if err != nil {
		return "max value not convertible to an int"
	}

	if i > m {
		return fmt.Sprintf("higher than %d", m)
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
		return "not a valid option"
	}

	return ""
}

func checkNum(str string) string {
	_, err := strconv.Atoi(str)
	if err != nil {
		return "input not convertible to an int"
	}

	return ""
}
