## Points and Circles

You can't really draw a point, but you can draw a tiny circle instead:

```c
color black
circle 50, 50, 2
fill
```

## Lines

Strictly speaking, we draw line segments. Just specify start and end coordinates:

```c
from 10, 15
to 35, 50
stroke
```

## Curves

You can draw quadratic or cubic bezier curves. Specify start point, control point(s), and end point. For example:

```c
from 10, 15
bezier 10, 50, 35, 50
stroke
```

Add one more control point for a cubic bezier:

```c
from 10, 15
bezier 10, 40, 10, 60, 35, 50
stroke
```

## Ellipses

Similar to circles, specify center, horizontal radius, and vertical radius:

```c
ellipse 50, 50, 10, 20
fill
```

## Rectangles

Specify top-left corner and width/height:

```c
rectangle 5, 5, 90, 40
stroke
```

You can also draw rounded rectangles by adding a corner radius parameter:

```c
rectangle 5, 5, 90, 40, 10
stroke
```

## Arcs

Another way to draw curves is to draw an arc. First specify the circle parameters (center and radius), then the start and end angles (positive = clockwise).

```c
arc 50, 50, 40, 45, 180
stroke
```

## Ellipse Arcs

You can also draw arcs on ellipses. Similar to arcs, specify ellipse parameters first, then angles:

```c
earc 50, 50, 40, 15, 0, 150
stroke
```

> Notice the drawing action is "stroke". Try changing it to "fill" and see what happens - which part gets filled?

## Regular Polygons

Zu Chongzhi used the polygon method to calculate pi. Conversely, how do we draw a regular polygon?

Draw a circle, divide it into n equal parts, connect the division points to form a regular n-gon.

So we just need to specify the number of sides and the circumcircle parameters:

```c
polygon 5, 50, 50, 40
stroke
```

This draws a regular pentagon.

> You can also add a rotation angle. Positive = clockwise, negative = counter-clockwise.
> The polygon above is flat on the bottom. Let's rotate it:

```c
polygon 5, 50, 50, 40, 20
stroke
```

## Text

You can draw simple text like coordinates. Only supports letters, numbers and common symbols - no Chinese characters.

Specify the top-left position and text content:

```c
color "black"
text 100, 100, "Hello!"
```

## Rendering

You can set line style (solid/dashed), line width, and color for all shapes. You can also rotate shapes by rotating the canvas with push/pop operations. TODO
