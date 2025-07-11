package data

import (
	"editGo/editor"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestNewEmptyFile(t *testing.T) {
	fm, _ := NewEmptyFile()
	if fm.FilePath != "" {
		t.Errorf("expected empty path, got %s", fm.FilePath)
	}
	if fm.Buffer.LineCount() != 1 || len(fm.Buffer.GetLine(0)) != 0 {
		t.Errorf("expected buffer with one empty line")
	}
}

func TestSaveAsAndLoad(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.txt")

	fm, _ := NewEmptyFile()
	fm.Buffer.Lines = [][]rune{
		[]rune("hello"),
		[]rune("world"),
	}

	err := fm.SaveAs(filePath)
	if err != nil {
		t.Fatalf("SaveAs failed: %v", err)
	}

	// Confirm file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Errorf("expected file to exist after save")
	}

	// Now load it back
	loaded, err := NewFile(filePath)
	if err != nil {
		t.Fatalf("NewFile failed: %v", err)
	}

	got := getStringLines(loaded.Buffer)
	want := []string{"hello", "world"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Loaded buffer mismatch. Got %v, want %v", got, want)
	}
}

func TestSaveWithoutPath(t *testing.T) {
	fm, _ := NewEmptyFile()
	err := fm.Save()
	if err == nil {
		t.Errorf("expected error when saving without path")
	}
}

// Helper
func getStringLines(buf *editor.TextBuffer) []string {
	lines := []string{}
	for i := 0; i < buf.LineCount(); i++ {
		lines = append(lines, string(buf.GetLine(i)))
	}
	return lines
}
