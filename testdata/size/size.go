package size

// just a simple function to check test-finder
// copied from https://blog.golang.org/cover
func size(a int) string {
	switch {
	case a < 0:
		return "negative"
	case a == 0:
		return "zero"
	case a < 10:
		return "small"
	case a < 100:
		return "big"
	case a < 1000:
		return "huge"
	}
	return "enormous"
}

func isEnormous(a int) bool {
	return size(a) == "enormous"
}

func isNegative(a int) bool {
	return size(a) == "negative"
}
