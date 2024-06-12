package main

import "strings"

type Lang struct {
	Name string
	Extensions []string
	TotalCount int
}
// newLang initializes a Lang instance with default counts
func newLang(name string, extensions []string) Lang {
	return Lang{
		Name:         name,
		Extensions:   extensions,
		TotalCount:   0,
	}
}

var langs = map[string]Lang{
	"Go":         newLang("Golang", []string{".go"}),
	"Python":     newLang("Python", []string{".py"}),
	"Ruby":       newLang("Ruby", []string{".rb"}),
	"JavaScript": newLang("JavaScript", []string{".js", ".cjs", ".mjs", ".jsx"}),
	"Java":       newLang("Java", []string{".java"}),
	"C":          newLang("C", []string{".c"}),
	"C++":        newLang("C++", []string{".cpp", ".cc"}),
	"Rust":       newLang("Rust", []string{".rs"}),
	"TypeScript": newLang("TypeScript", []string{".ts", ".tsx"}),
	"Shell":      newLang("Shell", []string{".sh"}),
	"Swift":      newLang("Swift", []string{".swift"}),
	"Kotlin":     newLang("Kotlin", []string{".kt"}),
	"PHP":        newLang("PHP", []string{".php"}),
	"HTML":       newLang("HTML", []string{".html"}),
	"CSS":        newLang("CSS", []string{".css", ".scss", ".sass", ".less"}),
	"SQL":        newLang("SQL", []string{".sql"}),
	"R":          newLang("R", []string{".r"}),
	"Scala":      newLang("Scala", []string{".scala"}),
	"Perl":       newLang("Perl", []string{".pl"}),
	"Lua":        newLang("Lua", []string{".lua"}),
	"Objective-C":newLang("Objective-C", []string{".m"}),
	"Assembly":   newLang("Assembly", []string{".asm"}),
	"Vim script": newLang("Vim script", []string{".vim"}),
	"Groovy":     newLang("Groovy", []string{".groovy"}),
	"PowerShell": newLang("PowerShell", []string{".ps1"}),
	"Racket":     newLang("Racket", []string{".rkt"}),
	"OCaml":      newLang("OCaml", []string{".ml"}),
	"Julia":      newLang("Julia", []string{".jl"}),
	"Scheme":     newLang("Scheme", []string{".scm"}),
	"Markdown":   newLang("Markdown", []string{".md", ".mdx"}),
	"TeX":        newLang("TeX", []string{".tex"}),
	"JSON":       newLang("JSON", []string{".json"}),
	"YAML":       newLang("YAML", []string{".yaml", ".yml"}),
	"XML":        newLang("XML", []string{".xml"}),
}

func getLangByExt(ext string) string {
	for lang, exts := range langs {
		for _, e := range exts.Extensions {
			if e == ext {
				return lang
			}
		}
	}
	return "Unknown"
}

func IsValidFile(file, ext string) bool {

	correctSuffix := strings.HasSuffix(file, ext)
	for _, e := range langs[getLangByExt(ext)].Extensions {
		if e == ext {
			return true && correctSuffix
		}
	}

	return false
}


func GetLang(ext string) Lang {
	return langs[getLangByExt(ext)]
}

