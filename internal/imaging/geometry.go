package imaging

import (
	"image"
	"image/color"
	"math"
)

type ResizeAlgorithm string

const (
	ResizeNearest  ResizeAlgorithm = "nearest"
	ResizeBilinear ResizeAlgorithm = "bilinear"
	ResizeBicubic  ResizeAlgorithm = "bicubic"
	ResizeLanczos3 ResizeAlgorithm = "lanczos3"
)

func Resize(img image.Image, width, height int, algo ResizeAlgorithm) (*image.RGBA, error) {
	if width <= 0 || height <= 0 {
		return nil, &ParameterError{"resize dimensions must be positive"}
	}
	src := toRGBA(img)
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	sb := src.Bounds()
	sw, sh := sb.Dx(), sb.Dy()
	if sw == width && sh == height {
		copy(dst.Pix, src.Pix)
		return dst, nil
	}
	for y := 0; y < height; y++ {
		sy := (float64(y)+0.5)*float64(sh)/float64(height) - 0.5
		for x := 0; x < width; x++ {
			sx := (float64(x)+0.5)*float64(sw)/float64(width) - 0.5
			c := sampleRGBA(src, sx, sy, algo)
			dst.SetRGBA(x, y, c)
		}
	}
	return dst, nil
}

func ResizeKeepAspect(src image.Image, maxW, maxH int, algo ResizeAlgorithm) (*image.RGBA, error) {
	if maxW <= 0 && maxH <= 0 {
		return nil, &ParameterError{"at least one target dimension must be positive"}
	}
	b := src.Bounds()
	w, h := b.Dx(), b.Dy()
	scale := math.Inf(1)
	if maxW > 0 {
		scale = math.Min(scale, float64(maxW)/float64(w))
	}
	if maxH > 0 {
		scale = math.Min(scale, float64(maxH)/float64(h))
	}
	if scale <= 0 || math.IsInf(scale, 0) {
		return nil, &ParameterError{"invalid keep-aspect dimensions"}
	}
	if scale > 1 {
		scale = 1
	}
	nw := max(1, int(math.Round(float64(w)*scale)))
	nh := max(1, int(math.Round(float64(h)*scale)))
	return Resize(src, nw, nh, algo)
}

func Crop(img image.Image, x, y, width, height int) (*image.RGBA, error) {
	if width <= 0 || height <= 0 {
		return nil, &ParameterError{"crop dimensions must be positive"}
	}
	b := img.Bounds()
	r := image.Rect(x, y, x+width, y+height)
	if !r.In(b) {
		return nil, &ParameterError{"crop rectangle is outside image bounds"}
	}
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	for yy := 0; yy < height; yy++ {
		for xx := 0; xx < width; xx++ {
			dst.Set(xx, yy, img.At(x+xx, y+yy))
		}
	}
	return dst, nil
}

func Flip(img image.Image, horizontal, vertical bool) (*image.RGBA, error) {
	if !horizontal && !vertical {
		return nil, &ParameterError{"choose --horizontal and/or --vertical"}
	}
	src := toRGBA(img)
	b := src.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	for y := 0; y < b.Dy(); y++ {
		for x := 0; x < b.Dx(); x++ {
			sx, sy := x, y
			if horizontal {
				sx = b.Dx() - 1 - x
			}
			if vertical {
				sy = b.Dy() - 1 - y
			}
			dst.SetRGBA(x, y, src.RGBAAt(sx+b.Min.X, sy+b.Min.Y))
		}
	}
	return dst, nil
}

func Rotate(img image.Image, degrees float64, algo ResizeAlgorithm) (*image.RGBA, error) {
	if math.IsNaN(degrees) || math.IsInf(degrees, 0) {
		return nil, &ParameterError{"rotation angle must be finite"}
	}
	deg := math.Mod(degrees, 360)
	if deg < 0 {
		deg += 360
	}
	if almostEqual(deg, 0) {
		return toRGBA(img), nil
	}
	if almostEqual(deg, 90) {
		return rotate90(img), nil
	}
	if almostEqual(deg, 180) {
		r := rotate90(img)
		return rotate90(r), nil
	}
	if almostEqual(deg, 270) {
		r := rotate90(img)
		r = rotate90(r)
		r = rotate90(r)
		return r, nil
	}

	src := toRGBA(img)
	w, h := src.Bounds().Dx(), src.Bounds().Dy()
	rad := deg * math.Pi / 180
	cosA, sinA := math.Cos(rad), math.Sin(rad)
	ow := int(math.Ceil(math.Abs(float64(w)*cosA) + math.Abs(float64(h)*sinA)))
	oh := int(math.Ceil(math.Abs(float64(w)*sinA) + math.Abs(float64(h)*cosA)))
	dst := image.NewRGBA(image.Rect(0, 0, ow, oh))
	cx, cy := float64(w-1)/2, float64(h-1)/2
	dcx, dcy := float64(ow-1)/2, float64(oh-1)/2
	for y := 0; y < oh; y++ {
		for x := 0; x < ow; x++ {
			dx, dy := float64(x)-dcx, float64(y)-dcy
			sx := dx*cosA + dy*sinA + cx
			sy := -dx*sinA + dy*cosA + cy
			if sx >= 0 && sx <= float64(w-1) && sy >= 0 && sy <= float64(h-1) {
				dst.SetRGBA(x, y, sampleRGBA(src, sx, sy, algo))
			} else {
				dst.SetRGBA(x, y, transparentRGBA())
			}
		}
	}
	return dst, nil
}

