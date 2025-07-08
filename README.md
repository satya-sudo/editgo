# Go Text Editor (Terminal-Based)

A terminal-based text editor written in **Go**, designed to explore **data structures** like `Stack` and `Trie`, and utilize **Go concurrency** patterns (channels, goroutines). The editor provides essential editing features and operates within the terminal using the **Bubbletea** framework.

---

## 🧠 Project Goals

* Build a **minimal, fast, local-first text editor**
* Practice structuring a Go project cleanly
* Use core data structures: `Stack` (Undo), `Trie` (Search)
* Apply Go concepts: `Goroutines`, `Channels`
* Start with a terminal UI (TUI) using `Bubbletea`, optionally extend to GUI later

---

## ✨ Core Features

| Feature                | Description                                                  | Implementation Logic                      |
| ---------------------- | ------------------------------------------------------------ | ----------------------------------------- |
| Open File              | Load text content from a file into the buffer                | `os.ReadFile()`, parse into 2D rune slice |
| Edit Text              | Insert, delete, newline, backspace                           | Buffer as 2D rune slice (`[][]rune`)      |
| Save File              | Save current buffer content to file                          | `os.WriteFile()`                          |
| Undo/Redo              | Revert or re-apply changes                                   | `UndoStack` and `RedoStack` storing diffs |
| Cursor Movement        | Arrow key movement and bounds checks                         | `cursorX`, `cursorY` vars with logic      |
| New File               | Start editing a new unsaved buffer                           | State management                          |
| Autosave               | Background goroutine that saves periodically                 | `time.Ticker` + `chan SaveSignal`         |
| Quit Prompt            | Ask to save before exiting unsaved changes                   | Buffer dirty flag + confirm popup         |
| Trie-Based Suggestions | Show live suggestions as user types based on words in buffer | Trie for prefix search, dynamic updates   |

---

## 🧩 Data Structures

### 🔁 Stack – Undo/Redo

Used to track user actions like insert/delete.
Each operation is stored as a diff struct:

```go
type EditAction struct {
    Line, Col int
    Char      rune
    Action    string // "insert" or "delete"
}
```

### 🌲 Trie – Word Search

All words from the buffer are added to a Trie:

* Allows fast prefix search (e.g., "pri" → Print, Println)
* Stores optional metadata: word → line numbers

```go
type TrieNode struct {
    Children map[rune]*TrieNode
    IsEnd    bool
    Lines    map[int]struct{}
}
```

### 🔄 Channels – Autosave

Used for background operations like saving:

* A channel listens for manual or timed triggers

```go
autosaveCh := make(chan struct{})
go func() {
    for range autosaveCh {
        SaveFile(buffer)
    }
}()
```

---

## 🖥️ Terminal UI (Bubbletea)

The UI is built using [Bubbletea](https://github.com/charmbracelet/bubbletea), and consists of:

* **Editor View**: Main area for text display/editing
* **Command Mode**: For actions like `:save`, `:quit`, `:find`
* **Search Suggestions**: Popup panel for Trie-based results
* **Status Bar**: Shows current file, cursor position, dirty flag

---

## 🗂️ Project Structure

```
texteditor/
├── cmd/
│   └── app.go            # Bubbletea app entrypoint and loop
├── editor/
│   ├── buffer.go         # Text buffer: [][]rune, insert/delete
│   ├── cursor.go         # Cursor logic
│   ├── undo.go           # Undo/Redo stacks
│   ├── autosave.go       # Autosave goroutine
│   └── search_trie.go    # Trie structure for word search
├── data/
│   └── fileio.go         # File open/save logic
├── ui/
│   └── render.go         # UI helpers, text rendering, status bar
├── internal/             # (Optional) internal helpers/utilities
├── main.go               # Application entrypoint
```

---

## 🔄 Example User Flow

1. Launch editor with `go run main.go`
2. Load existing file or start with empty buffer
3. Edit text using keyboard (char keys, arrows, backspace, Enter)
4. Autosave runs in the background every 5s
5. Press `Ctrl+Z` to undo, `Ctrl+Y` to redo
6. Press `/` to search — starts live Trie-based suggestions
7. Press `:w` to save or `:q` to quit

---

## 🧪 Future Ideas

* Tabs for multiple file buffers
* Syntax highlighting (Go, Markdown, etc.)
* Clipboard integration (copy, paste)
* Configurable keybindings
* Plugin system
* Encrypted notes (AES)

---

## 💪 Dependencies

* [Bubbletea](https://github.com/charmbracelet/bubbletea) - TUI framework
* [Lipgloss](https://github.com/charmbracelet/lipgloss) - Styling (optional)
* Go 1.21+

---

## 👨‍💻 Author

Built by [Satyam Shree](https://github.com/satya-sudo) in stealth mode during mandatory office days 😎
