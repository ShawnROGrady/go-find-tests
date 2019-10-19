package subtests

// simple function to check test finder
// will be adding sub tests to confirm correct behaviour
func length(s string) string {
	a := len(s)
	switch {
	case a == 0:
		return "empty"
	case a < 10:
		return "short"
	case a < 100:
		return "long"
	case a < 1000:
		return "very long"
	}
	return "a novel"
}

func isEmpty(s string) bool {
	return length(s) == "empty"
}

func isShort(s string) bool {
	return length(s) == "short"
}
