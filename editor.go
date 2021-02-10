package owasp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"

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

func (me *Editor) Load(filename string) error {
	fh, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fh.Close()
	return me.Import(fh)
}

func (me *Editor) Save(filename string) error {
	fh, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer fh.Close()
	return me.TidyExport(fh)
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

func (me *Editor) SaveReport(filename string) error {
	fh, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer fh.Close()
	return me.WriteReport(fh)
}

// WriteReport writes a markdown report
func (me *Editor) WriteReport(w io.Writer) error {
	p, err := nexus.NewPrinter(w)
	p.Println("# ISVS Report")
	p.Println("https://github.com/OWASP/IoT-Security-Verification-Standard-ISVS/")

	p.Println("## Summary")
	p.Println()
	p.Println("- L1:", me.Stats(me.list(1)))
	p.Println("- L2:", me.Stats(me.list(2)))
	p.Println("- L3:", me.Stats(me.list(3)))

	p.Println()
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

func (me *Editor) Stats(entries []Entry) string {
	var num int
	var verified int
	var applicable int
	for _, e := range entries {
		num++
		if e.Applicable {
			applicable++
		}
		if e.Verified {
			verified++
		}
	}
	return fmt.Sprintf("%d/%d applicable (total %d)", verified, applicable, num)
}

func (me *Editor) list(level int) []Entry {
	res := make([]Entry, 0)
	if level < 1 || level > 3 {
		panic(fmt.Errorf("no such level %v", level))
	}
	for _, e := range me.entries {
		switch {
		case level == 1 && e.L1:
			res = append(res, e)
		case level == 2 && e.L2:
			res = append(res, e)
		case level == 3 && e.L3:
			res = append(res, e)
		}
	}
	return res
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
