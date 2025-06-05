package main

import (
	_ "embed"
	"html/template"
)

//go:embed goreleaser.yml.tmpl
var goreleaserTemplateSource string
var goreleaserTemplate = template.Must(
	template.New("goreleaser.yml").
		Delims("⦃", "⦄").
		Parse(goreleaserTemplateSource),
)
