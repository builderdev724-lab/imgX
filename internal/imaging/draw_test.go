package imaging

import (
	"image/color"
	"testing"
)

func TestParseColor(t *testing.T) {
	c, err := ParseColor("#12abEF")
	if err != nil {
		t.Fatal(err)
	}
	if c.R != 0x12 || c.G != 0xab || c.B != 0xef || c.A != 255 {
		t.Fatalf("got %#v", c)
	}
	c, err = ParseColor("#01020380")
	if err != nil {
		t.Fatal(err)
	}
	if c.A != 0x80 {
		t.Fatalf("alpha: %d", c.A)
	}
	if _, err = ParseColor("not-a-color"); err == nil {
		t.Fatal("expected invalid color error")
	}
}

func TestBorder(t *testing.T) {
	m := testImage(3, 2)
	out, err := AddBorder(m, 2, rgbaBlack())
	if err != nil {
		t.Fatal(err)
	}
	if out.Bounds().Dx() != 7 || out.Bounds().Dy() != 6 {
		t.Fatalf("got %dx%d", out.Bounds().Dx(), out.Bounds().Dy())
	}
}

func rgbaBlack() color.RGBA { return color.RGBA{0, 0, 0, 255} }
