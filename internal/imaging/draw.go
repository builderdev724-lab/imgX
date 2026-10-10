package imaging

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
)

// ParseColor accepts #RGB, #RGBA, #RRGGBB and #RRGGBBAA hex values, with or
// without a leading '#', as well as the CSS colour names in cssNames.
func ParseColor(s string) (color.RGBA, error) {
	v := strings.TrimSpace(s)
	if v == "" {
		return color.RGBA{}, &ParameterError{"color must not be empty"}
	}

	if named, ok := cssNames[strings.ToLower(v)]; ok {
		return named, nil
	}

	hex := strings.TrimPrefix(v, "#")
	if len(hex) == 3 || len(hex) == 4 {
		// Expand the shorthand form so #1a2 becomes #11aa22.
		var sb strings.Builder
		for _, r := range hex {
			sb.WriteRune(r)
			sb.WriteRune(r)
		}
		hex = sb.String()
	}
	if len(hex) != 6 && len(hex) != 8 {
		return color.RGBA{}, &ParameterError{fmt.Sprintf("invalid color %q", s)}
	}

	parsed, err := strconv.ParseUint(hex, 16, 64)
	if err != nil {
		return color.RGBA{}, &ParameterError{fmt.Sprintf("invalid color %q", s)}
	}

	if len(hex) == 6 {
		return color.RGBA{
			R: uint8(parsed >> 16),
			G: uint8(parsed >> 8),
			B: uint8(parsed),
			A: 0xff,
		}, nil
	}
	return color.RGBA{
		R: uint8(parsed >> 24),
		G: uint8(parsed >> 16),
		B: uint8(parsed >> 8),
		A: uint8(parsed),
	}, nil
}

// AddBorder surrounds the image with a frame of the given width in pixels.
// A zero width is a no-op; a negative width is rejected.
func AddBorder(img image.Image, width int, c color.RGBA) (*image.RGBA, error) {
	if width < 0 {
		return nil, &ParameterError{"border width must not be negative"}
	}
	src := toRGBA(img)
	if width == 0 {
		return src, nil
	}
	b := src.Bounds()
	w, h := b.Dx(), b.Dy()
	dst := image.NewRGBA(image.Rect(0, 0, w+2*width, h+2*width))
	draw.Draw(dst, dst.Bounds(), &image.Uniform{C: c}, image.Point{}, draw.Src)
	draw.Draw(dst, image.Rect(width, width, width+w, width+h), src, b.Min, draw.Over)
	return dst, nil
}

// WatermarkPosition selects where Watermark anchors its text.
type WatermarkPosition string

const (
	PosTopLeft     WatermarkPosition = "top-left"
	PosTop         WatermarkPosition = "top"
	PosTopRight    WatermarkPosition = "top-right"
	PosLeft        WatermarkPosition = "left"
	PosCenter      WatermarkPosition = "center"
	PosRight       WatermarkPosition = "right"
	PosBottomLeft  WatermarkPosition = "bottom-left"
	PosBottom      WatermarkPosition = "bottom"
	PosBottomRight WatermarkPosition = "bottom-right"
)

var (
	fontOnce sync.Once
	fontFace *sfnt.Font
	fontErr  error
)

func watermarkFont() (*sfnt.Font, error) {
	fontOnce.Do(func() {
		fontFace, fontErr = sfnt.Parse(goregular.TTF)
	})
	return fontFace, fontErr
}

