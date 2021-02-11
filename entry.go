package owasp

import "fmt"

type Entry struct {
	L1          bool
	L2          bool
	L3          bool
	Description string
	ID          string
	Verified    bool
	Applicable  bool
}

func (me *Entry) checkbox() string {
	checkbox := "[ ]"
	if me.Verified {
		checkbox = "[x]"
	}
	return checkbox
}

func (me *Entry) shortString() string {
	return fmt.Sprintf("%s %s...", me.ID, me.shortDesc())
}

func (me *Entry) String() string {
	return fmt.Sprintf("%s %s", me.ID, me.Description)
}

func (me *Entry) shortDesc() string {
	if len(me.Description) < 80 {
		return me.Description
	}
	return me.Description[:80] + "..."
}
