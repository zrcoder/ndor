package goldmark

import (
	"bytes"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
	"github.com/zrcoder/ndor/pkg"
)

type HTMLRenderer struct{}

func (r *HTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindBlock, r.Render)
}

func (r *HTMLRenderer) Render(w util.BufWriter, src []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*Block)
	if !entering {
		w.WriteString("</img>")
		return ast.WalkContinue, nil
	}

	b := bytes.Buffer{}
	lines := n.Lines()
	for i := 0; i < lines.Len(); i++ {
		line := lines.At(i)
		b.Write(line.Value(src))
	}

	if b.Len() == 0 {
		return ast.WalkContinue, nil
	}

	data, lineErr := pkg.Run(0, 0, b.String())
	if lineErr != nil {
		_, err := w.WriteString(lineErr.Msg)
		return ast.WalkContinue, err
	}

	w.WriteString(`<img src=`)
	_, err := w.WriteString(data)
	w.WriteString(">")
	return ast.WalkContinue, err
}
