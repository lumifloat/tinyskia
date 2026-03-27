// Copyright 2006 The Android Open Source Project
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package shader

import (
	"math"

	"github.com/chewxy/math32"

	"github.com/lumifloat/tinyskia/color"
	"github.com/lumifloat/tinyskia/internal/scalar"
	"github.com/lumifloat/tinyskia/path"
	"github.com/lumifloat/tinyskia/pipeline"
)

type focalData struct {
	r1        float32 // r1 after mapping focal point to (0, 0)
	focalX    float32 // f
	isSwapped bool
}

func (fd *focalData) set(r0, r1 float32, matrix *path.Transform) bool {
	fd.isSwapped = false
	fd.focalX = r0 / (r0 - r1)

	if scalar.IsNearlyZero(fd.focalX - 1.0) {
		// swap r0, r1
		*matrix = matrix.PostTranslate(-1.0, 0.0).PostScale(-1.0, 1.0)
		r0, r1 = r1, r0

		fd.focalX = 0.0 // because r0 is now 0
		fd.isSwapped = true
	}

	// Map {focal point, (1, 0)} to {(0, 0), (1, 0)}
	from := [2]path.Point{{X: fd.focalX, Y: 0.0}, {X: 1.0, Y: 0.0}}
	to := [2]path.Point{{X: 0.0, Y: 0.0}, {X: 1.0, Y: 0.0}}

	focalMatrix, ok := tsFromPolyToPoly(from[0], from[1], to[0], to[1])
	if !ok {
		return false
	}

	*matrix = matrix.PostConcat(focalMatrix)
	fd.r1 = r1 / math32.Abs(1.0-fd.focalX) // focalMatrix has a scale of 1/(1-f).

	// The following transformations are just to accelerate the shader computation by saving
	// some arithmetic operations.
	if fd.isFocalOnCircle() {
		*matrix = matrix.PostScale(0.5, 0.5)
	} else {
		*matrix = matrix.PostScale(
			fd.r1/(fd.r1*fd.r1-1.0),
			1.0/math32.Sqrt(math32.Abs(fd.r1*fd.r1-1.0)),
		)
	}

	*matrix = matrix.PostScale(math32.Abs(1.0-fd.focalX), math32.Abs(1.0-fd.focalX)) // scale |1 - f|

	return true
}

func (fd focalData) isFocalOnCircle() bool {
	return scalar.IsNearlyZero(1.0 - fd.r1)
}

func (fd focalData) isWellBehaved() bool {
	return !fd.isFocalOnCircle() && fd.r1 > 1.0
}

func (fd focalData) isNativelyFocal() bool {
	return scalar.IsNearlyZero(fd.focalX)
}

type radialGradientType int

const (
	radialTypeRadial radialGradientType = iota
	radialTypeStrip
	radialTypeFocal
)

type gradientType struct {
	kind     radialGradientType
	radius1  float32
	radius2  float32
	scaledR0 float32
	focal    focalData
}

// A 2-point conical gradient shader.
type RadialGradient struct {
	base         *Gradient
	gradientType gradientType
}

