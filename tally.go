package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/charmbracelet/glamour"
)

//go:embed help.md
var helpText string

func Tally(args []string) {
	root, option, err := parseArgs(args)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if option.showHelp {
		formattedHelp, err := glamour.Render(helpText, "dark")
		if err != nil {
			fmt.Println("Error:", err)
		}
		fmt.Println(formattedHelp)
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
	if argsLength < 1 || argsLength > 4 {
		return "", Option{}, fmt.Errorf("usage: tally <directory|repo> [--blame | --remote | --html | --help]")
	}

	root := resolveRootDirectoryFromArgs(args)
	option := *NewOption()

	for _, arg := range args[1:] {
		if arg == "--help" || arg == "-h" {
			option.showHelp = true
			return "", option, nil
		}
	}

	if argsLength >= 3 && (contains(args[1:], "--remote") || contains(args[1:], "-r")) {
		option.remote = true
	}
	if argsLength >= 3 && contains(args[1:], "--html") {
		option.html = true
	}
	if argsLength >= 3 && (contains(args[1:], "--blame") || contains(args[1:], "-b")) {
		option.blame = true
	}

	return root, option, nil
}

// Helper function to check if a flag is in the arguments
func contains(args []string, flag string) bool {
	for _, arg := range args {
		if arg == flag {
			return true
		}
	}
	return false
}

func processLocalDirectory(root string, option Option) {
	if !isPathOk(root) {
		return
	}
	// rootBase := filepath.Base(root)
	languages, err := TallyDirectory(root)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if option.html {
		HtmlDisplay(root, languages, option)
		return
	}
	Display(languages, option)
}

func processRemoteRepo(repoName string, option Option) {
	loc, err := TallyRemoteRepo(repoName)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	languages := generateLangArrayFromLocByLangs(loc)

	if option.html {
		HtmlDisplay(repoName, languages, option)
		return
	}

	Display(languages, option)
}

func countLine(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

func Scan(path string) (int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	return countLine(reader)
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

func countLanguageInDir(lang Language, root string) (Language, error) {
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
		lang, err := countLanguageInDir(language, root)
		if err != nil {
			return recognisedLanguages, err
		}
		if lang.TotalCount != 0 {
			SetEncounteredLangs(lang)
		}
	}
	return sortByTotalLines(GetEncounteredLangs()), nil
}
