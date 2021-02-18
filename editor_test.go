package owasp

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestEditor(t *testing.T) {
	ed := NewEditor().UnderTest(t)

	filename := "testdata/OWASP_ISVS-1.0RC.json"
	ed.mustLoad(filename)

	ed.shouldSetVerified("1.3.1", true)
	man := Manual{
		How:  "Using hardware...",
		When: "2022-01-01",
		By:   "John Doe",
	}
	if err := ed.SetManuallyVerified("2.1.1", true, man); err != nil {
		t.Fatal(err)
	}

	if err := ed.SetVerified("no such", true); err == nil {
		t.Fatal("SetVerified should fail")
	}

	var buf bytes.Buffer
	ed.mustTidyExport(&buf)

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

func ExampleEditor_WriteReport() {
	ed := NewEditor()
	_ = ed.Load("testdata/asvsx.json")

	_ = ed.SetApplicableBy(`1\.1\.\d*`, true)
	_ = ed.SetVerified("1.1.1", true)
	_ = ed.SetVerified("1.3.1", false)

	man := Manual{
		How:  "Latest threatmodel design change was updated on ...",
		When: "2021-02-18",
		By:   "John Doe",
	}
	_ = ed.SetManuallyVerified("1.1.2", true, man)

	ed.Save("testdata/asvsx.json")
	ed.NewReport("Report ASVS").Save("example_report.md")

	// output:
}

func ExampleMustSetVerifiedNow() {
	MustSetVerifiedNow("1.1.1", "testdata/asvsx.json", false)
	// output:
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
	ed := NewEditor().UnderTest(t)
	ed.entries = entries

	ed.mustSave("checklist/asvs.json")
}

func Test_convert_isvs(t *testing.T) {
	ed := NewEditor().UnderTest(t)

	ed.mustLoad("testdata/OWASP_ISVS-1.0RC.json")
	ed.mustSave("checklist/isvs.json")
}
