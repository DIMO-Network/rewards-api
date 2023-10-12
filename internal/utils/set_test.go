package utils

import "testing"

func TestNewWithElements(t *testing.T) {
	s := NewSet("b", "d")

	if !s.Contains("b") {
		t.Errorf("set should contain \"b\", but we claim it doesn't")
	}

	if !s.Contains("d") {
		t.Errorf("set should not contain \"d\", but we claim it does")
	}
}

func TestAddContains(t *testing.T) {
	s := NewSet[string]()
	s.Add("a")
	s.Add("b")

	if !s.Contains("a") {
		t.Errorf("set should contain \"a\", but we claim it doesn't")
	}

	if s.Contains("c") {
		t.Errorf("set should not contain \"c\", but we claim it does")
	}
}

func TestRemove(t *testing.T) {
	s := NewSet[string]()
	s.Add("a")
	s.Add("b")

	s.Remove("a")

	if s.Contains("a") {
		t.Errorf("set should not contain the removed element \"a\", but it does")
	}
}

func TestSlice(t *testing.T) {
	s := NewSet[string]()
	s.Add("a")
	s.Add("b")

	v := s.Slice()

	if l := len(v); l != 2 {
		t.Errorf("slice length should be 2, but was %d", l)
	}
}

func TestLen(t *testing.T) {
	s := NewSet[string]()
	s.Add("a")
	s.Add("b")

	if l := s.Len(); l != 2 {
		t.Errorf("set should have length 2, but was %d", l)
	}
}
