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
func (dc *Context) SetFillStylePattern(p pattern) {
	dc.fillStyle = p
}

// SetFillStyleGradient to change the fill style to a gradient.
func (dc *Context) SetFillStyleGradient(g gradient) {
	dc.fillStyle = g
}

// GetFillStyle returns the current fill style.
func (dc *Context) GetFillStyle() style {
	return dc.fillStyle
}

// GetFillStyleSolidColor returns the current fill style solid color.
func (dc *Context) GetFillStyleSolidColor() color.Color {
	if solid, ok := dc.fillStyle.(*solidColor); ok {
		return solid.color
	}
	return nil
}

// GetFillStylePattern returns the current fill style pattern.
func (dc *Context) GetFillStylePattern() pattern {
	if pattern, ok := dc.fillStyle.(pattern); ok {
		return pattern
	}
	return nil
}

// GetFillStyleGradient returns the current fill style gradient.
func (dc *Context) GetFillStyleGradient() gradient {
	if gradient, ok := dc.fillStyle.(gradient); ok {
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
func (dc *Context) SetStrokeStylePattern(p pattern) {
	dc.strokeStyle = p
}

// SetStrokeStyleGradient to change the stroke style to a gradient.
func (dc *Context) SetStrokeStyleGradient(g gradient) {
	dc.strokeStyle = g
}

// GetStrokeStyle returns the current stroke style.
func (dc *Context) GetStrokeStyle() style {
	return dc.strokeStyle
}

// GetStrokeStyleSolidColor returns the current stroke style solid color.
func (dc *Context) GetStrokeStyleSolidColor() color.Color {
	if solid, ok := dc.strokeStyle.(*solidColor); ok {
		return solid.color
	}
	return nil
}

// GetStrokeStylePattern returns the current stroke style pattern.
func (dc *Context) GetStrokeStylePattern() pattern {
	if pattern, ok := dc.strokeStyle.(pattern); ok {
		return pattern
	}
	return nil
}

// GetStrokeStyleGradient returns the current stroke style gradient.
func (dc *Context) GetStrokeStyleGradient() gradient {
	if gradient, ok := dc.strokeStyle.(gradient); ok {
		return gradient
	}
	return nil
}

type solidColor struct {
	color color.Color
}

// CreateSolidColor returns a solid color style.
func (dc *Context) CreateSolidColor(c color.Color) *solidColor {
	return &solidColor{color: c}
}

func (p *solidColor) style() {}

type PatternRepeatOp int

const (
	PatternRepeatBoth PatternRepeatOp = iota
	PatternRepeatX
	PatternRepeatY
	PatternRepeatNone
)

type imagePattern struct {
	im        image.Image
	op        PatternRepeatOp
	transform *matrix
}

// CreatePattern returns a CanvasPattern object that uses the given image and
// repeats in the direction(s) given by the repetition argument.
// The allowed values for repetition are repeat (both directions),
// repeat-x (horizontal only), repeat-y (vertical only), and no-repeat (neither).
// If the repetition argument is empty, the value repeat is used.
func (dc *Context) CreatePattern(im image.Image, op PatternRepeatOp) *imagePattern {
	return &imagePattern{im: im, op: op}
}

func (p *imagePattern) style() {}

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

type linearGradient struct {
	x0, y0, x1, y1 float64
	stops          stops
}

// CreateLinearGradient returns a CanvasGradient object that represents
// a linear gradient that paints along the line given by the coordinates represented by the arguments.
func (dc *Context) CreateLinearGradient(x0, y0, x1, y1 float64) gradient {
	g := &linearGradient{
		x0: x0, y0: y0,
		x1: x1, y1: y1,
	}
	return g
}

func (g *linearGradient) style() {}

type circle struct {
	x, y, r float64
}

type radialGradient struct {
	c0, c1, cd circle
	stops      stops
}

// CreateRadialGradient returns a CanvasGradient object that represents
// a radial gradient that paints along the cone given by the circles represented by the arguments.
func (dc *Context) CreateRadialGradient(x0, y0, r0, x1, y1, r1 float64) gradient {
	c0 := circle{x0, y0, r0}
	c1 := circle{x1, y1, r1}
	cd := circle{x1 - x0, y1 - y0, r1 - r0}
	g := &radialGradient{
		c0: c0,
		c1: c1,
		cd: cd,
	}
	return g
}

func (g *radialGradient) style() {}

type conicGradient struct {
	x, y       float64
	startAngle float64
	stops      stops
}

// CreateConicGradient returns a CanvasGradient object that represents
// a conic gradient that paints clockwise along the rotation around the center represented by the arguments.
func (dc *Context) CreateConicGradient(startAngle, x, y float64) gradient {
	g := &conicGradient{
		x: x, y: y,
		startAngle: startAngle,
	}
	return g
}

func (g *conicGradient) style() {}
