package scrapper

import (
	"log"
	"net/http"
	"sync"

	"lms/internal/models"

	tea "github.com/charmbracelet/bubbletea"
)

// StartDataFetching starts the data fetching process for all courses
func StartDataFetching(client *http.Client, courses []models.Course) tea.Cmd {
	return func() tea.Msg {
		attendanceByCourse := make(map[string][]models.Attendance)
		assignmentByCourse := make(map[string][]models.Assignment)
		vplByCourse := make(map[string][]models.VPL)

		var mu sync.Mutex
		var wg sync.WaitGroup

		for _, course := range courses {
			wg.Add(1)
			go func(course models.Course) {
				defer wg.Done()
				attendanceURL, err := FindAttendanceURL(client, course.URL)
				if err == nil && attendanceURL != "" {
					log.Println("Attendance URL found:", attendanceURL)
					attendanceRecords, err := ScrapeAttendance(client, attendanceURL)
					if err == nil {
						mu.Lock()
						attendanceByCourse[course.Name] = attendanceRecords
						mu.Unlock()
						log.Println("Attendance records fetched:", len(attendanceRecords))
					}
				} else {
					log.Println("Attendance URL not found")
				}

				assignments, err := FindAssignmentsInCourse(client, course)
				if err == nil {
					var assignWg sync.WaitGroup
					for k := range assignments {
						assignWg.Go(func() {
							AssignmentDetailsStatusAndDueDate(client, &assignments[k])
						})
					}
					assignWg.Wait()

					mu.Lock()
					assignmentByCourse[course.Name] = assignments
					mu.Unlock()
					log.Println("Assignments fetched:", len(assignments))
				} else {
					log.Println("Error fetching assignments")
				}

				vpls, err := FindVPLInCourse(client, course)
				if err == nil {
					var vplWg sync.WaitGroup
					for k := range vpls {
						vplWg.Go(func() {
							VPLDetailsStatusAndDueDate(client, &vpls[k])
						})
					}
					vplWg.Wait()

					mu.Lock()
					vplByCourse[course.Name] = vpls
					mu.Unlock()
					log.Println("VPLs fetched:", len(vpls))
				} else {
					log.Println("Error fetching VPLs")
				}
			}(course)
		}
		wg.Wait()
		return models.DataLoadedMsg{Attendance: attendanceByCourse, Assignment: assignmentByCourse, VPL: vplByCourse}
	}
}
