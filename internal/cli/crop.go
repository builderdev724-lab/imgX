package cli

import (
	"github.com/priyojeet/imgx/internal/imaging"
	"github.com/spf13/cobra"
)

func init() {
	var x, y, w, h int
	var format string
	cmd := &cobra.Command{Use: "crop <input> <output>", Short: "Crop a rectangular region", Args: cobra.ExactArgs(2), RunE: func(cmd *cobra.Command, args []string) error {
		img, _, err := load(args[0])
		if err != nil {
			return err
		}
		out, err := imaging.Crop(img, x, y, w, h)
		if err != nil {
			return err
		}
		f, err := save(out, args[1], format, 90, -1, true)
		if err != nil {
			return err
		}
		success(cmd, "cropped -> %s", f)
		return nil
	}}
	cmd.Flags().IntVar(&x, "x", 0, "left coordinate")
	cmd.Flags().IntVar(&y, "y", 0, "top coordinate")
	cmd.Flags().IntVar(&w, "w", 0, "crop width")
	cmd.Flags().IntVar(&h, "h", 0, "crop height")
	cmd.Flags().StringVarP(&format, "format", "f", "", "output format override")
	rootCmd.AddCommand(cmd)
}
