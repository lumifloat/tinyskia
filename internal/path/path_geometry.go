// Copyright 2006 The Android Open Source Project
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// ! A collection of functions to work with Bezier paths.
// !
// ! Mainly for internal use. Do not rely on it!
package path

import (
	"github.com/chewxy/math32"

	"github.com/lumifloat/tinyskia/internal/numeric/normalized"
	"github.com/lumifloat/tinyskia/internal/numeric/scalar"
	"github.com/lumifloat/tinyskia/internal/numeric/wide"
)

// QuadCoeff is use for : eval(t) == A * t^2 + B * t + C
type QuadCoeff struct {
	A, B, C wide.F32x2
}

func NewQuadCoeffFromPoints(points [3]Point) QuadCoeff {
	c := points[0].ToF32x2()
	p1 := points[1].ToF32x2()
	p2 := points[2].ToF32x2()
	b := times2(p1.Sub(c))
	a := p2.Sub(times2(p1)).Add(c)

	return QuadCoeff{A: a, B: b, C: c}
}

func (q QuadCoeff) Eval(t wide.F32x2) wide.F32x2 {
	// (a*t + b)*t + c
	return (q.A.Mul(t).Add(q.B)).Mul(t).Add(q.C)
}

type CubicCoeff struct {
	A, B, C, D wide.F32x2
}

func NewCubicCoeffFromPoints(points [4]Point) CubicCoeff {
	p0 := points[0].ToF32x2()
	p1 := points[1].ToF32x2()
	p2 := points[2].ToF32x2()
	p3 := points[3].ToF32x2()
	three := wide.Splat(3.0)

	// a: p3 + 3*(p1 - p2) - p0
	a := p3.Add(three.Mul(p1.Sub(p2))).Sub(p0)
	// b: 3*(p2 - 2*p1 + p0)
	b := three.Mul(p2.Sub(times2(p1)).Add(p0))
	// c: 3*(p1 - p0)
	c := three.Mul(p1.Sub(p0))
	// d: p0
	d := p0

	return CubicCoeff{A: a, B: b, C: c, D: d}
}

func (q CubicCoeff) Eval(t wide.F32x2) wide.F32x2 {
	// ((a*t + b)*t + c)*t + d
	return q.A.Mul(t).Add(q.B).Mul(t).Add(q.C).Mul(t).Add(q.D)
}

func ChopQuadAt(src [3]Point, t normalized.NormalizedF32Exclusive, dst *[5]Point) {
	p0 := src[0].ToF32x2()
	p1 := src[1].ToF32x2()
	p2 := src[2].ToF32x2()
	tt := wide.Splat(t.Get())

	p01 := interp(p0, p1, tt)
	p12 := interp(p1, p2, tt)

	dst[0] = NewPointFromF32x2(p0)
	dst[1] = NewPointFromF32x2(p01)
	dst[2] = NewPointFromF32x2(interp(p01, p12, tt))
	dst[3] = NewPointFromF32x2(p12)
	dst[4] = NewPointFromF32x2(p2)
}

// FindUnitQuadRoots from Numerical Recipes in C.
//
// Q = -1/2 (B + sign(B) sqrt[B*B - 4*A*C])
// x1 = Q / A
// x2 = C / Q
func FindUnitQuadRoots(a, b, c float32, roots *[3]normalized.NormalizedF32Exclusive) int {
	if a == 0.0 {
		if r, ok := ValidUnitDivide(-c, b); ok {
			roots[0] = r
			return 1
		}
		return 0
	}

	dr := b*b - 4.0*a*c
	if dr < 0.0 {
		return 0
	}
	dr = math32.Sqrt(dr)

	if math32.IsInf(dr, 0) || math32.IsNaN(dr) {
		return 0
	}

	var q float32
	if b < 0.0 {
		q = -(b - dr) / 2.0
	} else {
		q = -(b + dr) / 2.0
	}

	offset := 0
	if r, ok := ValidUnitDivide(q, a); ok {
		roots[offset] = r
		offset++
	}

	if r, ok := ValidUnitDivide(c, q); ok {
		roots[offset] = r
		offset++
	}

	if offset == 2 {
		if roots[0] > roots[1] {
			roots[0], roots[1] = roots[1], roots[0]
		} else if roots[0] == roots[1] {
			offset--
		}
	}

	return offset
}

