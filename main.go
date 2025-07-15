package main

import (
	"editGo/cmd/app"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"log"
	"os"
)

func init() {
	logFile, err := os.OpenFile("editor.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Could not open log file:", err)
		os.Exit(1)
	}
	log.SetOutput(logFile)
}

func main() {
	var model app.Model
	if len(os.Args) > 1 {
		model = app.NewModel(os.Args[1])
	} else {
		model = app.NewModel("")
	}
	p := tea.NewProgram(model)
	if err := p.Start(); err != nil {
		panic(err)
	}
}
