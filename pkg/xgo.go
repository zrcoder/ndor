package pkg

import (
	"strconv"
	"strings"

	"github.com/zrcoder/ndor/internal"
	_ "github.com/zrcoder/ndor/internal/exported/github.com/zrcoder/ndor/api"

	"github.com/goplus/ixgo"
	_ "github.com/goplus/ixgo/pkg/math"
	_ "github.com/goplus/ixgo/pkg/math/rand"
	_ "github.com/goplus/ixgo/xgobuild"
)

const (
	xgofileName = "main.xgo"

	preCodes = `
	import (
		math "math"
		rand "math/rand"

		. "github.com/zrcoder/ndor/api"
	)

	var (
		_ = math.Pi
		_ rand.Source
	)

	const preserve = Preserve

	var (
		color  = Color
		dash   = Dash
		fill   = Fill
		stroke = Stroke
		clip   = Clip
		pop    = Pop
		push   = Push
		clear  = Clear
	)
	`
)

func xgoRun(code string) *internal.LineError {
	code = preCodes + code
	_, err := ixgo.RunFile(xgofileName, code, nil, 0)
	return parseXGoErr(err)
}

// err like : ./main.xgo:1:1: undefined: vv
func parseXGoErr(err error) *internal.LineError {
	if err == nil {
		return nil
	}
	msg := err.Error()
	i := strings.Index(msg, xgofileName)
	if i == -1 {
		return &internal.LineError{Number: -1, Msg: msg}
	}
	errUnexpected := &internal.LineError{Number: -1, Msg: "unexpected internal error"}
	msg = msg[i+len(xgofileName):]
	if len(msg) == 0 || msg[0] != ':' {
		return errUnexpected
	}
	msg = msg[1:]
	i = strings.Index(msg, ":")
	if i == -1 {
		return errUnexpected
	}
	n, err := strconv.Atoi(msg[:i])
	if err != nil {
		return &internal.LineError{Number: -1, Msg: err.Error()}
	}
	msg = msg[i+1:]
	i = strings.Index(msg, ":")
	if i == -1 {
		return errUnexpected
	}
	msg = strings.TrimSpace(msg[i+1:])
	i = strings.Index(msg, "\n")
	if i != -1 {
		msg = msg[:i]
	}
	return &internal.LineError{
		Number: n - strings.Count(preCodes, "\n"),
		Msg:    msg,
	}
}
