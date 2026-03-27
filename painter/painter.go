// Copyright 2006 The Android Open Source Project
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package painter

import (
	"github.com/lumifloat/tinyskia/blend"
	"github.com/lumifloat/tinyskia/color"
	"github.com/lumifloat/tinyskia/shader"
)

// FillRule is a path filling rule.
type FillRule int

const (
	// Winding specifies that "inside" is computed by a non-zero sum of signed edge crossings.
	FillRuleWinding FillRule = iota
	// EvenOdd specifies that "inside" is computed by an odd number of edge crossings.
	FillRuleEvenOdd
)

// Paint controls how a shape should be painted.
type Paint struct {
	// A paint shader.
	Shader shader.Shader
	// Paint blending mode.
	BlendMode blend.BlendMode
	// Enables anti-aliased painting.
	AntiAlias bool
	// Colorspace for blending.
	Colorspace color.ColorSpace
	// Forces the high quality/precision rendering pipeline.
	ForceHQPipeline bool
}

// DefaultPaint returns a Paint with default values.
func DefaultPaint() Paint {
	return Paint{
		Shader:          shader.NewSolidColor(color.ColorBlack),
		BlendMode:       blend.BlendModeSourceOver,
		AntiAlias:       true,
		Colorspace:      color.ColorSpaceLinear,
		ForceHQPipeline: false,
	}
}

// SetColor sets a paint source to a solid color.
func (p *Paint) SetColor(color color.Color) {
	p.Shader = shader.NewSolidColor(color)
}

// SetColorRGBA8 sets a paint source to a solid color using RGBA8 values.
func (p *Paint) SetColorRGBA8(r, g, b, a uint8) {
	p.SetColor(color.ColorFromRGBA8(r, g, b, a))
}

// IsSolidColor checks that the paint source is a solid color.
func (p *Paint) IsSolidColor() bool {
	_, ok := p.Shader.(*shader.SolidColor)
	return ok
}
