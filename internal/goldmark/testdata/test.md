# ndor examples

## beziel

```ndor
type Point struct {
    x, y float64
}

func points(ps []Point) {
    for p <- ps {
        circle p.x, p.y, 0.5
    }
    color "white"
    fill preserve
    color "black"
    linew 4
    stroke
}

func connect(ps []Point) {
    if len(ps) < 2 {
        return
    }
    from ps[0].x, ps[0].y
    for _, p := range ps {
        to p.x, p.y
    }
    color "#FF2D00"
    linew 8
    stroke
}

func bezierPoints(ps []Point) {
    if len(ps) < 3 {
        return
    }
    for len(ps) >= 3 {
        bezier ps[0].x, ps[0].y, ps[1].x, ps[1].y, ps[2].x, ps[2].y
        ps = ps[2:]
    }
    color "#3E606F"
    linew 16
    fill preserve
    color "black"
    stroke
}

color "white"
clear
scale 40, 40

ps := []Point{{-10, 0}, {-5, -10}, {0, 0}, {5, 10}, {10, 0}}
connect ps
bezierPoints ps
points ps
```

## clip

```ndor
circle -150, 0, 300
clip
circle 150, 0, 300
clip
rectangle -500, -500, 1000, 1000
color "black"
fill
```

## crisp

```ndor
const (
    half  = 500
    minor = 10
    major = 100
)

const (
    first = -half + minor
    top   = -half + major
)

func row(y int) {
    f := float64(y)
    from -half, f
    to half, f
}

func col(x int) {
    f := float64(x)
    from f, -half
    to f, half
}

linew 1

for x <- first:half:minor {
    col x
}
for y <- first:half:minor {
    row y
}
color 0, 0, 0, 55
stroke

for x <- top:half:major {
    col x
}
for y <- top:half:major {
    row y
}

color 0, 0, 0, 120
stroke
```

## fan

```ndor
scale 40, 40
from -10, 0
bezier -8, -8, 8, 8, 10, 0
color "orange"
linew 8
fill preserve
dash 1, 10
color "lightblue"
stroke
```

## flower

```ndor
const (
    dx, dy = 400, 80
    delta  = 30
    all    = 180
)

n := all / delta
degree := 0.0

color 0, 0, 0, 50

for i <- :n {
    push
    rotate 0, 0, degree
    ellipse 0, 0, dx, dy
    fill
    pop
    degree += delta
}
```

## geometry

```ndor
rectangle -250, -250, 500, 500
stroke

dash 5
rectangle -250, 190, 120, 60
rectangle 190, 130, 60, 120
rectangle 130, -250, 120, 60
rectangle -250, -250, 60, 120
stroke

dash // cacel dashes
from -130, 190
to 190, 130
to 130, -190
to -190, -130
to -130, 190
stroke

color "black"
text -260, 190, "c"
text -130, 255, "r"
text 180, 255, "n-1-c"
text 255, 130, "r"
text 255, -190, "m-1-c"
text 120, -265, "n-1-r"
text -195, -265, "c"
text -275, -130, "m-1-r"
```

## hellow world

```ndor
circle 0, 0, 300
color "lightgreen"
fill
```

## line width

```ndor
color "white"
clear
color "black"
w := 0.1
for i <- 100:901:20 {
    x := float64(i) - 500
    from x+50, 500
    to x-50, -500
    linew w
    stroke
    w += 0.1
}
```

## lines

```ndor
const width, height = 1024, 1024

color "snow"
clear

for i <- :100 {
    x1 := (rand.float64 - 0.5) * width
    y1 := (rand.float64 - 0.5) * height
    x2 := (rand.float64 - 0.5) * width
    y2 := (rand.float64 - 0.5) * height
    r := 255 * rand.float64
    g := 255 * rand.float64
    b := 255 * rand.float64
    a := 255 * (rand.float64*0.5 + 0.5)
    w := rand.float64*4 + 1
    color r, g, b, a
    linew w
    from x1, y1
    to x2, y2
    stroke
}
```

## open fill

```ndor
for j <- :10 {
    for i <- :10 {
        x := (float64(i) - 4.5) * 100
        y := (float64(j) - 4.5) * 100
        from := 360 * rand.float64
        to := from + 180*rand.float64 + 180
        arc x, y, 40, from, to
    }
}

color "orange"
fill preserve
color "green"
linew 8
stroke preserve
color "red"
linew 4
stroke preserve
```

## rainbow

```ndor
color "white"
clear

colors := [
    "red",
    "orange",
    "yellow",
    "green",
    "blue",
    "violet",
    "indigo",
    "white",
]

r := 400.0

for c <- colors {
    color c
    arc 0, 0, r, 0, 180
    fill
    r -= 50
}
```

## spiral

```ndor
const s = 1024
const n = 2048

color "black"

for i <- :n {
    t := float64(i) / n
    d := t*s*0.4 + 10
    a := math.Pi * 2 * t * 20
    x := math.cos(a) * d
    y := math.sin(a) * d
    r := t * 8
    circle x, y, r
}

fill
```

## star

```ndor
type point struct {
    x, y float64
}

func genPoints(n int, x, y, r float64) []point {
    result := make([]point, n)
    for i, _ <- result {
        angle := float64(i)*2*math.Pi/float64(n) - math.Pi/2
        result[i] = point{
            x: x + r*math.Cos(angle),
            y: y + r*math.Sin(angle),
        }
    }
    return result
}

n := 5
points := genPoints(n, 0, 0, 400)
for i <- :n + 1 {
    index := (i * 2) % n
    p := points[index]
    if i == 0 {
        from p.x, p.y
    } else {
        to p.x, p.y
    }
}
color 0, 127, 0
fill preserve
color 0, 127, 0, 130
linew 16
stroke
```
