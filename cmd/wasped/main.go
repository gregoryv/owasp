package main

import (
	"log"
	"os"

	"github.com/gregoryv/cmdline"
	"github.com/gregoryv/owasp"
)

func main() {
	var (
		cli   = cmdline.NewParser(os.Args...)
		help  = cli.Flag("-h, --help")
		set   = cli.Option("-set").String("")
		unset = cli.Option("-unset").String("")

		rfile   = cli.Option("-r, --report").String("")
		title   = cli.Option("-t, --title").String("Report")
		shortna = cli.Flag("-s, --short-description-na")

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

	if set != "" {
		must(ed.SetVerified(set, true))
	}
	if unset != "" {
		must(ed.SetVerified(unset, false))
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
