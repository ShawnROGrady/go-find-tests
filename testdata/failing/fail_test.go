package failing

import "testing"

func TestSum(t *testing.T) {
	var (
		a = 1
		b = 2
		c = 3
	)
	if sum(a, b) != c {
		t.Errorf("Unexpected sum(%d, %d) (expected = %d, actual = %d)", a, b, c, sum(a, b))
	}
}
