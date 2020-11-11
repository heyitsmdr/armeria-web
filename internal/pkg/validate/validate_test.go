package validate

import "testing"

func check(results []ValidationResult, pass bool, t *testing.T) {
	for _, c := range results {
		if c.Result != pass {
			t.Errorf("validation returned unexpected result: %v %v", c.Checks, c.Errors)
		}
	}
}

func TestMulti(t *testing.T) {
	shouldPass := []ValidationResult{
		Check("5", "num|min:1|max:7"),
	}
	shouldFail := []ValidationResult{
		Check("10", "num|min:1|max:7"),
	}

	check(shouldPass, true, t)
	check(shouldFail, false, t)
}

func TestBool(t *testing.T) {
	shouldPass := []ValidationResult{
		Check("true", "bool"),
		Check("false", "bool"),
		Check("True", "bool"),
		Check("False", "bool"),
	}
	shouldFail := []ValidationResult{
		Check("test", "bool"),
	}

	check(shouldPass, true, t)
	check(shouldFail, false, t)
}

func TestMin(t *testing.T) {
	shouldPass := []ValidationResult{
		Check("2", "min:2"),
		Check("3", "min:2"),
	}
	shouldFail := []ValidationResult{
		Check("1", "min:2"),
		Check("abc", "min:2"),
	}

	check(shouldPass, true, t)
	check(shouldFail, false, t)
}

func TestMax(t *testing.T) {
	shouldPass := []ValidationResult{
		Check("4", "max:5"),
		Check("5", "max:5"),
	}
	shouldFail := []ValidationResult{
		Check("6", "max:4"),
		Check("abc", "max:4"),
	}

	check(shouldPass, true, t)
	check(shouldFail, false, t)
}

func TestIn(t *testing.T) {
	shouldPass := []ValidationResult{
		Check("cat", "in:cat,dog"),
	}
	shouldFail := []ValidationResult{
		Check("blue", "in:red,green"),
	}

	check(shouldPass, true, t)
	check(shouldFail, false, t)
}

func TestNum(t *testing.T) {
	shouldPass := []ValidationResult{
		Check("5", "num"),
	}
	shouldFail := []ValidationResult{
		Check("test", "num"),
		Check("5a", "num"),
	}

	check(shouldPass, true, t)
	check(shouldFail, false, t)
}
