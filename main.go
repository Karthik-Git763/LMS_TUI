package main

import (
	"fmt"
	"log"
	"os"

	"lms/internal/auth"
	"lms/internal/scrapper"
	"lms/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
)

const baseURL = "https://lmsug23.iiitkottayam.ac.in"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	username := os.Getenv("YOUR_USERNAME")
	password := os.Getenv("YOUR_PASSWORD")

	client := auth.CreateClient()

	token, err := auth.FetchToken(client, baseURL)
	if err != nil {
		log.Fatal(err)
	}

	err = auth.Login(client, baseURL, username, password, token)
	if err != nil {
		log.Fatal(err)
	}

	ok := auth.CheckLogin(client, baseURL)
	if !ok {
		log.Fatal("Login Failed")
	}
	// fmt.Println("Login Successful")

	courses := scrapper.FetchCourses(client, baseURL)

	attendanceByCourse := make(map[string][]scrapper.Attendance)
	assignmentByCourse := make(map[string][]scrapper.Assignment)
	vplByCourse := make(map[string][]scrapper.VPL)

	for i, course := range courses {
		// fmt.Println("\n Course: ", course.Name)

		attendanceURL, err := scrapper.FindAttendanceURL(client, course.URL)
		if err != nil {
			// fmt.Println("No attendance module")
			continue
		}
		courses[i].AttendanceURL = attendanceURL
		
		attendanceRecords, err := scrapper.ScrapeAttendance(client, attendanceURL)
		if err != nil {
			// fmt.Println("Error fetching attendance")
			continue
		}
		attendanceByCourse[course.Name] = attendanceRecords

		
		assignments, err := scrapper.FindAssignmentsInCourse(client, course)
		if err != nil {
			fmt.Println("Error fetching assignments")
			continue
		}
		for i := range assignments {
			scrapper.AssignementDetailsStatusAndDueDate(client, &assignments[i])
		}
		assignmentByCourse[course.Name] = assignments
		vpls, err := scrapper.FindVPLInCourse(client, course)
		if err != nil {
			fmt.Println("Error fetching VPLs")
			continue
		}
		for i := range vpls {
			scrapper.VPLDetailsStatusAndDueDate(client, &vpls[i])
		}
		vplByCourse[course.Name] = vpls
	}

	p := tea.NewProgram(tui.InitialModel(courses, attendanceByCourse, assignmentByCourse, vplByCourse), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
