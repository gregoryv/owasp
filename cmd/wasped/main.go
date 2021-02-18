package main

import (
	"log"
	"os"
	"strings"

	"github.com/gregoryv/cmdline"
	"github.com/gregoryv/owasp"
)

func main() {
	var (
		cli      = cmdline.NewParser(os.Args...)
		help     = cli.Flag("-h, --help")
		verify   = cli.Option("--verify", "ID to set as verified").String("")
		unverify = cli.Option("--unverify", "ID to set as unverified").String("")
		manstr   = cli.Option("-m, --manual",
			"Comma separated list of manual verification notes.",
			"Format: YYYY-MM-DD, By, How",
		).String("")

		rfile   = cli.Option("-r, --report", "Markdown file to save report to").String("")
		title   = cli.Option("-t, --title", "Title of the report").String("Report")
		shortna = cli.Option("-s, --short-description-na",
			"If given, not applicable requiremens will have a shortened description",
		).Bool()

		file = cli.Required("FILE").String()
	)
	log.SetFlags(0)

	switch {
	case help:
		cli.WriteUsageTo(os.Stdout)
		os.Exit(0)

	case !cli.Ok():
		log.Fatal(cli.Error())

	case file == "":
		log.Fatal("empty file, see --help")
	}

	ed := owasp.NewEditor()
	ed.Load(file)

	var man *owasp.Manual
	if manstr != "" {
		parts := strings.Split(manstr, ",")
		if len(parts) != 3 {
			log.Fatal("invalid --manual text, see --help")
		}
		man = &owasp.Manual{
			How:  parts[2],
			When: parts[1],
			By:   parts[0],
		}
	}
	var (
		id       string = verify
		verified bool   = true
	)
	if unverify != "" {
		id = unverify
		verified = false
	}

	switch {
	case man != nil:
		must(ed.SetManuallyVerified(id, verified, *man))

	default:
		must(ed.SetVerified(id, verified))
	}

	if rfile != "" {
		report := ed.NewReport(title)
		report.ShortDescriptionNA = shortna
		must(report.Save(rfile))
	}
	ed.Save(file)
}

func must(err error) {
	if err == nil {
		return
	}
	log.Fatal(err)
}
