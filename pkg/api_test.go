package pkg

import (
	"testing"

	"github.com/zrcoder/niudour/internal"
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
		want *internal.LineError
	}{
		{
			args: args{code: `context 0, 0
			color "unkown"`},
			want: &internal.LineError{Number: 2, Msg: "invalid parameters"},
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

func Test_parseErrorline(t *testing.T) {
	type args struct {
		oriCode string
		flag    string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			args: args{
				oriCode: `line 1
				line 2
				...
				color "unkown"
				...`,
				flag: "color",
			},
			want: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseErrorline(tt.args.oriCode, tt.args.flag); got != tt.want {
				t.Errorf("parseErrorline() = %v, want %v", got, tt.want)
			}
		})
	}
}
