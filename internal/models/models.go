package models

import "time"

type Credentials struct {
	Username 		string
	Password 		string
	Authenticated 	bool
}

type Course struct {
	Name 			string
	URL 			string
}

type Assignment struct {
	Title 		string
	OpenDate	string
	DueDate 	string
	Status 		string
	Grade		string
	CourseName 	string
	URL 		string
}

type Attendance struct {
	Date 		string
	Session 	string
	Status 		string
	AttendanceURL 	string
}

type VPL struct {
	Title      string
	CourseName string
	URL        string
	OpenDate   string
	DueDate    string
}

type CourseData struct {
	Attendance 	[]Attendance
	Assignment 	[]Assignment
	VPL       	[]VPL
}

type DataLoadedMsg struct {
	Attendance map[string][]Attendance
	Assignment map[string][]Assignment
	VPL		   map[string][]VPL
}

type Cache struct {
	Timestamp	time.Time
	Courses 	[]Course
	Data 		map[string]CourseData
}

type ProgressType int

const (
	CourseStarted ProgressType = iota
	CourseCompleted
	CourseError
)

type ProgressMsg struct {
	Course string
	Type   ProgressType
	Err    error
}
