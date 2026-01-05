package scrapper

import (
	// "fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Course struct {
	Name 			string
	URL 			string
	AttendanceURL 	string
}

func FetchCourses(client *http.Client, baseURL string) []Course {
	resp, err := client.Get(baseURL + "/my/")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var courses []Course
	
	doc.Find("a[href*='course/view.php']").Each(func (i int, s *goquery.Selection) {
		name := strings.TrimSpace(s.Text())
		link, exists := s.Attr("href")
		if exists && name != "" {
			courses = append(courses, Course{Name: name, URL: link})
			// fmt.Printf("Course: %s\nLink: %s\n", name, link)
		}
	})
	return courses
}