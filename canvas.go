// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package tinyskia

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"

	color2 "github.com/lumifloat/tinyskia/internal/core/color"
)

type Canvas struct {
	im      *image.RGBA
	context *Context
}

func NewCanvas(width, height int) *Canvas {
	c := &Canvas{im: image.NewRGBA(image.Rect(0, 0, width, height))}
	c.context = &Context{
		canvas: c,

		fillStyle:   &SolidColor{color: color.RGBA{0, 0, 0, 255}},
		strokeStyle: &SolidColor{color: color.RGBA{0, 0, 0, 255}},

		matrix: NewMatrixIdentity(),

		composite: CompositeOperationSourceOver,

		lineWidth:      1,
		lineCap:        LineCapButt,
		lineJoin:       LineJoinMiter,
		miterLimit:     10,
		lineDash:       []float64{},
		lineDashOffset: 0,

		path2d: NewPath2D(),

		font: FontAttr{
			Family: []string{string(FontGenericSansSerif)},
			Weight: FontWeightNormal,
			Style:  FontStyleNormal,
			Size:   10,
		},
		textAlign:   TextAlignStart,
		direction:   DirectionLTR,
		fontStretch: FontStretchNormal,
		fontVariant: FontVariantNormal,
		fontKerning: FontKerningAuto,

		antiAlias:       true,
		colorspace:      color2.ColorSpaceLinear,
		forceHQPipeline: true,
	}

	return c
}

func (c *Canvas) GetContext() *Context {
	return c.context
}

func (c *Canvas) SetWidth(width int) {
	c = NewCanvas(width, c.im.Rect.Dy())
}

func (c *Canvas) GetWidth() int {
	return c.im.Rect.Dx()
}

func (c *Canvas) SetHeight(height int) {
	c = NewCanvas(c.im.Rect.Dx(), height)
}

func (c *Canvas) GetHeight() int {
	return c.im.Rect.Dy()
}

// PngConfig contains configuration options for PNG encoding
type PngConfig struct {
	CompressionLevel png.CompressionLevel
}

// JpegConfig contains configuration options for JPEG encoding
type JpegConfig struct {
	Quality int // Quality ranges from 1-100
}

// ToBufferPNG encodes the canvas as a PNG buffer
func (c *Canvas) ToBufferPNG(config *PngConfig) ([]byte, error) {
	var buf bytes.Buffer
	err := c.WritePNG(&buf, config)
	return buf.Bytes(), err
}

// ToBufferJPEG encodes the canvas as a JPEG buffer
func (c *Canvas) ToBufferJPEG(config *JpegConfig) ([]byte, error) {
	var buf bytes.Buffer
	err := c.WriteJPEG(&buf, config)
	return buf.Bytes(), err
}

// ToBufferRaw returns the unencoded pixel data, top-to-bottom in BGRA format on little-endian systems
func (c *Canvas) ToBufferRaw() []byte {
	pix := make([]byte, len(c.im.Pix))
	copy(pix, c.im.Pix)
	return pix
}

// Image returns the canvas as an image.Image
func (c *Canvas) Image() image.Image {
	im := image.NewRGBA(c.im.Rect)
	copy(im.Pix, c.im.Pix)
	return im
}

// WritePNG writes the canvas as a PNG file to the given writer
func (c *Canvas) WritePNG(w io.Writer, config *PngConfig) error {
	if config != nil {
		encoder := &png.Encoder{
			CompressionLevel: config.CompressionLevel,
		}
		return encoder.Encode(w, c.im)
	} else {
		return png.Encode(w, c.im)
	}
}

// WriteJPEG writes the canvas as a JPEG file to the given writer
func (c *Canvas) WriteJPEG(w io.Writer, config *JpegConfig) error {
	if config != nil && config.Quality > 0 {
		opts := &jpeg.Options{
			Quality: config.Quality,
		}
		return jpeg.Encode(w, c.im, opts)
	} else {
		return jpeg.Encode(w, c.im, nil)
	}
}
