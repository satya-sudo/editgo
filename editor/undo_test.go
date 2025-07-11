package editor

import (
	"reflect"
	"testing"
)

func makeBufferWithLines(lines []string) *TextBuffer {
	tb := NewTextBuffer()
	tb.Lines = [][]rune{}
	for _, line := range lines {
		tb.Lines = append(tb.Lines, []rune(line))
	}
	return tb
}

func getStringLines(tb *TextBuffer) []string {
	lines := []string{}
	for _, r := range tb.Lines {
		lines = append(lines, string(r))
	}
	return lines
}

func TestUndoManager_Push(t *testing.T) {
	um := NewUndoManager()
	buffer := makeBufferWithLines([]string{"one", "two"})

	um.Push(buffer)

	if len(um.undoStack) != 1 {
		t.Errorf("Expected undo stack length 1, got %d", len(um.undoStack))
	}
}

func TestUndoManager_UndoRedo(t *testing.T) {
	um := NewUndoManager()
	buffer := makeBufferWithLines([]string{"line1"})

	// Push original
	um.Push(buffer)

	// Simulate edit
	buffer.Lines[0] = []rune("changed")

	// Undo should bring back "line1"
	um.Undo(buffer)
	if got := getStringLines(buffer); !reflect.DeepEqual(got, []string{"line1"}) {
		t.Errorf("Undo failed. Got %v", got)
	}

	// Redo should bring back "changed"
	um.Redo(buffer)
	if got := getStringLines(buffer); !reflect.DeepEqual(got, []string{"changed"}) {
		t.Errorf("Redo failed. Got %v", got)
	}
}
func TestUndoManager_EmptyUndoRedo(t *testing.T) {
	um := NewUndoManager()
	buffer := makeBufferWithLines([]string{"initial"})

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Expected no panic on empty undo/redo, but got %v", r)
		}
	}()

	um.Undo(buffer)
	um.Redo(buffer)
}

func TestUndoManager_Redo(t *testing.T) {
	um := NewUndoManager()
	buffer := makeBufferWithLines([]string{"original"})

	um.Push(buffer)                    // Push original
	buffer.Lines[0] = []rune("edited") // Edit 1

	um.Push(buffer)                   // Push edited
	buffer.Lines[0] = []rune("final") // Edit 2

	um.Undo(buffer) // back to "edited"
	if got := getStringLines(buffer); !reflect.DeepEqual(got, []string{"edited"}) {
		t.Errorf("Undo to edited failed. Got %v", got)
	}

	um.Redo(buffer) // back to "final"
	if got := getStringLines(buffer); !reflect.DeepEqual(got, []string{"final"}) {
		t.Errorf("Redo to final failed. Got %v", got)
	}
}
