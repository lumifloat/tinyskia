// Copyright 2006 The Android Open Source Project
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Package path provides a memory-efficient Bezier path container, path builder, path stroker and path dasher.
//
// Also provides some basic geometry types, but they will be moved to an external crate eventually.
//
// Note that all types use float32.
package path

import (
	"math"

	"github.com/chewxy/math32"

	"github.com/lumifloat/tinyskia/internal/scalar"
	"github.com/lumifloat/tinyskia/internal/wide"
)

// Point is a point.
//
// Doesn't guarantee to be finite.
type Point struct {
	X float32
	Y float32
}

// NewPointFromF32x2 creates a new Point from F32x2.
func NewPointFromF32x2(p wide.F32x2) Point {
	return Point{X: p[0], Y: p[1]}
}

// ToF32x2 converts the Point to F32x2.
func (p Point) ToF32x2() wide.F32x2 {
	return wide.F32x2{p.X, p.Y}
}

// IsZero returns true if x and y are both zero.
func (p Point) IsZero() bool {
	return p.X == 0.0 && p.Y == 0.0
}

// IsFinite returns true if both x and y are measurable values.
//
// Both values are other than infinities and NaN.
func (p Point) IsFinite() bool {
	return !math32.IsNaN(p.X) && !math32.IsInf(p.X, 0) && !math32.IsNaN(p.Y) && !math32.IsInf(p.Y, 0)
}

// AlmostEqual checks that two Points are almost equal.
func (p Point) AlmostEqual(other Point) bool {
	return !p.Sub(other).canNormalize()
}

// EqualsWithinTolerance checks that two Points are almost equal using the specified tolerance.
func (p Point) EqualsWithinTolerance(other Point, tolerance float32) bool {
	return scalar.IsNearlyZeroWithinTolerance(p.X-other.X, tolerance) &&
		scalar.IsNearlyZeroWithinTolerance(p.Y-other.Y, tolerance)
}

// WithNormalizeFrom scales (X, Y) so that Length() returns one, while preserving ratio of X to Y,
// if possible.
//
// If prior length is nearly zero, returns a zero point and false;
// otherwise returns the normalized point and true.
func (p Point) WithNormalizeFrom() (Point, bool) {
	return p.WithLengthFrom(1.0)
}

func (p Point) canNormalize() bool {
	return p.IsFinite() && (p.X != 0.0 || p.Y != 0.0)
}

// Length returns the Euclidean distance from origin.
func (p Point) Length() float32 {
	mag2 := p.X*p.X + p.Y*p.Y
	if !math32.IsInf(mag2, 0) && !math32.IsNaN(mag2) {
		return math32.Sqrt(mag2)
	}
	// In float32, the double-precision logic from Rust is inherent.
	return math32.Sqrt(p.X*p.X + p.Y*p.Y)
}

// WithLengthFrom returns a vector scaled so that its distance from origin returns length, if possible.
//
// If former length is nearly zero, returns a zero point and false;
// otherwise returns the new point and true.
func (p Point) WithLengthFrom(length float32) (Point, bool) {
	// our mag2 step overflowed to infinity, so use doubles instead.
	// much slower, but needed when x or y are very large, other wise we
	// divide by inf. and return (0,0) vector.
	xx := float64(p.X)
	yy := float64(p.Y)
	dmag := math.Sqrt(xx*xx + yy*yy)
	if dmag == 0 || math.IsNaN(dmag) || math.IsInf(dmag, 0) {
		return Point{}, false
	}

	dscale := float64(length) / dmag
	nx := p.X * float32(dscale)
	ny := p.Y * float32(dscale)

	if nx == 0.0 && ny == 0.0 {
		return Point{}, false
	}

	return Point{X: nx, Y: ny}, true
}

// Distance returns the Euclidean distance between two points.
func (p Point) Distance(other Point) float32 {
	return p.Sub(other).Length()
}

// Dot returns the dot product of two points.
func (p Point) Dot(other Point) float32 {
	return p.X*other.X + p.Y*other.Y
}

// Cross returns the cross product of vector and vec.
//
// Vector and vec form three-dimensional vectors with z-axis value equal to zero.
// The cross product is a three-dimensional vector with x-axis and y-axis values
// equal to zero. The cross product z-axis component is returned.
func (p Point) Cross(other Point) float32 {
	return p.X*other.Y - p.Y*other.X
}

func (p Point) DistanceToSqd(pt Point) float32 {
	dx := p.X - pt.X
	dy := p.Y - pt.Y
	return dx*dx + dy*dy
}

func (p Point) LengthSqd() float32 {
	return p.Dot(p)
}

// WithScaleFrom scales Point by scale and returns the result.
func (p Point) WithScaleFrom(scale float32) Point {
	return Point{X: p.X * scale, Y: p.Y * scale}
}

func (p Point) swapCoords() Point {
	return Point{X: p.Y, Y: p.X}
}

// WithRotateCWFrom rotates the point clockwise.
func (p Point) WithRotateCWFrom() Point {
	res := p.swapCoords()
	res.X = -res.X
	return res
}

// WithRotateCCWFrom rotates the point counter-clockwise.
func (p Point) WithRotateCCWFrom() Point {
	res := p.swapCoords()
	res.Y = -res.Y
	return res
}

// Neg returns the negative of the point.
func (p Point) Neg() Point {
	return Point{X: -p.X, Y: -p.Y}
}

// Add returns the sum of two points.
func (p Point) Add(other Point) Point {
	return Point{X: p.X + other.X, Y: p.Y + other.Y}
}

// Sub returns the difference of two points.
func (p Point) Sub(other Point) Point {
	return Point{X: p.X - other.X, Y: p.Y - other.Y}
}

// Mul returns the component-wise product of two points.
func (p Point) Mul(other Point) Point {
	return Point{X: p.X * other.X, Y: p.Y * other.Y}
}
