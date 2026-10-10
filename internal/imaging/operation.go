package imaging

import (
	"fmt"
	"image"
	"strconv"
	"strings"
)

type Operation struct {
	Name string
	Args map[string]string
}

func Apply(img image.Image, op Operation) (image.Image, error) {
	name := strings.ToLower(strings.TrimSpace(op.Name))
	arg := func(key string) string { return strings.TrimSpace(op.Args[key]) }
	intArg := func(key string, def int) (int, error) {
		v := arg(key)
		if v == "" {
			return def, nil
		}
		n, err := strconv.Atoi(v)
		if err != nil {
			return 0, fmt.Errorf("%s must be an integer", key)
		}
		return n, nil
	}
	floatArg := func(key string, def float64) (float64, error) {
		v := arg(key)
		if v == "" {
			return def, nil
		}
		n, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, fmt.Errorf("%s must be a number", key)
		}
		return n, nil
	}

	switch name {
	case "resize":
		w, err := intArg("w", 0)
		if err != nil {
			return nil, err
		}
		h, err := intArg("h", 0)
		if err != nil {
			return nil, err
		}
		algo := ResizeAlgorithm(arg("algo"))
		if algo == "" {
			algo = ResizeBilinear
		}
		if parseBool(arg("keep-aspect")) {
			return ResizeKeepAspect(img, w, h, algo)
		}
		return Resize(img, w, h, algo)
	case "crop":
		x, err := intArg("x", 0)
		if err != nil {
			return nil, err
		}
		y, err := intArg("y", 0)
		if err != nil {
			return nil, err
		}
		w, err := intArg("w", 0)
		if err != nil {
			return nil, err
		}
		h, err := intArg("h", 0)
		if err != nil {
			return nil, err
		}
		return Crop(img, x, y, w, h)
	case "rotate":
		deg, err := floatArg("deg", 0)
		if err != nil {
			return nil, err
		}
		algo := ResizeAlgorithm(arg("algo"))
		if algo == "" {
			algo = ResizeBilinear
		}
		return Rotate(img, deg, algo)
	case "flip":
		return Flip(img, parseBool(arg("horizontal")), parseBool(arg("vertical")))
	case "grayscale":
		return Grayscale(img), nil
	case "brightness":
		v, err := floatArg("value", 0)
		if err != nil {
			return nil, err
		}
		if v < -100 || v > 100 {
			return nil, &ParameterError{"brightness must be between -100 and 100"}
		}
		return Brightness(img, v), nil
	case "contrast":
		v, err := floatArg("value", 0)
		if err != nil {
			return nil, err
		}
		if v < -100 || v > 100 {
			return nil, &ParameterError{"contrast must be between -100 and 100"}
		}
		return Contrast(img, v), nil
	case "blur":
		r, err := intArg("radius", 0)
		if err != nil {
			return nil, err
		}
		return Blur(img, r)
	case "sharpen":
		v, err := floatArg("amount", 1)
		if err != nil {
			return nil, err
		}
		return Sharpen(img, v)
	case "border":
		w, err := intArg("width", 10)
		if err != nil {
			return nil, err
		}
		c, err := ParseColor(argDefault(op.Args, "color", "#000000"))
		if err != nil {
			return nil, err
		}
		return AddBorder(img, w, c)
	case "watermark":
		text := arg("text")
		pos := WatermarkPosition(arg("pos"))
		if pos == "" {
			pos = PosBottomRight
		}
		sz, err := floatArg("size", 32)
		if err != nil {
			return nil, err
		}
		opacity, err := floatArg("opacity", 70)
		if err != nil {
			return nil, err
		}
		c, err := ParseColor(argDefault(op.Args, "color", "#ffffff"))
		if err != nil {
			return nil, err
		}
		return Watermark(img, text, pos, sz, opacity, c)
	default:
		return nil, fmt.Errorf("unknown operation %q", name)
	}
}

func parseBool(v string) bool {
	v = strings.ToLower(strings.TrimSpace(v))
	return v == "1" || v == "true" || v == "yes" || v == "on"
}
func argDefault(m map[string]string, key, def string) string {
	if v := strings.TrimSpace(m[key]); v != "" {
		return v
	}
	return def
}
