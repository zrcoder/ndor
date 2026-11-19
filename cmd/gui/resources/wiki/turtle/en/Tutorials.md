## Square

Move forward then turn right 90 degrees, four times.

```c
color "lightgreen"
width 3
for range 4 {
    step 10
    right 90
}
```

## Star

Move forward and turn right 144 degrees, five times.

```c
color "orange"
for range 5 {
    step 10
    right 144
}
```

## Backward

Pass a negative number to `step`; with the pen down it draws in reverse.

```c
step 8
step -8
```

## Colored squares

Change colors in a loop to draw a ring of rotated squares.

```c
func squar() {
    step 5
    right 90
    step 5
    right 90
    step 5
    right 90
    step 5
}

colors := ["red", "green", "blue"]

for c <- colors {
    color c
    squar
    right 30
}
```
