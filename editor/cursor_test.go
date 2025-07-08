package editor

import "testing"

type mockBuffer struct {
	lines []string
}

func (m *mockBuffer) GetLine(line int) []rune {
	if line >= 0 && line < len(m.lines) {
		return []rune(m.lines[line])
	}
	return []rune{}
}
func (m *mockBuffer) LineCount() int {
	return len(m.lines)
}
func (m *mockBuffer) IsDirty() bool             { return false }
func (m *mockBuffer) SetDirty(bool)             {}
func (m *mockBuffer) InsertRune(int, int, rune) {}
func (m *mockBuffer) DeleteRune(int, int, rune) {}
func (m *mockBuffer) InsertNewLine(int, int)    {}
func (m *mockBuffer) MergeLine(int)             {}
func TestCursor_MoveLeft(t *testing.T) {
	buf := &mockBuffer{lines: []string{"hello", "world"}}
	cursor := NewCursor(1, 1)

	t.Run("Move left within line", func(t *testing.T) {
		cursor.MoveLeft(buf)
		x, y := cursor.GetPosition()
		if x != 0 || y != 1 {
			t.Errorf("Expected (0,1), got (%d,%d)", x, y)
		}
	})

	t.Run("Move left across line", func(t *testing.T) {
		cursor.MoveLeft(buf)
		x, y := cursor.GetPosition()
		if x != 5 || y != 0 {
			t.Errorf("Expected (5,0), got (%d,%d)", x, y)
		}
	})
}

func TestCursor_MoveRight(t *testing.T) {
	buf := &mockBuffer{lines: []string{"hi", "there"}}
	cursor := NewCursor(2, 0)

	t.Run("Move right across line", func(t *testing.T) {
		cursor.MoveRight(buf)
		x, y := cursor.GetPosition()
		if x != 0 || y != 1 {
			t.Errorf("Expected (0,1), got (%d,%d)", x, y)
		}
	})
}

func TestCursor_MoveUpDown(t *testing.T) {
	buf := &mockBuffer{lines: []string{"short", "muchlonger", "end"}}
	cursor := NewCursor(10, 1)

	t.Run("Move up and clamp X", func(t *testing.T) {
		cursor.MoveUp(buf)
		x, y := cursor.GetPosition()
		if x != 5 || y != 0 {
			t.Errorf("Expected (5,0), got (%d,%d)", x, y)
		}
	})

	t.Run("Move down and clamp X", func(t *testing.T) {
		cursor.MoveDown(buf)
		cursor.MoveDown(buf)
		x, y := cursor.GetPosition()
		if x != 3 || y != 2 {
			t.Errorf("Expected (3,2), got (%d,%d)", x, y)
		}
	})
}

func TestCursor_Clamp(t *testing.T) {
	buf := &mockBuffer{lines: []string{"one", "two"}}
	cursor := NewCursor(10, 5)

	t.Run("Clamp to last valid position", func(t *testing.T) {
		cursor.Clamp(buf)
		x, y := cursor.GetPosition()
		if y != 1 || x != 3 {
			t.Errorf("Expected clamped to (3,1), got (%d,%d)", x, y)
		}
	})
}
