// Package webpbin provides a Go wrapper for the WebP image compression tools.
// It allows for easy conversion of images to WebP format with various options
// including quality control, cropping, and different input/output methods.
package webpbin

import (
	"errors"
	"fmt"
	"image"
	"io"

	"github.com/belphemur/go-binwrapper"
)

// cropInfo represents the cropping parameters for an image.
type cropInfo struct {
	x      int // x-coordinate of the top-left corner
	y      int // y-coordinate of the top-left corner
	width  int // width of the crop area
	height int // height of the crop area
}

// CWebP wraps the cwebp command-line tool for compressing images to WebP format.
// It supports various input formats including PNG, JPEG, TIFF, WebP, and raw Y'CbCr samples.
// For more information, see: https://developers.google.com/speed/webp/docs/cwebp
type CWebP struct {
	*binwrapper.BinWrapper
	inputFile  string      // Path to the input image file
	inputImage image.Image // Input image as Go image.Image
	input      io.Reader   // Input as io.Reader
	outputFile string      // Path to the output WebP file
	output     io.Writer   // Output as io.Writer
	quality    int         // Compression quality (0-100)
	crop       *cropInfo   // Cropping parameters
}

// NewCWebP creates a new CWebP instance with the given options.
// It initializes the binary wrapper and sets default values.
// The quality is set to -1 by default, which means the default cwebp quality will be used.
func NewCWebP(optionFuncs ...OptionFunc) *CWebP {
	bin := &CWebP{
		BinWrapper: createBinWrapper(optionFuncs...),
		quality:    -1,
	}
	bin.ExecPath("cwebp")

	return bin
}

// Version returns the version of the cwebp binary.
// Returns the version string and any error encountered.
func (c *CWebP) Version() (string, error) {
	return version(c.BinWrapper)
}

// InputFile sets the input image file to convert.
// Any previous calls to Input or InputImage will be ignored.
// Returns the CWebP instance for method chaining.
func (c *CWebP) InputFile(file string) *CWebP {
	c.input = nil
	c.inputImage = nil
	c.inputFile = file
	return c
}

// Input sets the reader to convert.
// Any previous calls to InputFile or InputImage will be ignored.
// Returns the CWebP instance for method chaining.
func (c *CWebP) Input(reader io.Reader) *CWebP {
	c.inputFile = ""
	c.inputImage = nil
	c.input = reader
	return c
}

// InputImage sets the image to convert.
// Any previous calls to InputFile or Input will be ignored.
// Returns the CWebP instance for method chaining.
func (c *CWebP) InputImage(img image.Image) *CWebP {
	c.inputFile = ""
	c.input = nil
	c.inputImage = img
	return c
}

// OutputFile specifies the name of the output WebP file.
// Any previous call to Output will be ignored.
// Returns the CWebP instance for method chaining.
func (c *CWebP) OutputFile(file string) *CWebP {
	c.output = nil
	c.outputFile = file
	return c
}

// Output specifies the writer to write WebP file content.
// Any previous call to OutputFile will be ignored.
// Returns the CWebP instance for method chaining.
func (c *CWebP) Output(writer io.Writer) *CWebP {
	c.outputFile = ""
	c.output = writer
	return c
}

// Quality specifies the compression factor for RGB channels.
// The value must be between 0 and 100, where:
// - A small factor produces a smaller file with lower quality
// - A value of 100 achieves the best quality
// - The default is 75
// Returns the CWebP instance for method chaining.
func (c *CWebP) Quality(quality uint) *CWebP {
	if quality > 100 {
		quality = 100
	}

	c.quality = int(quality)
	return c
}

// Crop sets the cropping parameters for the source image.
// The cropping area must be fully contained within the source rectangle.
// Parameters:
//   - x: x-coordinate of the top-left corner
//   - y: y-coordinate of the top-left corner
//   - width: width of the crop area
//   - height: height of the crop area
//
// Returns the CWebP instance for method chaining.
func (c *CWebP) Crop(x, y, width, height int) *CWebP {
	c.crop = &cropInfo{x, y, width, height}
	return c
}

// Run executes the cwebp command with the specified parameters.
// Returns an error if the command fails or if input/output is not properly configured.
func (c *CWebP) Run() error {
	defer c.BinWrapper.Reset()

	if c.quality > -1 {
		c.Arg("-q", fmt.Sprintf("%d", c.quality))
	}

	if c.crop != nil {
		c.Arg("-crop", fmt.Sprintf("%d", c.crop.x), fmt.Sprintf("%d", c.crop.y), fmt.Sprintf("%d", c.crop.width), fmt.Sprintf("%d", c.crop.height))
	}

	output, err := c.getOutput()

	if err != nil {
		return err
	}

	c.Arg("-o", output)

	err = c.setInput()

	if err != nil {
		return err
	}

	if c.output != nil {
		c.SetStdOut(c.output)
	}

	err = c.BinWrapper.Run()

	if err != nil {
		return errors.New(err.Error() + ". " + string(c.StdErr()))
	}

	return nil
}

// Reset restores all parameters to their default values.
// Returns the CWebP instance for method chaining.
func (c *CWebP) Reset() *CWebP {
	c.crop = nil
	c.quality = -1
	return c
}

// setInput configures the input source for the cwebp command.
// Returns an error if no input source is defined.
func (c *CWebP) setInput() error {
	if c.input != nil {
		c.Arg("--").Arg("-")
		c.StdIn(c.input)
	} else if c.inputImage != nil {
		r, err := createReaderFromImage(c.inputImage)

		if err != nil {
			return err
		}

		c.Arg("--").Arg("-")
		c.StdIn(r)
	} else if c.inputFile != "" {
		c.Arg(c.inputFile)
	} else {
		return errors.New("Undefined input")
	}

	return nil
}

// getOutput determines the output destination for the cwebp command.
// Returns the output path and an error if no output destination is defined.
func (c *CWebP) getOutput() (string, error) {
	if c.output != nil {
		return "-", nil
	} else if c.outputFile != "" {
		return c.outputFile, nil
	} else {
		return "", errors.New("Undefined output")
	}
}
