// Package xgo interprets ndor XGo code.
package xgo

import (
	_ "github.com/zrcoder/ndor/internal/exported/github.com/zrcoder/ndor/api/ndor"
	"github.com/zrcoder/ndor/internal/xgoerr"

	_ "github.com/goplus/ixgo/pkg/math"
	_ "github.com/goplus/ixgo/pkg/math/rand"
	_ "github.com/goplus/ixgo/xgobuild"
)

const preCodes = `
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

// Run interprets code on top of the ndor prelude.
func Run(code string) *xgoerr.LineError { return xgoerr.Run(code, preCodes) }
