// Package xgo interprets ndor XGo code.
package xgo

import (
	"strconv"
	"strings"

	_ "github.com/zrcoder/ndor/internal/exported/github.com/zrcoder/ndor/api/ndor"
	"github.com/zrcoder/ndor/internal/ndor"

	"github.com/goplus/ixgo"
	_ "github.com/goplus/ixgo/pkg/math"
	_ "github.com/goplus/ixgo/pkg/math/rand"
	_ "github.com/goplus/ixgo/xgobuild"
)

const (
	xgoFileName = "main.xgo"

	preCodes = `
	import (
		math "math"
		rand "math/rand"

		. "github.com/zrcoder/ndor/api/ndor"
	)

	var (
		_ = math.Pi
		_ rand.Source
	)

	const preserve = Preserve

	var (
		color   = Color
		dash    = Dash
		fill    = Fill
		stroke  = Stroke
		clip    = Clip
		pop     = Pop
		push    = Push
		clear   = Clear
	)
	`
)

// Run interprets code on top of the ndor prelude.
func Run(code string) *ndor.LineError {
	code = preCodes + code
	_, err := ixgo.RunFile(xgoFileName, code, nil, 0)
	return parseXGoErr(err)
}

// err like : ./main.xgo:1:1: undefined: vv
func parseXGoErr(err error) *ndor.LineError {
	if err == nil {
		return nil
	}
	msg := err.Error()
	_, after, found := strings.Cut(msg, xgoFileName)
	if !found {
		return &ndor.LineError{Number: -1, Msg: msg}
	}
	errUnexpected := &ndor.LineError{Number: -1, Msg: "unexpected internal error"}
	if len(after) == 0 || after[0] != ':' {
		return errUnexpected
	}
	after = after[1:]
	before, after, found := strings.Cut(after, ":")
	if !found {
		return errUnexpected
	}
	n, err := strconv.Atoi(before)
	if err != nil {
		return &ndor.LineError{Number: -1, Msg: err.Error()}
	}
	_, after, found = strings.Cut(after, ":")
	if !found {
		return errUnexpected
	}
	after = strings.TrimSpace(after)
	if i := strings.Index(after, "\n"); i != -1 {
		after = after[:i]
	}
	return &ndor.LineError{
		Number: n - strings.Count(preCodes, "\n"),
		Msg:    after,
	}
}
