package editor

type TextBuffer struct {
	Lines [][]rune
	Dirty bool
}

func NewTextBuffer() *TextBuffer {
	return &TextBuffer{
		Lines: [][]rune{{}},
		Dirty: false,
	}
}

type Buffer interface {
	InsertRune(line, col int, ch rune)
	DeleteRune(line, col int, ch rune)
	InsertNewLine(line, col int)
	MergeLine(line int)

	GetLine(line int) []rune
	LineCount() int
	IsDirty() bool
	SetDirty(dirty bool)
}

func (buffer *TextBuffer) InsertRune(line, col int, ch rune) {
	if buffer.validateIndex(line, col) {
		lineRunes := buffer.Lines[line]
		lineRunes = append(lineRunes[:col], append([]rune{ch}, lineRunes[col:]...)...)
		buffer.Lines[line] = lineRunes
		buffer.SetDirty(true)
	}
}
func (buffer *TextBuffer) DeleteRune(line, col int, ch rune) {
	if buffer.validateIndex(line, col) {
		if col == 0 { // delete new line
			if line == 0 {
				return // top of the buffer
			}
			buffer.MergeLine(line - 1) // merge into previous line
			buffer.SetDirty(true)
			return
		}
		// base case
		lineRunes := buffer.Lines[line]
		lineRunes = append(lineRunes[:col-1], lineRunes[col:]...)
		buffer.Lines[line] = lineRunes
		buffer.SetDirty(true)
	}
}
func (buffer *TextBuffer) InsertNewLine(line, col int) {
	if line < 0 || line >= len(buffer.Lines) {
		return
	}
	current := buffer.Lines[line]

	if col > len(current) {
		col = len(current)
	}
	before := current[:col]
	after := current[col:]
	buffer.Lines[line] = before

	buffer.Lines = append(
		buffer.Lines[:line+1], append([][]rune{after}, buffer.Lines[line+1:]...)...)
	buffer.SetDirty(true)
}
func (buffer *TextBuffer) MergeLine(line int) {
	if line < 0 || line+1 >= len(buffer.Lines) {
		return
	}
	buffer.Lines[line] = append(buffer.Lines[line], buffer.Lines[line+1]...)
	buffer.Lines = append(buffer.Lines[:line+1], buffer.Lines[line+2:]...)
	buffer.SetDirty(true)
}
func (buffer *TextBuffer) GetLine(line int) []rune {
	if line >= 0 && line < len(buffer.Lines) {
		return buffer.Lines[line]
	}
	return []rune{} // should panic
}
func (buffer *TextBuffer) LineCount() int {
	return len(buffer.Lines)
}
func (buffer *TextBuffer) IsDirty() bool {
	return buffer.Dirty
}
func (buffer *TextBuffer) SetDirty(dirty bool) {
	buffer.Dirty = dirty
}
func (buffer *TextBuffer) validateIndex(line, col int) bool {
	if line < 0 || line >= len(buffer.Lines) {
		return false
	}
	if col < 0 || col > len(buffer.Lines[line]) {
		return false
	}
	return true
}
