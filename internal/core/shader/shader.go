// Copyright 2006 The Android Open Source Project
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package shader

import (
	"image/color"

	"github.com/lumifloat/tinyskia/internal/core/colorf"
	"github.com/lumifloat/tinyskia/internal/core/pipeline"
	"github.com/lumifloat/tinyskia/internal/path"
)

// A shader spreading mode.
type SpreadMode int

const (
	// Replicate the edge color if the shader draws outside of its
	// original bounds.
	SpreadModePad SpreadMode = iota
	// Repeat the shader's image horizontally and vertically, alternating
	// mirror images so that adjacent images always seam.
	SpreadModeReflect
	// Repeat the shader's image horizontally and vertically.
	SpreadModeRepeat
)

type Shader interface {
	IsOpaque() bool
	PushStages(cs colorf.ColorSpace, p *pipeline.RasterPipelineBuilder) bool
	Transform(ts path.Transform)
	// ApplyOpacity(opacity float32)
}

// SolidColor a solid color shader.
type SolidColor struct {
	color color.Color
}

func NewSolidColor(c color.Color) *SolidColor {
	return &SolidColor{color: c}
}

func (sc *SolidColor) Color() color.Color {
	return sc.color
}

func (sc *SolidColor) IsOpaque() bool {
	return colorf.IsOpaque(sc.color)
}

func (sc *SolidColor) PushStages(cs colorf.ColorSpace, p *pipeline.RasterPipelineBuilder) bool {
	expand := cs.ExpandColor(sc.color)
	c0 := color.RGBA64Model.Convert(expand).(color.RGBA64)
	c1 := colorf.RGBAFModel.Convert(expand).(colorf.RGBAF)
	p.PushUniformColor(pipeline.UniformColorCtx{
		R0: c0.R, G0: c0.G, B0: c0.B, A0: c0.A,
		R1: float32(c1.R), G1: float32(c1.G), B1: float32(c1.B), A1: float32(c1.A),
	})
	return true
}

func (sc *SolidColor) Transform(ts path.Transform) {
	// Solid color shaders don't need transform
}

func (sc *SolidColor) ApplyOpacity(opacity float32) {
	sc.color = colorf.ApplyOpacity(sc.color, opacity)
}
