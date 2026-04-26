// Copyright 2016 Michael Fogleman
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package tinyskia

import (
	"image"
	"image/png"
	"io"

	color2 "github.com/lumifloat/tinyskia/internal/core/color"
	"github.com/lumifloat/tinyskia/internal/path"
)

type FillRule int

const (
	FillRuleWinding FillRule = iota
	FillRuleEvenOdd
)

// Context is the main drawing context, similar to gg.Context.
// It maintains drawing state and provides a canvas-like API.
type Context struct {
	width  int
	height int
	im     *image.RGBA
	mask   *image.Alpha

	// FillStrokeStyles
	fillStyle   style
	strokeStyle style

	// Transform
	matrix *matrix

	// Composite
	composite CompositeOperation

	// PathDrawingStyles
	lineWidth      float64
	lineCap        LineCap
	lineJoin       LineJoin
	miterLimit     float64
	lineDash       []float64
	lineDashOffset float64

	// Path
	path2d *path2d

	fillRule FillRule

	// TextDrawingStyles
	lang        string
	font        *Font
	textAlign   TextAlign
	fontKerning FontKerning

	antiAlias       bool
	colorspace      color2.ColorSpace
	forceHQPipeline bool
	stack           []*Context
	contextLost     bool // Tracks if the rendering context was lost
}

func NewContext(width, height int) *Context {
	return NewContextForRGBA(image.NewRGBA(image.Rect(0, 0, width, height)))
}

// NewContextForImage creates a context from an existing image.Image.
// No copy is made.
func NewContextForImage(im image.Image) *Context {
	return NewContextForRGBA(imageToRGBA(im))
}

func imageToRGBA(im image.Image) *image.RGBA {
	bounds := im.Bounds()
	rgba := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rgba.Set(x, y, im.At(x, y))
		}
	}
	return rgba
}

func NewContextForRGBA(im *image.RGBA) *Context {
	bounds := im.Bounds()
	return &Context{
		path2d: NewPath2D(),

		width:           bounds.Dx(),
		height:          bounds.Dy(),
		im:              im,
		lineWidth:       1,
		lineCap:         LineCapRound,
		lineJoin:        LineJoinMiter,
		fillRule:        FillRuleWinding,
		matrix:          NewMatrixIdentity(),
		composite:       CompositeOperationSourceOver,
		antiAlias:       true,
		colorspace:      color2.ColorSpaceLinear,
		forceHQPipeline: true,
	}
}

// Image returns the image that has been drawn by this context.
func (dc *Context) Image() image.Image {
	return dc.im
}

// Width returns the width of the image in pixels.
func (dc *Context) Width() int {
	return dc.width
}

// Height returns the height of the image in pixels.
func (dc *Context) Height() int {
	return dc.height
}

// SavePNG encodes the image as a PNG and writes it to disk.
func (dc *Context) SavePNG(path string) error {
	return SavePNG(path, dc.Image())
}

// EncodePNG encodes the image as a PNG and writes it to the provided io.Writer.
func (dc *Context) EncodePNG(w io.Writer) error {
	return png.Encode(w, dc.Image())
}

// BeginPath resets the current default path.
func (dc *Context) BeginPath() {
	dc.path2d.builder = path.NewPathBuilder()
}

// SetAntiAlias enables or disables anti-aliasing.
func (dc *Context) SetAntiAlias(aa bool) {
	dc.antiAlias = aa
}

func (dc *Context) SetForceHQPipeline(force bool) {
	dc.forceHQPipeline = force
}
