package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gookit/color"
	"golang.org/x/crypto/ssh/terminal"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	table table.Model
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			return m, tea.Batch(
				tea.Printf("Let's go to %s!", m.table.SelectedRow()[0]),
			)
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return baseStyle.Render(m.table.View()) + "\n"
}

func generateTableRowFromLanguageArray(langs []Language) []table.Row {
	rows := []table.Row{}
	rowStyle := lipgloss.NewStyle()
	for _, lang := range langs {
		bulletWithColor := coloriser(lang.getColor())
		rows = append(rows, table.Row{
			rowStyle.Align(lipgloss.Left).Inline(true).Render(bulletWithColor.Render() + lang.Name),
			rowStyle.Align(lipgloss.Right).Inline(true).Render(strconv.Itoa(lang.FileCount)),
			rowStyle.Align(lipgloss.Right).Inline(true).Render(NumberToString(lang.TotalCount, ',')),
		})
	}
	return rows
}

func coloriser(color string) lipgloss.Style {
	bullet := "● " //• ● ○ ◌ ◍ ◎ ◉ ○ ◌ ◍ ◎ ◉ ● ●
	return lipgloss.NewStyle().Background(lipgloss.NoColor{}).
		Inline(true).SetString(bullet).Foreground(lipgloss.Color(color))
}

func minimalColorise(langColor string) color.RGBColor {
	return color.HEX(langColor)

}

func BuildTable(languages []Language) {
	columns := []table.Column{
		{Title: "Language", Width: 40},
		{Title: "files", Width: 5},
		{Title: "lines", Width: 20},
	}

	rows := generateTableRowFromLanguageArray(languages)

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(6),
	)
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true)
	s.Selected = s.Selected.
		// Background(lipgloss.Color("12")).
		Foreground(lipgloss.Color("#f1e05")).
		Bold(false)
	t.SetStyles(s)
	m := model{t}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
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
