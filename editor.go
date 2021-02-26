package owasp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
)

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
			me.entries[i].Manual = nil
			return nil
		}
	}
	return fmt.Errorf("id %s not found", id)
}

// ResetVerifiedBy resets the verified state of all entries matching pattern
func (me *Editor) ResetVerifiedBy(pattern string) error {
	var found bool
	for i, e := range me.entries {
		if found, _ = regexp.MatchString(pattern, e.ID); found {
			me.entries[i].Verified = false
			me.entries[i].Manual = nil
		}
	}

	if !found {
		return fmt.Errorf("%s no match", pattern)
	}
	return nil
}

// SetManuallyVerified sets the given entry as verified and applicable
// with manual notes.
func (me *Editor) SetManuallyVerified(id string, v bool, man Manual) error {
	for i, e := range me.entries {
		if e.ID == id {
			me.entries[i].Verified = v
			me.entries[i].Applicable = true
			me.entries[i].Manual = &man
			return nil
		}
	}
	return fmt.Errorf("id %s not found", id)
}

func (me *Editor) SetApplicableBy(pattern string, v bool) error {
	var found bool
	for i, e := range me.entries {
		if found, _ = regexp.MatchString(pattern, e.ID); found {
			me.entries[i].Applicable = v
		}
	}
	if !found {
		return fmt.Errorf("%s no match", pattern)
	}
	return nil
}

func (me *Editor) SetApplicable(id string, v bool) error {
	for i, e := range me.entries {
		if e.ID == id {
			me.entries[i].Applicable = v
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

// NewReport returns a new report from the loaded entries.
func (me *Editor) NewReport(title string) *Report {
	r := &Report{
		Title:              title,
		ShortDescriptionNA: true,
	}
	r.AddEntries(me.entries...)
	return r
}
