package models

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
	Session 	string
	Status 		string
	Date 		string
	AttendanceURL 	string
}

type VPL struct {
	Title      string
	CourseName string
	URL        string
	OpenDate   string
	DueDate    string
}


