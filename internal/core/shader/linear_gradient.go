// Copyright 2006 The Android Open Source Project
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package shader

import (
	"github.com/chewxy/math32"
	"github.com/lumifloat/tinyskia/internal/core/color"
	"github.com/lumifloat/tinyskia/internal/core/pipeline"
	"github.com/lumifloat/tinyskia/internal/numeric/scalar"
	"github.com/lumifloat/tinyskia/internal/path"
)

// A linear gradient shader.
type LinearGradient struct {
	base *Gradient
}

// NewLinearGradient creates a new linear gradient shader.
func NewLinearGradient(start path.Point, end path.Point, stops []GradientStop, mode SpreadMode, transform path.Transform) Shader {
	if len(stops) == 0 {
		return nil
	}

	if len(stops) == 1 {
		color := stops[0].color
		return &SolidColor{color: color}
	}

	length := end.Sub(start).Length()
	if math32.IsInf(float32(length), 0) || math32.IsNaN(float32(length)) {
		return nil
	}

	if scalar.IsNearlyZeroWithinTolerance(length, DEGENERATE_THRESHOLD) {
		// Degenerate gradient, the only tricky complication is when in clamp mode,
		// the limit of the gradient approaches two half planes of solid color
		// (first and last). However, they are divided by the line perpendicular
		// to the start and end path.Point, which becomes undefined once start and end
		// are exactly the same, so just use the end color for a stable solution.
		switch mode {
		case SpreadModePad:
			color := stops[len(stops)-1].color
			return &SolidColor{color: color}
		case SpreadModeReflect, SpreadModeRepeat:
			avgColor := averageGradientColor(stops)
			return &SolidColor{color: avgColor}
		}
	}

	if _, ok := transform.Invert(); !ok {
		return nil
	}

	unitTs, ok := pointsToUnitTs(start, end)
	if !ok {
		return nil
	}

	return &LinearGradient{
		base: NewGradient(stops, mode, transform, unitTs),
	}
}

func (lg *LinearGradient) IsOpaque() bool {
	return lg.base.colorsAreOpaque
}

func (lg *LinearGradient) PushStages(cs color.ColorSpace, p *pipeline.RasterPipelineBuilder) bool {
	return lg.base.PushStages(p, cs, func(p *pipeline.RasterPipelineBuilder) {}, func(p *pipeline.RasterPipelineBuilder) {})
}

func (lg *LinearGradient) Transform(ts path.Transform) {
	lg.base.transform = lg.base.transform.PostConcat(ts)
}

func (lg *LinearGradient) ApplyOpacity(opacity float32) {
	lg.base.ApplyOpacity(opacity)
}

func pointsToUnitTs(start path.Point, end path.Point) (path.Transform, bool) {
	vec := end.Sub(start)
	mag := vec.Length()
	var inv float32
	if mag != 0.0 {
		inv = 1.0 / mag
	} else {
		inv = 0.0
	}

	vec = vec.WithScaleFrom(inv)

	ts := tsFromSinCosAt(-vec.Y, vec.X, start.X, start.Y)
	ts = ts.PostTranslate(-start.X, -start.Y)
	ts = ts.PostScale(inv, inv)
	return ts, true
}

func averageGradientColor(points []GradientStop) color.Color {
	// The gradient is a piecewise linear interpolation between colors. For a given interval,
	// the integral between the two endpoints is 0.5 * (ci + cj) * (pj - pi), which provides that
	// intervals average color. The overall average color is thus the sum of each piece.
	var r, g, b, a float32

	for i := 0; i < len(points)-1; i++ {
		p0 := points[i]
		p1 := points[i+1]

		w := p1.position.Get() - p0.position.Get()
		// 0.5 * w * (c1 + c0)
		r += 0.5 * w * (p1.color.Red() + p0.color.Red())
		g += 0.5 * w * (p1.color.Green() + p0.color.Green())
		b += 0.5 * w * (p1.color.Blue() + p0.color.Blue())
		a += 0.5 * w * (p1.color.Alpha() + p0.color.Alpha())
	}

	// Now account for any implicit intervals at the start or end of the stop definitions
	if points[0].position.Get() > 0.0 {
		w := points[0].position.Get()
		r += w * points[0].color.Red()
		g += w * points[0].color.Green()
		b += w * points[0].color.Blue()
		a += w * points[0].color.Alpha()
	}

	lastIdx := len(points) - 1
	if points[lastIdx].position.Get() < 1.0 {
		w := 1.0 - points[lastIdx].position.Get()
		r += w * points[lastIdx].color.Red()
		g += w * points[lastIdx].color.Green()
		b += w * points[lastIdx].color.Blue()
		a += w * points[lastIdx].color.Alpha()
	}

	color, _ := color.ColorFromRGBA(r, g, b, a)
	return color
}

func tsFromSinCosAt(sin, cos, px, py float32) path.Transform {
	cosInv := 1.0 - cos
	return path.TransformFromRow(
		cos,
		sin,
		-sin,
		cos,
		sdot(sin, py, cosInv, px),
		sdot(-sin, px, cosInv, py),
	)
}

func sdot(a, b, c, d float32) float32 {
	return a*b + c*d
}
