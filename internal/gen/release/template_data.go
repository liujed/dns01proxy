package main

import "github.com/liujed/dns01proxy/internal/gen"

type TemplateData struct {
	Builds map[string]gen.Build
}

func MakeTemplateData() TemplateData {
	return TemplateData{
		Builds: map[string]gen.Build{},
	}
}
