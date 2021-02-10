package owasp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/gregoryv/nexus"
)

// https://github.com/OWASP/IoT-Security-Verification-Standard-ISVS

func NewEditor() *Editor {
	return &Editor{}
}

//go:generate gentut -t Editor -p owasp -in editor.go -w
type Editor struct {
	entries []Entry
}

// SetVerified sets the given entry as verified and applicable
func (me *Editor) SetVerified(id string) error {
	for i, e := range me.entries {
		if e.ID == id {
			me.entries[i].Verified = true
			me.entries[i].Applicable = true
			return nil
		}
	}
	return fmt.Errorf("id %s not found", id)
}

func (me *Editor) ImportFile(filename string) error {
	fh, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fh.Close()
	return me.Import(fh)
}

func (me *Editor) Import(r io.Reader) error {
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
	p.Println("https://github.com/OWASP/IoT-Security-Verification-Standard-ISVS/")

	p.Println("## Applicable")

	for _, e := range me.entries {
		if !e.Applicable {
			continue
		}
		p.Println(e.String())
	}

	p.Println()
	p.Println("## Not Applicable")
	for _, e := range me.entries {
		if e.Applicable {
			continue
		}
		p.Println(e.String())
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
	Applicable  bool
}

func (me *Entry) String() string {
	checkbox := "- [ ]"
	if me.Verified {
		checkbox = "- [x]"
	}
	return fmt.Sprintf("%s [%s](%s) %s...", checkbox, me.ID, me.link(), me.shortDesc())
}

func (me *Entry) link() string {
	switch me.ID[:3] {
	case "1.1":
		return base + "/blob/master/en/V1-IoT_Ecosystem_Requirements.md#application-and-ecosystem-design"
		// todo add the other cases
	default:
		return base
	}
}

func (me *Entry) shortDesc() string {
	if len(me.Description) < 80 {
		return me.Description
	}
	return me.Description[:80] + "..."
}

const base = "https://github.com/OWASP/IoT-Security-Verification-Standard-ISVS"
