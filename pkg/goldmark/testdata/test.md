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

const s = 1200
context s, s
color "white"
clear
translate s/2, s/2
scale 40, 40

ps := []Point{{-10, 0}, {-5, -10}, {0, 0}, {5, 10}, {10, 0}}
connect ps
bezierPoints ps
points ps
```

## clip

```ndor
context 1000, 1000
circle 350, 500, 300
clip
circle 650, 500, 300
clip
rectangle 0, 0, 1000, 1000
color "black"
fill
```

## crisp

```ndor
const (
    width = 1000
    height = 1000
    minor = 10
    major = 100
)

func row(y int) {
    f := float64(y)
    from 0, f
    to width, f
}

func col(x int) {
    f := float64(x)
    from f, 0
    to f, height
}

context width, height
linew 1

for x <- minor:width:minor {
    col x
}
for y <- minor:height:minor {
    row y
}
color 0, 0, 0, 55
stroke

for x <- major:width:major {
    col x
}
for y <- major:width:major {
    row y
}

color 0, 0, 0, 120
stroke
```

## fan

```ndor
context 1000, 1000

translate 500, 500
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
    x, y = 500, 420
    dx, dy = 400, 80
    delta = 30
    all = 180
)

n := all/delta
degree := 0.0

context 1000, 1000
color 0, 0, 0, 50

for i <- :n {
    push
    rotate x, y, degree
    ellipse x, y, dx, dy
    fill
    pop
    degree += 30
}
```

## geometry

```ndor
context 600, 600

rectangle 50, 50, 500, 500
stroke

dash 5
rectangle 50, 50, 120, 60
rectangle 490, 50, 60, 120
rectangle 430, 490, 120, 60
rectangle 50, 430, 60, 120
stroke

dash // cacel dashes
from 170, 110
to 490, 170
to 430, 490
to 110, 430
to 170, 110
stroke

color "black"
text 40, 110, "c"
text 170, 45, "r"
text 480, 45, "n-1-c"
text 555, 170, "r"
text 555, 490, "m-1-c"
text 420, 565, "n-1-r"
text 105, 565, "c"
text 25, 430, "m-1-r"
```

## hellow world

```ndor
context 800, 800
circle 400, 400, 300
color "lightgreen"
fill
```

## line width

```ndor
context 1000, 1000
color "white"
clear
color "black"
w := 0.1
for i <- 100:901:20 {
    x := float64(i)
    from x+50, 0
    to x-50, 1000
    linew w
    stroke
    w += 0.1
}
```

## lines

```ndor
const width, height = 1024, 1024

context width, height

color "snow"
clear

for i <- :100 {
    x1 := rand.float64 * width
    y1 := rand.float64 * height
    x2 := rand.float64 * width
    y2 := rand.float64 * height
    r := 255 * rand.float64
    g := 255 * rand.float64
    b := 255 * rand.float64
    a := 255*(rand.float64*0.5 + 0.5)
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
context 1000, 1000

for j <- :10 {
    for i <- :10 {
        x := float64(i)*100+50
        y := float64(j)*100+50
        from := 360*rand.float64
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
context 1024, 512
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

const x, y = 512, 512
const from, to = 0, -180
r := 400.0

for c <- colors {
    color c
    arc x, y, r, from, to
    fill
    r -= 50
}
```

## spiral

```ndor
const s = 1024
const n = 2048

context s, s
color "black"

for i <- :n {
    t := float64(i) / n
    d := t*s*0.4 + 10
    a := math.Pi*2*t*20
    x := s/2 + math.cos(a)*d
    y := s/2 + math.sin(a)*d
    r := t*8
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
            x: x + r*math.cos(angle),
            y: y + r*math.sin(angle),
        }
    }
    return result
}

context 1024, 1024
n := 5
points := genPoints(n, 512, 512, 400)
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