func ChopCubicAt2(src [4]Point, t normalized.NormalizedF32Exclusive, dst *[7]Point) {
	p0 := src[0].ToF32x2()
	p1 := src[1].ToF32x2()
	p2 := src[2].ToF32x2()
	p3 := src[3].ToF32x2()
	tt := wide.Splat(t.Get())

	ab := interp(p0, p1, tt)
	bc := interp(p1, p2, tt)
	cd := interp(p2, p3, tt)
	abc := interp(ab, bc, tt)
	bcd := interp(bc, cd, tt)
	abcd := interp(abc, bcd, tt)

	dst[0] = NewPointFromF32x2(p0)
	dst[1] = NewPointFromF32x2(ab)
	dst[2] = NewPointFromF32x2(abc)
	dst[3] = NewPointFromF32x2(abcd)
	dst[4] = NewPointFromF32x2(bcd)
	dst[5] = NewPointFromF32x2(cd)
	dst[6] = NewPointFromF32x2(p3)
}

// Quad'(t) = At + B, where
// A = 2(a - 2b + c)
// B = 2(b - a)
// Solve for t, only if it fits between 0 < t < 1
func findQuadExtrema(a, b, c float32) (normalized.NormalizedF32Exclusive, bool) {
	// At + B == 0
	// t = -B / A
	return ValidUnitDivide(a-b, a-b-b+c)
}

func ValidUnitDivide(numer, denom float32) (normalized.NormalizedF32Exclusive, bool) {
	if numer < 0.0 {
		numer = -numer
		denom = -denom
	}

	if denom == 0.0 || numer == 0.0 || numer >= denom {
		return 0, false
	}

	r := numer / denom
	return normalized.NewNormalizedF32Exclusive(r)
}

func interp(v0, v1, t wide.F32x2) wide.F32x2 {
	return v0.Add(v1.Sub(v0).Mul(t))
}

func times2(value wide.F32x2) wide.F32x2 {
	return value.Add(value)
}

// F(t)    = a (1 - t) ^ 2 + 2 b t (1 - t) + c t ^ 2
// F'(t)   = 2 (b - a) + 2 (a - 2b + c) t
// F”(t)  = 2 (a - 2b + c)
//
// A = 2 (b - a)
// B = 2 (a - 2b + c)
//
// Maximum curvature for a quadratic means solving
// Fx' Fx” + Fy' Fy” = 0
//
// t = - (Ax Bx + Ay By) / (Bx ^ 2 + By ^ 2)
func findQuadMaxCurvature(src [3]Point) normalized.NormalizedF32 {
	ax := src[1].X - src[0].X
	ay := src[1].Y - src[0].Y
	bx := src[0].X - src[1].X - src[1].X + src[2].X
	by := src[0].Y - src[1].Y - src[1].Y + src[2].Y

	numer := -(ax*bx + ay*by)
	denom := bx*bx + by*by
	if denom < 0.0 {
		numer = -numer
		denom = -denom
	}

	if numer <= 0.0 {
		return 0.0
	}

	if numer >= denom {
		return 1.0
	}

	v, _ := normalized.NewNormalizedF32(numer / denom)
	return v
}

func evalQuadAt(src [3]Point, t normalized.NormalizedF32) Point {
	return NewPointFromF32x2(NewQuadCoeffFromPoints(src).Eval(wide.Splat(t.Get())))
}

func evalQuadTangentAt(src [3]Point, tol normalized.NormalizedF32) Point {
	// The derivative equation is 2(b - a +(a - 2b +c)t). This returns a
	// zero tangent vector when t is 0 or 1, and the control point is equal
	// to the end point. In this case, use the quad end points to compute the tangent.
	if (tol == 0.0 && src[0] == src[1]) ||
		(tol == 1.0 && src[1] == src[2]) {
		return src[2].Sub(src[0])
	}

	p0 := src[0].ToF32x2()
	p1 := src[1].ToF32x2()
	p2 := src[2].ToF32x2()

	b := p1.Sub(p0)
	a := p2.Sub(p1).Sub(b)
	t := a.Mul(wide.Splat(tol.Get())).Add(b)

	return NewPointFromF32x2(times2(t))
}

// FindCubicMaxCurvature
// Looking for F' dot F” == 0
//
// A = b - a
// B = c - 2b + a
// C = d - 3c + 3b - a
//
// F' = 3Ct^2 + 6Bt + 3A
// F” = 6Ct + 6B
//
// F' dot F” -> CCt^3 + 3BCt^2 + (2BB + CA)t + AB
// Note: use float32 in place of the original NormalizedF32 type.
func FindCubicMaxCurvature(src [4]Point, tValues *[3]normalized.NormalizedF32) []normalized.NormalizedF32 {
	coeffX := formulateF1DotF2([4]float32{src[0].X, src[1].X, src[2].X, src[3].X})
	coeffY := formulateF1DotF2([4]float32{src[0].Y, src[1].Y, src[2].Y, src[3].Y})

	for i := 0; i < 4; i++ {
		coeffX[i] += coeffY[i]
	}

	length := solveCubicPoly(coeffX, tValues)
	return tValues[:length]
}

