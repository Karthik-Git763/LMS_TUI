package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	viewportTitleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		
		return lipgloss.NewStyle().Border(b).Padding(0, 1).Margin(0, 1)	
	}()
	
	contentStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		
		return lipgloss.NewStyle().Border(b).Padding(0, 1).Margin(0, 1)
	}()
)

func (m *Model) HeaderView() string {
	title := viewportTitleStyle.Render("LMS Terminal")
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

// func (m *Model) FooterView() string {
// 	info := contentStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
// 	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))
// 	return lipgloss.JoinHorizontal(lipgloss.Left, line, info)
// }