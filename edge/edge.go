// Copyright 2006 The Android Open Source Project
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package edge

import (
	"math/bits"

	"github.com/lumifloat/tinyskia/internal/fixed"
	"github.com/lumifloat/tinyskia/path"
)

// We store 1<<shift in a (signed) byte, so its maximum value is 1<<6 == 64.
//
// Note that this limits the number of lines we use to approximate a curve.
// If we need to increase this, we need to store curve_count in something
// larger than i8.
const maxCoeffShift = 6

// EdgeType represents the type of edge.
type EdgeType int

const (
	EdgeTypeLine EdgeType = iota
	EdgeTypeQuadratic
	EdgeTypeCubic
)

// Edge represents a path edge for rasterization.
type Edge struct {
	Type      EdgeType
	Line      *LineEdge
	Quadratic *QuadraticEdge
	Cubic     *CubicEdge
}

// AsLine returns the underlying LineEdge.
func (e *Edge) AsLine() *LineEdge {
	switch e.Type {
	case EdgeTypeLine:
		return e.Line
	case EdgeTypeQuadratic:
		return e.Quadratic.Line
	case EdgeTypeCubic:
		return e.Cubic.Line
	default:
		return nil
	}
}

// Update updates the edge and returns true if successful.
func (e *Edge) Update() bool {
	switch e.Type {
	case EdgeTypeLine:
		return false // Line edges don't have update logic
	case EdgeTypeQuadratic:
		return e.Quadratic.Update()
	case EdgeTypeCubic:
		return e.Cubic.Update()
	default:
		return false
	}
}

// LineEdge represents a line edge for rasterization.
type LineEdge struct {
	Prev    int // -1 means none
	Next    int // -1 means none
	X       fixed.FDot16
	DX      fixed.FDot16
	FirstY  int32
	LastY   int32
	Winding int8 // 1 or -1
}

// NewLineEdge creates a new LineEdge from two points.
func NewLineEdge(p0, p1 path.Point, shift int32) *LineEdge {
	scale := float32(int32(1) << (shift + 6))
	x0 := int32(p0.X * scale)
	y0 := int32(p0.Y * scale)
	x1 := int32(p1.X * scale)
	y1 := int32(p1.Y * scale)

	winding := int8(1)

	if y0 > y1 {
		x0, x1 = x1, x0
		y0, y1 = y1, y0
		winding = -1
	}

	top := fixed.FDot6Round(fixed.FDot6(y0))
	bottom := fixed.FDot6Round(fixed.FDot6(y1))

	// are we a zero-height line?
	if top == bottom {
		return nil
	}

	slope := fixed.FDot6DivToFDot16(fixed.FDot6(x1-x0), fixed.FDot6(y1-y0))
	dy := computeDy(top, y0)

	return &LineEdge{
		Next:    -1,
		Prev:    -1,
		X:       fixed.FDot6ToFDot16(fixed.FDot6(x0 + fixed.FDot16Mul(slope, fixed.FDot16(dy)))),
		DX:      slope,
		FirstY:  int32(top),
		LastY:   int32(bottom) - 1,
		Winding: winding,
	}
}

// IsVertical returns true if the edge is vertical.
func (e *LineEdge) IsVertical() bool {
	return e.DX == 0
}

func (e *LineEdge) update(x0, y0, x1, y1 fixed.FDot16) bool {
	y0 >>= 10
	y1 >>= 10

	top := fixed.FDot6Round(fixed.FDot6(y0))
	bottom := fixed.FDot6Round(fixed.FDot6(y1))

	// are we a zero-height line?
	if top == bottom {
		return false
	}

	x0 >>= 10
	x1 >>= 10

	slope := fixed.FDot6DivToFDot16(fixed.FDot6(x1-x0), fixed.FDot6(y1-y0))
	dy := computeDy(top, y0)

	e.X = fixed.FDot6ToFDot16(fixed.FDot6(x0) + fixed.FDot16Mul(slope, fixed.FDot16(dy)))
	e.DX = slope
	e.FirstY = int32(top)
	e.LastY = int32(bottom) - 1

	return true
}

// QuadraticEdge represents a quadratic Bezier edge.
type QuadraticEdge struct {
	Line       *LineEdge
	CurveCount int8
	curveShift uint8 // applied to all dx/ddx/dddx
	qx         fixed.FDot16
	qy         fixed.FDot16
	qdx        fixed.FDot16
	qdy        fixed.FDot16
	qddx       fixed.FDot16
	qddy       fixed.FDot16
	qLastX     fixed.FDot16
	qLastY     fixed.FDot16
}

// NewQuadraticEdge creates a new QuadraticEdge from control points.
func NewQuadraticEdge(points [3]path.Point, shift int32) *QuadraticEdge {
	quad := newQuadraticEdge2(points, shift)
	if quad == nil {
		return nil
	}
	// Match Rust behavior: call Update() and return nil if it fails
	if quad.Update() {
		return quad
	}
	return nil
}

