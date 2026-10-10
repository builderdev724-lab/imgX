package cli

import (
	"github.com/priyojeet/imgx/internal/imaging"
	"github.com/spf13/cobra"
	"image"
)

func init() {
	var deg float64
	var algo, format string
	cmd := &cobra.Command{Use: "rotate <input> <output>", Short: "Rotate by any angle in degrees", Args: cobra.ExactArgs(2), Example: `  imgx rotate photo.png rotated.png --deg 90
  imgx rotate photo.png rotated.png --deg 17.5 --algo bicubic`, RunE: func(cmd *cobra.Command, args []string) error {
		a, err := parseResizeAlgo(algo)
		if err != nil {
			return err
		}
		img, _, err := load(args[0])
		if err != nil {
			return err
		}
		var out image.Image
		out, err = imaging.Rotate(img, deg, a)
		if err != nil {
			return err
		}
		f, err := save(out, args[1], format, 90, -1, true)
		if err != nil {
			return err
		}
		success(cmd, "rotated %.2f degrees -> %s", deg, f)
		return nil
	}}
	cmd.Flags().Float64Var(&deg, "deg", 0, "rotation angle in degrees")
	cmd.Flags().StringVar(&algo, "algo", "bilinear", "nearest, bilinear, bicubic, or lanczos3")
	cmd.Flags().StringVarP(&format, "format", "f", "", "output format override")
	rootCmd.AddCommand(cmd)
}
