package cli

import (
	"fmt"
	"github.com/priyojeet/imgx/internal/imaging"
	"github.com/spf13/cobra"
)

func init() {
	var horizontal, vertical bool
	var format string
	cmd := &cobra.Command{Use: "flip <input> <output>", Short: "Flip horizontally, vertically, or both", Args: cobra.ExactArgs(2), RunE: func(cmd *cobra.Command, args []string) error {
		if !horizontal && !vertical {
			return fmt.Errorf("use --horizontal and/or --vertical")
		}
		img, _, err := load(args[0])
		if err != nil {
			return err
		}
		out, err := imaging.Flip(img, horizontal, vertical)
		if err != nil {
			return err
		}
		f, err := save(out, args[1], format, 90, -1, true)
		if err != nil {
			return err
		}
		success(cmd, "flipped -> %s", f)
		return nil
	}}
	cmd.Flags().BoolVar(&horizontal, "horizontal", false, "flip left-to-right")
	cmd.Flags().BoolVar(&vertical, "vertical", false, "flip top-to-bottom")
	cmd.Flags().StringVarP(&format, "format", "f", "", "output format override")
	rootCmd.AddCommand(cmd)
}
