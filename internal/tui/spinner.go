package tui

import (
	"lms/internal/models"

	"github.com/charmbracelet/bubbles/spinner"
)

func newSpinner() spinner.Model {
	s := spinner.New()
	s.Spinner = spinner.Ellipsis
	return s
}

func (m *Model) applyProgress(msg models.ProgressMsg) {
	switch msg.Type {
	case models.CourseStarted:
		m.active[msg.Course] = true
	case models.CourseCompleted:
		delete(m.active, msg.Course)
		m.done[msg.Course] = true
	case models.CourseError:
		delete(m.active, msg.Course)
		m.errors[msg.Course] = msg.Err
	}

	if len(m.active) == 0 && m.screen == progressScreen {
		m.screen = menuScreen
		m.cursor = 0
	}
}

func (m *Model) applyDataLoaded(msg models.DataLoadedMsg) {
	m.attendanceByCourse = msg.Attendance
	m.assignmentsByCourse = msg.Assignment
	m.vplsByCourse = msg.VPL

	if m.screen == progressScreen {
		m.screen = menuScreen
		m.cursor = 0
	}
}