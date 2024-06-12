package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const DOT = "."


func getLangExt(file string) string {
	ext := filepath.Ext(file)
	return DOT + strings.TrimPrefix(ext, ".")
}

func countLines(file string) (int, error) {
	f, err := os.Open(file)

	if err != nil { return 0, err }
	defer f.Close()
	reader := bufio.NewReader(f)

	lineCount := 0
	for {
		_, err := reader.ReadString('\n')
		if err != nil {
			if err == os.ErrClosed {
				break
			}
			if err != nil {
				return lineCount, nil
			}
		}
		lineCount++
	}

	return lineCount, nil
	// scanner := bufio.NewScanner(f)

	// lineCount := 0
	// for scanner.Scan() { lineCount++ }

	// return lineCount, scanner.Err()
}

func getAllLangsInCodebase(root string) []Lang {
	alreadyInserted := map[string]bool{}
	availableLangs := []Lang{}
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil { return err }
		if !info.IsDir() {
			ext := getLangExt(path)
			langName := getLangByExt(ext)
			lang := langs[langName]
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

	
	return availableLangs
}

func tallyCodebase(root string) ([]Lang, error) {
	langs := getAllLangsInCodebase(root)
	for _, lang := range langs {
		resultBylang, err := CalculateByLangFromRoot(lang, root)
		if err != nil {
			return nil, err
		}
		SetEncounteredLangs(resultBylang)
	}
	return GetEncounteredLangs(), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: tally <directory>")
		return
	}

	root := os.Args[1]
	totalLinesPerLangs, err := tallyCodebase(root)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	sortByTotalLines(totalLinesPerLangs)	
	// TODO: use charmbracelet glow or similar for better readability
	fmt.Printf("%-20s %s\n", "Language", "Total Lines")
	fmt.Println(strings.Repeat("-", 30))
	for _, lang := range totalLinesPerLangs {
		fmt.Printf("%-20s %d\n", lang.Name, lang.TotalCount)
	}
	fmt.Println(strings.Repeat("-", 30))
	
}

func sortByTotalLines(langs []Lang) {
	for i := 0; i < len(langs); i++ {
		for j := i + 1; j < len(langs); j++ {
			if langs[i].TotalCount < langs[j].TotalCount {
				langs[i], langs[j] = langs[j], langs[i]
			}
		}
	}
}