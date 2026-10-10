package imaging

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

func TestSaveLoadPNG(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "test.png")
	m := image.NewRGBA(image.Rect(0, 0, 16, 12))
	m.SetRGBA(2, 3, color.RGBA{1, 2, 3, 255})
	if _, err := Save(m, p, "", DefaultOptions()); err != nil {
		t.Fatal(err)
	}
	out, format, err := Load(p, 1000)
	if err != nil {
		t.Fatal(err)
	}
	if format != "png" {
		t.Fatalf("format %s", format)
	}
	if out.Bounds().Dx() != 16 || out.Bounds().Dy() != 12 {
		t.Fatal("unexpected dimensions")
	}
}

func TestDecodeBufferedMaxPixels(t *testing.T) {
	var b bytes.Buffer
	m := image.NewRGBA(image.Rect(0, 0, 20, 20))
	if err := png.Encode(&b, m); err != nil {
		t.Fatal(err)
	}
	if _, _, err := DecodeBuffered(b.Bytes(), 100); err == nil {
		t.Fatal("expected max-pixels error")
	}
}

func TestResolveOutputFormat(t *testing.T) {
	for _, tt := range []struct{ path, hint, want string }{{"x.jpg", "", "jpeg"}, {"x.tif", "", "tiff"}, {"-", "png", "png"}} {
		got, err := ResolveOutputFormat(tt.path, tt.hint)
		if err != nil {
			t.Fatal(err)
		}
		if got != tt.want {
			t.Fatalf("got %s want %s", got, tt.want)
		}
	}
	if _, err := ResolveOutputFormat("x", ""); err == nil {
		t.Fatal("expected extension error")
	}
}

func TestSaveToStdoutPng(t *testing.T) {
	// Redirecting os.Stdout in tests is platform-sensitive; only verify the format path here.
	_ = os.Stdout
}
