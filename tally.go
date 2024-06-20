package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

func Tally(args []string) {
	argsLength := len(args)
	if argsLength < 1 || argsLength > 3 {
		fmt.Println("Usage: tally <directory> [--blame]")
		return
	}

	root := resolveRootDirectoryFromArgs(args)
	if !isPathOk(root) {
		return
	}
	var option Option = *NewOption()
	if argsLength == 3 && args[2] == "--blame" || argsLength == 2 && args[1] == "--blame" {
		option.blame = true
	}
	talliedDir, err := TallyDirectory(root)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	MinimalDisplay(talliedDir, option)
	// BuildTable(talliedDir)

}

func Scan(path string) (int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
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
}

func countUp(lang Language, root string) (Language, error) {
	totalLine := 0
	fileCount := 0
	filesInfo := []os.FileInfo{}
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && hasExpectedLanguage(path, lang.Extensions) && !isIgnoredFile(path) && !isIgnoredDir(path) && !isDotedDir(path) {
			line, err := Scan(path)
			fileCount++
			filesInfo = append(filesInfo, info)
			if err != nil {
				return err
			}

			totalLine += line
		}

		lang.TotalCount = totalLine
		lang.FileCount = fileCount
		lang.Files = filesInfo

		return nil
	})

	if err != nil {
		return lang, err
	}
	return lang, nil
}

func TallyDirectory(root string) ([]Language, error) {
	recognisedLanguages := GetAllLanguagesInDir(root)
	for _, language := range recognisedLanguages {
		lang, err := countUp(language, root)
		if err != nil {
			return recognisedLanguages, err
		}
		if lang.TotalCount != 0 {
			SetEncounteredLangs(lang)
		}
	}
	return sortByTotalLines(GetEncounteredLangs()), nil
}
