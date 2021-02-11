package owasp

import (
	"fmt"
	"io"
	"os"

	"github.com/gregoryv/draw"
	"github.com/gregoryv/draw/design"
	"github.com/gregoryv/draw/shape"
	"github.com/gregoryv/nexus"
)

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
	p.Println(me.sumChart().Inline())
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
		p.Printf("- %s **%s** %s\n", checkbox(e), e.ID, e.Description)
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

func (me *Report) sumChart() *design.Diagram {
	var (
		d = design.NewDiagram()
	)
	width := 400

	draw.DefaultClassAttributes["green"] = `stroke="black" stroke-width="0" fill="#ccff99" fill-opacity="1.0"`
	draw.DefaultClassAttributes["blue"] = `stroke="black" stroke-width="0" fill="#99e6ff" fill-opacity="1.0"`
	draw.DefaultClassAttributes["gray"] = `stroke="black" stroke-width="0" fill="#e2e2e2" fill-opacity="1.0"`

	l1v, l1a, l1na := me.bar(me.entries, width)
	d.Place(l1v).At(20, 20)
	d.Place(l1a, l1na).RightOf(l1v, 0)

	return &d
}

func (me *Report) bar(entries []Entry, width int) (v, a, na *shape.Rect) {
	verified, applicable, total := me.stats(entries)
	v = shape.NewRect("")
	a = shape.NewRect("")
	na = shape.NewRect("")

	draw.DefaultClassAttributes["green"] = `stroke="black" stroke-width="0" fill="#ccff99" fill-opacity="1.0"`
	v.SetClass("green")
	v.SetWidth(part(verified, total, width))

	draw.DefaultClassAttributes["blue"] = `stroke="black" stroke-width="0" fill="#99e6ff" fill-opacity="1.0"`
	a.SetClass("blue")
	a.SetWidth(part((applicable - verified), total, width))

	draw.DefaultClassAttributes["gray"] = `stroke="black" stroke-width="0" fill="#e2e2e2" fill-opacity="1.0"`
	na.SetClass("gray")
	na.SetWidth(part((total - applicable), total, width))
	return
}

func part(a, b, c int) int {
	v := (float64(a) / float64(b)) * float64(c)
	return int(v)
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
	verified, applicable, total := me.stats(entries)
	return fmt.Sprintf("%d/%d applicable (total %d)", verified, applicable, total)
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

// ----------------------------------------

type Entry struct {
	L1          bool
	L2          bool
	L3          bool
	Description string
	ID          string
	Verified    bool
	Applicable  bool
}
