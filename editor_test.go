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

	filename := "OWASP_ISVS-1.0RC.json"
	ed.mustLoad(filename)

	ed.shouldSetVerified("1.3.1", true)

	if err := ed.SetVerified("no such", true); err == nil {
		t.Fatal("SetVerified should fail")
	}

	var buf bytes.Buffer
	ed.mustTidyExport(&buf)

	report := NewReport("Report ISVS")
	var rbuf bytes.Buffer
	ed.WriteReport(&rbuf, *report)

	got := rbuf.String()
	exp := []string{
		"4.3.4",
		"[ ] 5",
		"[x] 1.3.1",
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
	_ = ed.Load("OWASP_ISVS-1.0RC.json")

	_ = ed.SetApplicableBy(`1\.1\.\d*`)
	_ = ed.SetVerified("1.1.1", true)
	_ = ed.SetVerified("1.3.1", false)

	ed.Save("isvs.json")
	report := NewReport("Report ISVS")
	ed.SaveReport("example_report.md", *report)
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
	fh, _ := os.Open("ASVS-4.0.2.json")
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
	ed.entries = entries

	report := NewReport("Report ASVS")
	ed.SaveReport("asvs_example.md", *report)

	ed.Save("asvs.json")

}
