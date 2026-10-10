# Command reference

## `convert`

```text
imgx convert <input> <output> [--quality 1-100] [--format FORMAT] [--png-level LEVEL]
```

Formats: PNG, JPEG, GIF, BMP, TIFF output. WebP input is supported.

## `resize`

```text
imgx resize <input> <output> --w WIDTH --h HEIGHT [--algo ALGO]
imgx resize <input> <output> --w WIDTH [--h HEIGHT] --keep-aspect [--algo ALGO]
```

Algorithms: `nearest`, `bilinear`, `bicubic`, `lanczos3`.

## `crop`

```text
imgx crop <input> <output> --x X --y Y --w WIDTH --h HEIGHT
```

The rectangle must be entirely inside the source image.

## `rotate`

```text
imgx rotate <input> <output> --deg DEGREES [--algo ALGO]
```

Multiples of 90 degrees use exact pixel transforms. Other angles use inverse mapping and the selected interpolation algorithm.

## `flip`

```text
imgx flip <input> <output> --horizontal
imgx flip <input> <output> --vertical
```

Both flags may be used together for a 180-degree reflection.

## `grayscale`

```text
imgx grayscale <input> <output>
```

Uses a luminance-weighted RGB conversion.

## `brightness`

```text
imgx brightness <input> <output> --value -100..100
```

The value is a percentage of the full 8-bit channel range.

## `contrast`

```text
imgx contrast <input> <output> --value -100..100
```

The value adjusts contrast around the 127.5 midpoint.

## `blur`

```text
imgx blur <input> <output> --radius 1..50
```

Uses a generated Gaussian kernel and the generic convolution engine.

## `sharpen`

```text
imgx sharpen <input> <output> --amount 0..10
```

Uses a high-pass/edge-emphasis kernel.

## `border`

```text
imgx border <input> <output> --width 10 --color '#000000'
```

Colors may be CSS color names or `#RRGGBB` / `#RRGGBBAA`.

## `watermark`

```text
imgx watermark <input> <output> --text 'Hello' \
  --pos bottom-right --size 32 --opacity 70 --color '#ffffff'
```

Positions: `top-left`, `top`, `top-right`, `left`, `center`, `right`, `bottom-left`, `bottom`, `bottom-right`.

The font is embedded from Go's regular font package, so watermarking does not depend on fonts installed on the target machine.

## `compress`

```text
imgx compress <input> <output> --quality 82 --png-level 9 --tiff-lzw
```

JPEG uses `--quality`; PNG uses `--png-level`; TIFF uses LZW by default when enabled.

## `info`

```text
imgx info <input>
```

Displays format, file size, dimensions, aspect ratio, color model, and GIF frame count.

## `batch`

```text
imgx batch '<glob>' --op 'resize --w 1200 --h 1200 --keep-aspect' --out-dir ./out --jobs 8
```

The `--op` value is parsed by a small shell-like tokenizer, supporting single/double quotes and backslash escaping.

## Global flags

```text
--max-pixels N    maximum decoded pixel count; 0 disables the guard
--quiet, -q       suppress success/progress messages
--help, -h        command help
--version         version information
```
