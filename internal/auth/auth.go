package auth

import (
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func CreateClient() *http.Client {
	Jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}
	return &http.Client{
		Jar : Jar,
	}
}

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

func CheckLogin(client *http.Client, baseURL string) bool {
	resp, err := client.Get(baseURL + "/my/")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	
	return resp.StatusCode == http.StatusOK
}