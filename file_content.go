package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var EncounteredLangsArr = []Language{}

func SetEncounteredLangs(lang Language) {
	EncounteredLangsArr = append(EncounteredLangsArr, lang)
}

func GetEncounteredLangs() []Language {
	return EncounteredLangsArr
}

func hasExpectedLanguage(file string, expectedExt []string) bool {
	for _, ext := range expectedExt {
		if strings.HasSuffix(file, ext) {
			return true
		}
	}
	return false
}

func CountUp(lang Language, root string) (Language, error) {
	totalLine := 0
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil { return err }
		// if isIgnoredDir(path) { return filepath.SkipDir }





		if !info.IsDir() && hasExpectedLanguage(path, lang.Extensions) && !isIgnoredFile(path)  && !isIgnoredDir(path) {
			
			if lang.Extensions[0] == ".py" {
				fmt.Println("py file found:", path)
			}

			line, err := Scan(path)

			if err != nil {
				return err
			}
			
			totalLine += line
		}
			
		lang.TotalCount = totalLine
		return nil
		})

		if err != nil{
			return lang, err
		}
		return lang, nil
}

var ignoredDirs = []string{".git", ".idea", "node_modules", "vendor"}
var ignoreFiles = []string{".DS_Store", ".gitignore", ".gitkeep", "package-lock.json", "yarn.lock"}


func isIgnoredDir(dir string) bool {
	for _, d := range ignoredDirs {
		if strings.HasSuffix(dir, d) {
			return true
		}
	}
	return false
}

func isIgnoredFile(file string) bool {
	for _, f := range ignoreFiles {
		if strings.HasSuffix(file, f) {
			return true
		}
	}
	return false
}

