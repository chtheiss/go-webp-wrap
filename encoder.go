// package webpwrap provides a Go wrapper for the WebP image compression tools.
// It allows for easy conversion of images to WebP format with various options
// including quality control, cropping, and different input/output methods.
package webpwrap

import (
	"context"
	"image"
	"io"
)

// Encoder encodes image.Image into WebP format using cwebp.
// It provides control over the encoding quality and other parameters.
type Encoder struct {
	// Quality specifies the compression factor for RGB channels.
	// The value must be between 0 and 100, where:
	// - A small factor produces a smaller file with lower quality
	// - A value of 100 achieves the best quality
	// - The default is 75
	Quality uint
}

// Encode writes the Image m to w in WebP format.
// Any Image type may be encoded.
//
// Parameters:
//   - w: The io.Writer to write the encoded WebP data
//   - m: The image.Image to encode
//
// Returns:
//   - error: Any error encountered during encoding
func (e *Encoder) Encode(w io.Writer, m image.Image) error {
	return e.EncodeWithContext(context.Background(), w, m)
}

// EncodeWithContext writes the Image m to w in WebP format with context support.
// The context can be used to cancel the operation.
// Any Image type may be encoded.
//
// Parameters:
//   - ctx: The context for cancellation
//   - w: The io.Writer to write the encoded WebP data
//   - m: The image.Image to encode
//
// Returns:
//   - error: Any error encountered during encoding
func (e *Encoder) EncodeWithContext(ctx context.Context, w io.Writer, m image.Image) error {
	return NewCWebP().
		Quality(e.Quality).
		InputImage(m).
		Output(w).
		RunWithContext(ctx)
}

// Encode writes the Image m to w in WebP format using default settings.
// It is a convenience function that creates an Encoder with default quality (75).
// Any Image type may be encoded.
//
// Parameters:
//   - w: The io.Writer to write the encoded WebP data
//   - m: The image.Image to encode
//
// Returns:
//   - error: Any error encountered during encoding
func Encode(w io.Writer, m image.Image) error {
	return EncodeWithContext(context.Background(), w, m)
}

// EncodeWithContext writes the Image m to w in WebP format using default settings and context support.
// The context can be used to cancel the operation.
// It is a convenience function that creates an Encoder with default quality (75).
// Any Image type may be encoded.
//
// Parameters:
//   - ctx: The context for cancellation
//   - w: The io.Writer to write the encoded WebP data
//   - m: The image.Image to encode
//
// Returns:
//   - error: Any error encountered during encoding
func EncodeWithContext(ctx context.Context, w io.Writer, m image.Image) error {
	e := &Encoder{Quality: 75}
	return e.EncodeWithContext(ctx, w, m)
}
