package editor

import (
	"reflect"
	"testing"
)

func TestInsertRune(t *testing.T) {
	buf := NewTextBuffer()
	buf.InsertRune(0, 0, 'H')
	buf.InsertRune(0, 1, 'i')

	want := []rune("Hi")
	got := buf.GetLine(0)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("InsertRune failed: got %q, want %q", string(got), string(want))
	}
}

func TestDeleteRune(t *testing.T) {
	buf := NewTextBuffer()
	buf.InsertRune(0, 0, 'A')
	buf.InsertRune(0, 1, 'B')
	buf.DeleteRune(0, 2, 'B') // Deletes B

	want := []rune("A")
	got := buf.GetLine(0)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("DeleteRune failed: got %q, want %q", string(got), string(want))
	}
}

func TestInsertNewLine(t *testing.T) {
	buf := NewTextBuffer()
	buf.InsertRune(0, 0, 'H')
	buf.InsertRune(0, 1, 'i')
	buf.InsertNewLine(0, 1)

	if buf.LineCount() != 2 {
		t.Errorf("InsertNewLine failed: expected 2 lines, got %d", buf.LineCount())
	}

	if string(buf.GetLine(0)) != "H" || string(buf.GetLine(1)) != "i" {
		t.Errorf("InsertNewLine failed: line split incorrect")
	}
}

func TestMergeLine(t *testing.T) {
	buf := NewTextBuffer()
	buf.InsertRune(0, 0, 'A')
	buf.InsertNewLine(0, 1)
	buf.InsertRune(1, 0, 'B')
	buf.MergeLine(0)

	if buf.LineCount() != 1 {
		t.Errorf("MergeLine failed: expected 1 line, got %d", buf.LineCount())
	}

	if string(buf.GetLine(0)) != "AB" {
		t.Errorf("MergeLine failed: got %q", string(buf.GetLine(0)))
	}
}
