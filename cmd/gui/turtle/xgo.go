package turtle

import (
	"fmt"

	_ "github.com/zrcoder/ndor/cmd/gui/turtle/internal/exported/turtle/api"

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

	. "turtle/api"
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

// runCode compiles and runs the user's xgo code through ixgo.
// All turtle commands are appended to api.CmdQueue during execution.
// After execution, the caller drains CmdQueue frame by frame.
func runCode(code string) error {
	code = preCodes + code
	_, err := ixgo.RunFile(xgoFileName, code, nil, 0)
	if err != nil {
		return fmt.Errorf("xgo: %w", err)
	}
	return nil
}
