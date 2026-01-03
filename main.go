package main

import (
	"fmt"
	"log"
	"os"

	"lms/internal/auth"
	"lms/internal/scrapper"

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
	fmt.Println("Login Successful")
	
	courses := scrapper.FetchCourses(client, baseURL)
	
	for _, course := range courses {
		fmt.Println("\n Course: ", course.Name)
		
		attendanceURL, err := scrapper.FindAttendanceURL(client, course.URL)
	 	if err != nil {
			fmt.Println("No attendance module")
			continue
		}
		
		attendanceRecords, err := scrapper.ScrapeAttendance(client, attendanceURL)
		if err != nil {
			fmt.Println("Error fetching attendance")
			continue
		}
	
		attended, total := scrapper.CalculateAttendancePercentage(attendanceRecords)
		percent := float64(attended) / float64(total) * 100
		fmt.Printf("Attendance: %d / %d = (%.2f%%)\n", attended, total, percent)
	}
	
	for _, course := range courses {
		assignments, err := scrapper.FindAssignmentsInCourse(client, course)
		if err != nil {
			fmt.Println("Error fetching assignments")
			continue
		}
		for i := range assignments {
			scrapper.AssignementDetailsStatusAndDueDate(client, &assignments[i])
			fmt.Println("Assignment Details: ", assignments[i].Title)
			fmt.Println("Course: ", assignments[i].CourseName)
			fmt.Println("Status: ", assignments[i].Status)
			fmt.Println("Open Date: ", assignments[i].OpenDate)
			fmt.Println("Due Date: ", assignments[i].DueDate)
			fmt.Println("Grade: ", assignments[i].Grade)
			fmt.Println()
		}
	}
}
