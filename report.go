package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/gookit/color"
	"golang.org/x/crypto/ssh/terminal"
)

//go:embed styles.css
var css string

func minimalColorise(langColor string) color.RGBColor {
	return color.HEX(langColor)

}
func generateBar(summaries []Language, width int) {
	innerWidth := width
	totalLines := 0
	totalLines = SumLines(summaries)

	if totalLines == 0 {
		fmt.Fprintf(os.Stderr, " no code found in %s\n", "/path/to/dir")
		return
	}

	filled := 0

	for _, summary := range summaries {
		percent := (summary.TotalCount * innerWidth) / totalLines
		if percent == 0 {
			continue
		}

		if filled == 0 {
			fmt.Println()
			fmt.Print(" ")
		}
		filled += percent

		if strings.HasPrefix(summary.getColor(), "#") {
			col := color.HEX(summary.getColor(), true)
			col.Print(strings.Repeat(" ", percent))
		}
	}

	if filled != 0 {
		fmt.Print(strings.Repeat(" ", innerWidth-filled))
		fmt.Println()
		fmt.Println()
	}
}

func Display(langs []Language, opts Option) {
	width, _, err := terminal.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		width = 80
	}

	fmt.Println()
	for _, lang := range langs {
		space := width - (len(lang.Name) + 4) - len(fmt.Sprintf("%d", lang.TotalCount))
		inlay := color.HEX("#a1a1a1").Sprint(strings.Repeat(".", space-2))
		format := map[bool]string{true: " %s %s%s %d", false: " %s  %s%s %d\n"}[opts.blame]
		fmt.Printf(format, minimalColorise(lang.getColor()).Sprint("●"), lang.Name, inlay, lang.TotalCount)
		if opts.blame {
			for i, file := range lang.Files {
				graphChar := "└"
				if i < len(lang.Files)-1 {
					graphChar = "├"
				}
				fmt.Printf("\n %s %s", graphChar, file.Name())
				if i == len(lang.Files)-1 {
					fmt.Println()
				}

			}

		}
	}
	generateBar(langs, width)
}

func HtmlDisplay(root string, langs []Language, opts Option) {

	var output bytes.Buffer
	openedInBrowser := false
	root = map[bool]string{true: root, false: filepath.Base(root)}[opts.html]
	output.WriteString("<!doctype html>\n")
	output.WriteString(fmt.Sprintf("<html>\n<head>\n<title>Tally - %s</title>\n<style>\n%s</style>\n</head>\n \n", root, css))
	output.WriteString("<body>\n\n")

	var totalLines int
	for _, lang := range langs {
		totalLines += lang.TotalCount
	}
	var remainingLines = totalLines
	output.WriteString("<div aria-hidden class='bar'>")
	for _, lang := range langs {
		remainingLines -= lang.TotalCount
		output.WriteString(fmt.Sprintf("<div aria-hidden title=%s style=\"background-color:%s; flex-grow: %d\"></div>", lang.Name, lang.getColor(), lang.TotalCount))
	}

	if remainingLines > 0 {
		output.WriteString(fmt.Sprintf("<div aria-hidden title=\"Other languages\" style=\"background-color: gray; flex-grow: %d\"></div>", remainingLines))
	}
	output.WriteString("</div>")
	output.WriteString("<table>\n<colgroup><col /><col width=\"15%\" /><col width=\"15%\" /></colgroup>\n <th>Language</th><th>Lines</th><th>File Count</th>")

	for _, lang := range langs {
		output.WriteString(fmt.Sprintf("<tr><td><span style=\"color: %s\">●</span>&nbsp;%s</td><td>%d</td><td>%d</td></tr>",
			lang.getColor(), lang.Name, lang.TotalCount, lang.FileCount))
	}

	output.WriteString("</table>\n</body>\n</html>")
	tempFile := temp(output.String())
	if tempFile == nil {
		return
	}

	if err := openBrowser(tempFile.Name()); err != nil {
		fmt.Println("Error opening browser:", err)
		defer os.Remove(tempFile.Name())
	} else {
		openedInBrowser = true
	}

	if openedInBrowser {
		defer func() {
			time.Sleep(10 * time.Second)
			os.Remove(tempFile.Name())
			fmt.Println("Deleted temporary file:", tempFile.Name())
		}()
	}
}

func temp(in string) *os.File {
	tempFile, err := os.CreateTemp("", "tally-*.html")
	if err != nil {
		fmt.Println("Error creating temporary file:", err)
		return nil
	}
	defer tempFile.Close()

	_, err = tempFile.WriteString(in)
	if err != nil {
		fmt.Println("Error writing to temporary file:", err)
		return nil
	}

	return tempFile
}

func openBrowser(filePath string) error {
	var cmd string
	var args []string
	switch runtime.GOOS {
	case "windows":
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", filePath}
	case "darwin":
		cmd = "open"
		args = []string{filePath}
	default:
		cmd = "xdg-open"
		args = []string{filePath}
	}

	return exec.Command(cmd, args...).Start()
}
