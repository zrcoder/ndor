# Ndor examples

## Hello world (XGo)

```ndor
context 800, 800
circle 400, 400, 300
color "lightgreen"
fill
```

## Hello world (Go)

```ndor
func main() {
	context(800, 800)
	circle(400, 400, 300)
	color("lightgreen")
	fill()
}
```