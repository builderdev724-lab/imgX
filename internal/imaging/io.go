package imaging

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"

	// Registers the WebP decoder with image.Decode, which is how WebP
	// input is supported. WebP output is intentionally not supported.
	_ "golang.org/x/image/webp"
)

const DefaultMaxPixels = int64(100_000_000)

type Options struct {
	JPEGQuality int
	PNGLevel    png.CompressionLevel
	TIFFLZW     bool
}

func DefaultOptions() Options {
	return Options{JPEGQuality: 90, PNGLevel: png.DefaultCompression, TIFFLZW: true}
}

func Load(path string, maxPixels int64) (image.Image, string, error) {
	data, err := readAll(path)
	if err != nil {
		return nil, "", err
	}
	return DecodeBuffered(data, maxPixels)
}

func DecodeBuffered(data []byte, maxPixels int64) (image.Image, string, error) {
	if len(data) == 0 {
		return nil, "", errors.New("input is empty")
	}
	if maxPixels <= 0 {
		maxPixels = DefaultMaxPixels
	}
	cfg, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return nil, "", fmt.Errorf("decode image header: %w", err)
	}
	if cfg.Width <= 0 || cfg.Height <= 0 {
		return nil, "", errors.New("image has invalid dimensions")
	}
	if int64(cfg.Width) > maxPixels/int64(cfg.Height) {
		return nil, "", fmt.Errorf("image is too large: %dx%d exceeds max-pixels=%d", cfg.Width, cfg.Height, maxPixels)
	}
	img, format, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, "", fmt.Errorf("decode image: %w", err)
	}
	return img, format, nil
}

func readAll(path string) ([]byte, error) {
	if path == "-" {
		return io.ReadAll(os.Stdin)
	}
	return os.ReadFile(path)
}

func Save(img image.Image, path, formatHint string, opts Options) (string, error) {
	format, err := ResolveOutputFormat(path, formatHint)
	if err != nil {
		return "", err
	}
	w := io.Writer(os.Stdout)
	var f *os.File
	if path != "-" {
		dir := filepath.Dir(path)
		if dir != "." {
			if err := os.MkdirAll(dir, 0o755); err != nil {
				return "", fmt.Errorf("create output directory: %w", err)
			}
		}
		f, err = os.Create(path)
		if err != nil {
			return "", fmt.Errorf("create output: %w", err)
		}
		defer f.Close()
		w = f
	}
	if opts.JPEGQuality == 0 {
		opts.JPEGQuality = 90
	}
	if format == "jpeg" {
		if err := validateJPEGQuality(opts.JPEGQuality); err != nil {
			return "", err
		}
	}
	encImg := img
	if format == "jpeg" {
		encImg = flatten(img, color.RGBA{255, 255, 255, 255})
	}
	switch format {
	case "png":
		err = (&png.Encoder{CompressionLevel: opts.PNGLevel}).Encode(w, encImg)
	case "jpeg":
		err = jpeg.Encode(w, encImg, &jpeg.Options{Quality: opts.JPEGQuality})
	case "gif":
		err = gif.Encode(w, encImg, &gif.Options{NumColors: 256})
	case "bmp":
		err = bmp.Encode(w, encImg)
	case "tiff":
		// Note: x/image/tiff can *decode* LZW but its encoder only
		// implements uncompressed and Deflate output, so TIFFLZW selects
		// Deflate here. Requesting tiff.LZW returns "unsupported
		// compression" from the encoder.
		compression := tiff.Uncompressed
		if opts.TIFFLZW {
			compression = tiff.Deflate
		}
		err = tiff.Encode(w, encImg, &tiff.Options{Compression: compression})
	case "webp":
		return "", errors.New("WebP output is not supported yet; WebP input is supported. Keeping the encoder dependency-free and cgo-free makes cross-compilation reliable")
	default:
		return "", fmt.Errorf("unsupported output format %q", format)
	}
	if err != nil {
		return "", fmt.Errorf("encode %s: %w", format, err)
	}
	return format, nil
}

func ResolveOutputFormat(path, hint string) (string, error) {
	if strings.TrimSpace(hint) != "" {
		return normalizeFormat(hint)
	}
	if path == "-" {
		return "", errors.New("stdout output requires --format")
	}
	ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(path)), ".")
	if ext == "" {
		return "", errors.New("cannot infer output format without a file extension; use --format")
	}
	return normalizeFormat(ext)
}

func normalizeFormat(v string) (string, error) {
	v = strings.ToLower(strings.TrimSpace(strings.TrimPrefix(v, ".")))
	switch v {
	case "jpg", "jpeg":
		return "jpeg", nil
	case "png":
		return "png", nil
	case "gif":
		return "gif", nil
	case "bmp":
		return "bmp", nil
	case "tif", "tiff":
		return "tiff", nil
	case "webp":
		return "webp", nil
	default:
		return "", fmt.Errorf("unsupported image format %q (supported output: png, jpeg, gif, bmp, tiff)", v)
	}
}

func validateJPEGQuality(q int) error {
	if q < 1 || q > 100 {
		return fmt.Errorf("JPEG quality must be between 1 and 100, got %d", q)
	}
	return nil
}

func flatten(src image.Image, bg color.RGBA) image.Image {
	b := src.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	for y := 0; y < b.Dy(); y++ {
		for x := 0; x < b.Dx(); x++ {
			sr, sg, sb, sa := src.At(x+b.Min.X, y+b.Min.Y).RGBA()
			a := float64(sa) / 65535
			br, bgV, bb, ba := uint32(bg.R)*257, uint32(bg.G)*257, uint32(bg.B)*257, uint32(bg.A)*257
			baF := float64(ba) / 65535
			outA := a + baF*(1-a)
			if outA == 0 {
				dst.SetRGBA(x, y, color.RGBA{})
				continue
			}
			r := (float64(sr)*a + float64(br)*baF*(1-a)) / outA
			g := (float64(sg)*a + float64(bgV)*baF*(1-a)) / outA
			bV := (float64(sb)*a + float64(bb)*baF*(1-a)) / outA
			dst.SetRGBA(x, y, color.RGBA{uint8(r / 257), uint8(g / 257), uint8(bV / 257), uint8(outA * 255)})
		}
	}
	return dst
}
