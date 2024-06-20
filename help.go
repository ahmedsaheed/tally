package main

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode/utf8"
)

func GetAllLanguagesInDir(root string) []Language {
	alreadyInserted := map[string]bool{}
	availableLangs := []Language{}

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
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

	return availableLangs
}

func sortByTotalLines(langs []Language) []Language {
	for i := 0; i < len(langs); i++ {
		for j := i + 1; j < len(langs); j++ {
			if langs[i].TotalCount < langs[j].TotalCount {
				langs[i], langs[j] = langs[j], langs[i]
			}
		}
	}
	return langs
}

func SumLines(langs []Language) int {
	total := 0
	for _, lang := range langs {
		total += lang.TotalCount
	}
	return total
}

func NumberToString(n int, sep rune) string {

	s := strconv.Itoa(n)

	startOffset := 0
	if n < 0 {
		startOffset = 1
	}

	const groupLen = 3
	groups := (len(s) - startOffset - 1) / groupLen

	if groups == 0 {
		return s
	}

	sepLen := utf8.RuneLen(sep)
	sepBytes := make([]byte, sepLen)
	_ = utf8.EncodeRune(sepBytes, sep)

	buf := make([]byte, groups*(groupLen+sepLen)+len(s)-(groups*groupLen))

	startOffset += groupLen
	p := len(s)
	q := len(buf)
	for p > startOffset {
		p -= groupLen
		q -= groupLen
		copy(buf[q:q+groupLen], s[p:])
		q -= sepLen
		copy(buf[q:], sepBytes)
	}
	if q > 0 {
		copy(buf[:q], s)
	}
	return string(buf)
}

func generateLanguageColorFromLanguageColorMap(lang Language) string {
	if color, ok := LanguageColors[lang.Name]; ok {
		return color
	}
	return generateRandomAnsiColor()
}

func generateRandomAnsiColor() string {
	return fmt.Sprintf("%d", 16+randomInt(0, 199))
}

func randomInt(min, max int) int {
	return min + rand.Intn(max-min)
}

func fileExtensionFromPath(path string) string {
	return "." + strings.TrimPrefix(filepath.Ext(path), ".")
}

func isPathOk(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		fmt.Println("Error:", err)
		return false
	}
	if !stat.IsDir() {
		fmt.Printf("Error %s is not a directory\n", stat.Name())
		return false
	}
	return true
}

func getWD() string {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	return currentDir
}

func resolveRootDirectoryFromArgs(args []string) string {
	if len(args) == 1 {
		return getWD()
	}
	return args[1]
}

func Tally(args []string) {
	argsLength := len(args)
	if argsLength == 1 || argsLength == 2 {
		ROOT := resolveRootDirectoryFromArgs(args)
		if !isPathOk(ROOT) {
			return
		}
		talliedDir, err := TallyDirectory(ROOT)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		BuildTable(talliedDir)
	} else {
		fmt.Println("Usage: tally <directory>")
		return
	}
}
