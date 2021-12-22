package bebop

import "fmt"

// a locError is an error that occured at a specific file location
type locError struct {
	loc location
	err error
}

func (l locError) Unwrap() error {
	return l.err
}

func (l locError) Error() string {
	return fmt.Sprintf("[%d:%d] %s", l.loc.line, l.loc.lineChar, l.err)
}

// a location is a line and character-in-line pair
type location struct {
	lineChar int
	line     int
}

func (l *location) inc(i int) {
	l.lineChar += i
}

func (l *location) incLine() {
	l.line++
	l.lineChar = 0
}
