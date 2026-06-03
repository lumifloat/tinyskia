// Copyright 2006 The Android Open Source Project
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package tinyskia

import (
	"image"
	"image/color"
	"math"

	color2 "github.com/lumifloat/tinyskia/internal/core/color"
	"github.com/lumifloat/tinyskia/internal/core/shader"
	"github.com/lumifloat/tinyskia/internal/path"
)

type FillRule int

const (
	FillRuleWinding FillRule = iota
	FillRuleEvenOdd
)

// SetFillStyle to change the fill style.
func (dc *Context) SetFillStyle(style style) {
	dc.fillStyle = style
}

// SetFillStyleSolidColor to change the fill style to a solid color.
func (dc *Context) SetFillStyleSolidColor(c color.Color) {
	dc.fillStyle = dc.CreateSolidColor(c)
}

// SetFillStylePattern to change the fill style to a pattern.
func (dc *Context) SetFillStylePattern(p Pattern) {
	dc.fillStyle = p
}

// SetFillStyleGradient to change the fill style to a gradient.
func (dc *Context) SetFillStyleGradient(g Gradient) {
	dc.fillStyle = g
}

// GetFillStyle returns the current fill style.
func (dc *Context) GetFillStyle() style {
	return dc.fillStyle
}

// GetFillStyleSolidColor returns the current fill style solid color.
func (dc *Context) GetFillStyleSolidColor() color.Color {
	if solid, ok := dc.fillStyle.(*SolidColor); ok {
		return solid.color
	}
	return nil
}

// GetFillStylePattern returns the current fill style pattern.
func (dc *Context) GetFillStylePattern() Pattern {
	if pattern, ok := dc.fillStyle.(Pattern); ok {
		return pattern
	}
	return nil
}

// GetFillStyleGradient returns the current fill style gradient.
func (dc *Context) GetFillStyleGradient() Gradient {
	if gradient, ok := dc.fillStyle.(Gradient); ok {
		return gradient
	}
	return nil
}

// SetStrokeStyle to change the stroke style.
func (dc *Context) SetStrokeStyle(style style) {
	dc.strokeStyle = style
}

// SetStrokeStyleSolidColor to change the stroke style to a solid color.
func (dc *Context) SetStrokeStyleSolidColor(c color.Color) {
	dc.strokeStyle = dc.CreateSolidColor(c)
}

// SetStrokeStylePattern to change the stroke style to a pattern.
func (dc *Context) SetStrokeStylePattern(p Pattern) {
	dc.strokeStyle = p
}

// SetStrokeStyleGradient to change the stroke style to a gradient.
func (dc *Context) SetStrokeStyleGradient(g Gradient) {
	dc.strokeStyle = g
}

// GetStrokeStyle returns the current stroke style.
func (dc *Context) GetStrokeStyle() style {
	return dc.strokeStyle
}

// GetStrokeStyleSolidColor returns the current stroke style solid color.
func (dc *Context) GetStrokeStyleSolidColor() color.Color {
	if solid, ok := dc.strokeStyle.(*SolidColor); ok {
		return solid.color
	}
	return nil
}

// GetStrokeStylePattern returns the current stroke style pattern.
func (dc *Context) GetStrokeStylePattern() Pattern {
	if pattern, ok := dc.strokeStyle.(Pattern); ok {
		return pattern
	}
	return nil
}

// GetStrokeStyleGradient returns the current stroke style gradient.
func (dc *Context) GetStrokeStyleGradient() Gradient {
	if gradient, ok := dc.strokeStyle.(Gradient); ok {
		return gradient
	}
	return nil
}

type SolidColor struct {
	color color.Color
}

// CreateSolidColor returns a solid color style.
func (dc *Context) CreateSolidColor(c color.Color) *SolidColor {
	return &SolidColor{color: c}
}

func (p *SolidColor) style() {}

type PatternRepeatOp int

const (
	PatternRepeatBoth PatternRepeatOp = iota
	PatternRepeatX
	PatternRepeatY
	PatternRepeatNone
)

type ImagePattern struct {
	im        image.Image
	op        PatternRepeatOp
	transform *Matrix
}

// CreatePattern returns a CanvasPattern object that uses the given image and
// repeats in the direction(s) given by the repetition argument.
// The allowed values for repetition are repeat (both directions),
// repeat-x (horizontal only), repeat-y (vertical only), and no-repeat (neither).
// If the repetition argument is empty, the value repeat is used.
func (dc *Context) CreatePattern(im image.Image, op PatternRepeatOp) Pattern {
	return &ImagePattern{im: im, op: op}
}

func (p *ImagePattern) style() {}

type stop struct {
	pos   float64
	color color.Color
}

type stops []stop

// Len satisfies the Sort interface.
func (s stops) Len() int {
	return len(s)
}

// Less satisfies the Sort interface.
func (s stops) Less(i, j int) bool {
	return s[i].pos < s[j].pos
}

// Swap satisfies the Sort interface.
func (s stops) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type LinearGradient struct {
	x0, y0, x1, y1 float64
	stops          stops
}

// CreateLinearGradient returns a CanvasGradient object that represents
// a linear gradient that paints along the line given by the coordinates represented by the arguments.
func (dc *Context) CreateLinearGradient(x0, y0, x1, y1 float64) Gradient {
	g := &LinearGradient{
		x0: x0, y0: y0,
		x1: x1, y1: y1,
	}
	return g
}

