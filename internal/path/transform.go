// Copyright 2006 The Android Open Source Project
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package path

import (
	"github.com/chewxy/math32"

	"github.com/lumifloat/tinyskia/internal/numeric/scalar"
)

// Transform is an affine transformation matrix.
//
// Unlike other types, doesn't guarantee to be valid. This is Skia quirk.
// Meaning Transform(0, 0, 0, 0, 0, 0) is ok, while it's technically not.
// Non-finite values are also not an error.
type Transform struct {
	SX, KX, KY, SY, TX, TY float32
}

// NewTransformDefault creates a default Transform.
func NewTransformDefault() Transform {
	return Transform{
		SX: 1.0,
		KX: 0.0,
		KY: 0.0,
		SY: 1.0,
		TX: 0.0,
		TY: 0.0,
	}
}

// Identity creates an identity transform.
func Identity() Transform {
	return NewTransformDefault()
}

// NewTransform creates a new Transform.
//
// We are using column-major-column-vector matrix notation, therefore it's ky-kx, not kx-ky.
func NewTransform(sx, ky, kx, sy, tx, ty float32) Transform {
	return Transform{
		SX: sx,
		KX: kx,
		KY: ky,
		SY: sy,
		TX: tx,
		TY: ty,
	}
}

// NewTransformFromTranslate creates a new translating Transform.
func NewTransformFromTranslate(tx, ty float32) Transform {
	return NewTransform(1.0, 0.0, 0.0, 1.0, tx, ty)
}

// NewTransformFromScale creates a new scaling Transform.
func NewTransformFromScale(sx, sy float32) Transform {
	return NewTransform(sx, 0.0, 0.0, sy, 0.0, 0.0)
}

// NewTransformFromSkew creates a new skewing Transform.
func NewTransformFromSkew(kx, ky float32) Transform {
	return NewTransform(1.0, ky, kx, 1.0, 0.0, 0.0)
}

// NewTransformFromRotate creates a new rotating Transform.
//
// angle in degrees.
func NewTransformFromRotate(angle float32) Transform {
	v := angle * math32.Pi / 180.0
	a := math32.Cos(v)
	b := math32.Sin(v)
	c := -b
	d := a
	return NewTransform(a, b, c, d, 0.0, 0.0)
}

// NewTransformFromRotateAt creates a new rotating Transform at the specified position.
//
// angle in degrees.
func NewTransformFromRotateAt(angle, tx, ty float32) Transform {
	ts := Identity()
	ts = ts.PreTranslate(tx, ty)
	ts = ts.PreConcat(NewTransformFromRotate(angle))
	ts = ts.PreTranslate(-tx, -ty)
	return ts
}

// NewTransformFromSinCos
func NewTransformFromSinCos(sin, cos float32) Transform {
	return NewTransform(cos, sin, -sin, cos, 0.0, 0.0)
}

// TransformFromRow creates a transform from a row.
func TransformFromRow(a, b, c, d, tx, ty float32) Transform {
	return NewTransform(a, b, c, d, tx, ty)
}

// IsFinite checks that transform is finite.
func (t Transform) IsFinite() bool {
	return !math32.IsInf(t.SX, 0) && !math32.IsNaN(t.SX) &&
		!math32.IsInf(t.KY, 0) && !math32.IsNaN(t.KY) &&
		!math32.IsInf(t.KX, 0) && !math32.IsNaN(t.KX) &&
		!math32.IsInf(t.SY, 0) && !math32.IsNaN(t.SY) &&
		!math32.IsInf(t.TX, 0) && !math32.IsNaN(t.TX) &&
		!math32.IsInf(t.TY, 0) && !math32.IsNaN(t.TY)
}

// IsValid checks that transform is finite and has non-zero scale.
func (t Transform) IsValid() bool {
	if t.IsFinite() {
		sx, sy := t.GetScale()
		return !scalar.IsNearlyZeroWithinTolerance(sx, math32.SmallestNonzeroFloat32) &&
			!scalar.IsNearlyZeroWithinTolerance(sy, math32.SmallestNonzeroFloat32)
	}
	return false
}

// IsIdentity checks that transform is identity.
func (t Transform) IsIdentity() bool {
	return t == Identity()
}

// IsScale checks that transform is scale-only.
func (t Transform) IsScale() bool {
	return t.HasScale() && !t.HasSkew() && !t.HasTranslate()
}

// IsSkew checks that transform is skew-only.
func (t Transform) IsSkew() bool {
	return !t.HasScale() && t.HasSkew() && !t.HasTranslate()
}

// IsTranslate checks that transform is translate-only.
func (t Transform) IsTranslate() bool {
	return !t.HasScale() && !t.HasSkew() && t.HasTranslate()
}

// IsScaleTranslate checks that transform contains only scale and translate.
func (t Transform) IsScaleTranslate() bool {
	return (t.HasScale() || t.HasTranslate()) && !t.HasSkew()
}

// HasScale checks that transform contains a scale part.
func (t Transform) HasScale() bool {
	return t.SX != 1.0 || t.SY != 1.0
}

// HasSkew checks that transform contains a skew part.
func (t Transform) HasSkew() bool {
	return t.KX != 0.0 || t.KY != 0.0
}

// HasTranslate checks that transform contains a translate part.
func (t Transform) HasTranslate() bool {
	return t.TX != 0.0 || t.TY != 0.0
}

// GetScale returns transform's scale part.
func (t Transform) GetScale() (float32, float32) {
	xScale := math32.Sqrt(t.SX*t.SX + t.KX*t.KX)
	yScale := math32.Sqrt(t.KY*t.KY + t.SY*t.SY)
	return xScale, yScale
}

