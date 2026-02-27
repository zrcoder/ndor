package goldmark

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/yuin/goldmark"
)

func Test(t *testing.T) {
	md := goldmark.New(
		goldmark.WithExtensions(&Extender{}),
	)
	data, err := os.ReadFile(filepath.Join("testdata", "test.md"))
	if err != nil {
		t.Error(err)
	}
	out := &bytes.Buffer{}
	err = md.Convert(data, out)
	if err != nil {
		t.Error(err)
	} else {
		err = os.WriteFile(filepath.Join("testdata", "test.html"), out.Bytes(), 0o640)
		if err != nil {
			t.Error(err)
		}
	}
}
