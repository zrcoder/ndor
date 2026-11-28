| Operation | Description | Args | Details |
| --------- | ----------- | ---- | ------- |
| rectangle | Rectangle | 4 or 5 | Bottom-left x, y, width, height; 5th param is corner radius for rounded rect |
| circle | Circle | 3 | Center x, y and radius |
| ellipse | Ellipse | 4 | Center x, y, horizontal radius, vertical radius |
| polygon | Regular polygon | 4 or 5 | Sides, circumcircle center x, y, radius; 5th param rotates the polygon in degrees (positive = counter-clockwise) |
| from | Move pen to point | 2 | x, y coordinates |
| to | Draw line to point | 2 | End point x, y coordinates |
| bezier | Draw bezier curve to point | 4 or 6 | 2 or 3 points (x,y pairs) for quadratic/cubic bezier |
| arc | Arc | 5 | Center x, y, radius, start/end angles in degrees (positive = counter-clockwise) |
| earc | Ellipse arc | 6 | Center x, y, horizontal radius, vertical radius, start/end angles in degrees |
| dash | Set dash pattern | variable | No args = solid line; otherwise alternating on/off lengths, e.g. `dash 8, 4` |
| text | Text | 3 | x, y of the start (left end of the baseline) and text content |
| color | Set current color | 1 or 3 or 4 | Color name ("red") or hex (#f00, #ff00ff) or rgba values 0-255, alpha defaults to 255 |
| fill | Fill current path with current color | 0 or 1 | Clears the path afterwards; "preserve" keeps it |
| stroke | Stroke current path with current color | 0 or 1 | Clears the path afterwards; "preserve" keeps it |
| clip | Clip later drawing to current path | 0 | Clears the path; multiple clips intersect |
| clear | Drop what was drawn, use the current color as the picture background | 0 | |
| linew | Set line width | 1 | Line width |
| rotate | Rotate the coordinate system around a point | 3 | Pivot x, y and angle in degrees (positive = counter-clockwise); later drawing only, consecutive calls accumulate |
| translate | Translate the coordinate system | 2 | x, y offset of the new origin; later drawing only |
| scale | Scale the coordinate system | 2 | x and y scale factors; later drawing only |
| push | Save the current drawing state (push) | 0 | Saves the coordinate system, colors, dashes and line width |
| pop | Restore the latest saved state (pop) | 0 | Restores that state; does nothing on an empty stack and leaves existing drawings alone |

## Rotation and the state stack

`rotate x, y, degree` rotates the coordinate system around `(x, y)`; `degree` is in degrees, positive being counter-clockwise and negative clockwise. The pivot is given in the current coordinate system. Rotation affects only shapes and text drawn afterwards — anything already drawn stays put. Consecutive rotations accumulate.

`push` saves the current drawing state on a stack, and `pop` restores the most recent one. The saved state includes the coordinate system (rotation, scale and translation), fill/stroke/text colors, dash pattern and line width. The path being built is not saved, and drawings already made are not undone. Calls can nest; `pop` on an empty stack does nothing.

So use `push`/`pop` whenever you want to change something temporarily and restore it afterwards: colors, line width, dashes, or a temporary `rotate`/`scale`/`translate`. Rotation is just one case: the common pattern is `push` → `rotate` → draw → `pop`, so each rotation starts from the same baseline and settings do not leak into later drawing. See `internal/examples/ndor/rotate.xgo` for a rotation example.

