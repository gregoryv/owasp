# Changelog
All notable changes to this project will be documented in this file.

The format is based on http://keepachangelog.com/en/1.0.0/
and this project adheres to http://semver.org/spec/v2.0.0.html.

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
