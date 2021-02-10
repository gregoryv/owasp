package owasp

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestEditor(t *testing.T) {
	ed := NewEditor().UnderTest(t)

	filename := "OWASP_ISVS-1.0RC.json"
	ed.mustImportFile(filename)

	ed.shouldSetVerified("1.3.1")

	if err := ed.SetVerified("no such"); err == nil {
		t.Fatal("SetVerified should fail")
	}

	var buf bytes.Buffer
	ed.mustTidyExport(&buf)

	var report bytes.Buffer
	ed.WriteReport(&report)

	got := report.String()
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
	_ = ed.SetVerified("1.1.1")
	_ = ed.SetVerified("1.3.1")

	fh, err := os.Create("example_report.md")
	if err != nil {
		panic(err)
	}
	defer fh.Close()
	ed.WriteReport(fh)
	// output:
}
