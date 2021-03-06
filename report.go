package owasp

import (
	"fmt"
	"io"
	"os"

	"github.com/gregoryv/nexus"
)

type Report struct {
	entries            []Entry
	Title              string
	ShortDescriptionNA bool // Short description for non applicable
}

func (me *Report) AddEntries(v ...Entry) {
	me.entries = append(me.entries, v...)
}

// SaveAs saves the report as markdown to the given filename.
func (me *Report) SaveAs(filename string) error {
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
	v, a, na := me.stats(me.entries)
	p.Printf("%d applicable requirements of %d\n", a, na)
	p.Println()

	p.Println("- L1:", me.Stats(me.groupByLevel(L1)))
	p.Println("- L2:", me.Stats(me.groupByLevel(L2)))
	p.Println("- L3:", me.Stats(me.groupByLevel(L3)))
	p.Println()

	sumLine := fmt.Sprintf("%d requirements left to verify!", a-v)
	if v == a {
		p.Println("All requirements verified!")
	}
	p.Println(sumLine)
	p.Println()
	p.Println("## Applicable")
	for _, e := range me.entries {
		if !e.Applicable {
			continue
		}
		p.Printf("- %s **%s** %s\n", checkbox(e), e.ID, e.Description)
		if e.Manual != nil {
			p.Printf("  _MANUAL %s by %s: %s_\n",
				e.Manual.When, e.Manual.By, e.Manual.How,
			)
		}
	}

	p.Println()
	p.Println("## Not Applicable")
	for _, e := range me.entries {
		if e.Applicable {
			continue
		}
		desc := e.Description
		if me.ShortDescriptionNA {
			desc = maxString(e.Description, 80)
		}
		p.Printf("- %s %s\n", e.ID, desc)
	}

	return p.Written, *err
}

func checkbox(e Entry) string {
	checkbox := "[ ]"
	if e.Verified {
		checkbox = "[x]"
	}
	return checkbox
}

func maxString(s string, l int) string {
	if len(s) < l {
		return s
	}
	return s[:l] + "..."
}

func (me *Report) Stats(entries []Entry) string {
	verified, applicable, _ := me.stats(entries)
	return fmt.Sprintf("%d verified of %d", verified, applicable)
}

func (me *Report) stats(entries []Entry) (verified, applicable, total int) {
	for _, e := range entries {
		total++
		if e.Applicable {
			applicable++
		}
		if e.Verified {
			verified++
		}
	}
	return
}

// groupByLevel returns requirements only in the given level, excluding
// multilevel requirements
func (me *Report) groupByLevel(level Level) []Entry {
	res := make([]Entry, 0)
	for _, e := range me.entries {
		if e.IsLevel(level) {
			res = append(res, e)
		}
	}
	return res
}
