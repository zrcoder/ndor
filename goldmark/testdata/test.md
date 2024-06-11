# Test ndor

## 1. Flower

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

## 2. rainbow

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
