package tui

import (
	"fmt"
	"lms/internal/models"
	"lms/internal/scrapper"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Screen int

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type Model struct {
	screen 				Screen
	cursor 				int
	courses 			[]models.Course
	selectedCourse 		models.Course
	assignmentsByCourse map[string][]models.Assignment
	attendanceByCourse 	map[string][]models.Attendance
	vplsByCourse 		map[string][]models.VPL
	selectedURL 		string
	spinner				spinner.Model
	active				map[string]bool
	done				map[string]bool
	errors 				map[string]error
	table				table.Model
}


const (
	menuScreen Screen = iota
	progressScreen
	attendanceCourseScreen
	attendanceDetailsScreen
	assignmentCourseScreen
	assignmentDetailsScreen
	vplCourseScreen
	vplDetailsScreen
)


func InitialModel(courses []models.Course) Model {
	return Model{
		screen:              progressScreen,
		courses:             courses,
		attendanceByCourse:  make(map[string][]models.Attendance),
		assignmentsByCourse: make(map[string][]models.Assignment),
		vplsByCourse:         make(map[string][]models.VPL),
		spinner:			 newSpinner(),
		active:				 make(map[string]bool),
		done:				 make(map[string]bool),
		errors:				 make(map[string]error),
	}
}

func (m Model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case models.ProgressMsg:
		m.applyProgress(msg)
	case models.DataLoadedMsg:
		m.applyDataLoaded(msg)
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "o":
			if m.selectedURL != "" {
				models.OpenBrowser(m.selectedURL)
			}
		case "q":
			switch m.screen {
			case attendanceDetailsScreen:
				m.screen = attendanceCourseScreen
			case assignmentDetailsScreen:
				m.screen = assignmentCourseScreen
			case vplDetailsScreen:
				m.screen = vplCourseScreen
			default:
				m.screen = menuScreen
				m.cursor = 0
			}
		case "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			switch m.screen {
			case menuScreen, attendanceCourseScreen, assignmentCourseScreen, vplCourseScreen:
				if m.cursor > 0 {
					m.cursor--
				}
			case assignmentDetailsScreen:
				if m.cursor > 0 {
					m.cursor--
				}
			case vplDetailsScreen:
				if m.cursor > 0 {
					m.cursor--
				}
			case attendanceDetailsScreen:
				// no navigation needed
			}
		case "down", "j":
			switch m.screen {
			case menuScreen:
				if m.cursor < 3 {
					m.cursor++
				}
		 	case attendanceCourseScreen, assignmentCourseScreen, vplCourseScreen:
				if m.cursor < len(m.courses)-1 {
					m.cursor++
				}
			case assignmentDetailsScreen:
				assignments := m.assignmentsByCourse[m.selectedCourse.Name]
				if m.cursor < len(assignments)-1 {
					m.cursor++
				}
			case vplDetailsScreen:
				vpls := m.vplsByCourse[m.selectedCourse.Name]
				if m.cursor < len(vpls)-1 {
					m.cursor++
				}
			case attendanceDetailsScreen:
				// no navigation needed
			}
		case "enter":
			switch m.screen {
			case menuScreen:
				if m.cursor == 0 {
					m.screen = attendanceCourseScreen
				}
				if m.cursor == 1 {
					m.screen = assignmentCourseScreen
				}
				if m.cursor == 2 {
					m.screen = vplCourseScreen
				}
				if m.cursor == 3 {
					return m, tea.Quit
				}
			case attendanceCourseScreen:
				m.selectedCourse = m.courses[m.cursor]
				records := m.attendanceByCourse[m.selectedCourse.Name]
				if len(records) > 0 {
					m.selectedURL = records[0].AttendanceURL
				}
				// Build a table only for the selected course
				m.table = AttendanceTable(records)
				m.screen = attendanceDetailsScreen
				m.cursor = 0
			case assignmentCourseScreen:
				m.selectedCourse = m.courses[m.cursor]
				m.screen = assignmentDetailsScreen
				m.cursor = 0
			case vplCourseScreen:
				m.selectedCourse = m.courses[m.cursor]
				m.screen = vplDetailsScreen
				m.cursor = 0
			}
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	switch m.screen {
		case progressScreen:
			return fmt.Sprintf("LOADING COURSES%s", m.spinner.View())
		case menuScreen:
			return m.MenuView()
		case attendanceCourseScreen:
			return m.AttendanceCourseView()
		case attendanceDetailsScreen:
			return m.AttendanceDetailsView()
		case assignmentCourseScreen:
			return m.AssignmentCourseView()
		case assignmentDetailsScreen:
			return m.AssignmentDetailView()
		case vplCourseScreen:
			return m.VPLCourseView()
		case vplDetailsScreen:
			return m.VPLDetailsView()
	}
	return baseStyle.Render(m.table.View()) + "\n"
}

