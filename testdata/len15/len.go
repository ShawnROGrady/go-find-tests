package length15

// another simple function to check test finder
// the main goal is to create a pkg w/ a medium number of tests (15)
// to measure performance as the number of tests grow
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

func isLong(s string) bool {
	return length(s) == "long"
}
