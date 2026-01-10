package auth

import (
	"fmt"
	"lms/internal/models"
	"lms/internal/scrapper"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	tea "github.com/charmbracelet/bubbletea"
)

// CreateClient creates a new HTTP client with a cookie jar
func CreateClient() *http.Client {
	Jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}
	return &http.Client{
		Jar : Jar,
		Timeout: time.Second * 10,
	}
}

// FetchToken fetches the login token from the Moodle website
func FetchToken(client *http.Client, baseURL string) (string, error) {
	resp, err := client.Get(baseURL + "/login/index.php")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}
	
	token, exists := doc.Find("input[name='logintoken']").Attr("value")
	if !exists {
		return "", fmt.Errorf("token not found")
	}
	return token, nil
}

// Login: logs in the user
func Login(client *http.Client, baseURL, username, password, token string) error {
	data := url.Values {}
	data.Set("username", username)
	data.Set("password", password)
	data.Set("logintoken", token)
	
	req, err := http.NewRequest(
		"POST",
		baseURL + "/login/index.php",
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		return err
	}
	
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}

// CheckLogin checks if the user is logged in
func CheckLogin(client *http.Client, baseURL string) bool {
	resp, err := client.Get(baseURL + "/my/")
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	// If Moodle redirected back to the login page, authentication failed
	if strings.Contains(resp.Request.URL.Path, "login/index.php") {
		return false
	}

	return resp.StatusCode == http.StatusOK
}

// AuthenticateAndFetchData performs authentication and fetches courses
func AuthenticateAndFetchData(client *http.Client, baseURL, username, password string) tea.Msg {
	// Fetch login token
	token, err := FetchToken(client, baseURL)
	if err != nil {
		log.Println("Error fetching token:", err)
		return models.AuthResultMsg{Success: false, Error: "Failed to fetch login token"}
	}

	// Authenticate with credentials
	err = Login(client, baseURL, username, password, token)
	if err != nil {
		log.Println("Error logging in:", err)
		return models.AuthResultMsg{Success: false, Error: "Invalid username or password"}
	}

	// Verify login
	ok := CheckLogin(client, baseURL)
	if !ok {
		log.Println("Login check failed")
		return models.AuthResultMsg{Success: false, Error: "Login verification failed"}
	}

	// Fetch courses
	courses := scrapper.FetchCourses(client, baseURL)
	log.Println("Courses fetched:", len(courses))

	return models.AuthResultMsg{Success: true, Courses: courses}
}