// NewRadialGradient creates a new two-point conical gradient shader.
func NewRadialGradient(startPoint path.Point, startRadius float32, endPoint path.Point, endRadius float32, stops []GradientStop, mode SpreadMode, transform path.Transform) Shader {
	if startRadius < 0.0 || endRadius < 0.0 {
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

	diff := startPoint.Sub(endPoint)
	length := diff.Length()
	if math.IsInf(float64(length), 0) || math.IsNaN(float64(length)) {
		return nil
	}

	if scalar.IsNearlyZeroWithinTolerance(length, DEGENERATE_THRESHOLD) {
		if scalar.IsNearlyEqualWithinTolerance(startRadius, endRadius, DEGENERATE_THRESHOLD) {
			if mode == SpreadModePad && endRadius > DEGENERATE_THRESHOLD {
				startColor := stops[0].color
				endColor := stops[len(stops)-1].color
				newStops := []GradientStop{
					NewGradientStop(0.0, startColor),
					NewGradientStop(1.0, startColor),
					NewGradientStop(1.0, endColor),
				}
				return newRadialUnchecked(startPoint, endRadius, newStops, mode, transform)
			}
			return nil
		}

		if scalar.IsNearlyZeroWithinTolerance(startRadius, DEGENERATE_THRESHOLD) {
			return newRadialUnchecked(startPoint, endRadius, stops, mode, transform)
		}
	}

	return createRadial(startPoint, startRadius, endPoint, endRadius, stops, mode, transform)
}

func newRadialUnchecked(
	center path.Point,
	radius float32,
	stops []GradientStop,
	mode SpreadMode,
	transform path.Transform,
) Shader {
	inv := 1.0 / radius
	pointsToUnit := path.NewTransformFromTranslate(-center.X, -center.Y).PostScale(inv, inv)

	return &RadialGradient{
		base: NewGradient(stops, mode, transform, pointsToUnit),
		gradientType: gradientType{
			kind:    radialTypeRadial,
			radius1: 0.0,
			radius2: radius,
		},
	}
}

func (rg *RadialGradient) IsOpaque() bool {
	return rg.base.colorsAreOpaque
}

func (rg *RadialGradient) PushStages(cs color.ColorSpace, p *pipeline.RasterPipelineBuilder) bool {
	var p0, p1 float32
	switch rg.gradientType.kind {
	case radialTypeRadial:
		if rg.gradientType.radius1 == 0.0 {
			p0, p1 = 1.0, 0.0
		} else {
			dRadius := rg.gradientType.radius2 - rg.gradientType.radius1
			p0 = math32.Max(rg.gradientType.radius1, rg.gradientType.radius2) / dRadius
			p1 = -rg.gradientType.radius1 / dRadius
		}
	case radialTypeStrip:
		p0, p1 = float32(rg.gradientType.scaledR0*rg.gradientType.scaledR0), 0.0
	case radialTypeFocal:
		p0, p1 = 1.0/rg.gradientType.focal.r1, rg.gradientType.focal.focalX
	}

	p.Ctx.TwoPointConicalGradient = pipeline.TwoPointConicalGradientCtx{
		Mask: [8]uint32{},
		P0:   p0,
		P1:   p1,
	}

	return rg.base.PushStages(p, cs, func(p *pipeline.RasterPipelineBuilder) {
		switch rg.gradientType.kind {
		case radialTypeRadial:
			p.Push(pipeline.StageXYToRadius)
			if p0 != 1.0 || p1 != 0.0 {
				p.Push(pipeline.StageApplyConcentricScaleBias)
			}
		case radialTypeStrip:
			p.Push(pipeline.StageXYTo2PtConicalStrip)
			p.Push(pipeline.StageMask2PtConicalNan)
		case radialTypeFocal:
			fd := rg.gradientType.focal
			if fd.isFocalOnCircle() {
				p.Push(pipeline.StageXYTo2PtConicalFocalOnCircle)
			} else if fd.isWellBehaved() {
				p.Push(pipeline.StageXYTo2PtConicalWellBehaved)
			} else if fd.isSwapped || (1.0-fd.focalX) < 0.0 {
				p.Push(pipeline.StageXYTo2PtConicalSmaller)
			} else {
				p.Push(pipeline.StageXYTo2PtConicalGreater)
			}

			if !fd.isWellBehaved() {
				p.Push(pipeline.StageMask2PtConicalDegenerates)
			}
			if (1.0 - fd.focalX) < 0.0 {
				p.Push(pipeline.StageNegateX)
			}
			if !fd.isNativelyFocal() {
				p.Push(pipeline.StageAlter2PtConicalCompensateFocal)
			}
			if fd.isSwapped {
				p.Push(pipeline.StageAlter2PtConicalUnswap)
			}
		}
	}, func(p *pipeline.RasterPipelineBuilder) {
		switch rg.gradientType.kind {
		case radialTypeStrip:
			p.Push(pipeline.StageApplyVectorMask)
		case radialTypeFocal:
			if !rg.gradientType.focal.isWellBehaved() {
				p.Push(pipeline.StageApplyVectorMask)
			}
		}
	})
}

func (rg *RadialGradient) Transform(ts path.Transform) {
	rg.base.transform = rg.base.transform.PostConcat(ts)
}

func (rg *RadialGradient) ApplyOpacity(opacity float32) {
	rg.base.ApplyOpacity(opacity)
}

func createRadial(
	c0 path.Point, r0 float32,
	c1 path.Point, r1 float32,
	stops []GradientStop,
	mode SpreadMode,
	transform path.Transform,
) Shader {
	var gType gradientType
	var gMatrix path.Transform

	dCenterLen := c0.Sub(c1).Length()

	if scalar.IsNearlyZero(dCenterLen) {
		if scalar.IsNearlyZero(math32.Max(r0, r1)) || scalar.IsNearlyEqual(r0, r1) {
			return nil
		}
		scale := 1.0 / math32.Max(r0, r1)
		gMatrix = path.NewTransformFromTranslate(-c1.X, -c1.Y).PostScale(scale, scale)
		gType = gradientType{
			kind:    radialTypeRadial,
			radius1: r0,
			radius2: r1,
		}
	} else {
		var ok bool
		gMatrix, ok = mapToUnitX(c0, c1)
		if !ok {
			return nil
		}
		if scalar.IsNearlyZero(r0 - r1) {
			gType = gradientType{
				kind:     radialTypeStrip,
				scaledR0: r0 / dCenterLen,
			}
		} else {
			gType = gradientType{kind: radialTypeFocal}
			if !gType.focal.set(r0/dCenterLen, r1/dCenterLen, &gMatrix) {
				return nil
			}
		}
	}

	return &RadialGradient{
		base:         NewGradient(stops, mode, transform, gMatrix),
		gradientType: gType,
	}
}

func mapToUnitX(origin, xIsOne path.Point) (path.Transform, bool) {
	return tsFromPolyToPoly(
		origin, xIsOne,
		path.Point{X: 0.0, Y: 0.0}, path.Point{X: 1.0, Y: 0.0},
	)
}

func tsFromPolyToPoly(src1, src2, dst1, dst2 path.Point) (path.Transform, bool) {
	tmp := fromPoly2(src1, src2)
	res, ok := tmp.Invert()
	if !ok {
		return path.Transform{}, false
	}
	tmpDst := fromPoly2(dst1, dst2)
	return tmpDst.PreConcat(res), true
}

func fromPoly2(p0, p1 path.Point) path.Transform {
	return path.TransformFromRow(
		p1.Y-p0.Y,
		p0.X-p1.X,
		p1.X-p0.X,
		p1.Y-p0.Y,
		p0.X,
		p0.Y,
	)
}
