package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
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
	
	client := createClient()
	
	token, err := fetchToken(client)
	if err != nil {
		log.Fatal(err)
	}
	
	err = login(client, username, password, token)
	if err != nil {
		log.Fatal(err)
	}
	
	ok := checkLogin(client)
	if !ok {
		log.Fatal("Login Failed")
	}
	fmt.Println("Login Successful")
	fetchCourses(client)
}

func createClient() *http.Client {
	Jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}
	return &http.Client{
		Jar : Jar,
	}
}

func fetchToken(client *http.Client) (string, error) {
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

func login(client *http.Client, username, password, token string) error {
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

func checkLogin(client *http.Client) bool {
	resp, err := client.Get(baseURL + "/my/")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	
	return resp.StatusCode == http.StatusOK
}

func fetchCourses(client *http.Client) {
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