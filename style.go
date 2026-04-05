// Copyright 2016 Michael Fogleman
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package tinyskia

import (
	"image"
	"image/color"
	"sort"

	color2 "github.com/lumifloat/tinyskia/internal/core/color"
	"github.com/lumifloat/tinyskia/internal/core/shader"
	"github.com/lumifloat/tinyskia/internal/path"
)

type Style interface {
	style()
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

type Gradient interface {
	Style
	AddColorStop(offset float64, color color.Color)
}

type RepeatOp int

const (
	RepeatBoth RepeatOp = iota
	RepeatX
	RepeatY
	RepeatNone
)

type linearGradient struct {
	x0, y0, x1, y1 float64
	stops          stops
}

func NewLinearGradient(x0, y0, x1, y1 float64) Gradient {
	g := &linearGradient{
		x0: x0, y0: y0,
		x1: x1, y1: y1,
	}
	return g
}

func (g *linearGradient) style() {}

func (g *linearGradient) AddColorStop(offset float64, color color.Color) {
	g.stops = append(g.stops, stop{pos: offset, color: color})
	sort.Sort(g.stops)
}

type circle struct {
	x, y, r float64
}

type radialGradient struct {
	c0, c1, cd circle
	stops      stops
}

func NewRadialGradient(x0, y0, r0, x1, y1, r1 float64) Gradient {
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

func (g *radialGradient) AddColorStop(offset float64, color color.Color) {
	g.stops = append(g.stops, stop{pos: offset, color: color})
	sort.Sort(g.stops)
}

type conicGradient struct {
	cx, cy float64
	deg    float64
	stops  stops
}

func NewConicGradient(cx, cy, deg float64) Gradient {
	g := &conicGradient{
		cx:  cx,
		cy:  cy,
		deg: deg,
	}
	return g
}

func (g *conicGradient) style() {}

func (g *conicGradient) AddColorStop(offset float64, color color.Color) {
	g.stops = append(g.stops, stop{pos: offset, color: color})
	sort.Sort(g.stops)
}

type solidPattern struct {
	color color.Color
}

func NewSolidPattern(c color.Color) Style {
	return &solidPattern{color: c}
}

func (p *solidPattern) style() {}

type surfacePattern struct {
	im image.Image
	op RepeatOp
}

func NewSurfacePattern(im image.Image, op RepeatOp) Style {
	return &surfacePattern{im: im, op: op}
}

func (p *surfacePattern) style() {}

func toShader(style Style, transform path.Transform) shader.Shader {
	switch s := style.(type) {
	case *linearGradient:
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
	case *radialGradient:
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
	case *conicGradient:
		stops := make([]shader.GradientStop, len(s.stops))
		for i, s := range s.stops {
			r, gb, b, a := s.color.RGBA()
			stops[i] = shader.NewGradientStop(
				float32(s.pos),
				color2.ColorFromRGBA8(uint8(r>>8), uint8(gb>>8), uint8(b>>8), uint8(a>>8)),
			)
		}

		center := path.Point{X: float32(s.cx), Y: float32(s.cy)}
		startAngle := float32(s.deg)
		endAngle := startAngle + 360.0
		return shader.NewSweepGradient(center, startAngle, endAngle, stops, shader.SpreadModePad, transform)
	case *solidPattern:
		r, g, b, a := s.color.RGBA()
		return shader.NewSolidColor(color2.ColorFromRGBA8(uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a>>8)))
	case *surfacePattern:
		return imageToPatternShader(s.im, s.op, transform)
	default:
		return shader.NewSolidColor(color2.ColorFromRGBA8(0, 0, 0, 255))
	}
}

func imageToPatternShader(im image.Image, op RepeatOp, transform path.Transform) shader.Shader {
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
	case RepeatBoth:
		spreadMode = shader.SpreadModeRepeat
	case RepeatX:
		// tinyskia 不支持单向重复，使用 Repeat 作为近似
		spreadMode = shader.SpreadModeRepeat
	case RepeatY:
		// tinyskia 不支持单向重复，使用 Repeat 作为近似
		spreadMode = shader.SpreadModeRepeat
	case RepeatNone:
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
