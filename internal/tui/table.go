package tui

import (
	"lms/internal/models"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

func AttendanceTable(records []models.Attendance) table.Model {
	columns := []table.Column{
		{Title: "Date", Width: 30},
		{Title: "Session", Width: 25},
		{Title: "Status", Width: 10},
	}
	rows := []table.Row{}
	for _, attendance := range records {
		rows = append(rows, table.Row{
			attendance.Date,
			attendance.Session,
			attendance.Status,
		})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return t
}