func (m Model) MenuView() string {
	options := []string{"Attendance", "Assignments", "VPL", "Exit"}
	
	s := "LMS Terminal Client\n\n"
	
	for i, option := range options {
		cursor := " "
		if i == m.cursor {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, option)
	}
	s += "\n↑/↓ or j/k to move • Enter to select • Ctrl+c to exit\n"
	return s
}

func (m Model) AttendanceCourseView() string {
	s := "Attendance -> Select Course\n\n"
	
	for i, course := range m.courses {
		cursor := ""
		if m.cursor == i {
			cursor = ">"
		}
		s += cursor + " " + course.Name + "\n" 
	}
	s += "\n↑/↓ navigate • Enter select • q to back • Ctrl+c to exit\n"
	return s
}

func (m Model) AttendanceDetailsView() string {
	records := m.attendanceByCourse[m.selectedCourse.Name]
	
	s := "Attendance -> " + m.selectedCourse.Name + "\n\n"
	
	if len(records) == 0 {
		s += "No attendance records found for this course.\n\n"
		s += "q to back • Ctrl+c to exit\n"
		return s
	}
	
	attended, total := scrapper.CalculateAttendancePercentage(records)
	percent := float64(attended) / float64(total) * 100
	warning := ""
	if percent < 80 {
		warning = "Warning: Low attendance!"
	}
	// Render the prepared table instead of plain lines
	s += baseStyle.Render(m.table.View()) + "\n"
	
	s += fmt.Sprintf("Overall: %d/%d (%.2f%%)\n%s\n", attended, total, percent, warning)
	if len(records) > 0 {
		m.selectedURL = records[0].AttendanceURL
		s += fmt.Sprintf("\nq to back • Ctrl+c to exit • o open in browser • %s\n", m.selectedURL)
	}
	return s
}

func (m Model) AssignmentCourseView() string {
	s := "Assignments -> Select Course\n\n"
	
	for i, course := range m.courses {
		cursor := ""
		if i == m.cursor {
			cursor = ">"
		}
		s += cursor + " " + course.Name + "\n"
	}
	s += "↑/↓ navigate • Enter select • q to back • Ctrl+c to exit\n"
	return s
}

func (m Model) AssignmentDetailView() string {
	assignments := m.assignmentsByCourse[m.selectedCourse.Name]
	
	s := "Assignments -> " + m.selectedCourse.Name + "\n\n"
	
	for i, assignment := range assignments {
		cursor := ""
		if i == m.cursor {
			cursor = ">"
			// m.selectedURL = assignment.URL
		}
		s += cursor + " " + assignment.Title + "\nDue: " + assignment.DueDate + "\nStatus: " + assignment.Status + "\nGrade: " + assignment.Grade + "\n"
	}
	s += "\nq to back • Ctrl+c to exit • o open in browser\n"
	return s
}

func (m Model) VPLCourseView() string {
	s := "VPL -> Select Course\n\n"
	
	for i, course := range m.courses {
		cursor := ""
		if i == m.cursor {
			cursor = ">"
		}
		s += cursor + " " + course.Name + "\n"
	}
	s += "↑/↓ navigate • Enter select • q to back • Ctrl+c to exit\n"
	return s
}

func (m Model) VPLDetailsView() string {
	vpls := m.vplsByCourse[m.selectedCourse.Name]
	
	s := "VPL -> " + m.selectedCourse.Name + "\n\n"
	
	for i, vpl := range vpls {
		cursor := ""
		if i == m.cursor {
			cursor = ">"
			m.selectedURL = vpl.URL
		}
		s += cursor + " " + vpl.Title + "\nDue: " + vpl.DueDate + "\n"
	}
	s += "\nq to back • Ctrl+c to exit • o open in browser\n"
	return s
}

