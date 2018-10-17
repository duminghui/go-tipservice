// Package main provides ...
package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

//go get -u github.com/kevinburke/go-bindata/...
//go:generate go-bindata -o=assets.go templates

type tmplValueMap map[string]interface{}

var msgTmpl *template.Template

func loadTemplates() {
	dir, err := ioutil.TempDir("", "piebottemplates")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)
	err = RestoreAssets(dir, "templates")
	if err != nil {
		panic(err)
	}
	log.Infof("Templates Path:%s", dir)
	funcMap := template.FuncMap{
		"join": strings.Join,
	}
	tmpl, err := template.New("main").Funcs(funcMap).ParseGlob(filepath.Join(dir, "templates/*.tmpl"))
	if err != nil {
		panic(err)
	}
	msgTmpl = tmpl
	// for _, t := range tmpl.Templates() {
	// 	log.Infof("Template Name:%s", t.Name())
	// }
}

func msgFromTmpl(tmplName string, data interface{}) string {
	buf := new(strings.Builder)
	err := msgTmpl.ExecuteTemplate(buf, tmplName, data)
	if err != nil {
		return err.Error()
	}
	return buf.String()
}
