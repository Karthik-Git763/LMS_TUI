package main

import (
	"fmt"
	"log"
	"os"

	"lms/internal/auth"
	"lms/internal/scrapper"

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
	
	client := auth.CreateClient()
	
	token, err := auth.FetchToken(client, baseURL)
	if err != nil {
		log.Fatal(err)
	}
	
	err = auth.Login(client, baseURL, username, password, token)
	if err != nil {
		log.Fatal(err)
	}
	
	ok := auth.CheckLogin(client, baseURL)
	if !ok {
		log.Fatal("Login Failed")
	}
	fmt.Println("Login Successful")
	scrapper.FetchCourses(client, baseURL)
}
