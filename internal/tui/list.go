package tui

import (
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

var (
	// Title styling - blue/purple background
	titleStyle = lipgloss.NewStyle().
		Background(lipgloss.Color("62")).
		Foreground(lipgloss.Color("255")).
		Bold(true).
		Padding(0, 1)
	
	// Selected item styling - pink/magenta left border
	selectedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("170")).
		BorderLeft(true).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("170")).
		PaddingLeft(1)
	
	// Normal item styling
	normalStyle = lipgloss.NewStyle().
		PaddingLeft(2)
	
	// Description styling - gray color
	descStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#444444"))
)

type Item struct {
	title       string
	desc string
}

func (i Item) Title() string {return i.title}
func (i Item) Description() string {return i.desc}
func (i Item) FilterValue() string {return i.title}