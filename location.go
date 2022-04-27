package bebop

import "fmt"

// a locError is an error that occurred at a specific file location
type locError struct {
	loc location
	err error
}

func (l locError) Unwrap() error {
	return l.err
}

func (l locError) Error() string {
	return fmt.Sprintf("%s %s", l.loc.String(), l.err)
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

func (l location) String() string {
	return fmt.Sprintf("[%d:%d]", l.line, l.lineChar)
}
