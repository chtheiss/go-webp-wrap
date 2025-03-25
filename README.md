# WebP Encoder/Decoder for Golang

[![](https://img.shields.io/badge/docs-godoc-blue.svg)](https://godoc.org/github.com/chtheiss/go-webp-wrap)
[![](https://circleci.com/gh/chtheiss/go-webp-wrap.png?circle-token=ebaa6a739ac4dc96dcb167e0700dcc699409f672)](https://circleci.com/gh/chtheiss/go-webp-wrap)

WebP Encoder/Decoder for Golang based on official libwebp distribution

## Install

```go get -u github.com/chtheiss/go-webp-wrap```



## Example of usage

```go
package main

import (
	"image"
	"image/color"
	"log"
	"os"
	"github.com/chtheiss/go-webp-wrap"
)

func main() {
	const width, height = 256, 256

	// Create a colored image of the given width and height.
	img := image.NewNRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.NRGBA{
				R: uint8((x + y) & 255),
				G: uint8((x + y) << 1 & 255),
				B: uint8((x + y) << 2 & 255),
				A: 255,
			})
		}
	}

	f, err := os.Create("image.webp")
	if err != nil {
		log.Fatal(err)
	}

	if err := webpbin.Encode(f, img); err != nil {
		f.Close()
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}
```

## CWebP

CWebP is a wrapper for *cwebp* command line tool.

Example to convert image.png to image.webp:

```go
err := webpbin.NewCWebP().
		Quality(80).
		InputFile("image.png").
		OutputFile("image.webp").
		Run()
```

## DWebP

DWebP is a wrapper for *dwebp* command line tool.

Example to convert image.webp to image.png:

```go
err := webpbin.NewDWebP().
		InputFile("image.webp").
		OutputFile("image.png").
		Run()
```
