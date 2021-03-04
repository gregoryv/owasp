package owasp

import "testing"

func Example_workWithTheEditor() {
	ed := NewEditor()
	ed.Load("checklist/asvs.json")
	ed.SetApplicableBy(`1\.1\.\d*`, true)

	verifyASVS(t, ed)

	man := Manual{
		How:  "Latest threatmodel design change was updated on ...",
		When: "2021-02-18",
		By:   "John Doe",
	}
	_ = ed.SetManuallyVerified("1.1.2", true, man)

	ed.NewReport("Report ASVS").SaveAs("example_report.md")
	// output:
}

var t *testing.T

func verifyASVS(t *testing.T, ed *Editor) {
	// write your tests here and check of requirements
	ed.SetVerified("1.1.1", true)
	ed.SetVerified("1.1.2", false)
}
