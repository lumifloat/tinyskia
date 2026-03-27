// Copyright 2006 The Android Open Source Project
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package shader

import (
	"github.com/chewxy/math32"
	"github.com/lumifloat/tinyskia/color"
	"github.com/lumifloat/tinyskia/internal/scalar"
	"github.com/lumifloat/tinyskia/path"
	"github.com/lumifloat/tinyskia/pipeline"
)

// A radial gradient.
type SweepGradient struct {
	base *Gradient
	t0   float32
	t1   float32
}

// NewSweepGradient creates a new sweep gradient shader.
func NewSweepGradient(center path.Point, startAngle float32, endAngle float32, stops []GradientStop, mode SpreadMode, transform path.Transform) Shader {
	if math32.IsInf(startAngle, 0) || math32.IsNaN(startAngle) ||
		math32.IsInf(endAngle, 0) || math32.IsNaN(endAngle) ||
		startAngle > endAngle {
		return nil
	}

	if len(stops) == 0 {
		return nil
	}
	if len(stops) == 1 {
		color := stops[0].color
		return &SolidColor{color: color}
	}

	if _, ok := transform.Invert(); !ok {
		return nil
	}

	if scalar.IsNearlyEqualWithinTolerance(startAngle, endAngle, DEGENERATE_THRESHOLD) {
		if mode == SpreadModePad && endAngle > DEGENERATE_THRESHOLD {
			// In this case, the first color is repeated from 0 to the angle, then a hardstop
			// switches to the last color (all other colors are compressed to the infinitely
			// thin interpolation region).
			frontColor := stops[0].color
			backColor := stops[len(stops)-1].color
			newStops := []GradientStop{
				NewGradientStop(0.0, frontColor),
				NewGradientStop(1.0, frontColor),
				NewGradientStop(1.0, backColor),
			}
			return NewSweepGradient(center, 0.0, endAngle, newStops, mode, transform)
		}
		return nil
	}

	if startAngle <= 0.0 && endAngle >= 360.0 {
		mode = SpreadModePad
	}

	t0 := startAngle / 360.0
	t1 := endAngle / 360.0

	return &SweepGradient{
		base: NewGradient(
			stops,
			mode,
			transform,
			path.NewTransformFromTranslate(-center.X, -center.Y),
		),
		t0: t0,
		t1: t1,
	}
}

func (sg *SweepGradient) IsOpaque() bool {
	return sg.base.colorsAreOpaque
}

func (sg *SweepGradient) PushStages(cs color.ColorSpace, p *pipeline.RasterPipelineBuilder) bool {
	scale := float32(1.0) / (sg.t1 - sg.t0)
	bias := -scale * sg.t0

	p.Ctx.TwoPointConicalGradient.P0 = scale
	p.Ctx.TwoPointConicalGradient.P1 = bias

	sg.base.PushStages(
		p,
		cs,
		func(p *pipeline.RasterPipelineBuilder) {
			p.Push(pipeline.StageXYToUnitAngle)
			if scale != 1.0 || bias != 0.0 {
				p.Push(pipeline.StageApplyConcentricScaleBias)
			}
		},
		func(p *pipeline.RasterPipelineBuilder) {},
	)

	return true
}

func (sg *SweepGradient) Transform(ts path.Transform) {
	sg.base.transform = sg.base.transform.PostConcat(ts)
}

func (sg *SweepGradient) ApplyOpacity(opacity float32) {
	sg.base.ApplyOpacity(opacity)
}
