package owasp

import "fmt"

type Entry struct {
	L1          bool
	L2          bool
	L3          bool
	Description string
	ID          string

	Applicable bool
	Verified   bool

	*Manual `json:"Manual,omitempty"`
}

func (me Entry) String() string {
	return fmt.Sprintf("L%v %s", me.Level(), me.ID)
}

// Level returns the highest level for this entry
func (me *Entry) Level() Level {
	switch {
	case me.L3:
		return L3
	case me.L2:
		return L2
	case me.L1:
		return L1
	}
	panic("Bad entry, level not set!")
}

func (me *Entry) IsLevel(level Level) bool {
	switch {

	case level == L1 && me.L1:
		return true

	case level == L2 && me.L2:
		return true

	case level == L3 && me.L3:
		return true

	}
	return false
}

// Manual describes manual verification
type Manual struct {
	How  string // How it was done
	When string // yyyy-mm-dd
	By   string
}