func newQuadraticEdge2(points [3]path.Point, shift int32) *QuadraticEdge {
	scale := float32(int32(1) << (shift + 6))
	x0 := int32(points[0].X * scale)
	y0 := int32(points[0].Y * scale)
	x1 := int32(points[1].X * scale)
	y1 := int32(points[1].Y * scale)
	x2 := int32(points[2].X * scale)
	y2 := int32(points[2].Y * scale)

	winding := int8(1)
	if y0 > y2 {
		x0, x2 = x2, x0
		y0, y2 = y2, y0
		winding = -1
	}

	// This ensures the control point is monotonic in Y after sorting

	top := fixed.FDot6Round(fixed.FDot6(y0))
	bottom := fixed.FDot6Round(fixed.FDot6(y2))

	// are we a zero-height quad (line)?
	if top == bottom {
		return nil
	}

	// compute number of steps needed (1 << shift)
	{
		dx := fixed.FDot6(((int32(uint32(x1)<<1) - x0 - x2) >> 2))
		dy := fixed.FDot6(((int32(uint32(y1)<<1) - y0 - y2) >> 2))
		// This is a little confusing:
		// before this line, shift is the scale up factor for AA;
		// after this line, shift is the fCurveShift.
		shift = diffToShift(dx, dy, shift)
	}

	// need at least 1 subdivision for our bias trick
	if shift == 0 {
		shift = 1
	} else if shift > maxCoeffShift {
		shift = maxCoeffShift
	}

	curveCount := int8(int32(1) << shift)

	// We want to reformulate into polynomial form, to make it clear how we
	// should forward-difference.
	//
	// p0 (1 - t)^2 + p1 t(1 - t) + p2 t^2 ==> At^2 + Bt + C
	//
	// A = p0 - 2p1 + p2
	// B = 2(p1 - p0)
	// C = p0
	//
	// Our caller must have constrained our inputs (p0..p2) to all fit into
	// 16.16. However, as seen above, we sometimes compute values that can be
	// larger (e.g. B = 2*(p1 - p0)). To guard against overflow, we will store
	// A and B at 1/2 of their actual value, and just apply a 2x scale during
	// application in updateQuadratic(). Hence we store (shift - 1) in
	// curve_shift.

	curveShift := uint8(shift - 1)

	a := fdot6ToFixedDiv2(fixed.FDot6(x0 - x1 - x1 + x2)) // 1/2 the real value
	b := fixed.FDot6ToFDot16(fixed.FDot6(x1 - x0))        // 1/2 the real value

	qx := fixed.FDot6ToFDot16(fixed.FDot6(x0))
	qdx := b + (a >> shift)  // biased by shift
	qddx := a >> (shift - 1) // biased by shift

	a = fdot6ToFixedDiv2(fixed.FDot6(y0 - y1 - y1 + y2)) // 1/2 the real value
	b = fixed.FDot6ToFDot16(fixed.FDot6(y1 - y0))        // 1/2 the real value

	qy := fixed.FDot6ToFDot16(fixed.FDot6(y0))
	qdy := b + (a >> shift)  // biased by shift
	qddy := a >> (shift - 1) // biased by shift

	qLastX := fixed.FDot16(x2 << 10)
	qLastY := fixed.FDot16(y2 << 10)

	quad := &QuadraticEdge{
		Line: &LineEdge{
			Next:    -1,
			Prev:    -1,
			X:       0,
			DX:      0,
			FirstY:  0,
			LastY:   0,
			Winding: winding,
		},
		CurveCount: curveCount,
		curveShift: curveShift,
		qx:         qx,
		qy:         qy,
		qdx:        qdx,
		qdy:        qdy,
		qddx:       qddx,
		qddy:       qddy,
		qLastX:     qLastX,
		qLastY:     qLastY,
	}
	// Return the quadratic edge without calling Update()
	// Caller should manually call Update() when needed
	return quad
}

// Update updates the quadratic edge and returns true if successful.
func (q *QuadraticEdge) Update() bool {
	var success bool
	count := q.CurveCount
	oldx := q.qx
	oldy := q.qy
	dx := q.qdx
	dy := q.qdy
	var newx, newy fixed.FDot16
	shift := q.curveShift

	// If count is invalid, return false to indicate failure
	if count <= 0 {
		return false
	}

	for {
		count--
		if count > 0 {
			newx = oldx + (dx >> shift)
			dx += q.qddx
			newy = oldy + (dy >> shift)
			dy += q.qddy
		} else {
			// last segment
			newx = q.qLastX
			newy = q.qLastY
		}
		success = q.Line.update(oldx, oldy, newx, newy)
		oldx = newx
		oldy = newy

		if count == 0 || success {
			break
		}
	}

	q.qx = newx
	q.qy = newy
	q.qdx = dx
	q.qdy = dy
	q.CurveCount = count

	return success
}

