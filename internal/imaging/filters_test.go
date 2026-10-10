package imaging

import (
	"image/color"
	"testing"
)

func TestGrayscale(t *testing.T) {
	m := testImage(1, 1)
	m.SetRGBA(0, 0, color.RGBA{255, 0, 0, 255})
	out := Grayscale(m)
	c := out.RGBAAt(0, 0)
	if c.R != c.G || c.G != c.B {
		t.Fatalf("not grayscale: %v", c)
	}
}

func TestBrightness(t *testing.T) {
	m := testImage(1, 1)
	m.SetRGBA(0, 0, color.RGBA{100, 100, 100, 255})
	out := Brightness(m, 10)
	c := out.RGBAAt(0, 0)
	if c.R <= 100 {
		t.Fatalf("brightness did not increase: %v", c)
	}
}

func TestContrast(t *testing.T) {
	m := testImage(1, 1)
	m.SetRGBA(0, 0, color.RGBA{200, 200, 200, 255})
	out := Contrast(m, 20)
	c := out.RGBAAt(0, 0)
	if c.R <= 200 {
		t.Fatalf("contrast did not increase bright pixel: %v", c)
	}
}

func TestBlur(t *testing.T) {
	m := testImage(5, 5)
	out, err := Blur(m, 2)
	if err != nil {
		t.Fatal(err)
	}
	if out.Bounds() != m.Bounds() {
		t.Fatal("blur changed bounds")
	}
}

func TestSharpen(t *testing.T) {
	m := testImage(5, 5)
	out, err := Sharpen(m, 1)
	if err != nil {
		t.Fatal(err)
	}
	if out.Bounds() != m.Bounds() {
		t.Fatal("sharpen changed bounds")
	}
}