// Watermark draws text over the image. size is the font size in points, and
// opacity is a 0-100 percentage applied to the colour's alpha channel.
func Watermark(img image.Image, text string, pos WatermarkPosition, size, opacity float64, c color.RGBA) (*image.RGBA, error) {
	if strings.TrimSpace(text) == "" {
		return nil, &ParameterError{"watermark text must not be empty"}
	}
	if size <= 0 || math.IsNaN(size) || math.IsInf(size, 0) {
		return nil, &ParameterError{"watermark size must be a positive number"}
	}
	if opacity < 0 {
		opacity = 0
	}
	if opacity > 100 {
		opacity = 100
	}

	f, err := watermarkFont()
	if err != nil {
		return nil, &ParameterError{"watermark font could not be loaded"}
	}

	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return nil, &ParameterError{"watermark face could not be created"}
	}
	defer face.Close()

	metrics := face.Metrics()
	ascent := float64(metrics.Ascent) / 64
	descent := float64(metrics.Descent) / 64
	advance := font.MeasureString(face, text)
	if advance <= 0 {
		return nil, &ParameterError{"watermark text could not be measured"}
	}
	textW := float64(advance) / 64
	textH := ascent + descent

	dst := toRGBA(img)
	imgW, imgH := float64(dst.Bounds().Dx()), float64(dst.Bounds().Dy())

	margin := math.Max(4, size*0.4)
	left := margin
	center := (imgW - textW) / 2
	right := imgW - margin - textW

	top := margin + ascent
	middle := (imgH-textH)/2 + ascent
	bottom := imgH - margin - descent

	var x, y float64
	switch WatermarkPosition(strings.ToLower(strings.TrimSpace(string(pos)))) {
	case PosTopLeft:
		x, y = left, top
	case PosTop:
		x, y = center, top
	case PosTopRight:
		x, y = right, top
	case PosLeft:
		x, y = left, middle
	case PosCenter:
		x, y = center, middle
	case PosRight:
		x, y = right, middle
	case PosBottomLeft:
		x, y = left, bottom
	case PosBottom:
		x, y = center, bottom
	case PosBottomRight, "":
		x, y = right, bottom
	default:
		return nil, &ParameterError{fmt.Sprintf("unknown watermark position %q", pos)}
	}

	// Keep the text inside the canvas even when it is wider or taller.
	x = math.Max(margin, math.Min(x, imgW-margin-textW))
	y = math.Max(ascent, math.Min(y, imgH-descent))

	baselineX := int(math.Round(x))
	baselineY := int(math.Round(y))
	ink := color.RGBA{R: c.R, G: c.G, B: c.B, A: uint8(math.Round(clampFloat(opacity/100, 0, 1) * 255))}

	rect := image.Rect(
		baselineX,
		baselineY-int(math.Round(ascent)),
		baselineX+int(math.Ceil(textW)),
		baselineY+int(math.Ceil(descent)),
	)
	draw.Draw(dst, rect.Intersect(dst.Bounds()), &image.Uniform{C: ink}, image.Point{}, draw.Over)

	d := font.Drawer{
		Dst:  dst,
		Src:  image.NewUniform(ink),
		Face: face,
		Dot:  fixed.Point26_6{X: fixed.I(baselineX), Y: fixed.I(baselineY)},
	}
	d.DrawString(text)
	return dst, nil
}

// cssNames holds the CSS colour keywords accepted by ParseColor.
var cssNames = map[string]color.RGBA{
	"transparent": {0, 0, 0, 0},
	"black":       {0, 0, 0, 255},
	"silver":      {192, 192, 192, 255},
	"gray":        {128, 128, 128, 255},
	"grey":        {128, 128, 128, 255},
	"white":       {255, 255, 255, 255},
	"maroon":      {128, 0, 0, 255},
	"red":         {255, 0, 0, 255},
	"purple":      {128, 0, 128, 255},
	"fuchsia":     {255, 0, 255, 255},
	"magenta":     {255, 0, 255, 255},
	"green":       {0, 128, 0, 255},
	"lime":        {0, 255, 0, 255},
	"olive":       {128, 128, 0, 255},
	"yellow":      {255, 255, 0, 255},
	"navy":        {0, 0, 128, 255},
	"blue":        {0, 0, 255, 255},
	"teal":        {0, 128, 128, 255},
	"aqua":        {0, 255, 255, 255},
	"cyan":        {0, 255, 255, 255},
	"orange":      {255, 165, 0, 255},
	"pink":        {255, 192, 203, 255},
	"brown":       {165, 42, 42, 255},
	"gold":        {255, 215, 0, 255},
	"indigo":      {75, 0, 130, 255},
	"violet":      {238, 130, 238, 255},
	"salmon":      {250, 128, 114, 255},
	"khaki":       {240, 230, 140, 255},
	"coral":       {255, 127, 80, 255},
	"crimson":     {220, 20, 60, 255},
	"tomato":      {255, 99, 71, 255},
	"orchid":      {218, 112, 214, 255},
	"plum":        {221, 160, 221, 255},
	"tan":         {210, 180, 140, 255},
	"beige":       {245, 245, 220, 255},
	"ivory":       {255, 255, 240, 255},
	"lavender":    {230, 230, 250, 255},
	"turquoise":   {64, 224, 208, 255},
	"azure":       {240, 255, 255, 255},
	"chocolate":   {210, 105, 30, 255},
	"sienna":      {160, 82, 45, 255},
	"peru":        {205, 133, 63, 255},
	"linen":       {250, 240, 230, 255},
	"wheat":       {245, 222, 179, 255},
	"darkgray":    {169, 169, 169, 255},
	"darkgrey":    {169, 169, 169, 255},
	"lightgray":   {211, 211, 211, 255},
	"lightgrey":   {211, 211, 211, 255},
	"dimgray":     {105, 105, 105, 255},
	"dimgrey":     {105, 105, 105, 255},
	"whitesmoke":  {245, 245, 245, 255},
}
