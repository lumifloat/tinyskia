// Copyright 2006 The Android Open Source Project
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package shader

import (
	"github.com/lumifloat/tinyskia/internal/normalized"
	"github.com/lumifloat/tinyskia/internal/scalar"
	"github.com/lumifloat/tinyskia/path"
	"github.com/lumifloat/tinyskia/pipeline"

	"github.com/lumifloat/tinyskia/color"
)

// The default SCALAR_NEARLY_ZERO threshold of .0024 is too big and causes regressions for svg
// gradients defined in the wild.
const DEGENERATE_THRESHOLD float32 = 1.0 / (1 << 15)

// GradientStop is a gradient point.
type GradientStop struct {
	position normalized.NormalizedF32
	color    color.Color
}

// NewGradientStop creates a new gradient point.
// `position` will be clamped to a 0..=1 range.
func NewGradientStop(position float32, color color.Color) GradientStop {
	return GradientStop{
		position: normalized.NewNormalizedF32WithClamped(position),
		color:    color,
	}
}

type Gradient struct {
	stops           []GradientStop
	tileMode        SpreadMode
	transform       path.Transform
	pointsToUnit    path.Transform
	colorsAreOpaque bool
	hasUniformStops bool
}

func NewGradient(stops []GradientStop, tileMode SpreadMode, transform path.Transform, pointsToUnit path.Transform) *Gradient {
	if len(stops) <= 1 {
		return nil
	}

	// Note: we let the caller skip the first and/or last position.
	// i.e. pos[0] = 0.3, pos[1] = 0.7
	// In these cases, we insert dummy entries to ensure that the final data
	// will be bracketed by [0, 1].
	// i.e. our_pos[0] = 0, our_pos[1] = 0.3, our_pos[2] = 0.7, our_pos[3] = 1
	dummyFirst := stops[0].position.Get() != 0.0
	dummyLast := stops[len(stops)-1].position.Get() != 1.0

	var newStops []GradientStop
	// Now copy over the colors, adding the dummies as needed.
	if dummyFirst {
		newStops = append(newStops, NewGradientStop(0.0, stops[0].color))
	}
	newStops = append(newStops, stops...)
	if dummyLast {
		newStops = append(newStops, NewGradientStop(1.0, stops[len(stops)-1].color))
	}

	colorsAreOpaque := true
	for _, p := range newStops {
		if !p.color.IsOpaque() {
			colorsAreOpaque = false
			break
		}
	}

	// Pin the last value to 1.0, and make sure positions are monotonic.
	startIndex := 1
	if dummyFirst {
		startIndex = 0
	}

	var prev float32 = 0.0
	hasUniformStops := true
	uniformStep := newStops[startIndex].position.Get() - prev

	for i := startIndex; i < len(newStops); i++ {
		var curr float32
		if i+1 == len(newStops) {
			// The last one must be one.
			curr = 1.0
		} else {
			curr = scalar.Bound(newStops[i].position.Get(), prev, 1.0)
		}

		hasUniformStops = hasUniformStops && scalar.IsNearlyEqual(uniformStep, curr-prev)
		newStops[i].position = normalized.NewNormalizedF32WithClamped(curr)
		prev = curr
	}

	return &Gradient{
		stops:           newStops,
		tileMode:        tileMode,
		transform:       transform,
		pointsToUnit:    pointsToUnit,
		colorsAreOpaque: colorsAreOpaque,
		hasUniformStops: hasUniformStops,
	}
}

