package imaging

import (
	"bytes"
	"fmt"
	"image/gif"
	"os"
)

type Info struct {
	Path       string
	Format     string
	Width      int
	Height     int
	ColorModel string
	Frames     int
	SizeBytes  int64
}

func GetInfo(path string, maxPixels int64) (Info, error) {
	st, err := os.Stat(path)
	if err != nil {
		return Info{}, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return Info{}, err
	}
	img, format, err := DecodeBuffered(data, maxPixels)
	if err != nil {
		return Info{}, err
	}
	frames := 1
	if format == "gif" {
		if g, err := gif.DecodeAll(bytes.NewReader(data)); err == nil && len(g.Image) > 0 {
			frames = len(g.Image)
		}
	}
	return Info{Path: path, Format: format, Width: img.Bounds().Dx(), Height: img.Bounds().Dy(), ColorModel: fmt.Sprintf("%T", img.ColorModel()), Frames: frames, SizeBytes: st.Size()}, nil
}
