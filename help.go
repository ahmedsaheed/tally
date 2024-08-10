package main

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
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
			if isDockerfile(path) && ext == "." {
				langName = "Dockerfile"
				lang = Languages[langName]
			}
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
	sort.Slice(langs, func(i, j int) bool {
		return langs[i].TotalCount > langs[j].TotalCount
	})
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

func generateRandomAnsiColor() string {
	return fmt.Sprintf("%d", 16+randomInt(0, 199))
}

func randomInt(min, max int) int {
	return min + rand.Intn(max-min)
}

func isDockerfile(file string) bool { return strings.HasSuffix(file, "Dockerfile") }

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
	if len(args) == 1 || args[1] == "--blame" {
		return getWD()
	}
	return args[1]
}