// Looking for F' dot F” == 0
//
// A = b - a
// B = c - 2b + a
// C = d - 3c + 3b - a
//
// F' = 3Ct^2 + 6Bt + 3A
// F” = 6Ct + 6B
//
// F' dot F” -> CCt^3 + 3BCt^2 + (2BB + CA)t + AB
func formulateF1DotF2(src [4]float32) [4]float32 {
	a := src[1] - src[0]
	b := src[2] - 2.0*src[1] + src[0]
	c := src[3] + 3.0*(src[1]-src[2]) - src[0]

	return [4]float32{c * c, 3.0 * b * c, 2.0*b*b + c*a, a * b}
}

// Solve coeff(t) == 0, returning the number of roots that lie within 0 < t < 1.
// coeff[0]t^3 + coeff[1]t^2 + coeff[2]t + coeff[3]
//
// Eliminates repeated roots (so that all t_values are distinct, and are always
// in increasing order.
func solveCubicPoly(coeff [4]float32, tValues *[3]normalized.NormalizedF32) int {
	if scalar.IsNearlyZero(coeff[0]) {
		// we're just a quadratic
		tmpT := [3]normalized.NormalizedF32Exclusive{}
		count := FindUnitQuadRoots(coeff[1], coeff[2], coeff[3], &tmpT)
		for i := 0; i < count; i++ {
			tValues[i] = tmpT[i].ToNormalized()
		}
		return count
	}

	inva := 1.0 / coeff[0]
	a := coeff[1] * inva
	b := coeff[2] * inva
	c := coeff[3] * inva

	q := (a*a - b*3.0) / 9.0
	r := (2.0*a*a*a - 9.0*a*b + 27.0*c) / 54.0

	q3 := q * q * q
	r2MinusQ3 := r*r - q3
	adiv3 := a / 3.0

	if r2MinusQ3 < 0.0 {
		// we have 3 real roots
		// the divide/root can, due to finite precisions, be slightly outside of -1...1
		theta := math32.Acos(scalar.Bound(r/math32.Sqrt(q3), -1.0, 1.0))
		neg2RootQ := -2.0 * math32.Sqrt(q)

		tValues[0] = normalized.NewNormalizedF32WithClamped(neg2RootQ*math32.Cos(theta/3.0) - adiv3)
		tValues[1] = normalized.NewNormalizedF32WithClamped(neg2RootQ*math32.Cos((theta+2.0*math32.Pi)/3.0) - adiv3)
		tValues[2] = normalized.NewNormalizedF32WithClamped(neg2RootQ*math32.Cos((theta-2.0*math32.Pi)/3.0) - adiv3)

		// now sort the roots
		sortArray3(tValues)
		return collapseDuplicates3(tValues)
	} else {
		// we have 1 real root
		a := math32.Abs(r) + math32.Sqrt(r2MinusQ3)
		a = scalarCubeRoot(a)
		if r > 0.0 {
			a = -a
		}

		if a != 0.0 {
			a += q / a
		}

		tValues[0] = normalized.NewNormalizedF32WithClamped(a - adiv3)
		return 1
	}
}

func sortArray3(array *[3]normalized.NormalizedF32) {
	if array[0] > array[1] {
		array[0], array[1] = array[1], array[0]
	}
	if array[1] > array[2] {
		array[1], array[2] = array[2], array[1]
	}
	if array[0] > array[1] {
		array[0], array[1] = array[1], array[0]
	}
}

func collapseDuplicates3(array *[3]normalized.NormalizedF32) int {
	length := 3
	if array[1] == array[2] {
		length = 2
	}
	if array[0] == array[1] {
		length = 1
	}
	return length
}

func scalarCubeRoot(x float32) float32 {
	return math32.Pow(x, 0.3333333)
}

// This is SkEvalCubicAt split into three functions.
func evalCubicPosAt(src [4]Point, t normalized.NormalizedF32) Point {
	return NewPointFromF32x2(NewCubicCoeffFromPoints(src).Eval(wide.Splat(t.Get())))
}

// This is SkEvalCubicAt split into three functions.
func evalCubicTangentAt(src [4]Point, t normalized.NormalizedF32) Point {
	// The derivative equation returns a zero tangent vector when t is 0 or 1, and the
	// adjacent control point is equal to the end point. In this case, use the
	// next control point or the end points to compute the tangent.
	if (t == 0.0 && src[0] == src[1]) || (t == 1.0 && src[2] == src[3]) {
		var tangent Point
		if t == 0.0 {
			tangent = src[2].Sub(src[0])
		} else {
			tangent = src[3].Sub(src[1])
		}

		if tangent.X == 0.0 && tangent.Y == 0.0 {
			tangent = src[3].Sub(src[0])
		}
		return tangent
	}
	return evalCubicDerivative(src, t)
}

