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
	Entries []Entry
}

// SetApplicableByLevel sets the all entries with the given level
// as applicable
func (me *Editor) SetApplicableByLevel(level Level, appl bool) error {
	for i, e := range me.Entries {
		switch {
		case level == L1 && e.L1:
			me.Entries[i].Applicable = appl

		case level == L2 && e.L2:
			me.Entries[i].Applicable = appl

		case level == L3 && e.L3:
			me.Entries[i].Applicable = appl

		default:
			return fmt.Errorf("unmatched level %v", level)
		}
	}
	return nil
}

// SetVerified sets the given entry as verified and applicable
func (me *Editor) SetVerified(id string, v bool) error {
	for i, e := range me.Entries {
		if e.ID == id {
			me.Entries[i].Verified = v
			me.Entries[i].Applicable = true
			me.Entries[i].Manual = nil
			return nil
		}
	}
	return fmt.Errorf("id %s not found", id)
}

// ResetVerifiedBy resets the verified state of all entries matching pattern
func (me *Editor) ResetVerifiedBy(pattern string) error {
	for i, e := range me.Entries {
		if found, _ := regexp.MatchString(pattern, e.ID); found {
			me.Entries[i].Verified = false
			me.Entries[i].Manual = nil
		}
	}

	return nil
}

// SetManuallyVerified sets the given entry as verified and applicable
// with manual notes.
func (me *Editor) SetManuallyVerified(id string, v bool, man Manual) error {
	for i, e := range me.Entries {
		if e.ID == id {
			me.Entries[i].Verified = v
			me.Entries[i].Applicable = true
			me.Entries[i].Manual = &man
			return nil
		}
	}
	return fmt.Errorf("id %s not found", id)
}

// SetManuallyVerifiedBy sets all entries matching pattern as verified and applicable
// with manual notes.
func (me *Editor) SetManuallyVerifiedBy(pattern string, v bool, man Manual) error {
	for i, e := range me.Entries {
		if found, _ := regexp.MatchString(pattern, e.ID); found {
			me.Entries[i].Verified = v
			me.Entries[i].Applicable = true
			me.Entries[i].Manual = &man
		}
	}
	return nil
}

func (me *Editor) SetApplicableBy(pattern string, v bool) error {
	for i, e := range me.Entries {
		if found, _ := regexp.MatchString(pattern, e.ID); found {
			me.Entries[i].Applicable = v
		}
	}
	return nil
}

func (me *Editor) SetApplicable(id string, v bool) error {
	for i, e := range me.Entries {
		if e.ID == id {
			me.Entries[i].Applicable = v
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
	return json.NewDecoder(r).Decode(&me.Entries)
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
	return json.NewEncoder(w).Encode(me.Entries)
}

// NewReport returns a new report from the loaded entries.
func (me *Editor) NewReport(title string) *Report {
	r := &Report{
		Title:              title,
		ShortDescriptionNA: true,
	}
	r.AddEntries(me.Entries...)
	return r
}
