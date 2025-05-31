package main

import (
	_ "embed"
	"text/template"
)

//go:embed main.go.tmpl
var mainTemplateSource string
var mainTemplate = template.Must(
	template.New("main.go").Parse(mainTemplateSource),
)

//go:embed Makefile.tmpl
var makefileTemplateSource string
var makefileTemplate = template.Must(
	template.New("Makefile").Parse(makefileTemplateSource),
)
