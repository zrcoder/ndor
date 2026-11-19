// Package xgo interprets turtle XGo code.
package xgo

import (
	"fmt"

	_ "github.com/zrcoder/ndor/internal/exported/github.com/zrcoder/ndor/api/turtle"

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

	. "github.com/zrcoder/ndor/api/turtle"
)

var (
	_ = math.Pi
	_ rand.Source
)

var (
	step = Step
	left = Left
	right = Right
	color = PenColor
	width = PenWidth
	clear = Clear
	home = Home
	pen = Pen
)
`
)

// Run interprets code; turtle commands collect in the api CmdQueue.
func Run(code string) error {
	code = preCodes + code
	_, err := ixgo.RunFile(xgoFileName, code, nil, 0)
	if err != nil {
		return fmt.Errorf("xgo: %w", err)
	}
	return nil
}
