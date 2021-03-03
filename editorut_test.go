package owasp

import (
	"io"
	"testing"
)

type EditorUnderTest struct {
	*testing.T
	*Editor
}

func (me *Editor) UnderTest(t *testing.T) *EditorUnderTest {
	return &EditorUnderTest{T: t, Editor: me}
}

func (me *EditorUnderTest) shouldSetApplicableByLevel(level Level, appl bool) {
	err := me.SetApplicableByLevel(level, appl)
	if err != nil {
		me.T.Helper()
		me.T.Error(err)
	}
}

func (me *EditorUnderTest) mustSetApplicableByLevel(level Level, appl bool) {
	err := me.SetApplicableByLevel(level, appl)
	if err != nil {
		me.T.Helper()
		me.T.Fatal(err)
	}
}

func (me *EditorUnderTest) shouldSetVerified(id string, v bool) {
	err := me.SetVerified(id, v)
	if err != nil {
		me.T.Helper()
		me.T.Error(err)
	}
}

func (me *EditorUnderTest) mustSetVerified(id string, v bool) {
	err := me.SetVerified(id, v)
	if err != nil {
		me.T.Helper()
		me.T.Fatal(err)
	}
}

func (me *EditorUnderTest) shouldResetVerifiedBy(pattern string) {
	err := me.ResetVerifiedBy(pattern)
	if err != nil {
		me.T.Helper()
		me.T.Error(err)
	}
}

func (me *EditorUnderTest) mustResetVerifiedBy(pattern string) {
	err := me.ResetVerifiedBy(pattern)
	if err != nil {
		me.T.Helper()
		me.T.Fatal(err)
	}
}

func (me *EditorUnderTest) shouldSetManuallyVerified(id string, v bool, man Manual) {
	err := me.SetManuallyVerified(id, v, man)
	if err != nil {
		me.T.Helper()
		me.T.Error(err)
	}
}

func (me *EditorUnderTest) mustSetManuallyVerified(id string, v bool, man Manual) {
	err := me.SetManuallyVerified(id, v, man)
	if err != nil {
		me.T.Helper()
		me.T.Fatal(err)
	}
}

func (me *EditorUnderTest) shouldSetManuallyVerifiedBy(pattern string, v bool, man Manual) {
	err := me.SetManuallyVerifiedBy(pattern, v, man)
	if err != nil {
		me.T.Helper()
		me.T.Error(err)
	}
}

func (me *EditorUnderTest) mustSetManuallyVerifiedBy(pattern string, v bool, man Manual) {
	err := me.SetManuallyVerifiedBy(pattern, v, man)
	if err != nil {
		me.T.Helper()
		me.T.Fatal(err)
	}
}

func (me *EditorUnderTest) shouldSetApplicableBy(pattern string, v bool) {
	err := me.SetApplicableBy(pattern, v)
	if err != nil {
		me.T.Helper()
		me.T.Error(err)
	}
}

func (me *EditorUnderTest) mustSetApplicableBy(pattern string, v bool) {
	err := me.SetApplicableBy(pattern, v)
	if err != nil {
		me.T.Helper()
		me.T.Fatal(err)
	}
}

func (me *EditorUnderTest) shouldSetApplicable(id string, v bool) {
	err := me.SetApplicable(id, v)
	if err != nil {
		me.T.Helper()
		me.T.Error(err)
	}
}

func (me *EditorUnderTest) mustSetApplicable(id string, v bool) {
	err := me.SetApplicable(id, v)
	if err != nil {
		me.T.Helper()
		me.T.Fatal(err)
	}
}

func (me *EditorUnderTest) shouldLoad(filename string) {
	err := me.Load(filename)
	if err != nil {
		me.T.Helper()
		me.T.Error(err)
	}
}

func (me *EditorUnderTest) mustLoad(filename string) {
	err := me.Load(filename)
	if err != nil {
		me.T.Helper()
		me.T.Fatal(err)
	}
}

func (me *EditorUnderTest) shouldSave(filename string) {
	err := me.Save(filename)
	if err != nil {
		me.T.Helper()
		me.T.Error(err)
	}
}

func (me *EditorUnderTest) mustSave(filename string) {
	err := me.Save(filename)
	if err != nil {
		me.T.Helper()
		me.T.Fatal(err)
	}
}

func (me *EditorUnderTest) shouldImport(r io.Reader) {
	err := me.Import(r)
	if err != nil {
		me.T.Helper()
		me.T.Error(err)
	}
}

func (me *EditorUnderTest) mustImport(r io.Reader) {
	err := me.Import(r)
	if err != nil {
		me.T.Helper()
		me.T.Fatal(err)
	}
}

func (me *EditorUnderTest) shouldTidyExport(w io.Writer) {
	err := me.TidyExport(w)
	if err != nil {
		me.T.Helper()
		me.T.Error(err)
	}
}

func (me *EditorUnderTest) mustTidyExport(w io.Writer) {
	err := me.TidyExport(w)
	if err != nil {
		me.T.Helper()
		me.T.Fatal(err)
	}
}

func (me *EditorUnderTest) shouldExport(w io.Writer) {
	err := me.Export(w)
	if err != nil {
		me.T.Helper()
		me.T.Error(err)
	}
}

func (me *EditorUnderTest) mustExport(w io.Writer) {
	err := me.Export(w)
	if err != nil {
		me.T.Helper()
		me.T.Fatal(err)
	}
}
