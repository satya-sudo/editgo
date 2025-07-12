package app

import (
	"editGo/data"
	"editGo/editor"
	"editGo/ui"
	tea "github.com/charmbracelet/bubbletea"
	"time"
)

type Model struct {
	Buffer        *editor.TextBuffer
	Cursor        *editor.CursorPointer
	File          *data.FileManager
	UndoStack     *editor.UndoManager
	AutoSaver     *data.AutoSave
	StatusMessage string
}

func NewModel(filePath string) Model {
	var file *data.FileManager
	var err error

	if filePath != "" {
		file, err = data.NewFile(filePath)
	} else {
		file, err = data.NewEmptyFile()
	}
	if err != nil {
		panic(err)
	}

	buffer := file.Buffer
	cursor := editor.NewCursor(0, 0)
	undo := editor.NewUndoManager()
	auto := data.NewAutoSave(file, 5*time.Second)
	auto.Start()

	return Model{
		Buffer:    buffer,
		Cursor:    cursor,
		File:      file,
		UndoStack: undo,
		AutoSaver: auto,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case (len(msg.Runes) == 1 && msg.Type == tea.KeyRunes) || msg.String() == " ":
			var r rune
			if len(msg.Runes) == 1 {
				r = msg.Runes[0]
			} else {
				r = ' ' // fallback for space
			}
			m.UndoStack.Push(m.Buffer)
			m.Buffer.InsertRune(m.Cursor.Y, m.Cursor.X, r)
			m.Cursor.MoveRight(m.Buffer)
		case msg.Type == tea.KeyBackspace:
			if m.Cursor.X == 0 && m.Cursor.Y == 0 {
				break // at top-left, nothing to delete
			}

			m.UndoStack.Push(m.Buffer)

			if m.Cursor.X > 0 {
				m.Buffer.DeleteRune(m.Cursor.Y, m.Cursor.X, 0)
				m.Cursor.MoveLeft(m.Buffer)
			} else if m.Cursor.Y > 0 {
				// Merge with previous line
				prevLineLen := len(m.Buffer.GetLine(m.Cursor.Y - 1))
				m.Buffer.MergeLine(m.Cursor.Y - 1)
				m.Cursor.Y--
				m.Cursor.X = prevLineLen
			}
		case msg.Type == tea.KeyEnter:
			m.UndoStack.Push(m.Buffer)
			m.Buffer.InsertNewLine(m.Cursor.Y, m.Cursor.X)
			m.Cursor.Y++
			m.Cursor.X = 0
		case msg.Type == tea.KeyCtrlQ, msg.Type == tea.KeyCtrlC:
			m.AutoSaver.Stop()
			return m, tea.Quit
		case msg.Type == tea.KeyUp:
			m.Cursor.MoveUp(m.Buffer)
		case msg.Type == tea.KeyDown:
			m.Cursor.MoveDown(m.Buffer)
		case msg.Type == tea.KeyLeft:
			m.Cursor.MoveLeft(m.Buffer)
		case msg.Type == tea.KeyRight:
			m.Cursor.MoveRight(m.Buffer)
		case msg.Type == tea.KeyCtrlZ:
			m.UndoStack.Undo(m.Buffer)
			m.Cursor.Clamp(m.Buffer)
		case msg.Type == tea.KeyCtrlY:
			m.UndoStack.Redo(m.Buffer)
			m.Cursor.Clamp(m.Buffer)
		case msg.Type == tea.KeyCtrlS:
			if m.File.FilePath == "" {
				// Optional: add SaveAs prompt later
				return m, nil
			}

			err := m.File.Save()
			if err != nil {
				m.StatusMessage = "Error: " + err.Error()
			} else {
				m.StatusMessage = "Saved to: " + m.File.FilePath
			}
		}
	}
	return m, nil
}

func (m Model) View() string {
	return ui.RenderStatusBar(m.File.FilePath, m.Buffer.IsDirty(), m.Cursor.X, m.Cursor.Y) + "\n" +
		ui.RenderBuffer(m.Buffer.Lines, m.Cursor.X, m.Cursor.Y) + "\n" +
		ui.RenderHelpBar() + "\n" +
		ui.RenderStatusMessage(m.StatusMessage)
}
