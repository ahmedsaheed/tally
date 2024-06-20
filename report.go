package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/gookit/color"
	"golang.org/x/crypto/ssh/terminal"
)

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

func MinimalDisplay(langs []Language, opts Option) {
	width, _, err := terminal.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		width = 80
	}

	fmt.Println()
	for _, lang := range langs {
		space := width - (len(lang.Name) + 4) - len(fmt.Sprintf("%d", lang.TotalCount))
		inlay := color.HEX("#a1a1a1").Sprint(strings.Repeat(".", space-2))
		format := map[bool]string{true: " %s %s%s %d", false: " %s %s%s %d\n"}[opts.blame]
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
