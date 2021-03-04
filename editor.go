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

type Editor struct {
	Entries []Entry
}

// SetApplicableByLevel sets applicable value of entries by specific
// level.
func (me *Editor) SetApplicableByLevel(level Level, v bool) error {
	for i, e := range me.Entries {
		if e.IsLevel(level) {
			me.Entries[i].Applicable = v
		}
	}
	return nil
}

// ResetVerifiedBy same as SetVerifiedBy(pattern, false)
func (me *Editor) ResetVerifiedBy(pattern string) error {
	return me.SetVerifiedBy(pattern, false)
}

// SetVerified sets the given entry as verified. Returns error if id
// is not found or the entry is not applicable.
func (me *Editor) SetVerified(id string, v bool) error {
	for i, e := range me.Entries {
		if e.ID == id {
			if !e.Applicable {
				return fmt.Errorf("%v is not applicable", e.ID)
			}
			me.Entries[i].Verified = v
			me.Entries[i].Manual = nil
			return nil
		}
	}
	return fmt.Errorf("id %s not found", id)
}

// SetVerifiedBy sets the verified state of all entries where id
// matches the pattern. Returns error if matching entry is not
// applicable.
func (me *Editor) SetVerifiedBy(pattern string, v bool) error {
	for i, e := range me.Entries {
		if found, _ := regexp.MatchString(pattern, e.ID); found {
			if !e.Applicable {
				return fmt.Errorf("%v is not applicable", e.ID)
			}
			me.Entries[i].Verified = v
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
			if !e.Applicable {
				return fmt.Errorf("%v is not applicable", e.ID)
			}
			me.Entries[i].Verified = v
			me.Entries[i].Applicable = true
			me.Entries[i].Manual = &man
			return nil
		}
	}
	return fmt.Errorf("id %s not found", id)
}

// SetManuallyVerifiedBy sets all entries matching pattern as verified
// and applicable with manual notes.
func (me *Editor) SetManuallyVerifiedBy(pattern string, v bool, man Manual) error {
	for i, e := range me.Entries {
		if found, _ := regexp.MatchString(pattern, e.ID); found {
			if !e.Applicable {
				return fmt.Errorf("%v is not applicable", e.ID)
			}
			me.Entries[i].Verified = v
			me.Entries[i].Manual = &man
		}
	}
	return nil
}

func (me *Editor) SetApplicableBy(pattern string, v bool) error {
	if err := doesMatch(pattern, me.Entries); err != nil {
		return fmt.Errorf("SetApplicableBy: %w", err)
	}
	for i, e := range me.Entries {
		if found, _ := regexp.MatchString(pattern, e.ID); found {
			me.Entries[i].Applicable = v
		}
	}
	return nil
}

func doesMatch(pattern string, entries []Entry) error {
	for i := range entries {
		if found, _ := regexp.MatchString(pattern, entries[i].ID); found {
			return nil
		}
	}
	return fmt.Errorf("pattern %q does not match any entries", pattern)
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
