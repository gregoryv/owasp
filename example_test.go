package owasp

import "testing"

func Example_workWithTheEditor() {
	ed := NewEditor()
	ed.Load("testdata/asvs.json")

	// Reset all
	ed.Reset()

	// Select requirements that apply for your project
	ed.SetApplicable(`1.1.*`, true)

	// Verify, with tests
	verifyASVS(t, ed)

	// or manually
	man := Manual{
		How:  "Latest threatmodel design change was updated on ...",
		When: "2021-02-18",
		By:   "John Doe",
	}
	_ = ed.SetManuallyVerified("1.1.2", true, man)

	// generate a nice report
	ed.NewReport("Report ASVS").SaveAs("example_report.md")

	// Save result
	ed.SaveAs("testdata/asvs.json")

	// output:
}

var t *testing.T

func verifyASVS(t *testing.T, ed *Editor) {
	// write your tests here and check of requirements
	ed.SetVerified("1.1.1", true)
	ed.SetVerified("1.1.2", false)
}
