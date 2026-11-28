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

## Rendering

You can set line style (solid/dashed), line width, and color for all shapes. You can also rotate shapes by rotating the canvas with push/pop operations. TODO
