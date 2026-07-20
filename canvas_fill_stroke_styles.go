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
	"image/draw"
	"math"
	"sort"

	"github.com/lumifloat/tinyskia/internal/core/shader"
	"github.com/lumifloat/tinyskia/internal/path"
)

// SetFillStyle to change the fill style.
func (ctx *Context) SetFillStyle(style Style) {
	if style == nil {
		ctx.fillStyle = &SolidColor{color: color.RGBA{0, 0, 0, 255}}
	}
	ctx.fillStyle = style
}

// SetFillStyleSolidColor to change the fill style to a solid color.
func (ctx *Context) SetFillStyleSolidColor(c color.Color) {
	if s, err := ctx.CreateSolidColor(c); err == nil {
		ctx.fillStyle = s
	}
}

// SetFillStylePattern to change the fill style to a pattern.
func (ctx *Context) SetFillStylePattern(p Pattern) {
	ctx.SetFillStyle(p)
}

// SetFillStyleGradient to change the fill style to a gradient.
func (ctx *Context) SetFillStyleGradient(g Gradient) {
	ctx.SetFillStyle(g)
}

// GetFillStyle returns the current fill style.
func (ctx *Context) GetFillStyle() Style {
	return ctx.fillStyle
}

// GetFillStyleSolidColor returns the current fill style solid color.
func (ctx *Context) GetFillStyleSolidColor() color.Color {
	if solid, ok := ctx.fillStyle.(*SolidColor); ok {
		return solid.color
	}
	return nil
}

// GetFillStylePattern returns the current fill style pattern.
func (ctx *Context) GetFillStylePattern() Pattern {
	if pattern, ok := ctx.fillStyle.(Pattern); ok {
		return pattern
	}
	return nil
}

// GetFillStyleGradient returns the current fill style gradient.
func (ctx *Context) GetFillStyleGradient() Gradient {
	if gradient, ok := ctx.fillStyle.(Gradient); ok {
		return gradient
	}
	return nil
}

// SetStrokeStyle to change the stroke style.
func (ctx *Context) SetStrokeStyle(style Style) {
	if style == nil {
		ctx.strokeStyle = &SolidColor{color: color.RGBA{0, 0, 0, 255}}
	}
	ctx.strokeStyle = style
}

// SetStrokeStyleSolidColor to change the stroke style to a solid color.
func (ctx *Context) SetStrokeStyleSolidColor(c color.Color) {
	if s, err := ctx.CreateSolidColor(c); err == nil {
		ctx.strokeStyle = s
	}
}

// SetStrokeStylePattern to change the stroke style to a pattern.
func (ctx *Context) SetStrokeStylePattern(p Pattern) {
	ctx.SetStrokeStyle(p)
}

// SetStrokeStyleGradient to change the stroke style to a gradient.
func (ctx *Context) SetStrokeStyleGradient(g Gradient) {
	ctx.SetStrokeStyle(g)
}

// GetStrokeStyle returns the current stroke style.
func (ctx *Context) GetStrokeStyle() Style {
	return ctx.strokeStyle
}

// GetStrokeStyleSolidColor returns the current stroke style solid color.
func (ctx *Context) GetStrokeStyleSolidColor() color.Color {
	if solid, ok := ctx.strokeStyle.(*SolidColor); ok {
		return solid.color
	}
	return nil
}

// GetStrokeStylePattern returns the current stroke style pattern.
func (ctx *Context) GetStrokeStylePattern() Pattern {
	if pattern, ok := ctx.strokeStyle.(Pattern); ok {
		return pattern
	}
	return nil
}

// GetStrokeStyleGradient returns the current stroke style gradient.
func (ctx *Context) GetStrokeStyleGradient() Gradient {
	if gradient, ok := ctx.strokeStyle.(Gradient); ok {
		return gradient
	}
	return nil
}

type Style interface {
	style() shader.Shader
}

type SolidColor struct {
	color color.Color
}

// CreateSolidColor returns a solid color style.
func (ctx *Context) CreateSolidColor(c color.Color) (*SolidColor, error) {
	return &SolidColor{color: c}, nil
}

func (p *SolidColor) style() shader.Shader {
	return shader.NewSolidColor(p.color)
}

type RepeatMode string

const (
	RepeatModeRepeat   = "repeat"
	RepeatModeRepeatX  = "repeat-x"
	RepeatModeRepeatY  = "repeat-y"
	RepeatModeNoRepeat = "no-repeat"
)

