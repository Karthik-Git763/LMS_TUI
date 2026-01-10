package tui

import (
	"fmt"
	"net/http"

	"lms/internal/auth"
	"lms/internal/models"
	"lms/internal/scrapper"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
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
	content 			string
	ready				bool
	viewport 			viewport.Model
	courselist 			list.Model
	assignmentList 		list.Model
	vplList				list.Model
	focusIndex			int
	inputs				[]textinput.Model
	cursorMode			cursor.Mode
	authError			string
	client				*http.Client
	baseURL				string
}


const (
	progressScreen Screen = iota
	inputScreen
	menuScreen
	attendanceCourseScreen
	attendanceDetailsScreen
	assignmentCourseScreen
	assignmentDetailsScreen
	vplCourseScreen
	vplDetailsScreen
)

func InitialModel(courses []models.Course, p *tea.Program, client *http.Client, baseURL string) Model {
	delegate := list.NewDefaultDelegate()
	
	courselist := list.New([]list.Item{}, delegate, 0, 0)
	courselist.SetShowHelp(true)
	courselist.SetShowStatusBar(true)
	courselist.SetFilteringEnabled(true)
	courselist.SetShowPagination(true)
	courselist.Paginator.ActiveDot = "●"
	courselist.Paginator.InactiveDot = "○"
	courselist.SetShowTitle(true)
	DisableListQuit(&courselist)
	
	// Create assignment list with help and status enabled
	assignmentList := list.New([]list.Item{}, delegate, 0, 0)
	assignmentList.SetShowHelp(true)
	assignmentList.SetFilteringEnabled(true)
	assignmentList.SetShowPagination(true)
	assignmentList.Paginator.ActiveDot = "●"
	assignmentList.Paginator.InactiveDot = "○"
	assignmentList.SetShowTitle(true)
	DisableListQuit(&assignmentList)
	
	// Create VPL list with help and status enabled
	vplList := list.New([]list.Item{}, delegate, 0, 0)
	vplList.SetShowHelp(true)
	vplList.SetShowStatusBar(true)
	vplList.SetFilteringEnabled(true)
	vplList.SetShowPagination(true)
	vplList.Paginator.ActiveDot = "●"
	vplList.Paginator.InactiveDot = "○"
	vplList.SetShowTitle(true)
	DisableListQuit(&vplList)
	
	inputs := make([]textinput.Model, 2)
	var t textinput.Model
	for i := range inputs {
		t =  textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 100
		t.Width = 60
		
		switch i {
			case 0:
				t.Prompt = "Username: "
				t.Placeholder = "Enter your username"
				t.Focus()
				t.PromptStyle = focusedStyle
				t.TextStyle = focusedStyle
			case 1:
				t.Prompt = "Password: "
				t.Placeholder = "Enter your password"
				t.EchoMode = textinput.EchoPassword
				t.EchoCharacter = '●'
				t.PromptStyle = blurredStyle
				t.TextStyle = blurredStyle
		}
		inputs[i] = t
	}
	
	return Model{
		screen:              inputScreen,
		courses:             courses,
		attendanceByCourse:  make(map[string][]models.Attendance),
		assignmentsByCourse: make(map[string][]models.Assignment),
		vplsByCourse:        make(map[string][]models.VPL),
		spinner:             newSpinner(),
		active:              make(map[string]bool),
		done:                make(map[string]bool),
		errors:              make(map[string]error),
		courselist:          courselist,
		assignmentList:      assignmentList,
		vplList:             vplList,
		inputs:              inputs,
		authError:           "",
		client:              client,
		baseURL:             baseURL,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, textinput.Blink)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case models.ProgressMsg:
		m.applyProgress(msg)
	case models.DataLoadedMsg:
		m.applyDataLoaded(msg)
	case models.Credentials:
		// Handle authentication with provided credentials
		return m, func() tea.Msg {
			return auth.AuthenticateAndFetchData(m.client, m.baseURL, msg.Username, msg.Password)
		}
	case models.AuthResultMsg:
		if msg.Success {
			m.courses = msg.Courses
			m.screen = progressScreen
			m.authError = ""
			m.selectedURL = ""
			// Start fetching data for all courses
			return m, scrapper.StartDataFetching(m.client, m.courses)
		} else {
			m.authError = msg.Error
			m.inputs[0].Reset()
			m.inputs[1].Reset()
			m.inputs[0].Focus()
			m.focusIndex = 0
			return m, nil
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "o", "O":
			if m.selectedURL != "" && (m.screen == attendanceDetailsScreen || m.screen == assignmentDetailsScreen || m.screen == vplDetailsScreen) {
				models.OpenBrowser(m.selectedURL)
			}
		case "q", "Q":
			if m.screen != inputScreen {
				m.selectedURL = ""
				switch m.screen {
				case attendanceCourseScreen:
					m.screen = menuScreen
					m.cursor = 0
				case attendanceDetailsScreen:
					m.screen = attendanceCourseScreen
				case assignmentCourseScreen:
					m.screen = menuScreen
					m.cursor = 0
				case assignmentDetailsScreen:
					m.screen = assignmentCourseScreen
				case vplCourseScreen:
					m.screen = menuScreen
					m.cursor = 0
				case vplDetailsScreen:
					m.screen = vplCourseScreen
				}
				return m, nil
			}
		case "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			switch m.screen {
			case menuScreen:
				if m.cursor > 0 {
					m.cursor--
				}
			case attendanceDetailsScreen, assignmentCourseScreen, vplCourseScreen, assignmentDetailsScreen, vplDetailsScreen, attendanceCourseScreen:
				// no navigation needed // List navigation
			}
		case "down", "j":
			switch m.screen {
			case menuScreen:
				if m.cursor < 3 {
					m.cursor++
				} 
			case attendanceDetailsScreen, assignmentDetailsScreen, vplDetailsScreen, assignmentCourseScreen, vplCourseScreen, attendanceCourseScreen:
				// no navigation needed
			}
		case "ctrl+r":
			m.cursorMode++
			if m.cursorMode > cursor.CursorHide {
				m.cursorMode = cursor.CursorBlink
			}
			cmds := make([]tea.Cmd, len(m.inputs))
			for i := range m.inputs {
				cmds[i] = m.inputs[i].Cursor.SetMode(m.cursorMode)
			}
			return m, tea.Batch(cmds...)
		case "tab", "shift+tab":
			s := msg.String()
			switch s {
				case "tab":
					m.focusIndex++
				case "shift+tab":
					m.focusIndex--
			}
			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}
			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i < len(m.inputs); i++ {
				if i == m.focusIndex {
					// Set the focused state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove the focused state
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = blurredStyle
				m.inputs[i].TextStyle = blurredStyle
			}
		case "enter":
			switch m.screen {
			case inputScreen:
				// Submit credentials and move to menu screen
				if m.focusIndex == len(m.inputs) {
					username := m.inputs[0].Value()
					password := m.inputs[1].Value()
					if username == "" || password == "" {
						m.inputs[0].Placeholder = "Username is required"
						m.inputs[1].Placeholder = "Password is required"
						m.inputs[0].Focus()
						return m, nil
					}
					return m, func() tea.Msg {
						return models.Credentials{Username: username, Password: password}
					}
				}
				
				m.cursor = 0
			case menuScreen:
				m.selectedURL = ""
				if m.cursor == 0 {
					items := m.AttendanceCourseView()
					m.courselist.SetItems(items)
					m.courselist.Title = "Select Course → Attendance"
					m.screen = attendanceCourseScreen
				}
				if m.cursor == 1 {
					items := m.AssignmentCourseView()
					m.courselist.SetItems(items)
					m.courselist.Title = "Select Course → Assignments"
					m.screen = assignmentCourseScreen
				}
				if m.cursor == 2 {
					// Initialize courselist for VPL selection
					items := m.VPLCourseView()
					m.courselist.SetItems(items)
					m.courselist.Title = "Select Course → VPLs"
					m.screen = vplCourseScreen
				}
				if m.cursor == 3 {
					return m, tea.Quit
				}
			case attendanceCourseScreen:
				// Get selected course from courselist
				if selectedItem, ok := m.courselist.SelectedItem().(Item); ok {
					for _, course := range m.courses {
						if course.Name == selectedItem.ItemTitle {
							m.selectedCourse = course
							break
						}
					}
				}
				records := m.attendanceByCourse[m.selectedCourse.Name]
				m.selectedURL = ""
				if len(records) > 0 {
					m.selectedURL = records[0].AttendanceURL
				}
				// Build a table only for the selected course
				m.table = AttendanceTable(records)
				m.screen = attendanceDetailsScreen
				m.cursor = 0
			case assignmentCourseScreen:
				// Get selected course from courselist
				if selectedItem, ok := m.courselist.SelectedItem().(Item); ok {
					for _, course := range m.courses {
						if course.Name == selectedItem.ItemTitle {
							m.selectedCourse = course
							break
						}
					}
				}
				items := m.AssignmentListView(m.selectedCourse.Name)
				m.assignmentList.SetItems(items)
				m.assignmentList.Title = fmt.Sprintf("Assignments for %s", m.selectedCourse.Name)
				m.selectedURL = ""
				if len(items) > 0 {
					m.assignmentList.Select(0)
					if first, ok := m.assignmentList.SelectedItem().(Item); ok {
						m.selectedURL = first.ItemUrl
					}
				}
				m.screen = assignmentDetailsScreen
				m.cursor = 0
			case vplCourseScreen:
				// Get selected course from courselist
				if selectedItem, ok := m.courselist.SelectedItem().(Item); ok {
					for _, course := range m.courses {
						if course.Name == selectedItem.ItemTitle {
							m.selectedCourse = course
							break
						}
					}
				}
				items := m.VPLListView(m.selectedCourse.Name)
				m.vplList.SetItems(items)
				m.vplList.Title = fmt.Sprintf("VPLs for %s", m.selectedCourse.Name)
				m.selectedURL = ""
				if len(items) > 0 {
					m.vplList.Select(0)
					if first, ok := m.vplList.SelectedItem().(Item); ok {
						m.selectedURL = first.ItemUrl
					}
				}
				m.screen = vplDetailsScreen
				m.cursor = 0
			}
			return m, tea.Batch(cmds...)
		}
		case tea.WindowSizeMsg:
			headerHeight := lipgloss.Height(m.HeaderView())
			verticalMarginHeight := 2 * headerHeight
			
			// For course and assignment/VPL lists, use full terminal size since they render directly
			m.courselist.SetSize(msg.Width, msg.Height-verticalMarginHeight)
			m.assignmentList.SetSize(msg.Width, msg.Height- verticalMarginHeight)
			m.vplList.SetSize(msg.Width, msg.Height-verticalMarginHeight)
			if !m.ready {
				m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
				m.viewport.YPosition = headerHeight
				m.ready = true
			} else {
				m.viewport.Width = msg.Width
				m.viewport.Height = msg.Height - verticalMarginHeight
			}
	}
	cmd = m.updateInputs(msg)
	cmds = append(cmds, cmd)
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)
	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)
	
	// Update List based on current screen
	switch m.screen {
		case attendanceCourseScreen, assignmentCourseScreen, vplCourseScreen:
			m.courselist, cmd = m.courselist.Update(msg)
			cmds = append(cmds, cmd)
		case assignmentDetailsScreen:
			m.assignmentList, cmd = m.assignmentList.Update(msg)
			cmds = append(cmds, cmd)
			selectedItem, ok := m.assignmentList.SelectedItem().(Item)
			if ok {
				m.selectedURL = selectedItem.ItemUrl
			} 
		case vplDetailsScreen:
			m.vplList, cmd = m.vplList.Update(msg)
			cmds = append(cmds, cmd)
			selectedItem, ok := m.vplList.SelectedItem().(Item)
			if ok {
				m.selectedURL = selectedItem.ItemUrl
			}
	}
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if !m.ready {
		return "\n Initializing..."
	}

	var body string

	switch m.screen {
	case progressScreen:
		body = fmt.Sprintf("LOADING COURSES%s", m.spinner.View())
	case inputScreen:
		var b strings.Builder
		b.WriteString("\n")
		b.WriteString(titleStyle.Render("Enter Your Credentials"))
		b.WriteString("\n\n")
			// Display auth error if present
			if m.authError != "" {
				b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render("Error: " + m.authError))
				b.WriteString("\n\n")
			}
	
		
		for i := range m.inputs {
			b.WriteString(m.inputs[i].View())
			if i < len(m.inputs)-1 {
				b.WriteRune('\n')
			}
		}
		button := &blurredButton
		if m.focusIndex == len(m.inputs) {
			button = &focusedButton
		}
		fmt.Fprintf(&b, "\n\n%s\n\n", *button)
		b.WriteString(helpStyle.Render("Cursor mode is "))
		b.WriteString(cursorModeHelpStyle.Render(m.cursorMode.String()))
		b.WriteString(helpStyle.Render(" (ctrl+r to change style)"))
		
		return m.HeaderView() + "\n" + b.String()
	case menuScreen:
		body = m.MenuView()
	case attendanceCourseScreen:
		return m.HeaderView() + "\n" + m.courselist.View()
	case attendanceDetailsScreen:
		body = m.AttendanceDetailsView()
	case assignmentCourseScreen:
		return m.HeaderView() + "\n" + m.courselist.View()
	case assignmentDetailsScreen:
		return m.HeaderView() + "\n" + m.assignmentList.View()
	case vplCourseScreen:
		return m.HeaderView() + "\n" + m.courselist.View()
	case vplDetailsScreen:
		return m.HeaderView() + "\n" + m.vplList.View()
	default:
		body = ""
	}

	m.viewport.SetContent(body)

	return fmt.Sprintf("%s\n%s\n", m.HeaderView(), m.viewport.View())
}

