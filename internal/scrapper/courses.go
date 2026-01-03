package scrapper

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func FetchCourses(client *http.Client, baseURL string) {
	resp, err := client.Get(baseURL + "/my/")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("a[href*='course/view.php']").Each(func (i int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Text())
		link, _ := s.Attr("href")
		fmt.Printf("Course: %s\nLink: %s\n", title, link)
	})
}