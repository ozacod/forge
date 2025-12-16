package terminal

// ANSI terminal control escape sequences

// ShowCursor makes the cursor visible
const ShowCursor = "\033[?25h"

// HideCursor makes the cursor invisible
const HideCursor = "\033[?25l"

// ClearLine clears the current line and returns cursor to the beginning
// This is equivalent to "\r\033[2K"
const ClearLine = "\r\033[2K"

// ClearLineToEnd clears from cursor position to end of line
const ClearLineToEnd = "\033[0K"

// ClearLineFromStart clears from start of line to cursor position
const ClearLineFromStart = "\033[1K"
