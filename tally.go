package main

import (
	"bufio"
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gookit/color"
)

//go:embed help.txt
var helpText string

func Tally(args []string) {
	root, option, err := parseArgs(args)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if option.showHelp {
		fmt.Println(color.HEX("#00FFFF", false).Sprint("Tally 0.0.1"))
		fmt.Println(helpText)
		return
	}

	if option.remote {
		processRemoteRepo(root, option)
	} else {
		processLocalDirectory(root, option)
	}
}

func parseArgs(args []string) (string, Option, error) {
	argsLength := len(args)
	if argsLength < 1 || argsLength > 3 {
		return "", Option{}, fmt.Errorf("usage: tally <directory|repo> [--blame | --remote]")
	}
	root, option := resolveRootDirectoryFromArgs(args), *NewOption()

	if argsLength == 2 && args[1] == "--help" || argsLength == 3 && args[2] == "--help" {
		option.showHelp = true
		return "", option, nil
	}

	if argsLength == 3 && args[2] == "--remote" {
		option.remote = true
		return args[1], option, nil
	}
	if !isPathOk(root) {
		return "", Option{}, fmt.Errorf("error: invalid path")
	}
	if argsLength == 3 && args[2] == "--blame" || argsLength == 2 && args[1] == "--blame" {
		option.blame = true
	}
	return root, option, nil
}

func processLocalDirectory(root string, option Option) {
	if !isPathOk(root) {
		return
	}
	languages, err := TallyDirectory(root)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	MinimalDisplay(languages, option)
}

func processRemoteRepo(repoName string, option Option) {
	loc, err := TallyRemoteRepo(repoName)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	languages := generateLangArrayFromLocByLangs(loc)
	MinimalDisplay(languages, option)
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

type RemoteResponse struct {
	LOC        int            `json:"loc"`
	LOCByLangs map[string]int `json:"locByLangs"`
	Children   map[string]any `json:"children"`
}

func TallyRemoteRepo(repoName string) (map[string]int, error) {
	url := "https://ghloc.ifels.dev/" + repoName + "?match=!package-lock.json&pretty=false"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		var data RemoteResponse
		json.NewDecoder(resp.Body).Decode(&data)
		return data.LOCByLangs, nil
	}
	reason := map[bool]string{
		true:  "repository is either private or does not exist.",
		false: "unknown error occurred.",
	}[resp.StatusCode == 404 || resp.StatusCode == 400]
	return nil, fmt.Errorf("%s - %s", resp.Status, reason)
}

func generateLangArrayFromLocByLangs(locByLangs map[string]int) []Language {
	langs := []Language{}
	for ext, loc := range locByLangs {
		if lang, ok := Languages[lookupLangByExtension(ext)]; ok {
			lang.TotalCount = loc
			langs = append(langs, lang)
		}

	}
	return sortByTotalLines(langs)

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