func rotate90(img image.Image) *image.RGBA {
	src := toRGBA(img)
	w, h := src.Bounds().Dx(), src.Bounds().Dy()
	dst := image.NewRGBA(image.Rect(0, 0, h, w))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			dst.SetRGBA(h-1-y, x, src.RGBAAt(x+src.Bounds().Min.X, y+src.Bounds().Min.Y))
		}
	}
	return dst
}

func sampleRGBA(img *image.RGBA, x, y float64, algo ResizeAlgorithm) (c color.RGBA) {
	if algo == "" {
		algo = ResizeBilinear
	}
	switch algo {
	case ResizeNearest:
		return nearest(img, x, y)
	case ResizeBilinear:
		return bilinear(img, x, y)
	case ResizeBicubic:
		return bicubic(img, x, y)
	case ResizeLanczos3:
		return lanczos(img, x, y, 3)
	default:
		return bilinear(img, x, y)
	}
}

func nearest(img *image.RGBA, x, y float64) color.RGBA {
	return img.RGBAAt(clampInt(int(math.Round(x)), 0, img.Bounds().Dx()-1), clampInt(int(math.Round(y)), 0, img.Bounds().Dy()-1))
}

func bilinear(img *image.RGBA, x, y float64) color.RGBA {
	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	x0, y0 := int(math.Floor(x)), int(math.Floor(y))
	x1, y1 := x0+1, y0+1
	fx, fy := x-float64(x0), y-float64(y0)
	x0, x1 = clampInt(x0, 0, w-1), clampInt(x1, 0, w-1)
	y0, y1 = clampInt(y0, 0, h-1), clampInt(y1, 0, h-1)
	c00 := img.RGBAAt(x0, y0)
	c10 := img.RGBAAt(x1, y0)
	c01 := img.RGBAAt(x0, y1)
	c11 := img.RGBAAt(x1, y1)
	mix := func(a, b, t float64) float64 { return a + (b-a)*t }
	return color.RGBA{
		uint8(clampFloat(mix(mix(float64(c00.R), float64(c10.R), fx), mix(float64(c01.R), float64(c11.R), fx), fy), 0, 255)),
		uint8(clampFloat(mix(mix(float64(c00.G), float64(c10.G), fx), mix(float64(c01.G), float64(c11.G), fx), fy), 0, 255)),
		uint8(clampFloat(mix(mix(float64(c00.B), float64(c10.B), fx), mix(float64(c01.B), float64(c11.B), fx), fy), 0, 255)),
		uint8(clampFloat(mix(mix(float64(c00.A), float64(c10.A), fx), mix(float64(c01.A), float64(c11.A), fx), fy), 0, 255)),
	}
}

func bicubic(img *image.RGBA, x, y float64) color.RGBA {
	return kernelSample(img, x, y, 2, cubicKernel)
}

func lanczos(img *image.RGBA, x, y float64, a int) color.RGBA {
	return kernelSample(img, x, y, a, func(v float64) float64 {
		if v == 0 {
			return 1
		}
		if math.Abs(v) >= float64(a) {
			return 0
		}
		return sinc(v) * sinc(v/float64(a))
	})
}

func kernelSample(img *image.RGBA, x, y float64, radius int, kernel func(float64) float64) color.RGBA {
	var r, g, b, a, sum float64
	minX, maxX := int(math.Floor(x))-radius+1, int(math.Floor(x))+radius
	minY, maxY := int(math.Floor(y))-radius+1, int(math.Floor(y))+radius
	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	for yy := minY; yy <= maxY; yy++ {
		wy := kernel(float64(yy) - y)
		if wy == 0 {
			continue
		}
		cy := clampInt(yy, 0, h-1)
		for xx := minX; xx <= maxX; xx++ {
			wx := kernel(float64(xx) - x)
			wt := wx * wy
			if wt == 0 {
				continue
			}
			cx := clampInt(xx, 0, w-1)
			c := img.RGBAAt(cx, cy)
			r += float64(c.R) * wt
			g += float64(c.G) * wt
			b += float64(c.B) * wt
			a += float64(c.A) * wt
			sum += wt
		}
	}
	if sum == 0 {
		return bilinear(img, x, y)
	}
	return color.RGBA{uint8(clampFloat(r/sum, 0, 255)), uint8(clampFloat(g/sum, 0, 255)), uint8(clampFloat(b/sum, 0, 255)), uint8(clampFloat(a/sum, 0, 255))}
}

func cubicKernel(x float64) float64 {
	x = math.Abs(x)
	if x < 1 {
		return 1.5*x*x*x - 2.5*x*x + 1
	}
	if x < 2 {
		return -0.5*x*x*x + 2.5*x*x - 4*x + 2
	}
	return 0
}

func sinc(x float64) float64 {
	if x == 0 {
		return 1
	}
	p := math.Pi * x
	return math.Sin(p) / p
}

func toRGBA(img image.Image) *image.RGBA {
	b := img.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	for y := 0; y < b.Dy(); y++ {
		for x := 0; x < b.Dx(); x++ {
			dst.Set(x, y, img.At(x+b.Min.X, y+b.Min.Y))
		}
	}
	return dst
}

func clampInt(v, minV, maxV int) int {
	if v < minV {
		return minV
	}
	if v > maxV {
		return maxV
	}
	return v
}
func clampFloat(v, minV, maxV float64) float64 {
	if v < minV {
		return minV
	}
	if v > maxV {
		return maxV
	}
	return v
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func almostEqual(a, b float64) bool { return math.Abs(a-b) < 1e-9 }
func transparentRGBA() color.RGBA   { return color.RGBA{0, 0, 0, 0} }

type ParameterError struct{ Message string }

func (e *ParameterError) Error() string { return e.Message }
