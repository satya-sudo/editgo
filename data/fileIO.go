package data

import (
	"bufio"
	"editGo/editor"
	"fmt"
	"os"
	"path/filepath"
)

type FileManager struct {
	FilePath string
	Buffer   *editor.TextBuffer
}

func NewFile(filePath string) (*FileManager, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	lines := [][]rune{}
	for scanner.Scan() {
		lines = append(lines, []rune(scanner.Text()))
	}
	if len(lines) == 0 {
		lines = [][]rune{{}} // ensure buffer isn't empty
	}
	buffer := editor.NewTextBufferWithLines(lines)
	fm := &FileManager{
		FilePath: filePath,
		Buffer:   buffer,
	}
	return fm, nil
}

func NewEmptyFile(filePath string) (*FileManager, error) {
	buffer := editor.NewTextBuffer()
	fm := &FileManager{
		Buffer:   buffer,
		FilePath: filePath,
	}
	return fm, nil
}

func (fm *FileManager) IsNewFile() bool {
	return fm.FilePath == ""
}

func (fm *FileManager) Save() error {
	if fm.FilePath == "" {
		return fmt.Errorf("file path is empty")
	}
	return fm.SaveAs(fm.FilePath)
}

func (fm *FileManager) SaveAs(filePath string) error {
	if err := ensurePath(filePath); err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	for _, line := range fm.Buffer.Lines {
		_, err := file.WriteString(string(line) + "\n")
		if err != nil {
			return err
		}
	}
	fm.FilePath = filePath
	fm.Buffer.SetDirty(false)

	return nil
}

func ensurePath(path string) error {
	dir := filepath.Dir(path)
	return os.MkdirAll(dir, 0755)
}
