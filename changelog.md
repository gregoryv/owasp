# Changelog
All notable changes to this project will be documented in this file.

The format is based on http://keepachangelog.com/en/1.0.0/
and this project adheres to http://semver.org/spec/v2.0.0.html.

## [unreleased]

- Add Editor.NewReport
- Replace Editor.SaveReport with Report.Save
- Replace Editor.WriteReport with Report.WriteTo, matching the io.WriterTo interface
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
