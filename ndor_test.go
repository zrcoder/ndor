package ndor

import (
	"testing"

	"github.com/zrcoder/ndor/internal/xgoerr"
)

func TestRun(t *testing.T) {
	type args struct {
		width  int
		height int
		code   string
	}
	tests := []struct {
		name string
		args args
		want *xgoerr.LineError
	}{
		{
			args: args{code: `context 0, 0
			color "unkown"`},
			want: &xgoerr.LineError{Number: 2, Msg: "invalid parameters"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, got := Run(tt.args.width, tt.args.height, tt.args.code)
			if got.Number != tt.want.Number {
				t.Errorf("Run() got error number:%d, want:%d", got.Number, tt.want.Number)
			}
		})
	}
}
