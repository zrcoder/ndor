| Command | Description               | Args   | Notes                                                          |
| ------- | ------------------------- | ------ | -------------------------------------------------------------- |
| step    | Move forward by n steps   | 1      | 1 step = 1 grid cell; negative moves backward (`step -3`)      |
| left    | Turn left                 | 1      | Degrees                                                        |
| right   | Turn right                | 1      | Degrees                                                        |
| color   | Set pen color             | 1/3/4  | Color name ("red"), hex ("#f00"), or r, g, b[, a] (0–255)       |
| width   | Set pen width             | 1      | Line width                                                     |
| pen     | Lift or lower the pen     | 0 or 1 | `"up"` moves without drawing, `"down"` (default) draws         |
| clear   | Clear canvas and reset    | 0      |                                                                |
| home    | Return to origin          | 0      | Keeps the drawing                                              |

> `forward` and `backward` are merged into `step`: `step 3` moves forward 3, `step -3` moves back 3.