// PreScale pre-scales the current transform.
func (t Transform) PreScale(sx, sy float32) Transform {
	return t.PreConcat(NewTransformFromScale(sx, sy))
}

// PostScale post-scales the current transform.
func (t Transform) PostScale(sx, sy float32) Transform {
	return t.PostConcat(NewTransformFromScale(sx, sy))
}

// PreTranslate pre-translates the current transform.
func (t Transform) PreTranslate(tx, ty float32) Transform {
	return t.PreConcat(NewTransformFromTranslate(tx, ty))
}

// PostTranslate post-translates the current transform.
func (t Transform) PostTranslate(tx, ty float32) Transform {
	return t.PostConcat(NewTransformFromTranslate(tx, ty))
}

// PreRotate pre-rotates the current transform.
func (t Transform) PreRotate(angle float32) Transform {
	return t.PreConcat(NewTransformFromRotate(angle))
}

// PostRotate post-rotates the current transform.
func (t Transform) PostRotate(angle float32) Transform {
	return t.PostConcat(NewTransformFromRotate(angle))
}

// PreRotateAt pre-rotates the current transform by the specified position.
func (t Transform) PreRotateAt(angle, tx, ty float32) Transform {
	return t.PreConcat(NewTransformFromRotateAt(angle, tx, ty))
}

// PostRotateAt post-rotates the current transform by the specified position.
func (t Transform) PostRotateAt(angle, tx, ty float32) Transform {
	return t.PostConcat(NewTransformFromRotateAt(angle, tx, ty))
}

// PreConcat pre-concats the current transform.
func (t Transform) PreConcat(other Transform) Transform {
	return concat(t, other)
}

// PostConcat post-concats the current transform.
func (t Transform) PostConcat(other Transform) Transform {
	return concat(other, t)
}

// MapPoint transforms a point using the current transform.
func (t Transform) MapPoint(point *Point) {
	if t.IsIdentity() {
		return
	} else if t.IsTranslate() {
		point.X += t.TX
		point.Y += t.TY
	} else if t.IsScaleTranslate() {
		point.X = point.X*t.SX + t.TX
		point.Y = point.Y*t.SY + t.TY
	} else {
		x := point.X*t.SX + point.Y*t.KX + t.TX
		y := point.X*t.KY + point.Y*t.SY + t.TY
		point.X = x
		point.Y = y
	}
}

// MapPoints transforms a slice of points using the current transform.
func (t Transform) MapPoints(points []Point) {
	if len(points) == 0 {
		return
	}

	if t.IsIdentity() {
		return
	} else if t.IsTranslate() {
		for i := range points {
			points[i].X += t.TX
			points[i].Y += t.TY
		}
	} else if t.IsScaleTranslate() {
		for i := range points {
			points[i].X = points[i].X*t.SX + t.TX
			points[i].Y = points[i].Y*t.SY + t.TY
		}
	} else {
		for i := range points {
			x := points[i].X*t.SX + points[i].Y*t.KX + t.TX
			y := points[i].X*t.KY + points[i].Y*t.SY + t.TY
			points[i].X = x
			points[i].Y = y
		}
	}
}

// Invert returns an inverted transform.
func (t Transform) Invert() (Transform, bool) {
	if t.IsIdentity() {
		return t, true
	}
	return invert(t)
}

func invert(t Transform) (Transform, bool) {
	if t.IsScaleTranslate() {
		if t.HasScale() {
			invX := 1.0 / t.SX
			invY := 1.0 / t.SY
			return NewTransform(invX, 0.0, 0.0, invY, -t.TX*invX, -t.TY*invY), true
		} else {
			return NewTransformFromTranslate(-t.TX, -t.TY), true
		}
	}

	invDet, ok := invDeterminant(t)
	if !ok {
		return Transform{}, false
	}

	invTs := computeInv(t, invDet)
	if invTs.IsFinite() {
		return invTs, true
	}
	return Transform{}, false
}

func invDeterminant(t Transform) (float32, bool) {
	det := dcross(t.SX, t.SY, t.KX, t.KY)
	tolerance := scalar.SCALAR_NEARLY_ZERO * scalar.SCALAR_NEARLY_ZERO * scalar.SCALAR_NEARLY_ZERO
	if scalar.IsNearlyZeroWithinTolerance(det, tolerance) {
		return 0, false
	}
	return 1.0 / det, true
}

func computeInv(t Transform, invDet float32) Transform {
	return NewTransform(
		t.SY*invDet,
		-t.KY*invDet,
		-t.KX*invDet,
		t.SX*invDet,
		dcrossDScale(t.KX, t.TY, t.SY, t.TX, invDet),
		dcrossDScale(t.KY, t.TX, t.SX, t.TY, invDet),
	)
}

func dcross(a, b, c, d float32) float32 {
	return a*b - c*d
}

func dcrossDScale(a, b, c, d, scale float32) float32 {
	return dcross(a, b, c, d) * scale
}

func concat(a, b Transform) Transform {
	if a.IsIdentity() {
		return b
	} else if b.IsIdentity() {
		return a
	} else if !a.HasSkew() && !b.HasSkew() {
		return NewTransform(
			a.SX*b.SX,
			0.0,
			0.0,
			a.SY*b.SY,
			a.SX*b.TX+a.TX,
			a.SY*b.TY+a.TY,
		)
	} else {
		return NewTransform(
			a.SX*b.SX+a.KX*b.KY,
			a.KY*b.SX+a.SY*b.KY,
			a.SX*b.KX+a.KX*b.SY,
			a.KY*b.KX+a.SY*b.SY,
			(a.SX*b.TX+a.KX*b.TY)+a.TX,
			(a.KY*b.TX+a.SY*b.TY)+a.TY,
		)
	}
}