func (m Model) MenuView() string {
	options := []string{"Attendance", "Assignments", "VPL", "Exit"}
	
	s := titleStyle.Render("Select an option") + "\n\n"
	
	for i, option := range options {
		if i == m.cursor {
			s += selectedStyle.Render(option) + "\n"
		} else {
			s += normalStyle.Render(option) + "\n"
		}
	}
	s += "\n↑/↓ or j/k to move • Enter to select • Ctrl+c to exit\n"
	return s
}

func (m Model) AttendanceCourseView() []list.Item {
	items := make([]list.Item, len(m.courses))
	for i, course := range m.courses {
		// Calculate attendance for this course
		records := m.attendanceByCourse[course.Name]
		attended, total := scrapper.CalculateAttendancePercentage(records)
		var percent float64
		if total > 0 {
			percent = float64(attended) / float64(total) * 100
		}
		items[i] = Item{
			ItemTitle: course.Name,
			ItemDesc: fmt.Sprintf("Overall Attendance: %d/%d (%.2f%%)", attended, total, percent),
		}
	}
	return items
}

func (m Model) AttendanceDetailsView() string {
	records := m.attendanceByCourse[m.selectedCourse.Name]
	
	s := titleStyle.Render("Attendance -> " + m.selectedCourse.Name) + "\n\n"
	
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
	s += descStyle.Render(fmt.Sprintf("Overall: %d/%d (%.2f%%)\n%s\n", attended, total, percent, warning))
	if len(records) > 0 {
		m.selectedURL = records[0].AttendanceURL
		s += fmt.Sprintf("\nq to back • Ctrl+c to exit • o open in browser • %s\n", m.selectedURL)
	}
	return s
}

