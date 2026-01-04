package tui

import (
	"fmt"
	"lms/internal/scrapper"

	tea "github.com/charmbracelet/bubbletea"
)

type Screen int

const (
	menuScreen Screen = iota
	attendanceScreen
	assignmentsScreen
)

type Model struct {
	screen 		Screen
	cursor 		int
	courses 	[]scrapper.Course
	assignments []scrapper.Assignment
	attendance 	map[string]string
}

func InitialModel(courses []scrapper.Course) Model {
	return Model {
		screen: menuScreen,
		courses: courses,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
				case "ctrl+c", "q":
					return m, tea.Quit
				case "up", "k":
					if m.cursor > 0 {
						m.cursor--
					}
				case "down", "j":
					if m.cursor < len(m.courses)-1 {
						m.cursor++
					}
				case "enter":
					switch m.screen {
						case 0:
							m.screen = attendanceScreen
						case 1:
							m.screen = assignmentsScreen
						case 2:
							return m, tea.Quit
					}
			}
	}
	return m, nil
}

func (m Model) View() string {
	switch m.screen {
		case menuScreen:
			return m.MenuView()
		case attendanceScreen:
			return m.AttendanceView()
		case assignmentsScreen:
			return m.AssignmentsView()
	}
	return ""
}

func (m Model) MenuView() string {
	options := []string{"Attendance", "Assignments", "Exit"}
	
	s := "LMS Terminal Client\n\n"
	
	for i, option := range options {
		cursor := " "
		if i == m.cursor {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, option)
	}
	s += "\n↑/↓ or j/k to move • Enter to select • q to quit\n"
	return s
}

func (m Model) AttendanceView() string {
	s := "Attendance\n\n"
	
	for _, course := range m.courses {
		s += fmt.Sprintf("%s\n", course.Name)
	}
	s += "\nPress q to quit\n"
	return s
}

func (m Model) AssignmentsView() string {
	s := "Assignments\n\n"
	
	for _, assignment :=  range m.assignments {
		s += fmt.Sprintf("%s\n", assignment.Title)
	}
	s += "\nPress q to quit\n"
	return s
}