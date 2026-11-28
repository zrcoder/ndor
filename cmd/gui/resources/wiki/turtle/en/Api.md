| Command | Description             | Args   | Notes                                                     |
| ------- | ----------------------- | ------ | --------------------------------------------------------- |
| step    | Move forward by n steps | 1      | Fractional steps allowed; 1 step = 1 grid cell; negative moves backward (`step -3`) |
| left    | Turn left               | 1      | Degrees, counter-clockwise                                |
| right   | Turn right              | 1      | Degrees, clockwise                                        |
| color   | Set pen color           | 1/3/4  | Color name ("red"), hex ("#f00"), or r, g, b[, a] (0–255, a defaults to 255; alpha supported) |
| width   | Set pen width           | 1      | Line width in pixels, default 2; applies to later strokes |
| pen     | Lift or lower the pen   | 0 or 1 | `"up"` moves without drawing; `"down"` or no args lowers it (default) |
| clear   | Clear canvas and reset  | 0      | Turtle returns to the origin facing up; color and width reset |
| home    | Return to origin        | 0      | Also restores the initial heading; keeps the drawing and pen style |
