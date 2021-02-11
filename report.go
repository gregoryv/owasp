package owasp

import (
	"fmt"
	"io"
	"os"

	"github.com/gregoryv/nexus"
)

func NewReport(title string) *Report {
	return &Report{
		Title:              title,
		ShortDescriptionNA: true,
	}
}

type Report struct {
	entries            []Entry
	Title              string
	ShortDescriptionNA bool // true to shorten description for all non applicable
}

func (me *Report) AddEntries(v ...Entry) {
	me.entries = append(me.entries, v...)
}

// Save saves the report as markdown to the given filename.
func (me *Report) Save(filename string) error {
	fh, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer fh.Close()
	_, err = me.WriteTo(fh)
	return err
}

// WriteReport writes a markdown report
func (me *Report) WriteTo(w io.Writer) (int64, error) {
	p, err := nexus.NewPrinter(w)
	p.Println("#", me.Title)

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
		p.Printf("- %s **%s** %s\n", e.checkbox(), e.ID, e.Description)
	}

	p.Println()
	p.Println("## Not Applicable")
	for _, e := range me.entries {
		if e.Applicable {
			continue
		}
		desc := e.Description
		if me.ShortDescriptionNA {
			desc = e.shortDesc()
		}
		p.Printf("- %s %s\n", e.ID, desc)
	}

	return p.Written, *err
}

func (me *Report) Stats(entries []Entry) string {
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

func (me *Report) list(level int) []Entry {
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
