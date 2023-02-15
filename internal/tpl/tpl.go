package main

import (
	"fmt"
	"runtime"

	math "math"
	rand "math/rand"
)

var (
	pi = math.Pi
	_  rand.Source
)

const (
	preserve = "preserve"
)

func context(x, y any) {
	helpPrint("context", x, y)
}

func color(args ...any) {
	helpPrint("color", args...)
}

func rectangle(x, y, w, h any, r ...any) {
	if len(r) == 0 {
		helpPrint("rectangle", x, y, w, h)
	} else {
		helpPrint("rectangle", x, y, w, h, r[0])
	}
}

func circle(x, y, r any) {
	helpPrint("circle", x, y, r)
}

func clip() {
	helpPrint("clip")
}

func ellipse(x, y, rx, ry any) {
	helpPrint("ellipse", x, y, rx, ry)
}

func from(x, y any) {
	helpPrint("from", x, y)
}

func to(x, y any) {
	helpPrint("to", x, y)
}

func dash(args ...any) {
	helpPrint("dash", args...)
}

func lineW(width any) {
	helpPrint("lineW", width)
}

func text(x, y any, s string) {
	helpPrint("text", x, y, s)
}

func arc(x, y, r, angle1, angle2 any) {
	helpPrint("arc", x, y, r, angle1, angle2)
}

func earc(x, y, rx, ry, angle1, angle2 any) {
	helpPrint("earc", x, y, rx, ry, angle1, angle2)
}

func polygon(n any, x, y, r any, degree ...any) {
	if len(degree) == 0 {
		helpPrint("polygon", x, y, r)
	} else {
		helpPrint("polygon", x, y, r, degree[0])
	}
}

func translate(x, y any) {
	helpPrint("translate", x, y)
}

func scale(x, y any) {
	helpPrint("scale", x, y)
}

func rotate(x, y, angle any) {
	helpPrint("rotate", x, y, angle)
}

func bezier(x1, y1, x2, y2 any, p ...any) {
	if len(p) == 2 {
		helpPrint("bezier", x1, y1, x2, y2, p[0], p[1])
	} else {
		helpPrint("bezier", x1, y1, x2, y2)
	}
}

func fill(preserve ...any) {
	if len(preserve) == 0 {
		helpPrint("fill")
	} else {
		helpPrint("fill", preserve...)
	}
}

func stroke(preserve ...any) {
	if len(preserve) == 0 {
		helpPrint("stroke")
	} else {
		helpPrint("stroke", preserve...)
	}
}

func clear() {
	helpPrint("clear")
}

func push() {
	helpPrint("push")
}

func pop() {
	helpPrint("pop")
}

func helpPrint(name string, args ...any) {
	fmt.Print(name)
	for _, arg := range args {
		fmt.Print(" ")
		fmt.Print(arg)
	}
	_, _, line, _ := runtime.Caller(2)
	fmt.Println(" //", line)
}
