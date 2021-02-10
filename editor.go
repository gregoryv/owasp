package owasp

import (
	"bytes"
	"encoding/json"
	"io"
	"os"

	"github.com/gregoryv/nexus"
)

// https://github.com/OWASP/IoT-Security-Verification-Standard-ISVS

type Editor struct {
	entries []Entry
}

func (me *Editor) ImportFile(filename string) error {
	fh, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fh.Close()
	return me.ImportOWASP(fh)
}

func (me *Editor) ImportOWASP(r io.Reader) error {
	return json.NewDecoder(r).Decode(&me.entries)
}

func (me *Editor) TidyExport(w io.Writer) error {
	var buf bytes.Buffer
	me.Export(&buf)

	var tidy bytes.Buffer
	err := json.Indent(&tidy, buf.Bytes(), "", "  ")
	if err != nil {
		return err
	}

	_, err = w.Write(tidy.Bytes())
	return err
}

func (me *Editor) Export(w io.Writer) error {
	return json.NewEncoder(w).Encode(me.entries)
}

// WriteReport writes a markdown report
func (me *Editor) WriteReport(w io.Writer) error {
	p, err := nexus.NewPrinter(w)
	p.Println("# ISVS Report")

	for _, e := range me.entries {
		checkbox := "[ ]"
		if e.Verified {
			checkbox = "[x]"
		}
		p.Println(checkbox, e.ID)
	}

	return *err
}

type Entry struct {
	L1          bool
	L2          bool
	L3          bool
	Description string
	ID          string
	Verified    bool
}
