package ndor

import "testing"

func TestImageSize(t *testing.T) {
	tests := []struct {
		name    string
		draw    func(c *Context)
		w, h    int
		wantErr bool
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
			name: "too large",
			draw: func(c *Context) {
				c.DrawCircle(0, 0, 5000)
				c.Fill()
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewContext()
			tt.draw(c)
			img, err := c.Image()
			if tt.wantErr {
				if err == nil {
					t.Fatal("Image() succeeded, want an error")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			b := img.Bounds()
			if b.Dx() != tt.w || b.Dy() != tt.h {
				t.Errorf("image size = %dx%d, want %dx%d", b.Dx(), b.Dy(), tt.w, tt.h)
			}
		})
	}
}
