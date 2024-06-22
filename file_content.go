package main

import "strings"

var EncounteredLangsArr = []Language{}

func SetEncounteredLangs(lang Language) {
	EncounteredLangsArr = append(EncounteredLangsArr, lang)
}

func GetEncounteredLangs() []Language {
	return EncounteredLangsArr
}

func hasExpectedLanguage(file string, expectedExt []string) bool {
	for _, ext := range expectedExt {
		if strings.HasSuffix(file, ext) || isDockerfile(file) {
			return true
		}
	}
	return false
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
