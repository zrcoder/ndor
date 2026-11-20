// Package xgoerr runs XGo code and turns ixgo errors into line errors
// shared by ndor and turtle.
package xgoerr

import (
	"errors"
	"go/scanner"
	"strconv"
	"strings"

	"github.com/goplus/gogen"
	"github.com/goplus/ixgo"
)

const (
	// xgoFile is the file name ixgo reports errors about.
	xgoFile = "main.xgo"

	// lineDirective makes ixgo report lines in the user's code, not in
	// the prelude.
	lineDirective = "\n//line " + xgoFile + ":1\n"
)

// LineError is a compile or runtime error with its line in the user's code.
type LineError struct {
	Number int
	Flag   string
	Msg    string
}

// Error formats the message with the line, when known.
func (e *LineError) Error() string {
	if e.Number > 0 {
		return "Line " + strconv.Itoa(e.Number) + ": " + e.Msg
	}
	return e.Msg
}

// Run compiles code after prelude.
func Run(code, prelude string) *LineError {
	_, err := ixgo.RunFile(xgoFile, prelude+lineDirective+code, nil, 0)
	if err == nil {
		return nil
	}
	return New(err)
}

// New returns err as a LineError, reading the line and message from the
// ixgo error itself.
func New(err error) *LineError {
	if codeErr, ok := errors.AsType[*gogen.CodeError](err); ok {
		return &LineError{Number: codeErr.Fset.Position(codeErr.Pos).Line, Msg: firstLine(codeErr.Msg)}
	}
	var list scanner.ErrorList
	if errors.As(err, &list) && len(list) > 0 {
		return &LineError{Number: list[0].Pos.Line, Msg: firstLine(list[0].Msg)}
	}
	return &LineError{Number: -1, Msg: firstLine(err.Error())}
}

// FindLine returns the line of the first statement whose leading word is
// flag, or -1 when none does.
func FindLine(code, flag string) int {
	j := 0
	for line := range strings.SplitSeq(code, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			j++
			continue
		}
		before, _, found := strings.Cut(line, " ")
		if !found {
			before, _, found = strings.Cut(line, "(")
		}
		if found {
			line = before
		}
		if strings.ToLower(line[:1])+line[1:] == flag {
			return j + 1
		}
		j++
	}
	return -1
}

func firstLine(s string) string {
	s, _, _ = strings.Cut(s, "\n")
	return s
}
