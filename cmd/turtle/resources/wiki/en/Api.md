| Operation | Description | Args | Details |
| --------- | ----------- | ---- | ------- |
| context | Set canvas size | 2 | Canvas width and height |
| color | Set current color | 1 or 3 or 4 | Color name ("red") or hex (#f00, #ff00ff) or rgba values 0-255, alpha defaults to 255 |
| linew | Set line width | 1 | Line width |
| turtle | Place the turtle | 2 | Turtle position, canvas center as origin, y axis pointing up |
| forward | Move forward and draw | 1 | Pixels |
| back | Move backward | 1 | Pixels |
| left | Turn left | 1 | Degrees |
| right | Turn right | 1 | Degrees |
| angle | Set absolute heading | 1 | Degrees, 0 = up, positive = clockwise |
| jump | Lift pen and move to point | 2 | Pixel coordinates |
| penup | Lift pen | 0 | Move without drawing |
| pendown | Put pen down | 0 | Move while drawing |
| home | Return to origin facing up, pen down | 0 | |

Turtle graphics are driven by xgo syntax, with no custom language. Running replays the turtle trail step by step as an animation over a coordinate grid.
