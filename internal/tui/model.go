// Package tui defines the terminal user interface for the calendar application using the Bubble Tea framework.
package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/jermartinz/calendar-cli/internal/calendar"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("240"))

	cellStyle = lipgloss.NewStyle()

	todayStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("229")).
			Background(lipgloss.Color("110")).
			Bold(true)

	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("222"))

	borderStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240"))
)

type model struct {
	year  int
	month time.Month
}

func InitialModel() model {
	now := time.Now()
	return model{
		year:  now.Year(),
		month: now.Month(),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "left":
			m.month--
			if m.month < time.January {
				m.month = time.December
				m.year--
			}
		case "right":
			m.month++
			if m.month > time.December {
				m.month = time.January
				m.year++
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	titleText := fmt.Sprintf("%s %d", m.month.String(), m.year)
	calendarWidth := 56 // 7 days × 8 characters = 56
	title := titleStyle.Width(calendarWidth).Align(lipgloss.Center).Render(titleText)

	var b strings.Builder

	// Each day has 8 characters: 2 spaces + name (3 chars) + 3 spaces
	days := []string{"  Sun   ", "  Mon   ", "  Tue   ", "  Wed   ", "  Thu   ", "  Fri   ", "  Sat   "}
	for _, day := range days {
		b.WriteString(headerStyle.Render(day))
	}
	b.WriteString("\n\n")

	// Calendar Gen
	weeks := calendar.CalendarGen(m.year, m.month)

	for i, week := range weeks {
		for _, cell := range week {
			// Center numbers in 8 characters
			// For 1 digit: 3 spaces + number + 4 spaces = 8
			// For 2 digits: 3 spaces + number + 3 spaces = 8
			var cellText string
			if cell.Day < 10 {
				cellText = fmt.Sprintf("   %d    ", cell.Day)
			} else {
				cellText = fmt.Sprintf("   %d   ", cell.Day)
			}

			if cell.IsToday {
				b.WriteString(todayStyle.Render(cellText))
			} else if !cell.IsCurrentMonth {
				b.WriteString(dimStyle.Render(cellText))
			} else {
				b.WriteString(cellStyle.Render(cellText))
			}
		}
		b.WriteString("\n")
		// Add empty line between weeks for extra height
		if i < len(weeks)-1 {
			b.WriteString("\n")
		}
	}

	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Render("\n← → change month | q exit")

	return title + "\n" + borderStyle.Render(b.String()) + help + "\n"
}
