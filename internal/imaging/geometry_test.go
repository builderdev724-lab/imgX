package imaging

import (
	"image"
	"image/color"
	"testing"
)

func testImage(w, h int) *image.RGBA {
	m := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			m.SetRGBA(x, y, color.RGBA{uint8(x * 40), uint8(y * 40), uint8((x + y) * 20), 255})
		}
	}
	return m
}

func TestResizeDimensions(t *testing.T) {
	m := testImage(4, 2)
	out, err := Resize(m, 8, 4, ResizeBilinear)
	if err != nil {
		t.Fatal(err)
	}
	if got := out.Bounds().Size(); got != (image.Point{X: 8, Y: 4}) {
		t.Fatalf("got %v", got)
	}
}

func TestResizeKeepAspect(t *testing.T) {
	m := testImage(400, 200)
	out, err := ResizeKeepAspect(m, 100, 100, ResizeBilinear)
	if err != nil {
		t.Fatal(err)
	}
	if out.Bounds().Dx() != 100 || out.Bounds().Dy() != 50 {
		t.Fatalf("got %dx%d", out.Bounds().Dx(), out.Bounds().Dy())
	}
}

func TestCrop(t *testing.T) {
	m := testImage(10, 8)
	out, err := Crop(m, 2, 3, 4, 2)
	if err != nil {
		t.Fatal(err)
	}
	if out.Bounds().Dx() != 4 || out.Bounds().Dy() != 2 {
		t.Fatal("unexpected crop size")
	}
	if got := out.RGBAAt(0, 0); got != m.RGBAAt(2, 3) {
		t.Fatalf("pixel mismatch: %v != %v", got, m.RGBAAt(2, 3))
	}
}

func TestFlip(t *testing.T) {
	m := testImage(3, 1)
	out, err := Flip(m, true, false)
	if err != nil {
		t.Fatal(err)
	}
	if out.RGBAAt(0, 0) != m.RGBAAt(2, 0) {
		t.Fatal("horizontal flip mismatch")
	}
}

func TestRotate90(t *testing.T) {
	m := testImage(4, 2)
	out, err := Rotate(m, 90, ResizeBilinear)
	if err != nil {
		t.Fatal(err)
	}
	if out.Bounds().Dx() != 2 || out.Bounds().Dy() != 4 {
		t.Fatalf("got %dx%d", out.Bounds().Dx(), out.Bounds().Dy())
	}
}
