package internal

import (
	"errors"
	"fmt"
	"image"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/fogleman/gg"
)

const (
	commentPrefix = "//"
)

var ErrEmptyInput = errors.New("empty input")

type Painter struct {
	context *gg.Context
	lines   []lineInfo
	err     error
}

type lineInfo struct {
	Num    int
	Oper   operation
	Params []string
}

type operation func(painter *Painter, line lineInfo) error

var operations = map[string]operation{
	"context":   context,
	"rectangle": rectanger,
	"circle":    circle,
	"ellipse":   ellipse,
	"from":      moveTo,
	"to":        lineTo,
	"dash":      dash,
	"color":     color,
	"clip":      clip,
	"clear":     clear,
	"lineW":     setLineWidth,
	"text":      text,
	"arc":       arc,
	"earc":      ellipseArc,
	"polygon":   regularPolygon,
	"fill":      fill,
	"stroke":    stroke,
	"bezier":    bezier,
	"translate": translate,
	"scale":     scale,
	"rotate":    rotate,
	"push":      push,
	"pop":       pop,

	"画布":   context,
	"矩形":   rectanger,
	"圆":    circle,
	"椭圆":   ellipse,
	"从":    moveTo,
	"到":    lineTo,
	"虚线":   dash,
	"颜色":   color,
	"裁剪":   clip,
	"清除":   clear,
	"线宽":   setLineWidth,
	"文字":   text,
	"弧":    arc,
	"椭圆弧":  ellipseArc,
	"正多边形": regularPolygon,
	"填充":   fill,
	"描边":   stroke,
	"贝塞尔":  bezier,
	// "translate": translate,
	"伸缩": scale,
	"旋转": rotate,
	"载入": push,
	"载出": pop,
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func NewPainter(body string, preLines int) *Painter {
	res := &Painter{}
	lines := strings.Split(body, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 || strings.HasPrefix(line, commentPrefix) {
			continue
		}
		line, lineNum := parseLineNumber(line, preLines)
		arr := strings.Fields(line)
		oper, ok := operations[arr[0]]
		if !ok {
			res.err = genError(lineNum, "uncognized operation")
			return res
		}
		res.lines = append(res.lines, lineInfo{Num: lineNum, Oper: oper, Params: arr[1:]})
	}
	if len(res.lines) == 0 {
		res.err = ErrEmptyInput
	}
	return res
}

func parseLineNumber(line string, preLines int) (string, int) {
	i := strings.LastIndex(line, commentPrefix)
	if i == -1 {
		return line, -1
	}
	res := line[i:]
	res = strings.TrimPrefix(res, commentPrefix)
	res = strings.TrimSpace(res)
	n, _ := strconv.Atoi(res)
	return line[:i], n - preLines
}

func (p *Painter) Draw(w, h int) (image.Image, error) {
	if p.err != nil {
		return nil, p.err
	}
	p.context = gg.NewContext(int(w), int(h))
	var err error
	for _, line := range p.lines {
		err = line.Oper(p, line)
		if err != nil {
			fmt.Println(line.Num, err)
			return nil, &LineError{
				Msg:    err.Error(),
				Number: line.Num,
			}
		}
	}
	return p.context.Image(), nil
}
