package main

import (
	"fmt"
	"log"
	"os"
	"sync"

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

	f, err := tea.LogToFile("./log.txt", "Log: ")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer f.Close()
	log.SetOutput(f)

	courses := scrapper.FetchCourses(client, baseURL)

	attendanceByCourse := make(map[string][]scrapper.Attendance)
	assignmentByCourse := make(map[string][]scrapper.Assignment)
	vplByCourse := make(map[string][]scrapper.VPL)
	
	var mu sync.Mutex
	
	var wg sync.WaitGroup

	for i, course := range courses {
		// fmt.Println("\n Course: ", course.Name)
		wg.Go(func() {
			attendanceURL, err := scrapper.FindAttendanceURL(client, course.URL)
			if err == nil && attendanceURL != "" {
				mu.Lock()
				courses[i].AttendanceURL = attendanceURL
				mu.Unlock()
				log.Println("Attendance URL found:", attendanceURL)
				attendanceRecords, err := scrapper.ScrapeAttendance(client, attendanceURL)
				if err == nil {
					mu.Lock()
					attendanceByCourse[course.Name] = attendanceRecords
					mu.Unlock()
					log.Println("Attendance records fetched:", len(attendanceRecords))
				}
			}

			assignments, err := scrapper.FindAssignmentsInCourse(client, course)
			if err == nil {
				for k := range assignments {
					wg.Go(func() {
						scrapper.AssignementDetailsStatusAndDueDate(client, &assignments[k])
					})
				}
				
				mu.Lock()
				assignmentByCourse[course.Name] = assignments
				mu.Unlock()
				log.Println("Assignments fetched:", len(assignments))
			} else {
				fmt.Println("Error fetching assignments")
			}

			vpls, err := scrapper.FindVPLInCourse(client, course)
			if err == nil {
				for k := range vpls {
					wg.Go(func() {
						scrapper.VPLDetailsStatusAndDueDate(client, &vpls[k])
					})
				}
				
				mu.Lock()
				vplByCourse[course.Name] = vpls
				mu.Unlock()
				log.Println("VPLs fetched:", len(vpls))
			} else {
				fmt.Println("Error fetching VPLs")
			}
		})
	}
	wg.Wait()

	p := tea.NewProgram(tui.InitialModel(courses, attendanceByCourse, assignmentByCourse, vplByCourse), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
