| Operation | Description | Args | Details |
| --------- | ----------- | ---- | ------- |
| context | Set canvas size | 2 | Canvas width and height |
| rectangle | Rectangle | 4 or 5 | Left-top x, y, width, height; 5th param is corner radius for rounded rect |
| circle | Circle | 3 | Center x, y and radius |
| ellipse | Ellipse | 4 | Center x, y, horizontal radius, vertical radius |
| polygon | Regular polygon | 4 or 5 | Sides, circumcircle center x, y, radius; 5th param rotates the polygon by that angle (positive = clockwise) |
| from | Move pen to point | 2 | x, y coordinates |
| to | Draw line to point | 2 | End point x, y coordinates |
| bezier | Draw bezier curve to point | 4 or 6 | 2 or 3 points (x,y pairs) for quadratic/cubic bezier |
| arc | Arc | 5 | Center x, y, radius, start angle, end angle (positive = clockwise) |
| earc | Ellipse arc | 6 | Center x, y, horizontal radius, vertical radius, start angle, end angle |
| dash | Set dash pattern | variable | No args = solid line; otherwise dash lengths |
| text | Text | 3 | x, y position and text content |
| color | Set current color | 1 or 3 or 4 | Color name ("red") or hex (#f00, #ff00ff) or rgba values 0-255, alpha defaults to 255 |
| fill | Fill with current color | 0 or 1 | Optional "preserve" arg to keep path after fill |
| stroke | Stroke with current color | 0 or 1 | Optional "preserve" arg to keep path after stroke |
| clip | Clip to path | 0 | |
| clear | Clear canvas with current color | 0 | |
| linew | Set line width | 1 | Line width |
| rotate | Rotate around point | 3 | Point x, y and rotation angle (positive = clockwise) |
| push | Save canvas state | 0 | |
| pop | Restore canvas state | 0 | |
