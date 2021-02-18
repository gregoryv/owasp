package owasp

// MustSetVerifiedNow loads and sets the verified flag. Panics on errors.
func MustSetVerifiedNow(id, filename string, v bool) {
	ed := NewEditor()
	if err := ed.Load(filename); err != nil {
		panic(err)
	}
	if err := ed.SetVerified(id, v); err != nil {
		panic(err)
	}
	if err := ed.Save(filename); err != nil {
		panic(err)
	}
}
