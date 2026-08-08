package hexempire

import "strconv"

// mapNumberInputMaxDigits matches the random map button's range
// (rand.Intn(999999) in game.go never produces more than 6 digits).
const mapNumberInputMaxDigits = 6

// MapNumberInput is a small single-line text editor for the map number
// button: it starts pre-filled with the current map number, and supports
// inserting/deleting at an arbitrary cursor position rather than only at
// the end. Its zero value is a valid "not currently editing" state.
type MapNumberInput struct {
	active bool
	digits []rune
	cursor int
}

// Begin starts editing, pre-filled with currentMapNumber and the cursor
// at the end, so a click-to-edit doesn't wipe out the existing value.
func (input *MapNumberInput) Begin(currentMapNumber int) {
	input.active = true
	input.digits = []rune(strconv.Itoa(currentMapNumber))
	input.cursor = len(input.digits)
}

func (input *MapNumberInput) Cancel() {
	input.active = false
	input.digits = nil
	input.cursor = 0
}

func (input *MapNumberInput) Active() bool {
	return input.active
}

func (input *MapNumberInput) Text() string {
	return string(input.digits)
}

// CursorIndex is a rune offset into Text(), in [0, len(Text())].
func (input *MapNumberInput) CursorIndex() int {
	return input.cursor
}

// InsertDigit inserts a typed digit at the cursor and advances it,
// ignoring non-digit runes and input past the length cap.
func (input *MapNumberInput) InsertDigit(r rune) {
	if r < '0' || r > '9' || len(input.digits) >= mapNumberInputMaxDigits {
		return
	}
	updated := make([]rune, 0, len(input.digits)+1)
	updated = append(updated, input.digits[:input.cursor]...)
	updated = append(updated, r)
	updated = append(updated, input.digits[input.cursor:]...)
	input.digits = updated
	input.cursor++
}

// DeleteBefore removes the digit immediately before the cursor (backspace).
func (input *MapNumberInput) DeleteBefore() {
	if input.cursor == 0 {
		return
	}
	input.digits = append(input.digits[:input.cursor-1], input.digits[input.cursor:]...)
	input.cursor--
}

// DeleteAfter removes the digit immediately after the cursor (forward delete/Del key).
func (input *MapNumberInput) DeleteAfter() {
	if input.cursor >= len(input.digits) {
		return
	}
	input.digits = append(input.digits[:input.cursor], input.digits[input.cursor+1:]...)
}

func (input *MapNumberInput) MoveLeft() {
	if input.cursor > 0 {
		input.cursor--
	}
}

func (input *MapNumberInput) MoveRight() {
	if input.cursor < len(input.digits) {
		input.cursor++
	}
}

func (input *MapNumberInput) MoveToStart() {
	input.cursor = 0
}

func (input *MapNumberInput) MoveToEnd() {
	input.cursor = len(input.digits)
}

// Confirm parses the typed digits into a map number. ok is false for an
// empty or unparseable buffer, leaving the input untouched so the player
// can keep editing.
func (input *MapNumberInput) Confirm() (mapNumber int, ok bool) {
	if len(input.digits) == 0 {
		return 0, false
	}
	number, err := strconv.Atoi(string(input.digits))
	if err != nil {
		return 0, false
	}
	return number, true
}
