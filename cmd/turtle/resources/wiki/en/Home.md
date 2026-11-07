# Turtle Graphics

A small app to guide kids into programming. You control a turtle with code to draw on the canvas.

## How to use

Draw turtle graphics in a 2D coordinate system with code.

`turtle` supports the [XGo](https://xgo.dev) programming language.

## Hello Turtle

```c
// 600×600 canvas
context 600, 600

// Set color to red, line width to 4
color "red"
linew 4

// Lift pen and move to start
penup
forward -150
right 90
forward -150
left 90
pendown

// Draw a square
for i <- :4 {
    forward 300
    right 90
}
```

> Lines starting with `"//"` or blocks between `"/*"` and `"*/"` are comments and ignored.
>
> `x` axis goes left to right, `y` axis goes bottom to top (canvas center is origin).
