package owasp

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"
)

func TestEditor_SetVerifiedBy(t *testing.T) {
	ed := NewEditor()
	ed.Entries = []Entry{
		{ID: "1.1.1", Applicable: true},
		{ID: "2.2.2", Applicable: true},
	}
	err := ed.SetVerifiedBy(`1.*`, true)
	if err != nil {
		t.Fatal(err)
	}
	if !ed.Entries[0].Verified {
		t.Error("Verified field not set")
	}
	if ed.Entries[1].Verified {
		t.Error("Verified field set on wrong entry")
	}
}

// ----------------------------------------

func TestEditor_SetManuallyVerifiedBy(t *testing.T) {
	ed := NewEditor()
	ed.Entries = []Entry{
		{ID: "1.1.1", Applicable: true},
		{ID: "2.2.2", Applicable: true},
	}
	err := ed.SetManuallyVerifiedBy(`1.*`, true, Manual{})
	if err != nil {
		t.Fatal(err)
	}
	if !ed.Entries[0].Verified {
		t.Error("Verified field not set")
	}
	if ed.Entries[0].Manual == nil {
		t.Error("Manual field not set")
	}
	if ed.Entries[1].Verified {
		t.Error("Verified field set on wrong entry")
	}
}

func TestEditor_SetManuallyVerifiedBy_fails(t *testing.T) {
	ed := NewEditor()
	ed.Entries = []Entry{
		{ID: "1.1.1", Applicable: true},
		{ID: "1.1.2", Applicable: true},
		{ID: "2.1.1"},
	}

	t.Run("when no entries match", func(t *testing.T) {
		err := ed.SetManuallyVerifiedBy(`3.*`, true, Manual{})
		if err == nil {
			t.Fail()
		}
	})

	t.Run("when not applicable", func(t *testing.T) {
		err := ed.SetManuallyVerifiedBy(`2.*`, true, Manual{})
		if err == nil {
			t.Fail()
		}
	})
}

// ----------------------------------------

func TestEditor_SetApplicableByLevel(t *testing.T) {
	ed := NewEditor()
	ed.Entries = []Entry{
		{L1: true},
		{L2: true},
		{L3: true},
	}
	if err := ed.SetApplicableByLevel(L1, true); err != nil {
		t.Error(err)
	}
	if err := ed.SetApplicableByLevel(L2, true); err != nil {
		t.Error(err)
	}
	if err := ed.SetApplicableByLevel(L3, true); err != nil {
		t.Error(err)
	}
	if err := ed.SetApplicableByLevel(0, true); err == nil {
		t.Error("did not fail for level 0")
	}
}

func TestEditor_SetApplicable(t *testing.T) {
	ed := NewEditor()
	ed.Entries = []Entry{
		{ID: "1.1.1"},
	}
	if err := ed.SetApplicable("1.1.1", true); err != nil {
		t.Error(err)
	}
	if !ed.Entries[0].Applicable {
		t.Error("did not set Applicable field")
	}
}

func TestEditor_SetApplicable_fails(t *testing.T) {
	ed := NewEditor()
	if err := ed.SetApplicable("1.1.1", true); err == nil {
		t.Error("did not fail for unknown id")
	}
}

func TestEditor_TidyExport(t *testing.T) {
	ed := NewEditor()
	ed.Entries = []Entry{
		{ID: "1.1.1", Applicable: true, Verified: true,
			Manual: &Manual{
				How:  "Using hardware...",
				When: "2022-01-01",
				By:   "John Doe",
			},
		},
		{ID: "1.1.2", Applicable: true},
		{ID: "2.1.1"},
	}

	var buf bytes.Buffer
	must(t, ed.TidyExport(&buf))

	var got []Entry
	json.NewDecoder(&buf).Decode(&got)
	if len(got) != 3 {
		t.Error(got)
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

	must(t, ed.SaveAs("checklist/asvs.json"))
}

func Test_convert_isvs(t *testing.T) {
	ed := NewEditor()

	must(t, ed.Load("testdata/OWASP_ISVS-1.0RC.json"))
	must(t, ed.SaveAs("checklist/isvs.json"))
}

// ----------------------------------------

func must(t *testing.T, err error) {
	if err != nil {
		t.Helper()
		t.Fatal(err)
	}
}
