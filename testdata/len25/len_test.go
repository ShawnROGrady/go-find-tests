package length25

import (
	"strings"
	"testing"
)

func TestEmptyStringIsEmpty(t *testing.T) {
	s := ""
	if !isEmpty(s) {
		t.Errorf("isEmpty('%s') unexpectedly false", s)
	}
}

func TestShortStringIsEmpty(t *testing.T) {
	s := "hello"
	if isEmpty(s) {
		t.Errorf("isEmpty('%s') unexpectedly true", s)
	}
}

func TestLongStringIsEmpty(t *testing.T) {
	s := "hello world!"
	if isEmpty(s) {
		t.Errorf("isEmpty('%s') unexpectedly true", s)
	}
}

func TestVeryLongStringIsEmpty(t *testing.T) {
	s := strings.Repeat("a", 500)
	if isEmpty(s) {
		t.Errorf("isEmpty(s) unexpectedly true with len(s) = %d", len(s))
	}
}

func TestNovelIsEmpty(t *testing.T) {
	s := strings.Repeat("a", 5000)
	if isEmpty(s) {
		t.Errorf("isEmpty(s) unexpectedly true with len(s) = %d", len(s))
	}
}

func TestEmptyStringIsShort(t *testing.T) {
	s := ""
	if isShort(s) {
		t.Errorf("isShort('%s') unexpectedly true", s)
	}
}

func TestShortStringIsShort(t *testing.T) {
	s := "hello"
	if !isShort(s) {
		t.Errorf("isShort('%s') unexpectedly false", s)
	}
}

func TestLongStringIsShort(t *testing.T) {
	s := "hello world!"
	if isShort(s) {
		t.Errorf("isShort('%s') unexpectedly true", s)
	}
}

func TestVeryLongStringIsShort(t *testing.T) {
	s := strings.Repeat("a", 500)
	if isShort(s) {
		t.Errorf("isShort(s) unexpectedly true with len(s) = %d", len(s))
	}
}

func TestNovelIsShort(t *testing.T) {
	s := strings.Repeat("a", 5000)
	if isShort(s) {
		t.Errorf("isShort(s) unexpectedly true with len(s) = %d", len(s))
	}
}

func TestEmptyStringIsLong(t *testing.T) {
	s := ""
	if isLong(s) {
		t.Errorf("isLong('%s') unexpectedly true", s)
	}
}

func TestShortStringIsLong(t *testing.T) {
	s := "hello"
	if isLong(s) {
		t.Errorf("isLong('%s') unexpectedly true", s)
	}
}

func TestLongStringIsLong(t *testing.T) {
	s := "hello world!"
	if !isLong(s) {
		t.Errorf("isLong('%s') unexpectedly false", s)
	}
}

func TestVeryLongStringIsLong(t *testing.T) {
	s := strings.Repeat("a", 500)
	if isLong(s) {
		t.Errorf("isLong(s) unexpectedly true with len(s) = %d", len(s))
	}
}

func TestNovelIsLong(t *testing.T) {
	s := strings.Repeat("a", 5000)
	if isLong(s) {
		t.Errorf("isLong(s) unexpectedly true with len(s) = %d", len(s))
	}
}

func TestEmptyStringIsVeryLong(t *testing.T) {
	s := ""
	if isVeryLong(s) {
		t.Errorf("isVeryLong('%s') unexpectedly true", s)
	}
}

func TestShortStringIsVeryLong(t *testing.T) {
	s := "hello"
	if isVeryLong(s) {
		t.Errorf("isVeryLong('%s') unexpectedly true", s)
	}
}

func TestLongStringIsVeryLong(t *testing.T) {
	s := "hello world!"
	if isVeryLong(s) {
		t.Errorf("isVeryLong('%s') unexpectedly true", s)
	}
}

func TestVeryLongStringIsVeryLong(t *testing.T) {
	s := strings.Repeat("a", 500)
	if !isVeryLong(s) {
		t.Errorf("isVeryLong(s) unexpectedly false with len(s) = %d", len(s))
	}
}

func TestNovelIsVeryLong(t *testing.T) {
	s := strings.Repeat("a", 5000)
	if isVeryLong(s) {
		t.Errorf("isVeryLong(s) unexpectedly true with len(s) = %d", len(s))
	}
}

func TestEmptyStringIsNovel(t *testing.T) {
	s := ""
	if isNovel(s) {
		t.Errorf("isNovel('%s') unexpectedly true", s)
	}
}

func TestShortStringIsNovel(t *testing.T) {
	s := "hello"
	if isNovel(s) {
		t.Errorf("isNovel('%s') unexpectedly true", s)
	}
}

func TestLongStringIsNovel(t *testing.T) {
	s := "hello world!"
	if isNovel(s) {
		t.Errorf("isNovel('%s') unexpectedly true", s)
	}
}

func TestVeryLongStringIsNovel(t *testing.T) {
	s := strings.Repeat("a", 500)
	if isNovel(s) {
		t.Errorf("isNovel(s) unexpectedly true with len(s) = %d", len(s))
	}
}

func TestNovelIsNovel(t *testing.T) {
	s := strings.Repeat("a", 5000)
	if !isNovel(s) {
		t.Errorf("isNovel(s) unexpectedly false with len(s) = %d", len(s))
	}
}
