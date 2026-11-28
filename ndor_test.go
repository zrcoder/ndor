package ndor

import (
	"testing"

	"github.com/zrcoder/ndor/internal/xgoerr"
)

func TestRun(t *testing.T) {
	type args struct {
		code string
	}
	tests := []struct {
		name string
		args args
		want *xgoerr.LineError
	}{
		{
			args: args{code: `color "unkown"`},
			want: &xgoerr.LineError{Number: 1, Msg: "invalid parameters"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, got := Run(tt.args.code)
			if got.Number != tt.want.Number {
				t.Errorf("Run() got error number:%d, want:%d", got.Number, tt.want.Number)
			}
		})
	}
}
