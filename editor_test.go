package owasp

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestEditor_SetApplicableBy_fails_on_no_match(t *testing.T) {
	ed := NewEditor()
	err := ed.SetApplicableBy(".*", true)
	if err == nil {
		t.Fail()
	}
}

func TestEditor_set_non_applicable(t *testing.T) {
	ed := NewEditor()
	filename := "testdata/OWASP_ISVS-1.0RC.json"
	if err := ed.Load(filename); err != nil {
		t.Fatal(err)
	}

	if err := ed.SetVerified("1.2.1", true); err == nil {
		t.Fatal("SetVerified should fail")
	}
	if err := ed.SetVerifiedBy("1.2.*", true); err == nil {
		t.Fatal("SetVerifiedBy should fail")
	}
	if err := ed.SetManuallyVerified("1.2.1", true, Manual{}); err == nil {
		t.Fatal("SetManuallyVerified should fail")
	}
}

func TestEditor_SetApplicableByLevel(t *testing.T) {
	ed := NewEditor()
	filename := "checklist/asvs.json"
	if err := ed.Load(filename); err != nil {
		t.Fatal(err)
	}

	if err := ed.SetApplicableByLevel(L1, true); err != nil {
		t.Error(err)
	}
	if err := ed.SetApplicableByLevel(L2, true); err != nil {
		t.Error(err)
	}
	if err := ed.SetApplicableByLevel(L3, false); err != nil {
		t.Error(err)
	}
}

func TestEditor(t *testing.T) {
	ed := NewEditor()
	must(t, ed.Load("testdata/OWASP_ISVS-1.0RC.json"))

	must(t, ed.SetApplicable("1.3.1", true))
	must(t, ed.SetVerified("1.3.1", true))

	man := Manual{
		How:  "Using hardware...",
		When: "2022-01-01",
		By:   "John Doe",
	}
	must(t, ed.SetApplicable("2.1.1", true))
	if err := ed.SetManuallyVerified("2.1.1", true, man); err != nil {
		t.Fatal(err)
	}

	if err := ed.SetVerified("no such", true); err == nil {
		t.Fatal("SetVerified should fail")
	}

	var buf bytes.Buffer
	must(t, ed.TidyExport(&buf))

	var rbuf bytes.Buffer
	ed.NewReport("Report ISVS").WriteTo(&rbuf)

	got := rbuf.String()
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

// ----------------------------------------

func Test_convert_original_asvs_to_checklist(t *testing.T) {
	var asvs struct {
		Requirements []struct {
			Items []struct {
				Items []struct {
					ShortCode   string
					Description string
					L1          struct {
						Required bool
					}
					L2 struct {
						Required bool
					}
					L3 struct {
						Required bool
					}
				}
			}
		}
	}

	// Load original
	fh, _ := os.Open("testdata/ASVS-4.0.2.json")
	defer fh.Close()
	json.NewDecoder(fh).Decode(&asvs)

	// convert to entries
	entries := make([]Entry, 0)
	for _, req := range asvs.Requirements {
		for _, item := range req.Items {
			for _, item := range item.Items {
				e := Entry{
					ID:          item.ShortCode[1:],
					Description: item.Description,
				}
				if item.L1.Required {
					e.L1 = true
				}
				if item.L2.Required {
					e.L2 = true
				}
				if item.L3.Required {
					e.L3 = true
				}
				entries = append(entries, e)
			}
		}
	}
	ed := NewEditor()
	ed.Entries = entries

	must(t, ed.Save("checklist/asvs.json"))
}

func Test_convert_isvs(t *testing.T) {
	ed := NewEditor()

	must(t, ed.Load("testdata/OWASP_ISVS-1.0RC.json"))
	must(t, ed.Save("checklist/isvs.json"))
}

// ----------------------------------------

func must(t *testing.T, err error) {
	if err != nil {
		t.Helper()
		t.Fatal(err)
	}
}
