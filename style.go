// Copyright 2016 Michael Fogleman
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package tinyskia

import (
	"image"

	color2 "github.com/lumifloat/tinyskia/internal/core/color"
	"github.com/lumifloat/tinyskia/internal/core/shader"
	"github.com/lumifloat/tinyskia/internal/path"
)

type style interface {
	style()
}

func toShader(style style, transform path.Transform) shader.Shader {
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

		center := path.Point{X: float32(s.x), Y: float32(s.y)}
		startAngle := float32(s.startAngle)
		endAngle := startAngle + 360.0
		return shader.NewSweepGradient(center, startAngle, endAngle, stops, shader.SpreadModePad, transform)
	case *solidColor:
		r, g, b, a := s.color.RGBA()
		return shader.NewSolidColor(color2.ColorFromRGBA8(uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a>>8)))
	case *imagePattern:
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
