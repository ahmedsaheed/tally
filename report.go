package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
		bulletWithColor := coloriser(generateLanguageColorFromLanguageColorMap(lang))
		rows = append(rows, table.Row{
			rowStyle.Align(lipgloss.Left).Inline(true).Render(bulletWithColor.Render() + lang.Name),
			rowStyle.Align(lipgloss.Right).Inline(true).Render(strconv.Itoa(lang.FileCount)),
			rowStyle.Align(lipgloss.Right).Inline(true).Render(NumberToString(lang.TotalCount, ',')),
		})
	}
	return rows
}

func coloriser(color string) lipgloss.Style {
	bullet := "● "
	// diffBullets := "• ● ○ ◌ ◍ ◎ ◉ ○ ◌ ◍ ◎ ◉"
	return lipgloss.NewStyle().Background(lipgloss.NoColor{}).
		Inline(true).SetString(bullet)
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
		Background(lipgloss.Color("12")).
		Foreground(lipgloss.Color("#f1e05")).
		Bold(false)
	t.SetStyles(s)
	m := model{t}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
