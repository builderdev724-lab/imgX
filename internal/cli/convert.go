package cli

import (
	"fmt"
	"github.com/spf13/cobra"
)

var convertQuality int
var convertFormat string
var convertPNGLevel int

func init() {
	cmd := &cobra.Command{
		Use:   "convert <input> <output>",
		Short: "Convert an image between supported formats",
		Args:  cobra.ExactArgs(2),
		Example: `  imgx convert photo.png photo.jpg --quality 92
  imgx convert photo.webp photo.png
  cat photo.png | imgx convert - out.jpg`,
		RunE: func(cmd *cobra.Command, args []string) error {
			img, srcFormat, err := load(args[0])
			if err != nil {
				return err
			}
			fmtHint := convertFormat
			if fmtHint == "" && args[1] == "-" {
				return fmt.Errorf("stdout output requires --format")
			}
			format, err := save(img, args[1], fmtHint, convertQuality, convertPNGLevel, true)
			if err != nil {
				return err
			}
			success(cmd, "converted %s (%s) -> %s", args[0], srcFormat, format)
			return nil
		},
	}
	cmd.Flags().IntVar(&convertQuality, "quality", 90, "JPEG quality (1-100)")
	cmd.Flags().StringVarP(&convertFormat, "format", "f", "", "output format when writing to stdout or overriding the extension")
	cmd.Flags().IntVar(&convertPNGLevel, "png-level", -1, "PNG compression level (-1 default, 0-9 fixed level; -2 is Huffman-only)")
	rootCmd.AddCommand(cmd)
}
