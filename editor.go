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
	switch level {
	case L1, L2, L3:
	default:
		return fmt.Errorf("no such level %v", level)
	}
	for i, e := range me.Entries {
		if e.IsLevel(level) {
			me.Entries[i].Applicable = v
		}
	}
	return nil
}

// Reset same as calling ResetVerified and ResetApplicable
func (me *Editor) Reset() {
	me.ResetVerified()
	me.ResetApplicable()
}

// ResetVerified sets Verified field to false on all entries and Manual to nil.
func (me *Editor) ResetVerified() {
	for i := range me.Entries {
		me.Entries[i].Verified = false
		me.Entries[i].Manual = nil
	}
}

// ResetApplicable sets applicable field to false on all entries.
func (me *Editor) ResetApplicable() {
	for i := range me.Entries {
		me.Entries[i].Applicable = false
	}
}

// ----------------------------------------

// SetApplicable sets the applicable field of given entry. Returns
// error if no entry is found.
func (me *Editor) SetApplicable(id string, v bool) error {
	for i, e := range me.Entries {
		if e.ID == id {
			me.Entries[i].Applicable = v
			return nil
		}
	}
	return fmt.Errorf("id %s not found", id)
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

// SetManuallyVerified sets the given entry as verified and applicable
// with manual notes.
func (me *Editor) SetManuallyVerified(id string, v bool, man Manual) error {
	for i, e := range me.Entries {
		if e.ID == id {
			if !e.Applicable {
				return fmt.Errorf("%v is not applicable", e.ID)
			}
			me.Entries[i].Verified = v
			me.Entries[i].Applicable = true // todo remove
			me.Entries[i].Manual = &man
			return nil
		}
	}
	return fmt.Errorf("id %s not found", id)
}

// SetVerifiedBy sets the verified state of all entries where id
// matches the pattern. Returns error if matching entry is not
// applicable.
func (me *Editor) SetVerifiedBy(pattern string, v bool) error {
	if err := doesMatch(pattern, me.Entries); err != nil {
		return fmt.Errorf("SetVerifiedBy: %w", err)
	}
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

// SetManuallyVerifiedBy sets all entries matching pattern as verified
// and applicable with manual notes.
func (me *Editor) SetManuallyVerifiedBy(pattern string, v bool, man Manual) error {
	if err := doesMatch(pattern, me.Entries); err != nil {
		return fmt.Errorf("SetApplicableBy: %w", err)
	}
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

// ----------------------------------------

// Load entries from given json file.
func (me *Editor) Load(filename string) error {
	fh, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fh.Close()
	return me.Import(fh)
}

// SaveAs writes entries as a tidy json to the given filename.
func (me *Editor) SaveAs(filename string) error {
	fh, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer fh.Close()
	_, err = me.WriteTo(fh)
	return err
}

// Import entries from json
func (me *Editor) Import(r io.Reader) error {
	return json.NewDecoder(r).Decode(&me.Entries)
}

// WriteTo exports entries as tidy json to the given writer.
func (me *Editor) WriteTo(w io.Writer) (int64, error) {
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(me.Entries)

	var tidy bytes.Buffer
	_ = json.Indent(&tidy, buf.Bytes(), "", "  ")

	return io.Copy(w, &tidy)
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
