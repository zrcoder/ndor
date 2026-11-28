2021, the Year of the Ox, I welcomed my first child and started writing this project. The goal is to use this program to spark my child's interest in programming.

This is a practical and fun app where you can draw by writing code.

## Usage

We draw in a 2D coordinate system using `code`.

`ndor` uses the [XGo](https://xgo.dev) programming language.

## Coordinate system

As in `turtle`, the origin is at the center of the picture, `x` grows to the right and `y` grows upwards.

The picture size needs no setting: it is sized to fit everything you draw.

## Hello World

```c
// Set current color to light green
color "lightgreen"

// Define a circle with center and radius
circle 0, 0, 300

// Fill the circle with current color
fill
```

> Lines starting with `"//"` or blocks wrapped in `"/*"` and `"*/"` are comments and are ignored.
>
> `x` coordinate goes left to right, `y` coordinate goes bottom to top, and angles are counter-clockwise.
