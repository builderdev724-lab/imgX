package cli

import (
	"fmt"
	"os"

	"github.com/priyojeet/imgx/internal/version"
	"github.com/spf13/cobra"
)

var maxPixels int64
var quiet bool

var rootCmd = &cobra.Command{
	Use:           "imgx",
	Short:         "Classical image processing from the command line",
	Long:          "imgx is a fast, cross-platform image-processing CLI for conversion, transforms, filters, compositing, and batch workflows.",
	SilenceUsage:  true,
	SilenceErrors: true,
	Version:       version.Version,
}

func init() {
	rootCmd.PersistentFlags().Int64Var(&maxPixels, "max-pixels", 100_000_000, "maximum decoded pixels (0 disables the limit)")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "suppress success messages")
	rootCmd.SetOut(os.Stdout)
	rootCmd.SetErr(os.Stderr)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "imgx:", err)
		os.Exit(1)
	}
}

func success(cmd *cobra.Command, format string, args ...any) {
	if quiet {
		return
	}
	fmt.Fprintf(cmd.ErrOrStderr(), format+"\n", args...)
}