func (m Model) AssignmentCourseView() []list.Item {
	// Initialize courselist for assignment selection
	items := make([]list.Item, len(m.courses))
	for i, course := range m.courses {
		items[i] = Item{
			ItemTitle: course.Name,
			ItemDesc: fmt.Sprintf("Total Assignments: %d", len(m.assignmentsByCourse[course.Name])),
		}
	}
	return items
}

func (m Model) AssignmentListView(selectedCourse string) []list.Item {
	assignments := m.assignmentsByCourse[selectedCourse]
	items := make([]list.Item, len(assignments))
	for i, assignment := range assignments {
		items[i] = Item{
			ItemTitle: assignment.Title,
			ItemDesc: fmt.Sprintf("Open Date: %s | Due Date: %s | Status: %s | Grade: %s\n", assignment.OpenDate, assignment.DueDate, assignment.Status, assignment.Grade),
			ItemUrl: assignment.URL,
		}
	}
	return items
}

func (m Model) VPLCourseView() []list.Item {
	items := make([]list.Item, len(m.courses))
	for i, course := range m.courses {
		items[i] = Item{
			ItemTitle: course.Name,
			ItemDesc: fmt.Sprintf("Total VPLs: %d", len(m.vplsByCourse[course.Name])),
		}
	}
	return items
}


func (m Model) VPLListView(selectedCourse string) []list.Item {
	vpls := m.vplsByCourse[selectedCourse]
	items := make([]list.Item, len(vpls))
	for i, vpl := range vpls {
		items[i] = Item {
			ItemTitle: vpl.Title,
			ItemDesc: fmt.Sprintf("Open Date: %s | Due Date: %s", vpl.OpenDate, vpl.DueDate),
			ItemUrl: vpl.URL,
		}
	}
	return items
}

// DisableListQuit removes the default 'q' quit binding from a list so that
// 'q' can be handled by the parent model to navigate back instead of exiting.
func DisableListQuit(l *list.Model) {
	l.KeyMap.Quit.SetKeys()
	l.KeyMap.ForceQuit.SetKeys()
	// Add custom help text for 'q' to show it goes back
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(
				key.WithKeys("q"),
				key.WithHelp("q", "back"),
			),
		}
	}
}
