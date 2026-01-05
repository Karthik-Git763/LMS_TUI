package scrapper

import (
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type VPL struct {
	Title      string
	CourseName string
	URL        string
	OpenDate   string
	DueDate    string
}

func FindVPLInCourse(client *http.Client, course Course) ([]VPL, error) {
	resp, err := client.Get(course.URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var vpls []VPL

	doc.Find("a[href*='/mod/vpl/view.php']").Each(func(i int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Text())
		href, exists := s.Attr("href")
		if !exists || title == "" {
			return
		}
		vpls = append(vpls, VPL{
			Title:      title,
			CourseName: course.Name,
			URL:        href,
		})
	})
	return vpls, nil
}

func VPLDetailsStatusAndDueDate(client *http.Client, v *VPL) error {
    resp, err := client.Get(v.URL)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    doc, err := goquery.NewDocumentFromReader(resp.Body)
    if err != nil {
        return err
    }

    box := doc.Find("div.box.generalbox").First()
    boxText := strings.TrimSpace(box.Text())
    
    // Split by exact "Label: "
    for _, label := range []string{"Available from:", "Due date:"} {
        parts := strings.SplitN(boxText, label, 2)
        if len(parts) > 1 {
            labelLower := strings.ToLower(strings.TrimSuffix(label, ":"))
            
            // Take until next label or end
            nextStart := len(parts[1])
            for _, nextLabel := range []string{"Available from:", "Due date:", "Maximum number", "Type of"} {
                if nextLabel != label {
                    pos := strings.Index(parts[1], nextLabel)
                    if pos != -1 && pos < nextStart {
                        nextStart = pos
                    }
                }
            }
            value := strings.TrimSpace(parts[1][:nextStart])
            
            // fmt.Printf("Label: %s\nValue: %s\n", labelLower, value)
            
            switch {
            case strings.Contains(labelLower, "available"):
                v.OpenDate = value
            case strings.Contains(labelLower, "due"):
                v.DueDate = value
            }
        }
    }
    return nil
}