func evalCubicDerivative(src [4]Point, t normalized.NormalizedF32) Point {
	p0 := src[0].ToF32x2()
	p1 := src[1].ToF32x2()
	p2 := src[2].ToF32x2()
	p3 := src[3].ToF32x2()

	coeff := QuadCoeff{
		A: p3.Add(wide.Splat(3.0).Mul(p1.Sub(p2))).Sub(p0),
		B: times2(p2.Sub(times2(p1)).Add(p0)),
		C: p1.Sub(p0),
	}

	return NewPointFromF32x2(coeff.Eval(wide.Splat(t.Get())))
}

// Cubic'(t) = At^2 + Bt + C, where
// A = 3(-a + 3(b - c) + d)
// B = 6(a - 2b + c)
// C = 3(b - a)
// Solve for t, keeping only those that fit between 0 < t < 1
func findCubicExtrema(a, b, c, d float32, tValues *[3]normalized.NormalizedF32Exclusive) int {
	// we divide A,B,C by 3 to simplify
	aa := d - a + 3.0*(b-c)
	bb := 2.0 * (a - b - b + c)
	cc := b - a

	return FindUnitQuadRoots(aa, bb, cc, tValues)
}

// http://www.faculty.idc.ac.il/arik/quality/appendixA.html
//
// Inflection means that curvature is zero.
// Curvature is [F' x F”] / [F'^3]
// So we solve F'x X F”y - F'y X F”y == 0
// After some canceling of the cubic term, we get
// A = b - a
// B = c - 2b + a
// C = d - 3c + 3b - a
// (BxCy - ByCx)t^2 + (AxCy - AyCx)t + AxBy - AyBx == 0
func findCubicInflections(src [4]Point, tValues *[3]normalized.NormalizedF32Exclusive) []normalized.NormalizedF32Exclusive {
	ax := src[1].X - src[0].X
	ay := src[1].Y - src[0].Y
	bx := src[2].X - 2.0*src[1].X + src[0].X
	by := src[2].Y - 2.0*src[1].Y + src[0].Y
	cx := src[3].X + 3.0*(src[1].X-src[2].X) - src[0].X
	cy := src[3].Y + 3.0*(src[1].Y-src[2].Y) - src[0].Y

	length := FindUnitQuadRoots(
		float32(bx*cy-by*cx),
		float32(ax*cy-ay*cx),
		float32(ax*by-ay*bx),
		tValues,
	)

	return tValues[:length]
}

// Return location (in t) of cubic cusp, if there is one.
// Note that classify cubic code does not reliably return all cusp'd cubics, so
// it is not called here.
func findCubicCusp(src [4]Point) (normalized.NormalizedF32Exclusive, bool) {
	// When the adjacent control point matches the end point, it behaves as if
	// the cubic has a cusp: there's a point of max curvature where the derivative
	// goes to zero. Ideally, this would be where t is zero or one, but math32
	// error makes not so. It is not uncommon to create cubics this way; skip them.
	if src[0] == src[1] || src[2] == src[3] {
		return 0, false
	}

	// Cubics only have a cusp if the line segments formed by the control and end points cross.
	// Detect crossing if line ends are on opposite sides of plane formed by the other line.
	if onSameSide(src, 0, 2) || onSameSide(src, 2, 0) {
		return 0, false
	}

	// Cubics may have multiple points of maximum curvature, although at most only
	// one is a cusp.
	tValues := [3]normalized.NormalizedF32{}
	maxCrvature := FindCubicMaxCurvature(src, &tValues)
	for _, testT := range maxCrvature {
		if testT <= 0.0 || testT >= 1.0 {
			// no need to consider max curvature on the end
			continue
		}

		// A cusp is at the max curvature, and also has a derivative close to zero.
		// Choose the 'close to zero' meaning by comparing the derivative length
		// with the overall cubic size.
		dPt := evalCubicDerivative(src, testT)
		dPtMagnitude := dPt.LengthSqd()
		precision := calcCubicPrecision(src)
		if dPtMagnitude < precision {
			// All three max curvature t values may be close to the cusp;
			// return the first one.
			return normalized.NewNormalizedF32ExclusiveWithBounded(testT.Get()), true
		}
	}

	return 0, false
}

