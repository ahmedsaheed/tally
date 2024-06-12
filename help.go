package main

import (
	"os"
	"path/filepath"
	"strings"
)

func GetAllLanguagesInDir(root string) []Language {
	alreadyInserted := map[string]bool{}
	availableLangs := []Language{}
	
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil { return err }
		if !info.IsDir() {
			ext := fileExtensionFromPath(path)
			langName := lookupLangByExtension(ext)
			lang := Languages[langName]
			if len(lang.Name) != 0 && !alreadyInserted[lang.Name] {
				availableLangs = append(availableLangs, lang)
				alreadyInserted[lang.Name] = true
			}
		}
		return nil
	})

	if err != nil {
		return nil
	}

	// fmt.Println("available lang are:", availableLangs)
	
	return availableLangs
}

func sortByTotalLines(langs []Language) {
	for i := 0; i < len(langs); i++ {
		for j := i + 1; j < len(langs); j++ {
			if langs[i].TotalCount < langs[j].TotalCount {
				langs[i], langs[j] = langs[j], langs[i]
			}
		}
	}
}

func SumLines(langs []Language) int {
	total := 0
	for _, lang := range langs {
		total += lang.TotalCount
	}
	return total
}

func fileExtensionFromPath(path string) string { return "." + strings.TrimPrefix(filepath.Ext(path), ".") }
