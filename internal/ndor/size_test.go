package ndor

import (
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestDrawingBounds(t *testing.T) {
	tests := []struct {
		name string
		draw func(c *Context)
		w, h float32
	}{
		{
			name: "circle",
			draw: func(c *Context) {
				c.DrawCircle(0, 0, 300)
				c.Fill()
			},
			w: 600, h: 600,
		},
		{
			name: "stroke width",
			draw: func(c *Context) {
				c.MoveTo(-100, 0)
				c.LineTo(100, 0)
				c.SetLineWidth(10)
				c.Stroke()
			},
			w: 210, h: 10,
		},
		{
			name: "clips intersect",
			draw: func(c *Context) {
				c.DrawCircle(-150, 0, 300)
				c.Clip()
				c.DrawCircle(150, 0, 300)
				c.Clip()
				c.DrawRectangle(-500, -500, 1000, 1000)
				c.Fill()
			},
			w: 300, h: 600,
		},
		{
			name: "huge drawing needs no full-size canvas",
			draw: func(c *Context) {
				c.DrawCircle(0, 0, 5000)
				c.Fill()
			},
			w: 10000, h: 10000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewContext()
			tt.draw(c)
			b := c.Bounds()
			if b.W != tt.w || b.H != tt.h {
				t.Errorf("bounds = %gx%g, want %gx%g", b.W, b.H, tt.w, tt.h)
			}
		})
	}
}

func TestTextBoundsRespectClip(t *testing.T) {
	measure := func(clip bool) ui.Rect {
		c := NewContext()
		c.SetRGBA255(0, 0, 0, 255)
		if clip {
			c.DrawCircle(-200, 0, 100)
			c.Clip()
		}
		c.DrawString("WWWWWWWWWW", -200, 0)
		return c.Bounds()
	}
	full, clipped := measure(false), measure(true)
	if full.W <= 0 || full.H <= 13 {
		t.Fatalf("text bounds = %+v, want a real glyph box", full)
	}
	if clipped.W <= 0 || clipped.W >= full.W || clipped.H > 200 {
		t.Fatalf("clipped text bounds = %+v, full = %+v", clipped, full)
	}
}

func TestRenderScalesContent(t *testing.T) {
	c := NewContext()
	c.SetRGBA255(0, 0, 0, 255)
	c.DrawCircle(0, 0, 5000)
	c.Fill()
	img := ui.Render(func(ctx *ui.Context) {
		ui.Column(ctx).Fill().Draw(func(p *ui.Painter, r ui.Rect) { c.Render(p, r) })
	}, 80, 80, 1)
	if got := img.Bounds(); got.Dx() != 80 || got.Dy() != 80 {
		t.Fatalf("render size = %dx%d, want 80x80", got.Dx(), got.Dy())
	}
	for _, x := range []int{10, 20, 30, 40, 50, 60, 70} {
		if r, _, _, _ := img.RGBAAt(x, 40).RGBA(); r > 0x1010 {
			t.Fatalf("scaled drawing has a seam at x=%d: red=%04x", x, r)
		}
	}
}

func TestSVG(t *testing.T) {
	c := NewContext()
	c.SetRGBA255(255, 0, 0, 255)
	c.DrawCircle(0, 0, 10)
	c.Fill()
	c.SetRGBA255(0, 0, 255, 255)
	c.DrawString("hi <ink>", 0, 0)
	svg := string(c.SVG())
	for _, want := range []string{
		`<svg xmlns="http://www.w3.org/2000/svg"`,
		`fill="#ff0000"`,
		`font-family="system-ui, sans-serif"`,
		`fill="#0000ff"`,
		`hi &lt;ink&gt;`,
	} {
		if !strings.Contains(svg, want) {
			t.Errorf("SVG missing %q: %s", want, svg)
		}
	}
}

func TestSVGClip(t *testing.T) {
	c := NewContext()
	c.DrawCircle(-150, 0, 300)
	c.Clip()
	c.DrawCircle(150, 0, 300)
	c.Clip()
	c.DrawRectangle(-500, -500, 1000, 1000)
	c.Fill()
	svg := string(c.SVG())
	if got := strings.Count(svg, `<clipPath id="`); got != 2 {
		t.Errorf("clip defs = %d, want 2: %s", got, svg)
	}
	if got := strings.Count(svg, `<g clip-path="url(#ndor-clip-`); got != 2 {
		t.Errorf("clip groups = %d, want 2: %s", got, svg)
	}
}
