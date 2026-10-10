package cli

import (
	"fmt"
	"github.com/priyojeet/imgx/internal/imaging"
	"github.com/spf13/cobra"
	"image"
)

func init() {
	var w, h int
	var keep bool
	var algo string
	var format string
	cmd := &cobra.Command{Use: "resize <input> <output>", Short: "Resize an image", Args: cobra.ExactArgs(2), Example: `  imgx resize photo.jpg small.jpg --w 800 --h 600
  imgx resize photo.jpg small.jpg --w 800 --keep-aspect --algo lanczos3`, RunE: func(cmd *cobra.Command, args []string) error {
		if keep {
			if w <= 0 && h <= 0 {
				return fmt.Errorf("--keep-aspect requires --w and/or --h")
			}
		} else if w <= 0 || h <= 0 {
			return fmt.Errorf("--w and --h are required unless --keep-aspect is used")
		}
		a, err := parseResizeAlgo(algo)
		if err != nil {
			return err
		}
		img, _, err := load(args[0])
		if err != nil {
			return err
		}
		var out image.Image
		if keep {
			out, err = imaging.ResizeKeepAspect(img, w, h, a)
		} else {
			out, err = imaging.Resize(img, w, h, a)
		}
		if err != nil {
			return err
		}
		f, err := save(out, args[1], format, 90, -1, true)
		if err != nil {
			return err
		}
		success(cmd, "resized -> %s", f)
		return nil
	}}
	cmd.Flags().IntVarP(&w, "w", "W", 0, "target width")
	cmd.Flags().IntVarP(&h, "h", "H", 0, "target height")
	cmd.Flags().BoolVar(&keep, "keep-aspect", false, "fit inside the requested bounds while preserving aspect ratio")
	cmd.Flags().StringVar(&algo, "algo", "bilinear", "nearest, bilinear, bicubic, or lanczos3")
	cmd.Flags().StringVarP(&format, "format", "f", "", "output format override")
	rootCmd.AddCommand(cmd)
}
