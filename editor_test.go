package owasp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestEditor_SetApplicable(t *testing.T) {
	cases := []struct {
		entryID string
		input   string
	}{
		{"1.1.1", "1.1.1"},
		{"1.1.1", `^1\.1\.1$`},
		{"1.1.1", `1.1.*`},
		{"1.1.1", `*.1.*`},
		{"1.1.1", `1.*.1`},
	}
	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			ed := NewEditor()
			ed.Entries = []Entry{
				{ID: c.entryID},
			}
			if err := ed.SetApplicable(c.input, true); err != nil {
				t.Error(err)
			}
			if !ed.Entries[0].Applicable {
				t.Error("did not set Applicable field")
			}
		})
	}
}

func ExampleEditor_SetApplicable() {
	ed := NewEditor()
	ed.Entries = []Entry{
		{ID: "1.1.1"},
		{ID: "1.2.1"},
		{ID: "1.2.2"},
		{ID: "2.1.1"},
		{ID: "3.1.1"},
		{L1: true, ID: "4.1.1"},
	}
	ed.SetApplicable("1.1.1", true) // specific
	ed.SetApplicable("1.2.*", true) // readable expression
	ed.SetApplicable(`^2.*`, true)  // raw regexp
	ed.SetApplicable(L1, true)      // by level
	for _, e := range ed.Entries {
		fmt.Println(e.ID, e.Applicable)
	}
	// output:
	// 1.1.1 true
	// 1.2.1 true
	// 1.2.2 true
	// 2.1.1 true
	// 3.1.1 false
	// 4.1.1 true
}

func TestEditor_SetApplicable_fails(t *testing.T) {
	ed := NewEditor()
	if err := ed.SetApplicable("1.1.1", true); err == nil {
		t.Error("did not fail for unknown id")
	}
}

// ----------------------------------------

func TestEditor_SetVerified(t *testing.T) {
	cases := []struct {
		entryID string
		input   string
	}{
		{"1.1.1", "1.1.1"},
		{"1.1.1", `^1\.1\.1$`},
		{"1.1.1", `1.1.*`},
		{"1.1.1", `*.1.*`},
		{"1.1.1", `1.*.1`},
	}
	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			ed := NewEditor()
			ed.Entries = []Entry{
				{ID: c.entryID, Applicable: true},
			}
			if err := ed.SetVerified(c.input, true); err != nil {
				t.Error(err)
			}
			if !ed.Entries[0].Verified {
				t.Error("did not set Verified field")
			}
		})
	}
}

func TestEditor_SetVerified_fails(t *testing.T) {
	ed := NewEditor()
	ed.Entries = []Entry{
		{ID: "1.1.1", Applicable: true},
		{ID: "1.1.2"},
	}
	t.Run("when no entries match", func(t *testing.T) {
		err := ed.SetVerified(`3.1.1`, true)
		if err == nil {
			t.Error("when no entries match")
		}
	})
	t.Run("when not applicable", func(t *testing.T) {
		err := ed.SetVerified(`1.1.2`, true)
		if err == nil {
			t.Fail()
		}
	})
}

// ----------------------------------------

func TestEditor_Reset(t *testing.T) {
	ed := NewEditor()
	ed.Entries = []Entry{
		{Verified: true, Manual: &Manual{}, Applicable: true},
	}
	ed.Reset()
	if ed.Entries[0].Verified {
		t.Error("Verified field not reset")
	}
	if ed.Entries[0].Manual != nil {
		t.Error("Manual field not reset")
	}
	if ed.Entries[0].Applicable {
		t.Error("Applicable field not reset")
	}
}

func TestEditor_ResetVerified(t *testing.T) {
	ed := NewEditor()
	ed.Entries = []Entry{
		{Verified: true, Manual: &Manual{}},
	}
	ed.ResetVerified()
	if ed.Entries[0].Verified {
		t.Error("Verified field not reset")
	}
	if ed.Entries[0].Manual != nil {
		t.Error("Manual field not reset")
	}
}

func TestEditor_ResetApplicable(t *testing.T) {
	ed := NewEditor()
	ed.Entries = []Entry{
		{Applicable: true},
	}
	ed.ResetApplicable()
	if ed.Entries[0].Applicable {
		t.Fail()
	}
}

// ----------------------------------------

func TestEditor_Open_fails(t *testing.T) {
	ed := NewEditor()
	err := ed.Load("no such file")
	if err == nil {
		t.Fail()
	}
}

// ----------------------------------------

func TestEditor_SaveAs_fails(t *testing.T) {
	ed := NewEditor()
	ed.Entries = []Entry{
		{ID: "1.1.1"},
	}

	// create write protected file
	filename, cleanup := tmpFile(0000)
	defer cleanup()

	if err := ed.SaveAs(filename); err == nil {
		t.Error("no error when saving to write protected file")
	}
}

func tmpFile(mode os.FileMode) (filename string, cleanup func()) {
	tmpfile, err := ioutil.TempFile("", "owasp")
	if err != nil {
		log.Fatal(err)
	}
	filename = tmpfile.Name()
	cleanup = func() {
		os.Chmod(filename, 0644)
		os.Remove(filename) // clean up
	}
	os.Chmod(filename, mode)
	return
}

// ----------------------------------------

func TestEditor_WriteTo(t *testing.T) {
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
	_, err := ed.WriteTo(&buf)
	if err != nil {
		t.Fatal(err)
	}

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
		if t != nil {
			t.Helper()
			t.Fatal(err)
			return
		}
		panic(err)
	}
}
