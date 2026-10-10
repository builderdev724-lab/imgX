package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/priyojeet/imgx/internal/imaging"
	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use: "info <input>", Short: "Show image metadata and dimensions", Args: cobra.ExactArgs(1),
		Example: `  imgx info photo.jpg
  imgx info image.webp`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if args[0] == "-" {
				data, err := io.ReadAll(os.Stdin)
				if err != nil {
					return err
				}
				img, format, err := imaging.DecodeBuffered(data, maxPixels)
				if err != nil {
					return err
				}
				fmt.Fprintf(cmd.OutOrStdout(), "format: %s\nwidth: %d\nheight: %d\ncolor-model: %T\n", format, img.Bounds().Dx(), img.Bounds().Dy(), img.ColorModel())
				return nil
			}
			in, err := imaging.GetInfo(args[0], maxPixels)
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "file: %s\nformat: %s\nsize: %s (%d bytes)\ndimensions: %dx%d\naspect-ratio: %.4f\ncolor-model: %s\nframes: %d\n", filepath.Clean(in.Path), in.Format, formatBytes(in.SizeBytes), in.SizeBytes, in.Width, in.Height, float64(in.Width)/float64(in.Height), in.ColorModel, in.Frames)
			return nil
		},
	}
	rootCmd.AddCommand(cmd)
}

func formatBytes(n int64) string {
	units := []string{"B", "KiB", "MiB", "GiB", "TiB"}
	v := float64(n)
	i := 0
	for v >= 1024 && i < len(units)-1 {
		v /= 1024
		i++
	}
	return fmt.Sprintf("%.2f %s", v, units[i])
}
