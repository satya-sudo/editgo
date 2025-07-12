package ui

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"
	"os"
	"strings"
)

type KeyHelp struct {
	Key    string
	Action string
}

var HelpKeys = []KeyHelp{
	{"←/→/↑/↓", "Move"},
	{"Enter", "New Line"},
	{"Backspace", "Delete"},
	{"Ctrl+S", "Save"},
	{"Ctrl+Z", "Undo"},
	{"Ctrl+Y", "Redo"},
	{"Ctrl+Q", "Quit"},
}
var statusBarStyle = lipgloss.NewStyle().
	Background(lipgloss.Color("#44475a")).
	Foreground(lipgloss.Color("#f8f8f2")).
	Padding(0, 1)

var helpBarStyle = lipgloss.NewStyle().
	Background(lipgloss.Color("#282a36")). // slightly different dark shade
	Foreground(lipgloss.Color("#bd93f9")). // lavender for commands
	Padding(0, 1).
	Italic(true)

var cursorCharStyle = lipgloss.NewStyle().
	Background(lipgloss.Color("#f8f8f2")).
	Foreground(lipgloss.Color("#282a36")).
	Bold(true)

var statusMsgStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#00FF00")). // Green text
	Background(lipgloss.Color("#1A1A1A")). // Dark gray background
	Padding(0, 1).
	Width(100)

func RenderStatusMessage(msg string) string {
	if msg == "" {
		return ""
	}
	return statusMsgStyle.Render("Status: " + msg)
}
func RenderHelpBar() string {
	helpText := " Ctrl+S Save | Ctrl+O Open | Ctrl+Z Undo | Ctrl+Y Redo | Ctrl+C Quit "
	return helpBarStyle.Render(helpText)
}
func RenderStatusBar(filePath string, isDirty bool, cursorX, cursorY int) string {
	dirtyFlag := ""
	if isDirty {
		dirtyFlag = "✱"
	}

	if filePath == "" {
		filePath = "[No Name]"
	}

	status := fmt.Sprintf(" %s %s | Ln %d, Col %d ", filePath, dirtyFlag, cursorY+1, cursorX+1)
	return statusBarStyle.Render(status)
}

func RenderBuffer(lines [][]rune, cursorX, cursorY int) string {
	var out strings.Builder

	_, height, err := term.GetSize(os.Stdout.Fd())
	if err != nil {
		height = 20 // fallback
	}
	contentHeight := height - 5 // for status and help bar

	for y := 0; y < contentHeight; y++ {
		var line string
		if y < len(lines) {
			runes := lines[y]
			if y == cursorY {
				before := string(runes[:cursorX])
				cursorChar := " "
				after := ""

				if cursorX < len(runes) {
					cursorChar = string(runes[cursorX])
					after = string(runes[cursorX+1:])
				}
				line = before + cursorCharStyle.Render(cursorChar) + after
			} else {
				line = string(runes)
			}
		} else {
			line = "" // or "~"
		}
		out.WriteString(line + "\n")
	}

	return out.String()
}
