package cli

import (
	"github.com/priyojeet/imgx/internal/imaging"
	"github.com/spf13/cobra"
)

func init() {
	var width int
	var colorValue, format string
	cmd := &cobra.Command{Use: "border <input> <output>", Short: "Add a solid-color border", Args: cobra.ExactArgs(2), RunE: func(cmd *cobra.Command, args []string) error {
		c, err := imaging.ParseColor(colorValue)
		if err != nil {
			return err
		}
		img, _, err := load(args[0])
		if err != nil {
			return err
		}
		out, err := imaging.AddBorder(img, width, c)
		if err != nil {
			return err
		}
		f, err := save(out, args[1], format, 90, -1, true)
		if err != nil {
			return err
		}
		success(cmd, "border %dpx -> %s", width, f)
		return nil
	}}
	cmd.Flags().IntVar(&width, "width", 10, "border width in pixels")
	cmd.Flags().StringVar(&colorValue, "color", "#000000", "CSS color name or #RRGGBB/#RRGGBBAA")
	cmd.Flags().StringVarP(&format, "format", "f", "", "output format override")
	rootCmd.AddCommand(cmd)
}