func (g *LinearGradient) style() {}

type circle struct {
	x, y, r float64
}

type RadialGradient struct {
	c0, c1, cd circle
	stops      stops
}

// CreateRadialGradient returns a CanvasGradient object that represents
// a radial gradient that paints along the cone given by the circles represented by the arguments.
func (dc *Context) CreateRadialGradient(x0, y0, r0, x1, y1, r1 float64) Gradient {
	c0 := circle{x0, y0, r0}
	c1 := circle{x1, y1, r1}
	cd := circle{x1 - x0, y1 - y0, r1 - r0}
	g := &RadialGradient{
		c0: c0,
		c1: c1,
		cd: cd,
	}
	return g
}

func (g *RadialGradient) style() {}

type ConicGradient struct {
	x, y       float64
	startAngle float64
	stops      stops
}

// CreateConicGradient returns a CanvasGradient object that represents
// a conic gradient that paints clockwise along the rotation around the center represented by the arguments.
func (dc *Context) CreateConicGradient(startAngle, x, y float64) Gradient {
	g := &ConicGradient{
		x: x, y: y,
		startAngle: startAngle,
	}
	return g
}

func (g *ConicGradient) style() {}

type style interface {
	style()
}

func toShader(style style, transform path.Transform) shader.Shader {
	switch s := style.(type) {
	case *LinearGradient:
		stops := make([]shader.GradientStop, len(s.stops))
		for i, s := range s.stops {
			r, gb, b, a := s.color.RGBA()
			stops[i] = shader.NewGradientStop(
				float32(s.pos),
				color2.ColorFromRGBA8(uint8(r>>8), uint8(gb>>8), uint8(b>>8), uint8(a>>8)),
			)
		}

		p0 := path.Point{X: float32(s.x0), Y: float32(s.y0)}
		p1 := path.Point{X: float32(s.x1), Y: float32(s.y1)}
		return shader.NewLinearGradient(p0, p1, stops, shader.SpreadModePad, transform)
	case *RadialGradient:
		stops := make([]shader.GradientStop, len(s.stops))
		for i, s := range s.stops {
			r, gb, b, a := s.color.RGBA()
			stops[i] = shader.NewGradientStop(
				float32(s.pos),
				color2.ColorFromRGBA8(uint8(r>>8), uint8(gb>>8), uint8(b>>8), uint8(a>>8)),
			)
		}

		center0 := path.Point{X: float32(s.c0.x), Y: float32(s.c0.y)}
		center1 := path.Point{X: float32(s.c1.x), Y: float32(s.c1.y)}
		return shader.NewRadialGradient(center0, float32(s.c0.r), center1, float32(s.c1.r), stops, shader.SpreadModePad, transform)
	case *ConicGradient:
		stops := make([]shader.GradientStop, len(s.stops))
		for i, s := range s.stops {
			r, gb, b, a := s.color.RGBA()
			stops[i] = shader.NewGradientStop(
				float32(s.pos),
				color2.ColorFromRGBA8(uint8(r>>8), uint8(gb>>8), uint8(b>>8), uint8(a>>8)),
			)
		}

		center := path.Point{X: float32(s.x), Y: float32(s.y)}
		startAngle := float32(s.startAngle * 180.0 / math.Pi)
		transform = transform.PreRotateAt(startAngle, center.X, center.Y)
		return shader.NewSweepGradient(center, 0, 360, stops, shader.SpreadModePad, transform)
	case *SolidColor:
		r, g, b, a := s.color.RGBA()
		return shader.NewSolidColor(color2.ColorFromRGBA8(uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a>>8)))
	case *ImagePattern:
		var transform path.Transform
		if s.transform != nil {
			transform = s.transform.transform
		} else {
			transform = path.NewTransformDefault()
		}
		return imageToPatternShader(s.im, s.op, transform)
	default:
		return shader.NewSolidColor(color2.ColorFromRGBA8(0, 0, 0, 255))
	}
}

func imageToPatternShader(im image.Image, op PatternRepeatOp, transform path.Transform) shader.Shader {
	bounds := im.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	if width <= 0 || height <= 0 {
		return shader.NewSolidColor(color2.ColorFromRGBA8(0, 0, 0, 0))
	}

	data := make([]uint8, width*height*4)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c := im.At(bounds.Min.X+x, bounds.Min.Y+y)
			r, g, b, a := c.RGBA()
			offset := (y*width + x) * 4
			data[offset+0] = uint8(r >> 8)
			data[offset+1] = uint8(g >> 8)
			data[offset+2] = uint8(b >> 8)
			data[offset+3] = uint8(a >> 8)
		}
	}

	size, _ := path.NewIntSize(uint32(width), uint32(height))

	var spreadMode shader.SpreadMode
	switch op {
	case PatternRepeatBoth:
		spreadMode = shader.SpreadModeRepeat
	case PatternRepeatX:
		// tinyskia 不支持单向重复，使用 Repeat 作为近似
		spreadMode = shader.SpreadModeRepeat
	case PatternRepeatY:
		// tinyskia 不支持单向重复，使用 Repeat 作为近似
		spreadMode = shader.SpreadModeRepeat
	case PatternRepeatNone:
		spreadMode = shader.SpreadModePad
	default:
		spreadMode = shader.SpreadModeRepeat
	}

	return shader.NewPattern(
		data,
		size,
		spreadMode,
		shader.FilterQualityBilinear,
		1.0,
		transform,
	)
}
