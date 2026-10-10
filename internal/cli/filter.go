package cli

import (
	"fmt"
	"github.com/priyojeet/imgx/internal/imaging"
	"github.com/spf13/cobra"
	"image"
)

func init() {
	registerSimpleFilter("grayscale", "Convert to grayscale", func(img image.Image, _ float64) (image.Image, error) { return imaging.Grayscale(img), nil }, 0, -100, 100)
	registerSimpleFilter("brightness", "Adjust brightness (-100 to 100)", func(img image.Image, v float64) (image.Image, error) {
		if err := requireRange("brightness", v, -100, 100); err != nil {
			return nil, err
		}
		return imaging.Brightness(img, v), nil
	}, 0, -100, 100)
	registerSimpleFilter("contrast", "Adjust contrast (-100 to 100)", func(img image.Image, v float64) (image.Image, error) {
		if err := requireRange("contrast", v, -100, 100); err != nil {
			return nil, err
		}
		return imaging.Contrast(img, v), nil
	}, 0, -100, 100)
	{
		var radius int
		var format string
		cmd := &cobra.Command{Use: "blur <input> <output>", Short: "Apply a Gaussian blur", Args: cobra.ExactArgs(2), RunE: func(cmd *cobra.Command, args []string) error {
			img, _, err := load(args[0])
			if err != nil {
				return err
			}
			out, err := imaging.Blur(img, radius)
			if err != nil {
				return err
			}
			f, err := save(out, args[1], format, 90, -1, true)
			if err != nil {
				return err
			}
			success(cmd, "blurred -> %s", f)
			return nil
		}}
		cmd.Flags().IntVar(&radius, "radius", 5, "blur radius (1-50)")
		cmd.Flags().StringVarP(&format, "format", "f", "", "output format override")
		rootCmd.AddCommand(cmd)
	}
	{
		var amount float64
		var format string
		cmd := &cobra.Command{Use: "sharpen <input> <output>", Short: "Sharpen using an unsharp-style convolution", Args: cobra.ExactArgs(2), RunE: func(cmd *cobra.Command, args []string) error {
			img, _, err := load(args[0])
			if err != nil {
				return err
			}
			out, err := imaging.Sharpen(img, amount)
			if err != nil {
				return err
			}
			f, err := save(out, args[1], format, 90, -1, true)
			if err != nil {
				return err
			}
			success(cmd, "sharpened -> %s", f)
			return nil
		}}
		cmd.Flags().Float64Var(&amount, "amount", 1, "sharpen amount (0-10)")
		cmd.Flags().StringVarP(&format, "format", "f", "", "output format override")
		rootCmd.AddCommand(cmd)
	}
}

func registerSimpleFilter(name, short string, fn func(image.Image, float64) (image.Image, error), def, minV, maxV float64) {
	var value float64
	var format string
	cmd := &cobra.Command{Use: fmt.Sprintf("%s <input> <output>", name), Short: short, Args: cobra.ExactArgs(2), RunE: func(cmd *cobra.Command, args []string) error {
		img, _, err := load(args[0])
		if err != nil {
			return err
		}
		out, err := fn(img, value)
		if err != nil {
			return err
		}
		f, err := save(out, args[1], format, 90, -1, true)
		if err != nil {
			return err
		}
		success(cmd, fmt.Sprintf("%s -> %%s", name), f)
		return nil
	}}
	if name != "grayscale" {
		cmd.Flags().Float64Var(&value, "value", def, fmt.Sprintf("value (%v to %v)", minV, maxV))
	}
	cmd.Flags().StringVarP(&format, "format", "f", "", "output format override")
	rootCmd.AddCommand(cmd)
}
