package cli

import (
	"github.com/priyojeet/imgx/internal/imaging"
	"github.com/spf13/cobra"
)

func init() {
	var text, pos, colorValue, format string
	var size, opacity float64
	cmd := &cobra.Command{Use: "watermark <input> <output>", Short: "Draw text onto an image", Args: cobra.ExactArgs(2), Example: `  imgx watermark photo.jpg marked.jpg --text "Copyright 2026" --pos bottom-right
  imgx watermark photo.jpg marked.jpg --text "IMGX" --size 48 --opacity 60 --color "#ffffff"`, RunE: func(cmd *cobra.Command, args []string) error {
		c, err := imaging.ParseColor(colorValue)
		if err != nil {
			return err
		}
		img, _, err := load(args[0])
		if err != nil {
			return err
		}
		out, err := imaging.Watermark(img, text, imaging.WatermarkPosition(pos), size, opacity, c)
		if err != nil {
			return err
		}
		f, err := save(out, args[1], format, 90, -1, true)
		if err != nil {
			return err
		}
		success(cmd, "watermark -> %s", f)
		return nil
	}}
	cmd.Flags().StringVar(&text, "text", "", "watermark text")
	cmd.Flags().StringVar(&pos, "pos", "bottom-right", "top-left, top, top-right, left, center, right, bottom-left, bottom, bottom-right")
	cmd.Flags().Float64Var(&size, "size", 32, "font size")
	cmd.Flags().Float64Var(&opacity, "opacity", 70, "text opacity (0-100)")
	cmd.Flags().StringVar(&colorValue, "color", "#ffffff", "text color")
	cmd.Flags().StringVarP(&format, "format", "f", "", "output format override")
	_ = cmd.MarkFlagRequired("text")
	rootCmd.AddCommand(cmd)
}