// Returns true if both points src[testIndex], src[testIndex+1] are in the same half plane defined
// by the line segment src[lineIndex], src[lineIndex+1].
func onSameSide(src [4]Point, testIndex, lineIndex int) bool {
	origin := src[lineIndex]
	line := src[lineIndex+1].Sub(origin)
	var crosses [2]float32
	for i := 0; i < 2; i++ {
		testLine := src[testIndex+i].Sub(origin)
		crosses[i] = line.Cross(testLine)
	}
	return crosses[0]*crosses[1] >= 0.0
}

// Returns a constant proportional to the dimensions of the cubic.
// Constant found through experimentation -- maybe there's a better way....
func calcCubicPrecision(src [4]Point) float32 {
	return (src[1].DistanceToSqd(src[0]) +
		src[2].DistanceToSqd(src[1]) +
		src[3].DistanceToSqd(src[2])) * 1e-8
}

type conic struct {
	Points [3]Point
	Weight float32
}

func NewConic(pt0, pt1, pt2 Point, weight float32) conic {
	return conic{
		Points: [3]Point{pt0, pt1, pt2},
		Weight: weight,
	}
}

func NewConicFromPoints(points [3]Point, weight float32) conic {
	return conic{
		Points: points,
		Weight: weight,
	}
}

func (c conic) computeQuadPow2(tolerance float32) (int, bool) {
	if tolerance < 0.0 || math32.IsInf(tolerance, 0) || math32.IsNaN(tolerance) {
		return 0, false
	}

	if !c.Points[0].IsFinite() || !c.Points[1].IsFinite() || !c.Points[2].IsFinite() {
		return 0, false
	}

	// Limit the number of suggested quads to approximate a conic
	const MAX_CONIC_TO_QUAD_POW2 = 4

	// "High order approximation of conic sections by quadratic splines"
	// by Michael Floater, 1993
	a := c.Weight - 1.0
	k := a / (4.0 * (2.0 + a))
	x := k * (c.Points[0].X - 2.0*c.Points[1].X + c.Points[2].X)
	y := k * (c.Points[0].Y - 2.0*c.Points[1].Y + c.Points[2].Y)

	errorVal := math32.Sqrt(x*x + y*y)
	var pow2 = 0
	for i := 0; i < MAX_CONIC_TO_QUAD_POW2; i++ {
		if errorVal <= tolerance {
			break
		}
		errorVal *= 0.25
		pow2 += 1
	}

	// Unlike Skia, we always expect `pow2` to be at least 1.
	// Otherwise it produces ugly results.
	if pow2 < 1 {
		pow2 = 1
	}
	return pow2, true
}

// ChopIntoQuadsPow2
// Chop this conic into N quads, stored continuously in pts[], where
// N = 1 << pow2. The amount of storage needed is (1 + 2 * N)
func (c conic) ChopIntoQuadsPow2(pow2 int, points *[64]Point) int {
	if pow2 >= 5 {
		panic("pow2 too large")
	}

	points[0] = c.Points[0]
	subdivide(c, points[1:], pow2)

	quadCount := 1 << pow2
	ptCount := 2*quadCount + 1

	isAnyNonFinite := false
	for i := 0; i < ptCount; i++ {
		if !points[i].IsFinite() {
			isAnyNonFinite = true
			break
		}
	}

	if isAnyNonFinite {
		// if we generated a non-finite, pin ourselves to the middle of the hull,
		// as our first and last are already on the first/last pts of the hull.
		for i := 1; i < ptCount-1; i++ {
			points[i] = c.Points[1]
		}
	}

	return 1 << pow2
}

func (c conic) chop() (conic, conic) {
	sv := 1.0 / (1.0 + c.Weight)
	scale := wide.Splat(sv)
	newW := subdivideWeightValue(c.Weight)

	p0 := c.Points[0].ToF32x2()
	p1 := c.Points[1].ToF32x2()
	p2 := c.Points[2].ToF32x2()
	ww := wide.Splat(c.Weight)

	wp1 := ww.Mul(p1)
	// (p0 + times_2(wp1) + p2) * scale * f32x2{0.5, 0.5}
	m := (p0.Add(times2(wp1)).Add(p2)).Mul(scale).Mul(wide.F32x2{0.5, 0.5})
	mPt := NewPointFromF32x2(m)

	if !mPt.IsFinite() {
		wD := float64(c.Weight)
		w2 := wD * 2.0
		scaleHalf := 1.0 / (1.0 + wD) * 0.5
		mPt.X = float32((float64(c.Points[0].X) +
			w2*float64(c.Points[1].X) +
			float64(c.Points[2].X)) *
			scaleHalf)
		mPt.Y = float32((float64(c.Points[0].Y) +
			w2*float64(c.Points[1].Y) +
			float64(c.Points[2].Y)) *
			scaleHalf)
	}

	return conic{
			Points: [3]Point{c.Points[0], NewPointFromF32x2(p0.Add(wp1).Mul(scale)), mPt},
			Weight: newW,
		}, conic{
			Points: [3]Point{mPt, NewPointFromF32x2(wp1.Add(p2).Mul(scale)), c.Points[2]},
			Weight: newW,
		}
}

