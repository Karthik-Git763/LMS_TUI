package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"lms/internal/auth"
	"lms/internal/models"
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

	p := tea.NewProgram(tui.InitialModel(courses), tea.WithAltScreen())
	go func() {
		attendanceByCourse := make(map[string][]models.Attendance)
		assignmentByCourse := make(map[string][]models.Assignment)
		vplByCourse := make(map[string][]models.VPL)
		
		var mu sync.Mutex
		
		var wg sync.WaitGroup
	
		for _, course := range courses {
			// fmt.Println("\n Course: ", course.Name)
			p.Send(models.ProgressMsg{Course: course.Name, Type: models.CourseStarted})
			wg.Go(func() {
				attendanceURL, err := scrapper.FindAttendanceURL(client, course.URL)
				if err == nil && attendanceURL != "" {
					log.Println("Attendance URL found:", attendanceURL)
					attendanceRecords, err := scrapper.ScrapeAttendance(client, attendanceURL)
					if err == nil {
						mu.Lock()
						attendanceByCourse[course.Name] = attendanceRecords
						mu.Unlock()
						log.Println("Attendance records fetched:", len(attendanceRecords))
					}
				} else {
					log.Println("Attendance URL not found")
					p.Send(models.ProgressMsg{Course: course.Name, Err: err, Type: models.CourseError})
				}
	
				assignments, err := scrapper.FindAssignmentsInCourse(client, course)
				if err == nil {
					var assignWg sync.WaitGroup
					for k := range assignments {
						assignWg.Go(func() {
							scrapper.AssignmentDetailsStatusAndDueDate(client, &assignments[k])
						})
					}
					assignWg.Wait()
					
					mu.Lock()
					assignmentByCourse[course.Name] = assignments
					mu.Unlock()
					log.Println("Assignments fetched:", len(assignments))
				} else {
					log.Println("Error fetching assignments")
					p.Send(models.ProgressMsg{Course: course.Name, Err: err, Type: models.CourseError})
				}
	
				vpls, err := scrapper.FindVPLInCourse(client, course)
				if err == nil {
					var vplWg sync.WaitGroup
					for k := range vpls {
						vplWg.Go(func() {
							scrapper.VPLDetailsStatusAndDueDate(client, &vpls[k])
						})
					}
					vplWg.Wait()
					
					mu.Lock()
					vplByCourse[course.Name] = vpls
					mu.Unlock()
					log.Println("VPLs fetched:", len(vpls))
				} else {
					log.Println("Error fetching VPLs")
					p.Send(models.ProgressMsg{Course: course.Name, Err: err, Type: models.CourseError})
				}
				p.Send(models.ProgressMsg{Course: course.Name, Type: models.CourseCompleted})
			})
		}
		wg.Wait()
		p.Send(models.DataLoadedMsg{Attendance: attendanceByCourse, Assignment: assignmentByCourse, VPL: vplByCourse})
	}()
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
