2021, the Year of the Ox, I welcomed my first child and started writing this project. The goal is to use this program to spark my child's interest in programming.

This is a practical and fun app where you can draw by writing code.

## Usage

We draw in a 2D coordinate system using `code`.

`ndor` supports both [XGo](https://xgo.dev) and [Go](https://go.dev) programming languages.

## Hello World

```c
// 800×800 canvas, can be omitted, defaults to the size of the picture frame
context 800, 800

// Set current color to light green
color "lightgreen"

// Define a circle with center and radius
circle 400, 400, 300

// Fill the circle with current color
fill
```

> Lines starting with `"//"` or blocks wrapped in `"/*"` and `"*/"` are comments and are ignored.
>
> `x` coordinate goes left to right, `y` coordinate goes top to bottom
