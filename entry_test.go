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
