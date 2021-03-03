Package [owasp](https://pkg.go.dev/github.com/gregoryv/owasp) provides
an [OWASP](https://github.com/OWASP) checklist editor.

It was written to integrate ISVS and ASVS checklists with tests that
verify the requirements.

## Quick start

Install the editor

    go get -u github.com/gregoryv/cmd/wasped

Prepare a checklist, start of with [asvs.json](checklist/asvs.json)
or [isvs.json](checklist/isvs.json) found in this repository. Set the
Applicable field to true on each entry that is applicable to your
project.

When you have verified a requirement check it off with 

    $ wasped --verify "1.3.2" asvs.json

finally you can render a markdown report summarizing your progress

    $ wasped --report asvs_report.md --title "My ASVS report" asvs.json

## Automate verification in tests

```go
func Test_some_feature(t *testing.T) {
    // test code here
    // Then update the checklist	
    MustSetVerifiedNow("1.3.2", "asvs.json", true) // or false if failed
}
```

