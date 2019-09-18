package size

import "testing"

// few different tests with various amounts of coverage
// some copied from https://blog.golang.org/cover
type Test struct {
	in  int
	out string
}

var tests = []Test{
	{-1, "negative"},
	{5, "small"},
}

func TestSize(t *testing.T) {
	for i, test := range tests {
		size := size(test.in)
		if size != test.out {
			t.Errorf("#%d: size(%d)=%s; want %s", i, test.in, size, test.out)
		}
	}
}

func TestNegativeSize(t *testing.T) {
	for i := -10; i < 0; i++ {
		size := size(i)
		if size != "negative" {
			t.Errorf("size(%d)=%s; want %s", i, size, "negative")
		}
	}
}

func TestIsNegative(t *testing.T) {
	for i := -10; i < 0; i++ {
		if !isNegative(i) {
			t.Errorf("isNegative(%d) unexpectedly false", i)
		}
	}
}

func TestIsEnormous(t *testing.T) {
	if !isEnormous(1001) {
		t.Errorf("isEnormous(%d) unexpectedly false", 1001)
	}
}
