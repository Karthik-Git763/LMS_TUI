package scrapper

import (
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Assignment struct {
	Title 		string
	OpenDate	string
	DueDate 	string
	Status 		string
	Grade		string
	CourseName 	string
	URL 		string
}

func FindAssignmentsInCourse(client *http.Client, course Course) ([]Assignment, error) {
	resp, err := client.Get(course.URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	
	var assignments []Assignment
	
	doc.Find("a[href*='/mod/assign/view.php']").Each(func(i int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Text())
		href, exists := s.Attr("href")
		if !exists || title == "" {
			return 
		}
		// fmt.Println("Found Assignment URL", href)
		
		assignments = append(assignments, Assignment{
			Title: title,
			CourseName: course.Name,
			URL: href,
		})
	})
	return assignments, nil
}

func ParseLabelValues(label, value string, a *Assignment) {
	switch {
		case strings.Contains(label, "submission status"):
			a.Status = value
		case strings.Contains(label, "grading status"):
			a.Grade = value
		case strings.Contains(label, "opened"):
			a.OpenDate = value
		case strings.Contains(label, "due"):
			a.DueDate = value
	}
}

func AssignementDetailsStatusAndDueDate(client *http.Client, a *Assignment) error {
	resp, err := client.Get(a.URL)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}
	
	doc.Find("table.generaltable tr").Each(func(i int, row *goquery.Selection) {
		label := strings.ToLower(strings.TrimSpace(row.Find("th").Text()))
		value := strings.TrimSpace(row.Find("td").Text())
		// fmt.Println("Label:", label, "Value:", value)
		ParseLabelValues(label, value, a)
	})
	
	doc.Find("div.activity-dates div").Each(func(i int, s *goquery.Selection) {
		label := strings.ToLower(strings.TrimSpace(s.Find("strong").Text()))
		value := strings.TrimSpace(strings.ReplaceAll(s.Text(), s.Find("strong").Text(), ""))
		// fmt.Println("Label:", label, "Value:", value)
		ParseLabelValues(label, value, a)
	})
	return nil
}