// CubicEdge represents a cubic Bezier edge.
type CubicEdge struct {
	Line       *LineEdge
	CurveCount int8
	curveShift uint8 // applied to all dx/ddx/dddx except for dshift exception
	dshift     uint8 // applied to cdx and cdy
	cx         fixed.FDot16
	cy         fixed.FDot16
	cdx        fixed.FDot16
	cdy        fixed.FDot16
	cddx       fixed.FDot16
	cddy       fixed.FDot16
	cdddx      fixed.FDot16
	cdddy      fixed.FDot16
	cLastX     fixed.FDot16
	cLastY     fixed.FDot16
}

// NewCubicEdge creates a new CubicEdge from control points.
func NewCubicEdge(points [4]path.Point, shift int32) *CubicEdge {
	// Match Rust behavior: call Update() and return nil if it fails
	cubic := newCubicEdge2(points, shift, true)
	if cubic == nil {
		return nil
	}
	if cubic.Update() {
		return cubic
	}
	return nil
}

func newCubicEdge2(points [4]path.Point, shift int32, sortY bool) *CubicEdge {
	scale := float32(int32(1) << (shift + 6))
	x0 := int32(points[0].X * scale)
	y0 := int32(points[0].Y * scale)
	x1 := int32(points[1].X * scale)
	y1 := int32(points[1].Y * scale)
	x2 := int32(points[2].X * scale)
	y2 := int32(points[2].Y * scale)
	x3 := int32(points[3].X * scale)
	y3 := int32(points[3].Y * scale)

	winding := int8(1)
	if sortY && y0 > y3 {
		x0, x3 = x3, x0
		x1, x2 = x2, x1
		y0, y3 = y3, y0
		y1, y2 = y2, y1
		winding = -1
	}

	top := fixed.FDot6Round(fixed.FDot6(y0))
	bot := fixed.FDot6Round(fixed.FDot6(y3))

	// are we a zero-height cubic (line)?
	if sortY && top == bot {
		return nil
	}

	// compute number of steps needed (1 << shift)
	{
		// Can't use (center of curve - center of baseline), since center-of-curve
		// need not be the max delta from the baseline (it could even be coincident)
		// so we try just looking at the two off-curve points
		dx := cubicDeltaFromLine(fixed.FDot6(x0), fixed.FDot6(x1), fixed.FDot6(x2), fixed.FDot6(x3))
		dy := cubicDeltaFromLine(fixed.FDot6(y0), fixed.FDot6(y1), fixed.FDot6(y2), fixed.FDot6(y3))
		// add 1 (by observation)
		shift = diffToShift(dx, dy, 2) + 1
	}
	// need at least 1 subdivision for our bias trick
	if shift <= 0 {
		panic("shift should be > 0")
	}
	if shift > maxCoeffShift {
		shift = maxCoeffShift
	}

	// Since our incoming data is initially shifted down by 10 (or 8 in
	// antialias). That means the most we can shift up is 8. However, we
	// compute coefficients with a 3*, so the safest upshift is really 6
	upShift := int32(6) // largest safe value
	downShift := shift + upShift - 10
	if downShift < 0 {
		downShift = 0
		upShift = 10 - shift
	}

	curveCount := int8(^((1 << shift) - 1))
	curveShift := uint8(shift)
	dshift := uint8(downShift)

	b := fdot6UpShift(fixed.FDot6(3*(x1-x0)), upShift)
	c := fdot6UpShift(fixed.FDot6(3*(x0-x1-x1+x2)), upShift)
	d := fdot6UpShift(fixed.FDot6(x3+3*(x1-x2)-x0), upShift)

	cx := fixed.FDot6ToFDot16(fixed.FDot6(x0))
	cdx := b + (c >> shift) + (d >> (2 * shift)) // biased by shift
	cddx := 2*c + ((3 * d) >> (shift - 1))       // biased by 2*shift
	cdddx := (3 * d) >> (shift - 1)              // biased by 2*shift

	b = fdot6UpShift(fixed.FDot6(3*(y1-y0)), upShift)
	c = fdot6UpShift(fixed.FDot6(3*(y0-y1-y1+y2)), upShift)
	d = fdot6UpShift(fixed.FDot6(y3+3*(y1-y2)-y0), upShift)

	cy := fixed.FDot6ToFDot16(fixed.FDot6(y0))
	cdy := b + (c >> shift) + (d >> (2 * shift)) // biased by shift
	cddy := 2*c + ((3 * d) >> (shift - 1))       // biased by 2*shift
	cdddy := (3 * d) >> (shift - 1)              // biased by 2*shift

	cLastX := fixed.FDot16(x3 << 10)
	cLastY := fixed.FDot16(y3 << 10)

	return &CubicEdge{
		Line: &LineEdge{
			Next:    -1,
			Prev:    -1,
			X:       0,
			DX:      0,
			FirstY:  0,
			LastY:   0,
			Winding: winding,
		},
		CurveCount: int8(curveCount),
		curveShift: curveShift,
		dshift:     dshift,
		cx:         cx,
		cy:         cy,
		cdx:        cdx,
		cdy:        cdy,
		cddx:       cddx,
		cddy:       cddy,
		cdddx:      cdddx,
		cdddy:      cdddy,
		cLastX:     cLastX,
		cLastY:     cLastY,
	}
}