func BuildUnitArc(uStart, uStop Point, dir pathDirection, userTransform Transform, dst *[5]conic) []conic {
	// rotate by x,y so that u_start is (1.0)
	x := uStart.Dot(uStop)
	y := uStart.Cross(uStop)
	absY := math32.Abs(y)

	// check for (effectively) coincident vectors
	// this can happen if our angle is nearly 0 or nearly 180 (y == 0)
	// ... we use the dot-prod to distinguish between 0 and 180 (x > 0)
	if absY <= scalar.SCALAR_NEARLY_ZERO && x > 0.0 &&
		((y >= 0.0 && dir == pathDirectionCW) || (y <= 0.0 && dir == pathDirectionCCW)) {
		return nil
	}

	if dir == pathDirectionCCW {
		y = -y
	}

	// We decide to use 1-conic per quadrant of a circle. What quadrant does [xy] lie in?
	//      0 == [0  .. 90)
	//      1 == [90 ..180)
	//      2 == [180..270)
	//      3 == [270..360)
	//
	quadrant := 0
	if y == 0.0 {
		quadrant = 2 // 180
	} else if x == 0.0 {
		if y > 0.0 {
			quadrant = 1
		} else {
			quadrant = 3
		}
	} else {
		if y < 0.0 {
			quadrant += 2
		}
		if (x < 0.0) != (y < 0.0) {
			quadrant += 1
		}
	}

	quadrantPoints := []Point{
		{X: 1.0, Y: 0.0},
		{X: 1.0, Y: 1.0},
		{X: 0.0, Y: 1.0},
		{X: -1.0, Y: 1.0},
		{X: -1.0, Y: 0.0},
		{X: -1.0, Y: -1.0},
		{X: 0.0, Y: -1.0},
		{X: 1.0, Y: -1.0},
	}

	const quadrantWeight = scalar.SCALAR_ROOT_2_OVER_2

	conicCount := quadrant
	for i := 0; i < conicCount; i++ {
		dst[i] = NewConicFromPoints([3]Point{quadrantPoints[i*2], quadrantPoints[i*2+1], quadrantPoints[i*2+2]}, quadrantWeight)
	}

	// Now compute any remaining (sub-90-degree) arc for the last conic
	finalPt := Point{X: x, Y: y}
	lastQ := quadrantPoints[quadrant*2]
	dot := lastQ.Dot(finalPt)

	if dot < 1.0 {
		offCurve := Point{X: lastQ.X + x, Y: lastQ.Y + y}
		// compute the bisector vector, and then rescale to be the off-curve point.
		// we compute its length from cos(theta/2) = length / 1, using half-angle identity we get
		// length = sqrt(2 / (1 + cos(theta)). We already have cos() when to computed the dot.
		// This is nice, since our computed weight is cos(theta/2) as well!
		cosThetaOver2 := math32.Sqrt((1.0 + dot) / 2.0)

		offCurve, _ = offCurve.WithLengthFrom(1.0 / cosThetaOver2)
		if !lastQ.AlmostEqual(offCurve) {
			dst[conicCount] = NewConic(lastQ, offCurve, finalPt, cosThetaOver2)
			conicCount += 1
		}
	}

	// now handle counter-clockwise and the initial unitStart rotation
	transform := NewTransformFromSinCos(uStart.Y, uStart.X)
	if dir == pathDirectionCCW {
		transform = transform.PreScale(1.0, -1.0)
	}
	transform = transform.PostConcat(userTransform)

	for i := 0; i < conicCount; i++ {
		transform.MapPoints(dst[i].Points[:])
	}

	if conicCount == 0 {
		return nil
	}
	return dst[0:conicCount]
}

func subdivideWeightValue(w float32) float32 {
	return math32.Sqrt(0.5 + w*0.5)
}

