package cli

import "github.com/spf13/cobra"

func init() {
	var quality, pngLevel int
	var format string
	var tiffLZW bool
	cmd := &cobra.Command{Use: "compress <input> <output>", Short: "Re-encode an image with format-specific compression settings", Args: cobra.ExactArgs(2), Example: `  imgx compress photo.png optimized.png --png-level 9
  imgx compress photo.png optimized.jpg --quality 82`, RunE: func(cmd *cobra.Command, args []string) error {
		img, src, err := load(args[0])
		if err != nil {
			return err
		}
		f, err := save(img, args[1], format, quality, pngLevel, tiffLZW)
		if err != nil {
			return err
		}
		success(cmd, "compressed %s (%s) -> %s", args[0], src, f)
		return nil
	}}
	cmd.Flags().IntVar(&quality, "quality", 80, "JPEG quality (1-100)")
	cmd.Flags().IntVar(&pngLevel, "png-level", -1, "PNG compression level: -1 default, -2 Huffman-only, 0 none, 1-9 fixed")
	cmd.Flags().BoolVar(&tiffLZW, "tiff-lzw", true, "compress TIFF output (Deflate; LZW encoding is unavailable in x/image/tiff)")
	cmd.Flags().StringVarP(&format, "format", "f", "", "output format override")
	rootCmd.AddCommand(cmd)
}
