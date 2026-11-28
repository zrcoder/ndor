2021, the Year of the Ox, I welcomed my first child and started writing this project. The goal is to use this program to spark my child's interest in programming.

This is a practical and fun app where you can draw by writing code.

## Usage

We draw in a 2D coordinate system.

`ndor` uses the [XGo](https://xgo.dev) programming language. Write code in the editor on the right, then click the run button at the bottom right (or press `Ctrl+R`). The help button at the bottom left (or `Ctrl+H`) opens the docs and examples.

Commands such as `circle` and `rectangle` only add shapes to the current path; nothing is painted until you call `fill`, `stroke` or `clip`. `fill` and `stroke` clear the path afterwards unless you pass `preserve`.

## Hello World

```c
// Set current color to light green
color "lightgreen"

// Add a circle with center and radius to the current path
circle 0, 0, 300

// Fill the current path with the current color
fill
```

> Lines starting with `"//"` or blocks wrapped in `"/*"` and `"*/"` are comments and are ignored.
>
> `x` goes left to right, `y` goes bottom to top; angles are in degrees and positive means counter-clockwise.