func subdivide(src conic, points []Point, level int) []Point {
	if level == 0 {
		points[0] = src.Points[1]
		points[1] = src.Points[2]
		return points[2:]
	} else {
		dst0, dst1 := src.chop()
		startY := src.Points[0].Y
		endY := src.Points[2].Y

		if between(startY, src.Points[1].Y, endY) {
			midY := dst0.Points[2].Y
			if !between(startY, midY, endY) {
				// If the computed midpoint is outside the ends, move it to the closer one.
				closerY := startY
				if math32.Abs(midY-startY) >= math32.Abs(midY-endY) {
					closerY = endY
				}
				dst0.Points[2].Y = closerY
				dst1.Points[0].Y = closerY
			}

			if !between(startY, dst0.Points[1].Y, dst0.Points[2].Y) {
				// If the 1st control is not between the start and end, put it at the start.
				// This also reduces the quad to a line.
				dst0.Points[1].Y = startY
			}

			if !between(dst1.Points[0].Y, dst1.Points[1].Y, endY) {
				// If the 2nd control is not between the start and end, put it at the end.
				// This also reduces the quad to a line.
				dst1.Points[1].Y = endY
			}

			// Verify that all five points are in order.
		}

		level -= 1
		points = subdivide(dst0, points, level)
		return subdivide(dst1, points, level)
	}
}

// This was originally developed and tested for pathops: see SkOpTypes.h
// returns true if (a <= b <= c) || (a >= b >= c)
func between(a, b, c float32) bool {
	return (a-b)*(c-b) <= 0.0
}

type autoConicToQuads struct {
	Points [64]Point
	Len    int // the number of quads
}

func ComputeAutoConicToQuads(pt0, pt1, pt2 Point, weight float32) (autoConicToQuads, bool) {
	conic := NewConic(pt0, pt1, pt2, weight)
	pow2, ok := conic.computeQuadPow2(0.25)
	if !ok {
		return autoConicToQuads{}, false
	}
	points := [64]Point{}
	count := conic.ChopIntoQuadsPow2(pow2, &points)
	return autoConicToQuads{
		Points: points,
		Len:    count,
	}, true
}

// Returns 0 for 1 quad, and 1 for two quads, either way the answer is stored in dst[].
//
// Guarantees that the 1/2 quads will be monotonic.
func ChopQuadAtXExtrema(src [3]Point, dst *[5]Point) int {
	a := src[0].X
	b := src[1].X
	c := src[2].X

	if isNotMonotonic(a, b, c) {
		if tValue, ok := ValidUnitDivide(a-b, a-b-b+c); ok {
			var tmpDst [5]Point
			ChopQuadAt(src, tValue, &tmpDst)

			// flatten double quad extrema
			dst[0] = tmpDst[0]
			dst[1] = tmpDst[1]
			dst[2] = tmpDst[2]
			dst[3] = tmpDst[3]
			dst[4] = tmpDst[4]

			return 1
		}

		// if we get here, we need to force dst to be monotonic, even though
		// we couldn't compute a unit_divide value (probably underflow).
		if math32.Abs(a-b) < math32.Abs(b-c) {
			b = a
		} else {
			b = c
		}
	}

	dst[0] = Point{X: a, Y: src[0].Y}
	dst[1] = Point{X: b, Y: src[1].Y}
	dst[2] = Point{X: c, Y: src[2].Y}
	return 0
}

// Returns 0 for 1 quad, and 1 for two quads, either way the answer is stored in dst[].
//
// Guarantees that the 1/2 quads will be monotonic.
func ChopQuadAtYExtrema(src [3]Point, dst *[5]Point) int {
	a := src[0].Y
	b := src[1].Y
	c := src[2].Y

	if isNotMonotonic(a, b, c) {
		if tValue, ok := ValidUnitDivide(a-b, a-b-b+c); ok {
			var tmpDst [5]Point
			ChopQuadAt(src, tValue, &tmpDst)

			// flatten double quad extrema
			dst[0] = tmpDst[0]
			dst[1] = tmpDst[1]
			dst[2] = tmpDst[2]
			dst[3] = tmpDst[3]
			dst[4] = tmpDst[4]

			return 1
		}

		// if we get here, we need to force dst to be monotonic, even though
		// we couldn't compute a unit_divide value (probably underflow).
		if math32.Abs(a-b) < math32.Abs(b-c) {
			b = a
		} else {
			b = c
		}
	}

	dst[0] = Point{X: src[0].X, Y: a}
	dst[1] = Point{X: src[1].X, Y: b}
	dst[2] = Point{X: src[2].X, Y: c}
	return 0
}

func isNotMonotonic(a, b, c float32) bool {
	ab := a - b
	bc := b - c
	if ab < 0.0 {
		bc = -bc
	}

	return ab == 0.0 || bc < 0.0
}

func ChopCubicAtXExtrema(src [4]Point, dst *[10]Point) int {
	var tValues [3]normalized.NormalizedF32Exclusive
	roots := findCubicExtrema(src[0].X, src[1].X, src[2].X, src[3].X, &tValues)
	actualT := tValues[:roots]

	ChopCubicAt(src, actualT, dst[:])
	if len(actualT) > 0 {
		// we do some cleanup to ensure our X extrema are flat
		dst[2].X = dst[3].X
		dst[4].X = dst[3].X
		if len(actualT) == 2 {
			dst[5].X = dst[6].X
			dst[7].X = dst[6].X
		}
	}

	return len(actualT)
}

