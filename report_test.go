package owasp

import (
	"bytes"
	"strings"
	"testing"
)

func TestReport_WriteTo(t *testing.T) {
	report := Report{
		Title: "test report",
	}
	report.AddEntries(
		Entry{ID: "1.3.1", L1: true, Verified: true, Applicable: true,
			Manual: &Manual{
				How:  "Using hardware...",
				When: "2022-01-01",
				By:   "John Doe",
			},
		},
		Entry{ID: "4.3.4", L2: true},
		Entry{ID: "5.1.1", L3: true},
	)

	var buf bytes.Buffer
	report.WriteTo(&buf)
	got := buf.String()

	exp := []string{
		"4.3.4",
		"- 5",
		"[x] **1.3.1**",
		"Using",
		"John Doe",
		"2022-01-01",
	}
	for _, exp := range exp {
		if !strings.Contains(got, exp) {
			t.Log(got)
			t.Fatal("missing", exp)
		}
	}
}

func TestReport_SaveAs_fails(t *testing.T) {
	var report Report

	filename, cleanup := tmpFile(0000)
	defer cleanup()

	if err := report.SaveAs(filename); err == nil {
		t.Fail()
	}
}
