# Security policy

Please do not disclose image-processing vulnerabilities in a public issue before a fix is available.

For security-sensitive reports, use GitHub's private vulnerability reporting for this repository.

Important implementation detail: imgx validates image dimensions from `DecodeConfig` before decoding and exposes `--max-pixels` to limit decoded image size. This is a defense against accidental or hostile oversized-image allocations, not a complete sandbox.
