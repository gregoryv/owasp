package owasp

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestReport_ASVS_check(t *testing.T) {
	ed := NewEditor()
	ed.Load("checklist/asvs.json")
	ed.SetApplicable(L3, true)

	report := ed.NewReport("")

	var buf bytes.Buffer
	report.WriteTo(&buf)
	got := buf.String()

	exp := []string{
		fmt.Sprintf("L1: 0 verified of %v", 131),
		fmt.Sprintf("L2: 0 verified of %v", 267),
		fmt.Sprintf("L3: 0 verified of %v", 286),
	}
	for _, exp := range exp {
		if !strings.Contains(got, exp) {
			t.Log(got[:125])
			t.Fatal("missing", exp)
		}
	}
}

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

func TestReport_WriteTo_all_verified(t *testing.T) {
	report := Report{
		Title:              "full report",
		ShortDescriptionNA: true,
		entries: []Entry{
			Entry{
				ID: "1.1.1", L2: true, Applicable: true, Verified: true,
			},
		},
	}

	var buf bytes.Buffer
	report.WriteTo(&buf)
	got := buf.String()

	exp := "All requirements"
	if !strings.Contains(got, exp) {
		t.Log(got)
		t.Fatal("missing", exp)
	}
}

func TestReport_WriteTo_short_description(t *testing.T) {
	exp := "out of range"
	report := Report{
		Title:              "short report",
		ShortDescriptionNA: true,
		entries: []Entry{
			Entry{
				ID: "1.1.1", L2: true, Applicable: true, Verified: true,
				Description: strings.Repeat("x", 80) + "more here",
			},
			Entry{
				ID: "1.1.2", L2: true,
				Description: strings.Repeat("x", 80) + exp,
			},
			Entry{
				ID: "1.1.3", L2: true,
				Description: strings.Repeat("l", 79),
			},
		},
	}

	var buf bytes.Buffer
	report.WriteTo(&buf)
	got := buf.String()

	if strings.Contains(got, exp) {
		t.Log(got)
		t.Fatal(got)
	}
}

// ----------------------------------------

func TestReport_SaveAs_fails(t *testing.T) {
	var report Report

	filename, cleanup := tmpFile(0000)
	defer cleanup()

	if err := report.SaveAs(filename); err == nil {
		t.Fail()
	}
}
