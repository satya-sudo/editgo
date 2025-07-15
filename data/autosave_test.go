package data

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestAutoSave_DoesNotStartForEmptyPath(t *testing.T) {
	fm, _ := NewEmptyFile("")
	auto := NewAutoSave(fm, 100*time.Millisecond)
	auto.Start()

	// Should not autosave because file path is empty
	time.Sleep(200 * time.Millisecond)
	// No panic, no effect = success
	auto.Stop()
}

func TestAutoSave_SavesWhenDirty(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "autosave_test.txt")

	// Create file with one line
	fm, err := NewEmptyFile("")
	if err != nil {
		t.Fatalf("failed to create buffer: %v", err)
	}
	fm.FilePath = filePath
	fm.Buffer.Lines = [][]rune{[]rune("first")}
	fm.Buffer.SetDirty(true) // trigger autosave

	auto := NewAutoSave(fm, 200*time.Millisecond)
	auto.Start()

	// Wait enough time for autosave to trigger
	time.Sleep(300 * time.Millisecond)
	auto.Stop()

	// Check that file exists and content is saved
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("file not written: %v", err)
	}

	got := strings.TrimSpace(string(content))
	want := "first"
	if got != want {
		t.Errorf("autosave failed, got '%s', want '%s'", got, want)
	}
}

func TestAutoSave_Stop(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "stop_test.txt")

	fm, _ := NewEmptyFile("")
	fm.FilePath = filePath
	fm.Buffer.Lines = [][]rune{[]rune("stop test")}
	fm.Buffer.SetDirty(true)

	auto := NewAutoSave(fm, 100*time.Millisecond)
	auto.Start()
	auto.Stop()

	// Wait to confirm no goroutine is running after stop
	time.Sleep(200 * time.Millisecond)

	// Dirty buffer should still be dirty because autosave stopped early
	if !fm.Buffer.IsDirty() {
		t.Errorf("expected buffer to remain dirty after stop")
	}
}
