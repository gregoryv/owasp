Package [owasp](https://pkg.go.dev/github.com/gregoryv/owasp) provides
an [OWASP](https://github.com/OWASP) conformance editor.

It was written to integrate ISVS and ASVS checklists with tests that
verify the requirements.

In your test use it like

```go:
func Test_some_feature(t *testing.T) {
    // test code here
    // Then update the checklist
    MustSetVerifiedNow("1.3.2", "isvs.json", true) // or false if failed
}
```