// Given 4 points on a cubic bezier, chop it into 1, 2, 3 beziers such that
// the resulting beziers are monotonic in Y.
//
// This is called by the scan converter.
//
// Depending on what is returned, dst[] is treated as follows:
//
// - 0: dst[0..3] is the original cubic
// - 1: dst[0..3] and dst[3..6] are the two new cubics
// - 2: dst[0..3], dst[3..6], dst[6..9] are the three new cubics
func ChopCubicAtYExtrema(src [4]Point, dst *[10]Point) int {
	var tValues [3]normalized.NormalizedF32Exclusive
	roots := findCubicExtrema(src[0].Y, src[1].Y, src[2].Y, src[3].Y, &tValues)
	actualT := tValues[:roots]

	ChopCubicAt(src, actualT, dst[:])
	if len(actualT) > 0 {
		// we do some cleanup to ensure our Y extrema are flat
		dst[2].Y = dst[3].Y
		dst[4].Y = dst[3].Y
		if len(actualT) == 2 {
			dst[5].Y = dst[6].Y
			dst[7].Y = dst[6].Y
		}
	}

	return len(actualT)
}

// Cubic'(t) = At^2 + Bt + C, where
// A = 3(-a + 3(b - c) + d)
// B = 6(a - 2b + c)
// C = 3(b - a)
// Solve for t, keeping only those that fit between 0 < t < 1
// Note: This is a duplicate of the float32 version above, kept for compatibility.
// Deprecated: Use the float32 version at line 422 instead.
func findCubicExtremaFloat64(a, b, c, d float64, tValues *[3]normalized.NormalizedF32Exclusive) int {
	// we divide A,B,C by 3 to simplify
	aa := d - a + 3.0*(b-c)
	bb := 2.0 * (a - b - b + c)
	cc := b - a

	return FindUnitQuadRoots(float32(aa), float32(bb), float32(cc), tValues)
}

func ChopCubicAt(src [4]Point, tValues []normalized.NormalizedF32Exclusive, dst []Point) {
	if len(tValues) == 0 {
		// nothing to chop
		copy(dst[:4], src[:4])
	} else {
		t := tValues[0]
		tmp := [4]Point{}
		currentSrc := src

		dstOffset := 0
		for i := 0; i < len(tValues); i++ {
			var tmpDst [7]Point
			ChopCubicAt2(currentSrc, t, &tmpDst)

			// Copy the result to dst
			for j := 0; j < 7 && dstOffset+j < len(dst); j++ {
				dst[dstOffset+j] = tmpDst[j]
			}
			if i == len(tValues)-1 {
				break
			}

			dstOffset += 3
			// have src point to the remaining cubic (after the chop)
			tmp[0] = dst[dstOffset+0]
			tmp[1] = dst[dstOffset+1]
			tmp[2] = dst[dstOffset+2]
			tmp[3] = dst[dstOffset+3]
			currentSrc = tmp

			// watch out in case the renormalized t isn't in range
			if n, ok := ValidUnitDivide(float32(tValues[i+1].Get()-tValues[i].Get()), float32(1.0-tValues[i].Get())); ok {
				t = n
			} else {
				// if we can't, just create a degenerate cubic
				dst[dstOffset+4] = currentSrc[3]
				dst[dstOffset+5] = currentSrc[3]
				dst[dstOffset+6] = currentSrc[3]
				break
			}
		}
	}
}

func ChopCubicAtMaxCurvature(src [4]Point, tValues *[3]normalized.NormalizedF32, dst []Point) int {
	var roots [3]normalized.NormalizedF32
	foundRoots := FindCubicMaxCurvature(src, &roots)

	// Throw out values not inside 0..1.
	actualCount := 0
	for i := 0; i < len(foundRoots); i++ {
		root := foundRoots[i]
		if root.Get() > 0.0 && root.Get() < 1.0 {
			tValues[actualCount] = root
			actualCount++
		}
	}

	if actualCount == 0 {
		copy(dst[:4], src[:4])
	} else {
		// Convert normalized.NormalizedF32 to normalized.NormalizedF32Exclusive for ChopCubicAt
		var exclusiveTValues []normalized.NormalizedF32Exclusive
		for i := 0; i < actualCount; i++ {
			if excl, ok := normalized.NewNormalizedF32Exclusive(tValues[i].Get()); ok {
				exclusiveTValues = append(exclusiveTValues, excl)
			}
		}
		ChopCubicAt(src, exclusiveTValues, dst)
	}

	return actualCount + 1
}
