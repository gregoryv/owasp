package main

import (
	"log"
	"os"

	"github.com/gregoryv/cmdline"
	"github.com/gregoryv/owasp"
)

func main() {
	var (
		cli      = cmdline.NewParser(os.Args...)
		help     = cli.Flag("-h, --help")
		verify   = cli.Option("--verify", "ID to set as verified").String("")
		unverify = cli.Option("--unverify", "ID to set as unverified").String("")

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
		log.Fatal("empty file, see -help")
	}

	ed := owasp.NewEditor()
	ed.Load(file)

	if verify != "" {
		must(ed.SetVerified(verify, true))
	}
	if unverify != "" {
		must(ed.SetVerified(unverify, false))
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
