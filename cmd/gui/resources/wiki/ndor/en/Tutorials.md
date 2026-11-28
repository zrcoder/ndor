## Points and Circles

You can't really draw a point, but you can draw a tiny circle instead:

```c
color "black"
circle 0, 0, 2
fill
```

## Lines

Strictly speaking, we draw line segments. Just specify start and end coordinates:

```c
from -20, -15
to 35, 50
stroke
```

## Curves

You can draw quadratic or cubic bezier curves. Specify start point, control point(s), and end point. For example:

```c
from 0, 0
bezier 0, 35, 25, 35
stroke
```

Add one more control point for a cubic bezier:

```c
from 0, 0
bezier 0, 25, 0, 45, 25, 35
stroke
```

## Ellipses

Similar to circles, specify center, horizontal radius, and vertical radius:

```c
ellipse 0, 0, 10, 20
fill
```

## Rectangles

Specify bottom-left corner and width/height:

```c
rectangle 0, 0, 90, 40
stroke
```

You can also draw rounded rectangles by adding a corner radius parameter:

```c
rectangle 0, 0, 90, 40, 10
stroke
```

## Arcs

Another way to draw curves is to draw an arc. First specify the circle parameters (center and radius), then the start and end angles (positive = counter-clockwise).

```c
arc 0, 0, 40, 45, 180
stroke
```

## Ellipse Arcs

You can also draw arcs on ellipses. Similar to arcs, specify ellipse parameters first, then angles:

```c
earc 0, 0, 40, 15, 0, 150
stroke
```

> Notice the drawing action is "stroke". Try changing it to "fill" and see what happens - which part gets filled?

## Regular Polygons

Zu Chongzhi used the polygon method to calculate pi. Conversely, how do we draw a regular polygon?

Draw a circle, divide it into n equal parts, connect the division points to form a regular n-gon.

So we just need to specify the number of sides and the circumcircle parameters:

```c
polygon 5, 0, 0, 40
stroke
```

This draws a regular pentagon.

> You can also add a rotation angle. Positive = counter-clockwise, negative = clockwise.
> The polygon above is flat on top. Let's rotate it:

```c
polygon 5, 0, 0, 40, 20
stroke
```

## Text

You can draw simple text like coordinates. Only supports letters, numbers and common symbols - no Chinese characters.

Specify the start (left end of the baseline) and text content:

```c
color "black"
text 0, 0, "Hello!"
```

## Clipping

`clip` clips everything drawn afterwards to the current path, then clears the path. Multiple clips intersect. For example, two overlapping circles define the region, then a rectangle fills it:

```c
circle -150, 0, 300
clip
circle 150, 0, 300
clip
rectangle -500, -500, 1000, 1000
color "black"
fill
```

## Saving and restoring state

You can set colors and line width for all shapes. Use `dash` for line style: no arguments restores a solid line, and the numbers are alternating on/off lengths, e.g. `dash 8, 4`.

`push` saves the current drawing state on a stack and `pop` restores the most recent one. The saved state covers the coordinate system, colors, dashes and line width, but not the path being built or drawings already made. Calls can nest, and `pop` on an empty stack does nothing. Use it whenever you want to change settings temporarily and restore them afterwards: colors, line width, dashes, or a temporary rotation/scale/translation of the coordinate system.

`rotate x, y, degree` rotates the coordinate system around `(x, y)`. The angle is in degrees, positive being counter-clockwise and negative clockwise. It affects only shapes and text drawn afterwards; anything already drawn stays put, and consecutive rotations accumulate. To give several shapes their own rotation without affecting each other, `push` before rotating and `pop` after drawing.

The example below deliberately rotates only 5 times: six 60° steps make a full turn, so even a missing `pop` would leave the final orientation unchanged and the check would prove nothing.

```c
color "steelblue"
linew 2

for i <- :5 {
    push
    rotate 0, 0, float64(i)*60
    rectangle 40, -10, 160, 20
    stroke
    pop
}

// the coordinate system is restored; a non-square rectangle makes any leftover rotation obvious
color "orange"
rectangle -8, -25, 16, 50
fill
```
