// Package webpbin provides a Go wrapper for the WebP image compression tools.
// It allows for easy conversion of images to WebP format with various options
// including quality control, cropping, and different input/output methods.
package webpbin

import (
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
	return NewDWebP().Input(r).Run()
}
