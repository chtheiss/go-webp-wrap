// Package webpbin provides a Go wrapper for the WebP image compression tools.
// It allows for easy conversion of images to WebP format with various options
// including quality control, cropping, and different input/output methods.
package webpbin

import (
	"context"
	"fmt"
	"image"
	"io"
)

// Decode reads a WebP image from r and returns it as an image.Image.
// It is a convenience function that wraps the DWebP decoder.
//
// Parameters:
//   - r: The io.Reader containing the WebP image data
//
// Returns:
//   - image.Image: The decoded image
//   - error: Any error encountered during decoding
func Decode(r io.Reader) (image.Image, error) {
	return DecodeWithContext(context.Background(), r)
}

// DecodeWithContext reads a WebP image from r and returns it as an image.Image.
// The context can be used to cancel the operation.
// It is a convenience function that wraps the DWebP decoder.
//
// Parameters:
//   - ctx: The context for cancellation
//   - r: The io.Reader containing the WebP image data
//
// Returns:
//   - image.Image: The decoded image
//   - error: Any error encountered during decoding
func DecodeWithContext(ctx context.Context, r io.Reader) (image.Image, error) {
	img, err := NewDWebP().Input(r).RunWithContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to decode WebP image: %w", err)
	}
	return img, nil
}
