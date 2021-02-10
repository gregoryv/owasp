package edisvs_test

import (
	"bytes"
	"testing"

	"github.com/gregoryv/edisvs"
)

func Test_generate_isvs(t *testing.T) {
	var ed edisvs.Editor

	filename := "OWASP_ISVS-1.0RC.json"
	if err := ed.ImportOWASPFile(filename); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	if err := ed.Export(&buf); err != nil {
		t.Fatal(err)
	}

}
