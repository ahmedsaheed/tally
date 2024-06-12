package main

import "strings"

type Language struct {
	Name string
	Extensions []string
	TotalCount int
}

var Languages = map[string]Language{
	"Go":         newLanguage("Golang", []string{".go"}),
	"Python":     newLanguage("Python", []string{".py"}),
	"Ruby":       newLanguage("Ruby", []string{".rb"}),
	"JavaScript": newLanguage("JavaScript", []string{".js", ".cjs", ".mjs", ".jsx"}),
	"Java":       newLanguage("Java", []string{".java"}),
	"C":          newLanguage("C", []string{".c"}),
	"C++":        newLanguage("C++", []string{".cpp", ".cc"}),
	"Rust":       newLanguage("Rust", []string{".rs"}),
	"TypeScript": newLanguage("TypeScript", []string{".ts", ".tsx"}),
	"Bash":      newLanguage("Bash", []string{".sh"}),
	"Swift":      newLanguage("Swift", []string{".swift"}),
	"Kotlin":     newLanguage("Kotlin", []string{".kt"}),
	"PHP":        newLanguage("PHP", []string{".php"}),
	"HTML":       newLanguage("HTML", []string{".html"}),
	"CSS":        newLanguage("CSS", []string{".css", ".scss", ".sass", ".less"}),
	"SQL":        newLanguage("SQL", []string{".sql"}),
	"R":          newLanguage("R", []string{".r"}),
	"Scala":      newLanguage("Scala", []string{".scala"}),
	"Perl":       newLanguage("Perl", []string{".pl"}),
	"Lua":        newLanguage("Lua", []string{".lua"}),
	"Objective-C":newLanguage("Objective-C", []string{".m"}),
	"Assembly":   newLanguage("Assembly", []string{".asm"}),
	"Vim script": newLanguage("Vim script", []string{".vim"}),
	"Groovy":     newLanguage("Groovy", []string{".groovy"}),
	"Racket":     newLanguage("Racket", []string{".rkt"}),
	"OCaml":      newLanguage("OCaml", []string{".ml"}),
	"Julia":      newLanguage("Julia", []string{".jl"}),
	"Scheme":     newLanguage("Scheme", []string{".scm"}),
	"Markdown":   newLanguage("Markdown", []string{".md", ".mdx"}),
	"TeX":        newLanguage("TeX", []string{".tex"}),
	"JSON":       newLanguage("JSON", []string{".json"}),
	"YAML":       newLanguage("YAML", []string{".yaml", ".yml"}),
	"TOML":       newLanguage("TOML", []string{".toml"}),
	"Fish":       newLanguage("Fish", []string{".fish"}),
}

// newLanguage initializes a Lang instance with default counts
func newLanguage(name string, extensions []string) Language {
	return Language{
		Name:         name,
		Extensions:   extensions,
		TotalCount:   0,
	}
}

func lookupLangByExtension(ext string) string {
	for lang, exts := range Languages {
		for _, e := range exts.Extensions {
			if e == ext {
				return lang
			}
		}
	}
	return "Unknown"
}

func IsFileExtensionValid(file, ext string) bool {
	correctSuffix := strings.HasSuffix(file, ext)
	for _, e := range Languages[lookupLangByExtension(ext)].Extensions {
		if e == ext {
			return true && correctSuffix
		}
	}
	return false
}