func (r RepeatMode) convert() shader.SpreadMode {
	switch r {
	case RepeatModeRepeat:
		return shader.SpreadModeRepeat
	case RepeatModeRepeatX:
		return shader.SpreadModeRepeatX
	case RepeatModeRepeatY:
		return shader.SpreadModeRepeatY
	case RepeatModeNoRepeat:
		return shader.SpreadModeNoRepeat
	default:
		return shader.SpreadModeRepeat
	}
}

type ImagePattern struct {
	im        image.Image
	op        RepeatMode
	transform *Matrix
}

// CreatePattern returns a CanvasPattern object that uses the given image and
// repeats in the direction(s) given by the repetition argument.
// The allowed values for repetition are repeat (both directions),
// repeat-x (horizontal only), repeat-y (vertical only), and no-repeat (neither).
// If the repetition argument is empty, the value repeat is used.
func (ctx *Context) CreatePattern(im image.Image, op RepeatMode) (Pattern, error) {
	return &ImagePattern{im: im, op: op, transform: NewMatrixIdentity()}, nil
}

func (p *ImagePattern) style() shader.Shader {
	// TODO
	bounds := p.im.Bounds()
	rect := image.Rect(0, 0, bounds.Dx(), bounds.Dy())
	im := image.NewRGBA(rect)
	draw.Draw(im, rect, p.im, bounds.Min, draw.Src)
	return shader.NewPattern(im, p.op.convert(), shader.FilterQualityBilinear, 1.0, p.transform.transform)
}

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

func (s stops) convert() []shader.GradientStop {
	sort.Sort(s)
	stops := make([]shader.GradientStop, len(s))
	for i, stop := range s {
		stops[i] = shader.NewGradientStop(
			float32(stop.pos), stop.color,
		)
	}
	return stops
}

type LinearGradient struct {
	x0, y0, x1, y1 float64
	stops          stops
}

// CreateLinearGradient returns a CanvasGradient object that represents
// a linear gradient that paints along the line given by the coordinates represented by the arguments.
func (ctx *Context) CreateLinearGradient(x0, y0, x1, y1 float64) (Gradient, error) {
	g := &LinearGradient{
		x0: x0, y0: y0,
		x1: x1, y1: y1,
	}
	return g, nil
}

func (g *LinearGradient) style() shader.Shader {
	p0 := path.Point{X: float32(g.x0), Y: float32(g.y0)}
	p1 := path.Point{X: float32(g.x1), Y: float32(g.y1)}
	s := shader.NewLinearGradient(
		p0, p1, g.stops.convert(), shader.SpreadModePad, path.NewTransformDefault(),
	)
	return s
}

type circle struct {
	x, y, r float64
}

type RadialGradient struct {
	c0, c1, cd circle
	stops      stops
}

// CreateRadialGradient returns a CanvasGradient object that represents
// a radial gradient that paints along the cone given by the circles represented by the arguments.
func (ctx *Context) CreateRadialGradient(x0, y0, r0, x1, y1, r1 float64) (Gradient, error) {
	c0 := circle{x0, y0, r0}
	c1 := circle{x1, y1, r1}
	cd := circle{x1 - x0, y1 - y0, r1 - r0}
	g := &RadialGradient{
		c0: c0,
		c1: c1,
		cd: cd,
	}
	return g, nil
}

func (g *RadialGradient) style() shader.Shader {
	center0 := path.Point{X: float32(g.c0.x), Y: float32(g.c0.y)}
	center1 := path.Point{X: float32(g.c1.x), Y: float32(g.c1.y)}
	s := shader.NewRadialGradient(
		center0, float32(g.c0.r), center1, float32(g.c1.r), g.stops.convert(), shader.SpreadModePad, path.NewTransformDefault(),
	)
	return s
}

type ConicGradient struct {
	x, y       float64
	startAngle float64
	stops      stops
}

// CreateConicGradient returns a CanvasGradient object that represents
// a conic gradient that paints clockwise along the rotation around the center represented by the arguments.
func (ctx *Context) CreateConicGradient(startAngle, x, y float64) (Gradient, error) {
	g := &ConicGradient{
		x: x, y: y,
		startAngle: startAngle,
	}
	return g, nil
}

func (g *ConicGradient) style() shader.Shader {
	center := path.Point{X: float32(g.x), Y: float32(g.y)}
	startAngle := float32(g.startAngle * 180.0 / math.Pi)
	// TODO fix in internal
	tf := path.NewTransformFromRotateAt(startAngle, center.X, center.Y)
	s := shader.NewSweepGradient(
		center, 0, 360, g.stops.convert(), shader.SpreadModePad, tf,
	)
	return s
}
