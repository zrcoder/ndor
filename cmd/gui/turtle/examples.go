package turtle

// example 是 goplayspace 里的海龟示例指令（xgo 命令式语法）。
type example struct {
	name string
	code string
}

var examples = []example{
	{
		name: "Star（五角星）",
		code: `
for range 5 {
    step 10
    right 144
}`,
	},
	{
		name: "Tree（树）",
		code: `
func branch(size, ratio, angle float64, iterations int) {
    color "black"
    width size
    pen "down"
    step size
    pen "up"

    if iterations > 0 {
        left angle
        branch size*ratio, ratio, angle, iterations-1

        right angle*2
        branch size*ratio, ratio, angle, iterations-1

        left angle
    }

    left 180
    step size
    left 180
}

left 90
branch 5, 0.7, 30, 5
`,
	},
	{
		name: "Colored squares（彩色正方形）",
		code: `func squar() {
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
`,
	},
}