// Update updates the cubic edge and returns true if successful.
func (c *CubicEdge) Update() bool {
	var success bool
	count := c.CurveCount
	oldx := c.cx
	oldy := c.cy
	var newx, newy fixed.FDot16
	ddshift := c.curveShift
	dshift := c.dshift

	// If count is invalid, return false to indicate failure
	// Note: Cubic edges use negative count that increments to 0
	if count >= 0 {
		return false
	}

	for {
		count++
		if count < 0 {
			newx = oldx + (c.cdx >> dshift)
			c.cdx += c.cddx >> ddshift
			c.cddx += c.cdddx

			newy = oldy + (c.cdy >> dshift)
			c.cdy += c.cddy >> ddshift
			c.cddy += c.cdddy
		} else {
			// last segment
			newx = c.cLastX
			newy = c.cLastY
		}

		// we want to say debug_assert(oldy <= newy), but our finite fixedpoint
		// doesn't always achieve that, so we have to explicitly pin it here.
		if newy < oldy {
			newy = oldy
		}

		success = c.Line.update(oldx, oldy, newx, newy)
		oldx = newx
		oldy = newy

		if count == 0 || success {
			break
		}
	}

	c.cx = newx
	c.cy = newy
	c.CurveCount = count

	return success
}

// computeDy computes the Y delta for line rasterization.
// This correctly favors the lower-pixel when y0 is on a 1/2 pixel boundary.
func computeDy(top, y0 fixed.FDot6) fixed.FDot6 {
	return fixed.FDot6(int32(uint32(top)<<6)) + 32 - y0
}

// diffToShift converts a distance to a shift value.
func diffToShift(dx, dy fixed.FDot6, shiftAA int32) int32 {
	// cheap calc of distance from center of p0-p2 to the center of the curve
	dist := cheapDistance(dx, dy)

	// shift down dist (it is currently in dot6)
	// down by 3 should give us 1/8 pixel accuracy (assuming our dist is accurate...)
	// this is chosen by heuristic: make it as big as possible (to minimize segments)
	// ... but small enough so that our curves still look smooth
	// When shift > 0, we're using AA and everything is scaled up so we can
	// lower the accuracy.
	dist = (dist + (1 << (2 + shiftAA))) >> (3 + shiftAA)

	// each subdivision (shift value) cuts this dist (error) by 1/4
	return (32 - int32(bits.LeadingZeros32(uint32(dist)))) >> 1
}

// cheapDistance computes a cheap approximation of distance.
func cheapDistance(dx, dy fixed.FDot6) fixed.FDot6 {
	if dx < 0 {
		dx = -dx
	}
	if dy < 0 {
		dy = -dy
	}
	// return max + min/2
	if dx > dy {
		return dx + (dy >> 1)
	}
	return dy + (dx >> 1)
}

// fdot6ToFixedDiv2 converts FDot6 to fixed point divided by 2.
func fdot6ToFixedDiv2(value fixed.FDot6) fixed.FDot16 {
	// we want to return SkFDot6ToFixed(value >> 1), but we don't want to throw
	// away data in value, so just perform a modify up-shift
	return fixed.FDot16(int32(uint32(value) << (16 - 6 - 1)))
}

// fdot6UpShift shifts FDot6 up by the given amount and returns FDot16.
func fdot6UpShift(x fixed.FDot6, upShift int32) fixed.FDot16 {
	return fixed.FDot16(int32(uint32(x) << uint(upShift)))
}

// cubicDeltaFromLine computes the maximum deviation of a cubic from its baseline.
func cubicDeltaFromLine(a, b, c, d fixed.FDot6) fixed.FDot6 {
	// since our parameters may be negative, we don't use <<
	oneThird := ((a*8 - b*15 + 6*c + d) * 19) >> 9
	twoThird := ((a + 6*b - c*15 + d*8) * 19) >> 9

	absOneThird := oneThird
	if absOneThird < 0 {
		absOneThird = -absOneThird
	}
	absTwoThird := twoThird
	if absTwoThird < 0 {
		absTwoThird = -absTwoThird
	}

	if absOneThird > absTwoThird {
		return absOneThird
	}
	return absTwoThird
}
