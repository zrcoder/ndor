// export by github.com/goplus/igop/cmd/qexp

package api

import (
	q "github.com/zrcoder/niudour/api"

	"go/constant"
	"reflect"

	"github.com/goplus/igop"
)

func init() {
	igop.RegisterPackage(&igop.Package{
		Name: "api",
		Path: "github.com/zrcoder/niudour/api",
		Deps: map[string]string{
			"github.com/fogleman/gg":              "gg",
			"github.com/zrcoder/niudour/internal": "internal",
		},
		Interfaces: map[string]reflect.Type{},
		NamedTypes: map[string]reflect.Type{},
		AliasTypes: map[string]reflect.Type{},
		Vars:       map[string]reflect.Value{},
		Funcs: map[string]reflect.Value{
			"Arc":       reflect.ValueOf(q.Arc),
			"Bezier":    reflect.ValueOf(q.Bezier),
			"Circle":    reflect.ValueOf(q.Circle),
			"Clear":     reflect.ValueOf(q.Clear),
			"Clip":      reflect.ValueOf(q.Clip),
			"Color":     reflect.ValueOf(q.Color),
			"Context":   reflect.ValueOf(q.Context),
			"Dash":      reflect.ValueOf(q.Dash),
			"Earc":      reflect.ValueOf(q.Earc),
			"Ellipse":   reflect.ValueOf(q.Ellipse),
			"Fill":      reflect.ValueOf(q.Fill),
			"From":      reflect.ValueOf(q.From),
			"Linew":     reflect.ValueOf(q.Linew),
			"Polygon":   reflect.ValueOf(q.Polygon),
			"Pop":       reflect.ValueOf(q.Pop),
			"Push":      reflect.ValueOf(q.Push),
			"Rectangle": reflect.ValueOf(q.Rectangle),
			"Rotate":    reflect.ValueOf(q.Rotate),
			"Scale":     reflect.ValueOf(q.Scale),
			"Stroke":    reflect.ValueOf(q.Stroke),
			"Text":      reflect.ValueOf(q.Text),
			"To":        reflect.ValueOf(q.To),
			"Translate": reflect.ValueOf(q.Translate),
		},
		TypedConsts: map[string]igop.TypedConst{},
		UntypedConsts: map[string]igop.UntypedConst{
			"Preserve": {"untyped string", constant.MakeString(string(q.Preserve))},
		},
	})
}
