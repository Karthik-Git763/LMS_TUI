package scrapper

import (
	"lms/internal/models"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func FetchCourses(client *http.Client, baseURL string) []models.Course {
	resp, err := client.Get(baseURL + "/my/")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var courses []models.Course
	
	doc.Find("a[href*='course/view.php']").Each(func (i int, s *goquery.Selection) {
		name, _ := s.Attr("title")
		name = strings.TrimSpace(name)
		link, exists := s.Attr("href")
		if exists && name != "" {
			courses = append(courses, models.Course{Name: name, URL: link})
			// fmt.Printf("Course: %s\nLink: %s\n", name, link)
		}
	})
	return courses
}