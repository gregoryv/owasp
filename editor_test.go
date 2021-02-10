package edisvs_test

import (
	"bytes"
	"testing"

	"github.com/gregoryv/edisvs"
)

func TestEditor(t *testing.T) {
	var ed edisvs.Editor

	filename := "OWASP_ISVS-1.0RC.json"
	if err := ed.ImportOWASPFile(filename); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	if err := ed.TidyExport(&buf); err != nil {
		t.Fatal(err)
	}

	page := ed.Report()
	var report bytes.Buffer
	page.WriteTo(&report)
	t.Error(report.String())
}
