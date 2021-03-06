package owasp

import "testing"

func TestEntry_IsLevel(t *testing.T) {
	entry := Entry{
		L1: true,
		L2: true,
		L3: true,
	}
	if !entry.IsLevel(L1) {
		t.Error("false, expected true")
	}
	if !entry.IsLevel(L2) {
		t.Error("false, expected true")
	}
	if !entry.IsLevel(L3) {
		t.Error("false, expected L3 true")
	}
}

func TestEntry_String(t *testing.T) {

	cases := []struct {
		Entry
		exp string
	}{
		{Entry{ID: "1.1.1", L1: true}, "L1 1.1.1"},
		{Entry{ID: "1.1.1", L2: true}, "L2 1.1.1"},
		{Entry{ID: "1.1.1", L3: true}, "L3 1.1.1"},
		{Entry{ID: "1.1.1", L1: true, L3: true}, "L3 1.1.1"},
	}

	for _, c := range cases {
		got := c.Entry.String()
		if got != c.exp {
			t.Error(got, "expected", c.exp)
		}
	}

	defer func() {
		e := recover()
		if e == nil {
			t.Error("did not panic")
		}
	}()
	Entry{}.String()
}
