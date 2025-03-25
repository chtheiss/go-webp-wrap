// Package webpbin provides a Go wrapper for the WebP image compression tools.
// It allows for easy conversion of images to WebP format with various options
// including quality control, cropping, and different input/output methods.
package webpbin

import (
	"bytes"
	"errors"
	"image"
	"image/png"
	"io"

	"github.com/belphemur/go-binwrapper"
)

// DWebP wraps the dwebp command-line tool for decompressing WebP files into PNG format.
// It provides various options for input/output handling and supports both file and stream-based operations.
// For more information, see: https://developers.google.com/speed/webp/docs/dwebp
type DWebP struct {
	*binwrapper.BinWrapper
	inputFile  string    // Path to the input WebP file
	input      io.Reader // Input as io.Reader
	outputFile string    // Path to the output PNG file
	output     io.Writer // Output as io.Writer
}

// NewDWebP creates a new DWebP instance with the given options.
// It initializes the binary wrapper and sets up the dwebp executable.
func NewDWebP(optionFuncs ...OptionFunc) *DWebP {
	bin := &DWebP{
		BinWrapper: createBinWrapper(optionFuncs...),
	}
	bin.ExecPath("dwebp")

	return bin
}

// InputFile sets the WebP file to convert.
// Any previous calls to Input will be ignored.
// Returns the DWebP instance for method chaining.
func (c *DWebP) InputFile(file string) *DWebP {
	c.input = nil
	c.inputFile = file
	return c
}

// Input sets the reader to convert.
// Any previous calls to InputFile will be ignored.
// Returns the DWebP instance for method chaining.
func (c *DWebP) Input(reader io.Reader) *DWebP {
	c.inputFile = ""
	c.input = reader
	return c
}

// OutputFile specifies the name of the output PNG file.
// Any previous call to Output will be ignored.
// Returns the DWebP instance for method chaining.
func (c *DWebP) OutputFile(file string) *DWebP {
	c.output = nil
	c.outputFile = file
	return c
}

// Output specifies the writer to write PNG file content.
// Any previous call to OutputFile will be ignored.
// Returns the DWebP instance for method chaining.
func (c *DWebP) Output(writer io.Writer) *DWebP {
	c.outputFile = ""
	c.output = writer
	return c
}

// Version returns the version of the dwebp binary.
// Returns the version string and any error encountered.
func (c *DWebP) Version() (string, error) {
	return version(c.BinWrapper)
}

// Run executes the dwebp command with the specified parameters.
// Returns the decoded image and any error encountered during the process.
// If no output is specified, returns the decoded image as an image.Image.
// If an output is specified (file or writer), returns nil, nil.
func (c *DWebP) Run() (image.Image, error) {
	defer c.BinWrapper.Reset()

	output, err := c.getOutput()

	if err != nil {
		return nil, err
	}

	c.Arg("-o", output)

	err = c.setInput()

	if err != nil {
		return nil, err
	}

	if c.output != nil {
		c.SetStdOut(c.output)
	}

	err = c.BinWrapper.Run()

	if err != nil {
		return nil, errors.New(err.Error() + ". " + string(c.StdErr()))
	}

	if c.output == nil && c.outputFile == "" {
		return png.Decode(bytes.NewReader(c.BinWrapper.StdOut()))
	}

	return nil, nil
}

// setInput configures the input source for the dwebp command.
// Returns an error if no input source is defined.
func (c *DWebP) setInput() error {
	if c.input != nil {
		c.Arg("--").Arg("-")
		c.StdIn(c.input)
	} else if c.inputFile != "" {
		c.Arg(c.inputFile)
	} else {
		return errors.New("Undefined input")
	}

	return nil
}

// getOutput determines the output destination for the dwebp command.
// Returns the output path and an error if no output destination is defined.
func (c *DWebP) getOutput() (string, error) {
	if c.outputFile != "" {
		return c.outputFile, nil
	}

	return "-", nil
}
