package pkg

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/zrcoder/ndor/internal"
)

const (
	tests = `
context 100, 100
color "black"
//---
context 100, 100
color
//---
context 100, 100
Color
//---
Context 800, 800
Color "lightgreen"
Circle 400, 400, 300
Fill
//---
context 800, 800
color "lightgreen"
circle 400, 400, 300
fill
//---
context 1000, 1000
color "notExistColorName"
`
	sep = "//---\n"
)

func Test_xgoRun(t *testing.T) {
	codes := strings.Split(tests, sep)
	for i, code := range codes {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			err := xgoRun(code)
			if err != nil {
				t.Error("code:\n", code)
				t.Error(err)
			}
		})
	}
}

func Test_parseXGoErr(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want *internal.LineError
	}{
		{
			name: "./main.xgo:100:1: undefined: cc",
			err:  errors.New("./main.xgo:100:1: undefined: cc"),
			want: &internal.LineError{Number: 100 - strings.Count(preCodes, "\n"), Msg: "undefined: cc"},
		},
		{
			name: "multiple errors",
			err: errors.New(`main.xgo:126:2: undefined: abc
			main.xgo:127:1: undefined: vv123`),
			want: &internal.LineError{Number: 126 - strings.Count(preCodes, "\n"), Msg: "undefined: abc"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseXGoErr(tt.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\nparseXGoErr() = %v, \nwant %v", got, tt.want)
			}
		})
	}
}
