package main

import (
	// "fmt"
	"log"

	"lms/internal/auth"
	"lms/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
)

const baseURL = "https://lmsug23.iiitkottayam.ac.in"

func main() {
	f, err := tea.LogToFile("./log.txt", "Log: ")
	if err != nil {
		// fmt.Println(err.Error())
		// os.Exit(1)
	}
	defer f.Close()
	log.SetOutput(f)

	client := auth.CreateClient()

	m := tui.InitialModel(nil, nil, client, baseURL)
	
	p := tea.NewProgram(&m, tea.WithAltScreen())
	
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
