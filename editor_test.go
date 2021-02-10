package owasp_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/gregoryv/owasp"
)

func TestEditor(t *testing.T) {
	var ed owasp.Editor

	filename := "OWASP_ISVS-1.0RC.json"
	if err := ed.ImportFile(filename); err != nil {
		t.Fatal(err)
	}

	if err := ed.SetVerified("1.3.1"); err != nil {
		t.Fatal(err)
	}

	if err := ed.SetVerified("no such"); err == nil {
		t.Fatal("SetVerified should fail")
	}

	var buf bytes.Buffer
	if err := ed.TidyExport(&buf); err != nil {
		t.Fatal(err)
	}

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
