# Changelog
All notable changes to this project will be documented in this file.

The format is based on http://keepachangelog.com/en/1.0.0/
and this project adheres to http://semver.org/spec/v2.0.0.html.

## [unreleased]

- Add Report.ShowNonApplicable, set to false hides the non applicable section
- SetVerified uses same pattern as SetApplicable
- SetApplicable, SetVerified understands id values as patterns or by level
- Removed SetApplicableBy and SetApplicableByLevel, use SetApplicable
- Removed SetManuallyVerified and SetManuallyVerifiedBy, use SetVerified
- SetVerified fails if entry is not applicable

## [0.8.0] - 2021-03-05

- Replace Editor.TidyExport with WriteTo implementing io.WriterTo
- Remove Editor.Export
- Rename Editor.Save to Editor.SaveAs
- Remove MustSetVerifiedNow
- Add Entry.IsLevel

## [0.7.1] - 2021-03-03

- Fix SetApplicableByLevel

## [0.7.0] - 2021-03-03

- Methods that set verified field fail if the requirement is not applicable
- Add SetVerifiedBy and ResetVerifiedBy
- Use SetApplicable and SetApplicableBy to set value to true or false

## [0.6.1] - 2021-02-18

- Fix double entries in report

## [0.6.0] - 2021-02-18

- Add Editor.SetManuallyVerified and type Manual to describe how
- Renamed wasped options -set and -unset to --verify and --unverify

## [0.5.0] - 2021-02-11

- Summary shows more clearly number of requirements left to verify

## [0.4.1] - 2021-02-11
## [0.4.0] - 2021-02-11

- Add Editor.NewReport
- Replace Editor.SaveReport with Report.Save
- Replace Editor.WriteReport with Report.WriteTo, matching the
  io.WriterTo interface
- Report uses short description for non applicable entries

## [0.3.0] - 2021-02-10

- Add examples of using wasped
- Add func MustSetVerifiedNow for easy integration in tests

## [0.2.0] - 2021-02-10

- Add flag --title to cmd/wasped
- Add cmd/wasped
- SetVerified takes value for toggling

## [0.1.0] - 2021-02-10

- Render basic markdown report
- Edit OWASP checklist, tried on ISVS