func (g *Gradient) PushStages(
	p *pipeline.RasterPipelineBuilder,
	cs color.ColorSpace,
	pushStagesPre func(*pipeline.RasterPipelineBuilder),
	pushStagesPost func(*pipeline.RasterPipelineBuilder),
) bool {
	p.Push(pipeline.StageSeedShader)

	ts, ok := g.transform.Invert()
	if !ok {
		return false
	}
	ts = ts.PostConcat(g.pointsToUnit)
	p.PushTransform(ts)

	pushStagesPre(p)

	switch g.tileMode {
	case SpreadModeReflect:
		p.Push(pipeline.StageReflectX1)
	case SpreadModeRepeat:
		p.Push(pipeline.StageRepeatX1)
	case SpreadModePad:
		if g.hasUniformStops {
			// We clamp only when the stops are evenly spaced.
			// If not, there may be hard stops, and clamping ruins hard stops at 0 and/or 1.
			// In that case, we must make sure we're using the general "gradient" stage,
			// which is the only stage that will correctly handle unclamped t.
			p.Push(pipeline.StagePadX1)
		}
	}

	// The two-stop case with stops at 0 and 1.
	if len(g.stops) == 2 {
		c0 := cs.ExpandColor(g.stops[0].color)
		c1 := cs.ExpandColor(g.stops[1].color)

		p.Ctx.EvenlySpaced2StopGradient = pipeline.EvenlySpaced2StopGradientCtx{
			Factor: pipeline.NewGradientColor(
				c1.Red()-c0.Red(),
				c1.Green()-c0.Green(),
				c1.Blue()-c0.Blue(),
				c1.Alpha()-c0.Alpha(),
			),
			Bias: pipeline.GradientColor{R: c0.Red(), G: c0.Green(), B: c0.Blue(), A: c0.Alpha()},
		}

		p.Push(pipeline.StageEvenlySpaced2StopGradient)
	} else {
		// Unlike Skia, we do not support the `evenly_spaced_gradient` stage.
		// In our case, there is no performance difference.
		ctx := pipeline.GradientCtx{}

		// Note: In order to handle clamps in search, the search assumes
		// a stop conceptually placed at -inf.
		// Therefore, the max number of stops is `self.points.len()+1`.

		// Remove the dummy stops inserted by Gradient::new
		// because they are naturally handled by the search method.
		firstStop, lastStop := 0, 1
		if len(g.stops) > 2 {
			if g.stops[0].color != g.stops[1].color {
				firstStop = 0
			} else {
				firstStop = 1
			}

			length := len(g.stops)
			if g.stops[length-2].color != g.stops[length-1].color {
				lastStop = length - 1
			} else {
				lastStop = length - 2
			}
		}

		tL := g.stops[firstStop].position.Get()
		cLExpanded := cs.ExpandColor(g.stops[firstStop].color)
		cL := pipeline.GradientColor{R: cLExpanded.Red(), G: cLExpanded.Green(), B: cLExpanded.Blue(), A: cLExpanded.Alpha()}
		ctx.PushConstColor(cL)
		ctx.TValues = append(ctx.TValues, 0.0)

		// N.B. lastStop is the index of the last stop, not one after.
		for i := firstStop; i < lastStop; i++ {
			tR := g.stops[i+1].position.Get()
			cRExpanded := cs.ExpandColor(g.stops[i+1].color)
			cR := pipeline.GradientColor{R: cRExpanded.Red(), G: cRExpanded.Green(), B: cRExpanded.Blue(), A: cRExpanded.Alpha()}

			if tL < tR {
				// For each stop we calculate a bias B and a scale factor F, such that
				// for any t between stops n and n+1, the color we want is B[n] + F[n]*t.
				invT := 1.0 / (tR - tL)
				f := pipeline.NewGradientColor(
					(cR.R-cL.R)*invT,
					(cR.G-cL.G)*invT,
					(cR.B-cL.B)*invT,
					(cR.A-cL.A)*invT,
				)

				ctx.Factors = append(ctx.Factors, f)

				ctx.Biases = append(ctx.Biases, pipeline.NewGradientColor(
					cL.R-f.R*tL,
					cL.G-f.G*tL,
					cL.B-f.B*tL,
					cL.A-f.A*tL,
				))

				ctx.TValues = append(ctx.TValues, normalized.NewNormalizedF32WithClamped(tL))
			}

			tL = tR
			cL = cR
		}

		ctx.PushConstColor(cL)
		ctx.TValues = append(ctx.TValues, normalized.NewNormalizedF32WithClamped(tL))

		ctx.Len = len(ctx.Factors)

		// Fill with zeros until we have enough data.
		// Note: TValues for padding stops are set to 1.0 to prevent accidental matches
		for len(ctx.Factors) < 16 {
			ctx.Factors = append(ctx.Factors, pipeline.GradientColor{})
			ctx.Biases = append(ctx.Biases, pipeline.GradientColor{})
			//ctx.TValues[ctx.Len] = 1.0
			//ctx.Len++
		}

		p.Push(pipeline.StageGradient)
		p.Ctx.Gradient = ctx
	}

	if !g.colorsAreOpaque {
		p.Push(pipeline.StagePremultiply)
	}

	pushStagesPost(p)

	return true
}

func (g *Gradient) ApplyOpacity(opacity float32) {
	for i := range g.stops {
		g.stops[i].color.ApplyOpacity(opacity)
	}

	g.colorsAreOpaque = true
	for _, p := range g.stops {
		if !p.color.IsOpaque() {
			g.colorsAreOpaque = false
			break
		}
	}
}
