package scrapper

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type AttendanceRecord struct {
	Session string
	Status 	string
	Date 	string
}

func FindAttendanceURL(client *http.Client, courseURL string) (string, error) {
	resp, err := client.Get(courseURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}
	
	var attendanceURL string
	
	doc.Find("div.activityname a").EachWithBreak(func(i int, s *goquery.Selection) bool {
		text := strings.ToLower(strings.TrimSpace(s.Find("span.instancename").Text()))
		if strings.Contains(text, "attendance") {
			href, exists := s.Attr("href")
			if exists && strings.Contains(href, "/mod/attendance/") {
				u, err := url.Parse(href)
				if err != nil {
					return true
				}
				// Add / Override view=5 
				// Only for ug23 websites
				q := u.Query()
				q.Set("view", "5")
				u.RawQuery = q.Encode()
				
				attendanceURL = u.String()
				// fmt.Println("Attendance URL found:", attendanceURL)
				return false
			}
		}
		return true
	})
	if attendanceURL == "" {
		return "", fmt.Errorf("attendance URL not found")
	}
	return attendanceURL, nil
}

func ScrapeAttendance(client *http.Client, attendanceURL string) ([]AttendanceRecord, error) {
	resp, err := client.Get(attendanceURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	
	var attendanceRecords = []AttendanceRecord{}
	
	doc.Find("table.generaltable tbody tr").Each(func(i int, row *goquery.Selection) {
		cols := row.Find("td")
		if cols.Length() < 3 {
			return
		}
		
		record := AttendanceRecord {
			Session: strings.TrimSpace(cols.Eq(0).Text()),
			Date:  strings.TrimSpace(cols.Eq(1).Text()),
			Status:    strings.TrimSpace(cols.Eq(2).Text()),
		}
		attendanceRecords = append(attendanceRecords, record)
	})
	return attendanceRecords, nil
}

func CalculateAttendancePercentage(records []AttendanceRecord) (attended int, total int) {
	for _, record := range records {
		total++
		if record.Status == "Present" {
			attended++
		}
	}
	return
}