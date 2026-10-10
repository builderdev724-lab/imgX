package cli

import (
	"fmt"
	"image"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/priyojeet/imgx/internal/imaging"
	"github.com/spf13/cobra"
)

func init() {
	var outDir, opText, format string
	var jobs int
	var overwrite bool
	cmd := &cobra.Command{Use: "batch <glob>", Short: "Run one operation over matching images", Args: cobra.ExactArgs(1), Example: `  imgx batch "./photos/*.jpg" --op "resize --w 1200 --h 1200 --keep-aspect" --out-dir ./out
  imgx batch "./photos/*.png" --op "grayscale" --out-dir ./gray`, RunE: func(cmd *cobra.Command, args []string) error {
		if jobs <= 0 {
			jobs = runtime.NumCPU()
		}
		files, err := filepath.Glob(args[0])
		if err != nil {
			return err
		}
		if len(files) == 0 {
			return fmt.Errorf("pattern matched no files: %s", args[0])
		}
		op, err := parseOpText(opText)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(outDir, 0o755); err != nil {
			return err
		}
		jobs = minInt(jobs, len(files))
		sem := make(chan struct{}, jobs)
		var wg sync.WaitGroup
		var mu sync.Mutex
		var firstErr error
		done := 0
		for _, input := range files {
			input := input
			wg.Add(1)
			go func() {
				defer wg.Done()
				sem <- struct{}{}
				defer func() { <-sem }()
				if err := processBatchFile(input, outDir, format, op, overwrite); err != nil {
					mu.Lock()
					if firstErr == nil {
						firstErr = err
					}
					mu.Unlock()
					return
				}
				mu.Lock()
				done++
				if !quiet {
					fmt.Fprintf(cmd.ErrOrStderr(), "[%d/%d] %s\n", done, len(files), filepath.Base(input))
				}
				mu.Unlock()
			}()
		}
		wg.Wait()
		return firstErr
	}}
	cmd.Flags().StringVar(&opText, "op", "", "operation and arguments, e.g. 'resize --w 800 --h 600'")
	cmd.Flags().StringVar(&outDir, "out-dir", "./out", "output directory")
	cmd.Flags().StringVarP(&format, "format", "f", "", "override output format/extension")
	cmd.Flags().IntVar(&jobs, "jobs", runtime.NumCPU(), "parallel workers")
	cmd.Flags().BoolVar(&overwrite, "overwrite", false, "allow replacement of existing output files")
	_ = cmd.MarkFlagRequired("op")
	rootCmd.AddCommand(cmd)
}

func processBatchFile(input, outDir, format string, op imaging.Operation, overwrite bool) error {
	img, _, err := load(input)
	if err != nil {
		return fmt.Errorf("%s: %w", input, err)
	}
	name := filepath.Base(input)
	ext := filepath.Ext(name)
	if strings.TrimSpace(format) != "" {
		ext = "." + strings.TrimPrefix(strings.ToLower(format), ".")
	} else if v := op.Args["format"]; v != "" {
		format = v
		ext = "." + strings.TrimPrefix(strings.ToLower(format), ".")
	}
	output := filepath.Join(outDir, strings.TrimSuffix(name, filepath.Ext(name))+ext)
	if !overwrite {
		if _, err := os.Stat(output); err == nil {
			return fmt.Errorf("output exists (use --overwrite): %s", output)
		}
	}
	var out image.Image = img
	nameOp := strings.ToLower(op.Name)
	if nameOp != "compress" && nameOp != "convert" {
		out, err = imaging.Apply(img, op)
		if err != nil {
			return fmt.Errorf("%s: %w", input, err)
		}
	}
	quality := 90
	if v := op.Args["quality"]; v != "" {
		quality, err = strconv.Atoi(v)
		if err != nil {
			return fmt.Errorf("%s: invalid --quality %q", input, v)
		}
	}
	pngLevel := -1
	if v := op.Args["png-level"]; v != "" {
		pngLevel, err = strconv.Atoi(v)
		if err != nil {
			return fmt.Errorf("%s: invalid --png-level %q", input, v)
		}
	}
	_, err = save(out, output, format, quality, pngLevel, true)
	if err != nil {
		return fmt.Errorf("%s: %w", input, err)
	}
	return nil
}

func parseOpText(s string) (imaging.Operation, error) {
	tokens, err := shellLikeSplit(s)
	if err != nil {
		return imaging.Operation{}, err
	}
	if len(tokens) == 0 {
		return imaging.Operation{}, fmt.Errorf("--op cannot be empty")
	}
	op := imaging.Operation{Name: tokens[0], Args: map[string]string{}}
	for i := 1; i < len(tokens); i++ {
		t := tokens[i]
		if !strings.HasPrefix(t, "--") {
			return imaging.Operation{}, fmt.Errorf("expected --flag in batch op, got %q", t)
		}
		key := strings.TrimPrefix(t, "--")
		val := "true"
		if i+1 < len(tokens) && !strings.HasPrefix(tokens[i+1], "--") {
			val = tokens[i+1]
			i++
		}
		op.Args[key] = val
	}
	return op, nil
}

func shellLikeSplit(s string) ([]string, error) {
	var out []string
	var b strings.Builder
	var quote rune
	esc := false
	flush := func() {
		if b.Len() > 0 {
			out = append(out, b.String())
			b.Reset()
		}
	}
	for _, r := range s {
		if esc {
			b.WriteRune(r)
			esc = false
			continue
		}
		if r == '\\' && quote == 0 {
			esc = true
			continue
		}
		if quote != 0 {
			if r == quote {
				quote = 0
			} else {
				b.WriteRune(r)
			}
			continue
		}
		if r == '\'' || r == '"' {
			quote = r
			continue
		}
		if r == ' ' || r == '\t' || r == '\n' {
			flush()
			continue
		}
		b.WriteRune(r)
	}
	if esc {
		return nil, fmt.Errorf("unterminated escape in --op")
	}
	if quote != 0 {
		return nil, fmt.Errorf("unterminated quote in --op")
	}
	flush()
	return out, nil
}
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
