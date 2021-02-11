package owasp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
)

// https://github.com/OWASP/IoT-Security-Verification-Standard-ISVS

// MustSetVerifiedNow loads and sets the verified flag. Panics on errors.
func MustSetVerifiedNow(id, filename string, v bool) {
	ed := NewEditor()
	if err := ed.Load(filename); err != nil {
		panic(err)
	}
	if err := ed.SetVerified(id, v); err != nil {
		panic(err)
	}
	if err := ed.Save(filename); err != nil {
		panic(err)
	}
}

// NewEditor returns an empty editor. Use Load or Import methods to
// fill with entries.
func NewEditor() *Editor {
	return &Editor{}
}

//go:generate gentut -t Editor -p owasp -in editor.go -w
type Editor struct {
	entries []Entry
}

// SetVerified sets the given entry as verified and applicable
func (me *Editor) SetVerified(id string, v bool) error {
	for i, e := range me.entries {
		if e.ID == id {
			me.entries[i].Verified = v
			me.entries[i].Applicable = true
			return nil
		}
	}
	return fmt.Errorf("id %s not found", id)
}

func (me *Editor) SetApplicableBy(pattern string) error {
	var found bool
	for i, e := range me.entries {
		if found, _ = regexp.MatchString(pattern, e.ID); found {
			me.entries[i].Applicable = true
		}
	}
	if !found {
		return fmt.Errorf("%s no match", pattern)
	}
	return nil
}

func (me *Editor) SetApplicable(id string) error {
	for i, e := range me.entries {
		if e.ID == id {
			me.entries[i].Applicable = true
			return nil
		}
	}
	return fmt.Errorf("id %s not found", id)
}

// Load entries from given json file.
func (me *Editor) Load(filename string) error {
	fh, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fh.Close()
	return me.Import(fh)
}

// Save writes entries as a tidy json to the given filename.
func (me *Editor) Save(filename string) error {
	fh, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer fh.Close()
	return me.TidyExport(fh)
}

// Import entries from json
func (me *Editor) Import(r io.Reader) error {
	return json.NewDecoder(r).Decode(&me.entries)
}

// TidyExport exports entries as tidy json to the given writer.
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

// Export entries as json
func (me *Editor) Export(w io.Writer) error {
	return json.NewEncoder(w).Encode(me.entries)
}

// SaveReport saves entries as markdown to the given filename.
func (me *Editor) SaveReport(filename string, report Report) error {
	fh, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer fh.Close()
	return me.WriteReport(fh, report)
}

// WriteReport writes a markdown report
func (me *Editor) WriteReport(w io.Writer, report Report) error {
	report.entries = me.entries
	_, err := report.WriteTo(w)
	return err
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

func (me *Entry) checkbox() string {
	checkbox := "[ ]"
	if me.Verified {
		checkbox = "[x]"
	}
	return checkbox
}

func (me *Entry) shortString() string {
	return fmt.Sprintf("%s %s...", me.ID, me.shortDesc())
}

func (me *Entry) String() string {
	return fmt.Sprintf("%s %s", me.ID, me.Description)
}

func (me *Entry) shortDesc() string {
	if len(me.Description) < 80 {
		return me.Description
	}
	return me.Description[:80] + "..."
}
