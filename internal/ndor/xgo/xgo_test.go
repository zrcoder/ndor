package xgo

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/zrcoder/ndor/internal/xgoerr"
)

func TestRun(t *testing.T) {
	codes := []string{
		`context 100, 100
color "black"`,
		`context 100, 100
color`,
		`context 100, 100
Color`,
		`Context 800, 800
Color "lightgreen"
Circle 400, 400, 300
Fill`,
		`context 800, 800
color "lightgreen"
circle 400, 400, 300
fill`,
		`context 1000, 1000
color "notExistColorName"`,
	}
	for i, code := range codes {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			err := Run(code)
			if err != nil {
				t.Error("code:\n", code)
				t.Error(err)
			}
		})
	}
}

func TestRunErrors(t *testing.T) {
	tests := []struct {
		name string
		code string
		want *xgoerr.LineError
	}{
		{
			name: "undefined identifier",
			code: "context 100, 100\ncc",
			want: &xgoerr.LineError{Number: 2, Msg: "undefined: cc"},
		},
		{
			name: "syntax error",
			code: "context 100, 100\ncolor (",
			want: &xgoerr.LineError{Number: 2, Msg: "expected operand, found 'EOF'"},
		},
		{
			name: "too many arguments",
			code: "context 100, 100\ncircle 1, 2, 3, 4",
			want: &xgoerr.LineError{Number: 2, Msg: "too many arguments in call to ndor.Circle"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Run(tt.code); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Run() = %v, want %v", got, tt.want)
			}
		})
	}
}
