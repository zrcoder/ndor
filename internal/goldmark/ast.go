package goldmark

import "github.com/yuin/goldmark/ast"

type Block struct {
	ast.BaseBlock
}

func (n *Block) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}

var KindBlock = ast.NewNodeKind("NdorBlock")

func (n *Block) Kind() ast.NodeKind {
	return KindBlock
}
