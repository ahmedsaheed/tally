package main

import (
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
	fileCount := 0
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && hasExpectedLanguage(path, lang.Extensions) && !isIgnoredFile(path) && !isIgnoredDir(path) && !isDotedDir(path) {
			line, err := Scan(path)
			fileCount++
			if err != nil {
				return err
			}

			totalLine += line
		}

		lang.TotalCount = totalLine
		lang.FileCount = fileCount

		return nil
	})

	if err != nil {
		return lang, err
	}
	return lang, nil
}

var ignoredDirs = []string{".git", ".idea", "node_modules", "vendor"}
var ignoreFiles = []string{".DS_Store", ".gitignore", ".gitkeep", "package-lock.json", "yarn.lock"}

func isIgnoredDir(dir string) bool {
	for _, d := range ignoredDirs {
		if strings.Contains(dir, d) {
			return true
		}
	}
	return false
}

func isDotedDir(dir string) bool {
	firstSlash := strings.Index(dir, "/")
	if firstSlash == -1 {
		return false
	}

	secondSlash := strings.Index(dir[firstSlash+1:], "/")
	if secondSlash == -1 {
		return false
	}

	secondSlash += firstSlash + 1

	if len(dir) > secondSlash+1 && dir[secondSlash+1] == '.' {
		return true
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
