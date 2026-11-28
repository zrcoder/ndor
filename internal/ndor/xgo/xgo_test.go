package xgo

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/zrcoder/ndor/internal/ndor"
	"github.com/zrcoder/ndor/internal/xgoerr"
)

func TestRun(t *testing.T) {
	codes := []string{
		`color "black"`,
		`color`,
		`Color`,
		`Color "lightgreen"
Circle 0, 0, 300
Fill`,
		`color "lightgreen"
circle 0, 0, 300
fill`,
		`color "notExistColorName"`,
	}
	for i, code := range codes {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ndor.Init()
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
			code: "color \"black\"\ncc",
			want: &xgoerr.LineError{Number: 2, Msg: "undefined: cc"},
		},
		{
			name: "syntax error",
			code: "color \"black\"\ncolor (",
			want: &xgoerr.LineError{Number: 2, Msg: "expected operand, found 'EOF'"},
		},
		{
			name: "too many arguments",
			code: "color \"black\"\ncircle 1, 2, 3, 4",
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
