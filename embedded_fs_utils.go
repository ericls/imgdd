package main

import (
	"embed"
	"imgdd/buildflag"
	"io/fs"
	"os"
)

//go:embed httpserver/templates/*.gotmpl
var embeddedTemplates embed.FS

//go:embed web_client/dist/*
var embeddedStatic embed.FS

type moutingFS struct {
	Templates fs.FS
	Static    fs.FS
}

var MoutingFS moutingFS

func init() {
	if buildflag.Debug == "true" {
		MoutingFS.Templates = os.DirFS("httpserver/templates")
		MoutingFS.Static = os.DirFS("web_client/dist")
	} else {
		templatesFs, err := fs.Sub(embeddedTemplates, "httpserver/templates")
		if err != nil {
			println("Error loading embedded templates: ", err.Error())
			panic(err)
		}
		MoutingFS.Templates = templatesFs
		staticFs, err := fs.Sub(embeddedStatic, "web_client/dist")
		if err != nil {
			println("Error loading embedded templates: ", err.Error())
			panic(err)
		}
		MoutingFS.Static = staticFs
	}
}
