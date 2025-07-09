package editor

type EditState struct {
	Lines [][]rune
}

type UndoManager struct {
	undoStack []EditState
	redoStack []EditState
}

func NewUndoManager() *UndoManager {
	return &UndoManager{
		undoStack: make([]EditState, 0),
		redoStack: make([]EditState, 0),
	}
}

type Undo interface {
	Undo(buffer *TextBuffer)
	Push(buffer *TextBuffer)
	Redo(buffer *TextBuffer)
}

func (um *UndoManager) Push(buffer *TextBuffer) {
	lines := buffer.Lines
	copyLines := copyBuffer(lines)
	um.undoStack = append(um.undoStack, EditState{copyLines})
	um.redoStack = nil
}

func (um *UndoManager) Undo(buffer *TextBuffer) {
	if len(um.undoStack) == 0 {
		return
	}
	lastState := um.undoStack[len(um.undoStack)-1]
	currentState := copyBuffer(buffer.Lines)
	buffer.Lines = lastState.Lines
	um.undoStack = um.undoStack[:len(um.undoStack)-1]
	um.redoStack = append(um.redoStack, EditState{currentState})
	buffer.SetDirty(true)
}
func (um *UndoManager) Redo(buffer *TextBuffer) {
	if len(um.redoStack) == 0 {
		return
	}
	lastState := um.redoStack[len(um.redoStack)-1]
	currentState := copyBuffer(buffer.Lines)
	buffer.Lines = lastState.Lines
	um.undoStack = append(um.undoStack, EditState{currentState})
	um.redoStack = um.redoStack[:len(um.redoStack)-1]
	buffer.SetDirty(true)
}

func copyBuffer(buffer [][]rune) [][]rune {
	copied := make([][]rune, len(buffer))
	for i, line := range buffer {
		copied[i] = make([]rune, len(line))
		copy(copied[i], line)
	}
	return copied
}
