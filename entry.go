package owasp

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
