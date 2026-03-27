// Copyright 2006 The Android Open Source Project
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package shader

import (
	"github.com/lumifloat/tinyskia/color"
	"github.com/lumifloat/tinyskia/path"
	"github.com/lumifloat/tinyskia/pipeline"
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
	PushStages(cs color.ColorSpace, p *pipeline.RasterPipelineBuilder) bool
	Transform(ts path.Transform)
	ApplyOpacity(opacity float32)
}

// SolidColor a solid color shader.
type SolidColor struct {
	color color.Color
}

func NewSolidColor(color color.Color) *SolidColor {
	return &SolidColor{color: color}
}

func (sc *SolidColor) Color() color.Color {
	return sc.color
}

func (sc *SolidColor) IsOpaque() bool {
	return sc.color.IsOpaque()
}

func (sc *SolidColor) PushStages(cs color.ColorSpace, p *pipeline.RasterPipelineBuilder) bool {
	expanded := cs.ExpandColor(sc.color)
	premultiplied := expanded.Premultiply()
	p.PushUniformColor(premultiplied)
	return true
}

func (sc *SolidColor) Transform(ts path.Transform) {
	// Solid color shaders don't need transform
}

func (sc *SolidColor) ApplyOpacity(opacity float32) {
	sc.color.ApplyOpacity(opacity)
}
