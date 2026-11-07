package pkg

import (
	"github.com/zrcoder/ndor/internal"
)

// TurtlePlan is everything the GUI needs to replay one turtle drawing:
// canvas size and every pen-down segment with its color and width.
type TurtlePlan struct {
	W, H  int
	Steps []internal.TurtleStep
}

// TurtleGen runs a turtle snippet and returns its plan. It keeps the
// canvas gg.Context (so Color, fill and friends still work), but the
// visible output is only the turtle steps; ordinary gg drawing is not
// replayed.
func TurtleGen(width, height int, code string) (*TurtlePlan, *internal.LineError) {
	if err := gen(width, height, code); err != nil {
		return nil, err
	}
	steps := make([]internal.TurtleStep, len(internal.TurtleSteps))
	copy(steps, internal.TurtleSteps)
	return &TurtlePlan{W: width, H: height, Steps: steps}, nil
}
