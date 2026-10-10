package cli

import (
	"fmt"
	"image"
	"strconv"
	"strings"

	"github.com/priyojeet/imgx/internal/imaging"
)

func load(path string) (image.Image, string, error) { return imaging.Load(path, maxPixels) }

func save(img image.Image, output, format string, quality int, pngLevel int, tiffLZW bool) (string, error) {
	level := imaging.DefaultOptions().PNGLevel
	if pngLevel < -2 || pngLevel > 9 {
		return "", fmt.Errorf("png compression level must be between -2 and 9")
	}
	switch pngLevel {
	case -2:
		level = -2
	case -1:
		level = -1
	case 0:
		level = 0
	case 1:
		level = 1
	case 2:
		level = 2
	case 3:
		level = 3
	case 4:
		level = 4
	case 5:
		level = 5
	case 6:
		level = 6
	case 7:
		level = 7
	case 8:
		level = 8
	case 9:
		level = 9
	}
	return imaging.Save(img, output, format, imaging.Options{JPEGQuality: quality, PNGLevel: level, TIFFLZW: tiffLZW})
}

func parseResizeAlgo(v string) (imaging.ResizeAlgorithm, error) {
	a := imaging.ResizeAlgorithm(strings.ToLower(strings.TrimSpace(v)))
	switch a {
	case imaging.ResizeNearest, imaging.ResizeBilinear, imaging.ResizeBicubic, imaging.ResizeLanczos3:
		return a, nil
	default:
		return "", fmt.Errorf("unsupported resize algorithm %q (use nearest, bilinear, bicubic, or lanczos3)", v)
	}
}

func requireRange(name string, v, min, max float64) error {
	if v < min || v > max {
		return fmt.Errorf("%s must be between %v and %v", name, min, max)
	}
	return nil
}

func parseFloat(s string) (float64, error) { return strconv.ParseFloat(strings.TrimSpace(s), 64) }
