// export by github.com/goplus/ixgo/cmd/qexp

package api

import (
	q "github.com/zrcoder/ndor/cmd/gui/turtle/api"

	"reflect"

	"github.com/goplus/ixgo"
)

func init() {
	ixgo.RegisterPackage(&ixgo.Package{
		Name: "api",
		Path: "turtle/api",
		Deps: map[string]string{
			"errors":                        "errors",
			"golang.org/x/image/colornames": "colornames",
			"image/color":                   "color",
		},
		Interfaces: map[string]reflect.Type{},
		NamedTypes: map[string]reflect.Type{},
		AliasTypes: map[string]reflect.Type{},
		Vars: map[string]reflect.Value{
			"CmdQueue": reflect.ValueOf(&q.CmdQueue),
		},
		Funcs: map[string]reflect.Value{
			"Backward": reflect.ValueOf(q.Backward),
			"Clear":    reflect.ValueOf(q.Clear),
			"Forward":  reflect.ValueOf(q.Forward),
			"Home":     reflect.ValueOf(q.Home),
			"Left":     reflect.ValueOf(q.Left),
			"Pen":      reflect.ValueOf(q.Pen),
			"PenColor": reflect.ValueOf(q.PenColor),
			"PenWidth": reflect.ValueOf(q.PenWidth),
			"Right":    reflect.ValueOf(q.Right),
		},
		TypedConsts:   map[string]ixgo.TypedConst{},
		UntypedConsts: map[string]ixgo.UntypedConst{},
	})
}
