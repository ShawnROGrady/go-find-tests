package length15

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
