package internal

import (
	"log"
	"reflect"
	"testing"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func Test_parseGopPlayError(t *testing.T) {
	type args struct {
		ori      string
		preLines int
	}
	tests := []struct {
		name string
		args args
		want *LineError
	}{
		{
			name: ".gop",
			args: args{
				ori: `prog.gop:136:11: missing ',' in argument list (and 1 more errors)
				===> errors stack:
				gop.LoadDir(".", 0xc00008df20, false, false)
					/root/goplus/gengo.go:73 LoadDir(dir, conf, genTestPkg, prompt)`,
				preLines: 133,
			},
			want: &LineError{
				Number: 3,
				Msg:    `missing ',' in argument list (and 1 more errors)`,
			},
		},
		{
			name: ".go",
			args: args{
				ori: `prog.go:136:13: undefined: vv
				===> errors stack:
				gop.LoadDir(".", 0xc000073f20, false, false)
					/root/goplus/gengo.go:73 LoadDir(dir, conf, genTestPkg, prompt)`,
				preLines: 133,
			},
			want: &LineError{
				Number: 3,
				Msg:    `undefined: vv`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseGopPlayError(tt.args.ori, tt.args.preLines); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseGopPlayError() = %v, want %v", got, tt.want)
			}
		})
	}
}
