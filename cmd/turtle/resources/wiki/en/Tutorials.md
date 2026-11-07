## Draw a square

```c
context 400, 400
color "blue"
linew 3

penup
forward -100
right 90
forward -100
left 90
pendown

for i <- :4 {
    forward 200
    right 90
}
```

## Draw an equilateral triangle

```c
context 400, 400
color "green"
linew 3

penup
forward -80
right 30
pendown

for i <- :3 {
    forward 160
    right 120
}
```

## Draw a spiral

```c
context 600, 600
color "purple"
linew 2

for i <- :100 {
    forward float(i) * 3
    right 91
}
```

## Pen up / Pen down

Use `penup` to lift the pen so the turtle moves without drawing; use `pendown` to put it down again. `jump x, y` teleports the turtle to a coordinate without leaving a trail.
