package imaging

import (
	"image"
	"image/color"
	"math"
)

func Grayscale(img image.Image) *image.RGBA {
	src := toRGBA(img)
	dst := image.NewRGBA(src.Bounds())
	for y := 0; y < src.Bounds().Dy(); y++ {
		for x := 0; x < src.Bounds().Dx(); x++ {
			c := src.RGBAAt(x, y)
			v := uint8(clampFloat(0.2126*float64(c.R)+0.7152*float64(c.G)+0.0722*float64(c.B), 0, 255))
			dst.SetRGBA(x, y, color.RGBA{v, v, v, c.A})
		}
	}
	return dst
}

func Brightness(img image.Image, value float64) *image.RGBA {
	src := toRGBA(img)
	dst := image.NewRGBA(src.Bounds())
	delta := 255 * value / 100
	for y := 0; y < src.Bounds().Dy(); y++ {
		for x := 0; x < src.Bounds().Dx(); x++ {
			c := src.RGBAAt(x, y)
			dst.SetRGBA(x, y, color.RGBA{
				uint8(clampFloat(float64(c.R)+delta, 0, 255)),
				uint8(clampFloat(float64(c.G)+delta, 0, 255)),
				uint8(clampFloat(float64(c.B)+delta, 0, 255)), c.A,
			})
		}
	}
	return dst
}

func Contrast(img image.Image, value float64) *image.RGBA {
	src := toRGBA(img)
	dst := image.NewRGBA(src.Bounds())
	factor := (100 + value) / 100
	for y := 0; y < src.Bounds().Dy(); y++ {
		for x := 0; x < src.Bounds().Dx(); x++ {
			c := src.RGBAAt(x, y)
			dst.SetRGBA(x, y, color.RGBA{
				uint8(clampFloat((float64(c.R)-127.5)*factor+127.5, 0, 255)),
				uint8(clampFloat((float64(c.G)-127.5)*factor+127.5, 0, 255)),
				uint8(clampFloat((float64(c.B)-127.5)*factor+127.5, 0, 255)), c.A,
			})
		}
	}
	return dst
}

func Blur(img image.Image, radius int) (*image.RGBA, error) {
	if radius < 1 || radius > 50 {
		return nil, &ParameterError{"blur radius must be between 1 and 50"}
	}
	sigma := math.Max(float64(radius)/2, 0.5)
	size := radius*2 + 1
	kernel := make([]float64, size*size)
	var sum float64
	for y := -radius; y <= radius; y++ {
		for x := -radius; x <= radius; x++ {
			v := math.Exp(-(float64(x*x + y*y)) / (2 * sigma * sigma))
			kernel[(y+radius)*size+x+radius] = v
			sum += v
		}
	}
	for i := range kernel {
		kernel[i] /= sum
	}
	return Convolve(img, kernel, size)
}

func Sharpen(img image.Image, amount float64) (*image.RGBA, error) {
	if amount < 0 || amount > 10 {
		return nil, &ParameterError{"sharpen amount must be between 0 and 10"}
	}
	k := []float64{-amount, -amount, -amount, -amount, 1 + 8*amount, -amount, -amount, -amount, -amount}
	return Convolve(img, k, 3)
}

func Convolve(img image.Image, kernel []float64, size int) (*image.RGBA, error) {
	if size <= 0 || size%2 == 0 || len(kernel) != size*size {
		return nil, &ParameterError{"kernel must be a non-empty odd square"}
	}
	src := toRGBA(img)
	dst := image.NewRGBA(src.Bounds())
	r := size / 2
	for y := 0; y < src.Bounds().Dy(); y++ {
		for x := 0; x < src.Bounds().Dx(); x++ {
			var rr, gg, bb, aa float64
			for ky := -r; ky <= r; ky++ {
				for kx := -r; kx <= r; kx++ {
					cx, cy := clampInt(x+kx, 0, src.Bounds().Dx()-1), clampInt(y+ky, 0, src.Bounds().Dy()-1)
					c := src.RGBAAt(cx, cy)
					w := kernel[(ky+r)*size+kx+r]
					rr += float64(c.R) * w
					gg += float64(c.G) * w
					bb += float64(c.B) * w
					aa += float64(c.A) * w
				}
			}
			// Preserve alpha unless a kernel explicitly changes it.
			if aa == 0 {
				aa = float64(src.RGBAAt(x, y).A)
			}
			dst.SetRGBA(x, y, color.RGBA{uint8(clampFloat(rr, 0, 255)), uint8(clampFloat(gg, 0, 255)), uint8(clampFloat(bb, 0, 255)), uint8(clampFloat(aa, 0, 255))})
		}
	}
	return dst, nil
}
