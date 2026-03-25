// Copyright 2009 The Android Open Source Project
// Copyright 2020 Yevhenii Reizner
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package edge

import (
	"math"

	"github.com/chewxy/math32"
	"github.com/lumifloat/tinyskia/internal/normalized"
	"github.com/lumifloat/tinyskia/path"
)

// Max curvature in X and Y split cubic into 9 pieces, * (line + cubic).
const maxVerbs = 18

type ClippedEdges []PathEdge

type EdgeClipper struct {
	clip              path.Rect
	canCullToTheRight bool
	edges             ClippedEdges
}

func NewEdgeClipper(clip path.Rect, canCullToTheRight bool) *EdgeClipper {
	return &EdgeClipper{
		clip:              clip,
		canCullToTheRight: canCullToTheRight,
		edges:             make(ClippedEdges, 0, maxVerbs),
	}
}

func (c *EdgeClipper) ClipLine(p0, p1 path.Point) []PathEdge {
	if !p0.IsFinite() || !p1.IsFinite() {
		return nil
	}

	src := [2]path.Point{p0, p1}
	var points [path.MAX_POINTS]path.Point
	n := path.Clip(src, c.clip, c.canCullToTheRight, &points)

	if n == 0 {
		return nil
	}

	for i := 0; i < n-1; i++ {
		c.pushLine(points[i], points[i+1])
	}

	if len(c.edges) == 0 {
		return nil
	}
	return c.edges
}

func (c *EdgeClipper) pushLine(p0, p1 path.Point) {
	c.edges = append(c.edges, PathEdgeLineTo{P0: p0, P1: p1})
}

func (c *EdgeClipper) pushVLine(x, y0, y1 float32, reverse bool) {
	if reverse {
		y0, y1 = y1, y0
	}

	c.edges = append(c.edges, PathEdgeLineTo{
		P0: path.Point{X: x, Y: y0},
		P1: path.Point{X: x, Y: y1},
	})
}

func (c *EdgeClipper) ClipQuad(p0, p1, p2 path.Point) []PathEdge {
	points := [3]path.Point{p0, p1, p2}
	bounds := rectFromPoints(points[:])

	if !quickReject(bounds, c.clip) {
		var monoY [5]path.Point
		countY := path.ChopQuadAtYExtrema(points, &monoY)
		for y := 0; y <= int(countY); y++ {
			var monoX [5]path.Point
			yPoints := [3]path.Point{monoY[y*2], monoY[y*2+1], monoY[y*2+2]}
			countX := path.ChopQuadAtXExtrema(yPoints, &monoX)
			for x := 0; x <= int(countX); x++ {
				xPoints := [3]path.Point{monoX[x*2], monoX[x*2+1], monoX[x*2+2]}
				c.clipMonoQuad(xPoints)
			}
		}
	}

	if len(c.edges) == 0 {
		return nil
	}
	return c.edges
}

// src must be monotonic in X and Y
func (c *EdgeClipper) clipMonoQuad(src [3]path.Point) {
	var pts [3]path.Point
	reverse := sortIncreasingY(src[:], pts[:])

	// are we completely above or below
	if pts[2].Y <= c.clip.Top() || pts[0].Y >= c.clip.Bottom() {
		return
	}

	// Now chop so that pts is contained within clip in Y
	chopQuadInY(c.clip, &pts)

	if pts[0].X > pts[2].X {
		pts[0], pts[2] = pts[2], pts[0]
		reverse = !reverse
	}

	// Now chop in X has needed, and record the segments
	if pts[2].X <= c.clip.Left() {
		// wholly to the left
		c.pushVLine(c.clip.Left(), pts[0].Y, pts[2].Y, reverse)
		return
	}

	if pts[0].X >= c.clip.Right() {
		// wholly to the right
		if !c.canCullToTheRight {
			c.pushVLine(c.clip.Right(), pts[0].Y, pts[2].Y, reverse)
		}
		return
	}

	var t normalized.NormalizedF32Exclusive
	var tmp [5]path.Point

	// are we partially to the left
	if pts[0].X < c.clip.Left() {
		if t, ok := chopMonoQuadAtX(pts, c.clip.Left(), t); ok {
			path.ChopQuadAt(pts, t, &tmp)
			c.pushVLine(c.clip.Left(), tmp[0].Y, tmp[2].Y, reverse)
			// clamp to clean up imprecise numerics in the chop
			tmp[2].X = c.clip.Left()
			tmp[3].X = math32.Max(tmp[3].X, c.clip.Left())

			pts[0] = tmp[2]
			pts[1] = tmp[3]
		} else {
			// if chopMonoQuadAtY failed, then we may have hit inexact numerics
			// so we just clamp against the left
			c.pushVLine(c.clip.Left(), pts[0].Y, pts[2].Y, reverse)
			return
		}
	}

	// are we partially to the right
	if pts[2].X > c.clip.Right() {
		if t, ok := chopMonoQuadAtX(pts, c.clip.Right(), t); ok {
			path.ChopQuadAt(pts, t, &tmp)
			// clamp to clean up imprecise numerics in the chop
			tmp[1].X = math32.Min(tmp[1].X, c.clip.Right())
			tmp[2].X = c.clip.Right()

			c.pushQuad([3]path.Point{tmp[0], tmp[1], tmp[2]}, reverse)
			c.pushVLine(c.clip.Right(), tmp[2].Y, tmp[4].Y, reverse)
		} else {
			// if chopMonoQuadAtY failed, then we may have hit inexact numerics
			// so we just clamp against the right
			pts[1].X = math32.Min(pts[1].X, c.clip.Right())
			pts[2].X = math32.Min(pts[2].X, c.clip.Right())
			c.pushQuad(pts, reverse)
		}
	} else {
		// wholly inside the clip
		c.pushQuad(pts, reverse)
	}
}

func (c *EdgeClipper) pushQuad(pts [3]path.Point, reverse bool) {
	if reverse {
		c.edges = append(c.edges, PathEdgeQuadTo{P0: pts[2], P1: pts[1], P2: pts[0]})
	} else {
		c.edges = append(c.edges, PathEdgeQuadTo{P0: pts[0], P1: pts[1], P2: pts[2]})
	}
}

func (c *EdgeClipper) ClipCubic(p0, p1, p2, p3 path.Point) []PathEdge {
	points := [4]path.Point{p0, p1, p2, p3}
	bounds := rectFromPoints(points[:])

	// check if we're clipped out vertically
	if bounds.Bottom() > c.clip.Top() && bounds.Top() < c.clip.Bottom() {
		if tooBigForReliableFloatMath(bounds) {
			// can't safely clip the cubic, so we give up and draw a line (which we can safely clip)
			return c.ClipLine(p0, p3)
		} else {
			var monoY [10]path.Point
			countY := path.ChopCubicAtYExtrema(points, &monoY)
			for y := 0; y <= int(countY); y++ {
				var monoX [10]path.Point
				yPoints := [4]path.Point{monoY[y*3], monoY[y*3+1], monoY[y*3+2], monoY[y*3+3]}
				countX := path.ChopCubicAtXExtrema(yPoints, &monoX)
				for x := 0; x <= int(countX); x++ {
					xPoints := [4]path.Point{monoX[x*3], monoX[x*3+1], monoX[x*3+2], monoX[x*3+3]}
					c.clipMonoCubic(xPoints)
				}
			}
		}
	}

	if len(c.edges) == 0 {
		return nil
	}
	return c.edges
}

// src must be monotonic in X and Y
func (c *EdgeClipper) clipMonoCubic(src [4]path.Point) {
	var pts [4]path.Point
	reverse := sortIncreasingY(src[:], pts[:])

	// are we completely above or below
	if pts[3].Y <= c.clip.Top() || pts[0].Y >= c.clip.Bottom() {
		return
	}

	// Now chop so that pts is contained within clip in Y
	chopCubicInY(c.clip, &pts)

	if pts[0].X > pts[3].X {
		pts[0], pts[3] = pts[3], pts[0]
		pts[1], pts[2] = pts[2], pts[1]
		reverse = !reverse
	}

	// Now chop in X has needed, and record the segments
	if pts[3].X <= c.clip.Left() {
		// wholly to the left
		c.pushVLine(c.clip.Left(), pts[0].Y, pts[3].Y, reverse)
		return
	}

	if pts[0].X >= c.clip.Right() {
		// wholly to the right
		if !c.canCullToTheRight {
			c.pushVLine(c.clip.Right(), pts[0].Y, pts[3].Y, reverse)
		}
		return
	}

	// are we partially to the left
	if pts[0].X < c.clip.Left() {
		var tmp [7]path.Point
		chopMonoCubicAtX(pts, c.clip.Left(), &tmp)
		c.pushVLine(c.clip.Left(), tmp[0].Y, tmp[3].Y, reverse)

		tmp[3].X = c.clip.Left()
		tmp[4].X = math32.Max(tmp[4].X, c.clip.Left())

		pts[0] = tmp[3]
		pts[1] = tmp[4]
		pts[2] = tmp[5]
	}

	// are we partially to the right
	if pts[3].X > c.clip.Right() {
		var tmp [7]path.Point
		chopMonoCubicAtX(pts, c.clip.Right(), &tmp)
		tmp[3].X = c.clip.Right()
		tmp[2].X = math32.Min(tmp[2].X, c.clip.Right())

		c.pushCubic([4]path.Point{tmp[0], tmp[1], tmp[2], tmp[3]}, reverse)
		c.pushVLine(c.clip.Right(), tmp[3].Y, tmp[6].Y, reverse)
	} else {
		// wholly inside the clip
		c.pushCubic(pts, reverse)
	}
}

func (c *EdgeClipper) pushCubic(pts [4]path.Point, reverse bool) {
	if reverse {
		c.edges = append(c.edges, PathEdgeCubicTo{P0: pts[3], P1: pts[2], P2: pts[1], P3: pts[0]})
	} else {
		c.edges = append(c.edges, PathEdgeCubicTo{P0: pts[0], P1: pts[1], P2: pts[2], P3: pts[3]})
	}
}

type EdgeClipperIter struct {
	edgeIter          *PathEdgeIter
	clip              path.Rect
	canCullToTheRight bool
}

func NewEdgeClipperIter(p *path.Path, clip path.Rect, canCullToTheRight bool) *EdgeClipperIter {
	return &EdgeClipperIter{
		edgeIter:          NewPathEdgeIter(p),
		clip:              clip,
		canCullToTheRight: canCullToTheRight,
	}
}

func (iter *EdgeClipperIter) next() ([]PathEdge, bool) {
	for {
		edge, ok := iter.edgeIter.next()
		if !ok {
			break
		}
		clipper := NewEdgeClipper(iter.clip, iter.canCullToTheRight)

		switch v := edge.(type) {
		case PathEdgeLineTo:
			if edges := clipper.ClipLine(v.P0, v.P1); edges != nil {
				return edges, true
			}
		case PathEdgeQuadTo:
			if edges := clipper.ClipQuad(v.P0, v.P1, v.P2); edges != nil {
				return edges, true
			}
		case PathEdgeCubicTo:
			if edges := clipper.ClipCubic(v.P0, v.P1, v.P2, v.P3); edges != nil {
				return edges, true
			}
		}
	}
	return nil, false
}

func quickReject(bounds, clip path.Rect) bool {
	return bounds.Top() >= clip.Bottom() || bounds.Bottom() <= clip.Top()
}

func sortIncreasingY(src []path.Point, dst []path.Point) bool {
	if src[0].Y > src[len(src)-1].Y {
		for i := 0; i < len(src); i++ {
			dst[i] = src[len(src)-1-i]
		}
		return true
	}
	copy(dst, src)
	return false
}

func chopQuadInY(clip path.Rect, pts *[3]path.Point) {
	var t normalized.NormalizedF32Exclusive
	var tmp [5]path.Point

	if pts[0].Y < clip.Top() {
		if t, ok := chopMonoQuadAtY(*pts, clip.Top(), t); ok {
			path.ChopQuadAt(*pts, t, &tmp)
			tmp[2].Y = clip.Top()
			tmp[3].Y = math32.Max(tmp[3].Y, clip.Top())
			pts[0] = tmp[2]
			pts[1] = tmp[3]
		} else {
			for i := range pts {
				if pts[i].Y < clip.Top() {
					pts[i].Y = clip.Top()
				}
			}
		}
	}

	if pts[2].Y > clip.Bottom() {
		if t, ok := chopMonoQuadAtY(*pts, clip.Bottom(), t); ok {
			path.ChopQuadAt(*pts, t, &tmp)
			tmp[1].Y = math32.Min(tmp[1].Y, clip.Bottom())
			tmp[2].Y = clip.Bottom()
			pts[1] = tmp[1]
			pts[2] = tmp[2]
		} else {
			for i := range pts {
				if pts[i].Y > clip.Bottom() {
					pts[i].Y = clip.Bottom()
				}
			}
		}
	}
}

func chopMonoQuadAtX(pts [3]path.Point, x float32, t normalized.NormalizedF32Exclusive) (normalized.NormalizedF32Exclusive, bool) {
	return chopMonoQuadAt(pts[0].X, pts[1].X, pts[2].X, x, t)
}

func chopMonoQuadAtY(pts [3]path.Point, y float32, t normalized.NormalizedF32Exclusive) (normalized.NormalizedF32Exclusive, bool) {
	return chopMonoQuadAt(pts[0].Y, pts[1].Y, pts[2].Y, y, t)
}

func chopMonoQuadAt(c0, c1, c2, target float32, t normalized.NormalizedF32Exclusive) (normalized.NormalizedF32Exclusive, bool) {
	a := c0 - c1 - c1 + c2
	b := 2.0 * (c1 - c0)
	c := c0 - target

	var roots [3]normalized.NormalizedF32Exclusive
	count := findUnitQuadRoots(a, b, c, &roots)
	if count != 0 {
		return roots[0], true
	}
	return 0, false
}

// findUnitQuadRoots finds roots of quadratic equation.
func findUnitQuadRoots(a, b, c float32, roots *[3]normalized.NormalizedF32Exclusive) int {
	if a == 0.0 {
		if r, ok := validUnitDivide(-c, b); ok {
			roots[0] = r
			return 1
		}
		return 0
	}

	dr := b*b - 4.0*a*c
	if dr < 0.0 {
		return 0
	}
	dr = float32(math.Sqrt(float64(dr)))

	if math.IsInf(float64(dr), 0) || math.IsNaN(float64(dr)) {
		return 0
	}

	var q float32
	if b < 0.0 {
		q = -(b - dr) / 2.0
	} else {
		q = -(b + dr) / 2.0
	}

	offset := 0
	if r, ok := validUnitDivide(q, a); ok {
		roots[offset] = r
		offset++
	}

	if r, ok := validUnitDivide(c, q); ok {
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

// validUnitDivide validates and performs unit division.
func validUnitDivide(numer, denom float32) (normalized.NormalizedF32Exclusive, bool) {
	n := numer
	d := denom
	if n < 0.0 {
		n = -n
		d = -d
	}

	if d == 0.0 || n == 0.0 || n >= d {
		return 0, false
	}

	r := n / d
	return normalized.NewNormalizedF32Exclusive(r)
}

func tooBigForReliableFloatMath(r path.Rect) bool {
	const limit = float32(1 << 22)
	return r.Left() < -limit || r.Top() < -limit || r.Right() > limit || r.Bottom() > limit
}

func chopCubicInY(clip path.Rect, pts *[4]path.Point) {
	if pts[0].Y < clip.Top() {
		var tmp [7]path.Point
		chopMonoCubicAtY(*pts, clip.Top(), &tmp)

		if tmp[3].Y < clip.Top() && tmp[4].Y < clip.Top() && tmp[5].Y < clip.Top() {
			tmp2 := [4]path.Point{tmp[3], tmp[4], tmp[5], tmp[6]}
			chopMonoCubicAtY(tmp2, clip.Top(), &tmp)
		}

		tmp[3].Y = clip.Top()
		tmp[4].Y = math32.Max(tmp[4].Y, clip.Top())

		pts[0] = tmp[3]
		pts[1] = tmp[4]
		pts[2] = tmp[5]
	}

	if pts[3].Y > clip.Bottom() {
		var tmp [7]path.Point
		chopMonoCubicAtY(*pts, clip.Bottom(), &tmp)
		tmp[3].Y = clip.Bottom()
		tmp[2].Y = math32.Min(tmp[2].Y, clip.Bottom())

		pts[1] = tmp[1]
		pts[2] = tmp[2]
		pts[3] = tmp[3]
	}
}

func rectFromPoints(pts []path.Point) path.Rect {
	if len(pts) == 0 {
		return path.Rect{}
	}

	left := pts[0].X
	right := pts[0].X
	top := pts[0].Y
	bottom := pts[0].Y

	for _, p := range pts[1:] {
		if p.X < left {
			left = p.X
		}
		if p.X > right {
			right = p.X
		}
		if p.Y < top {
			top = p.Y
		}
		if p.Y > bottom {
			bottom = p.Y
		}
	}

	rect, _ := path.NewRectFromLTRB(left, top, right, bottom)
	return rect
}

func chopMonoCubicAtX(src [4]path.Point, x float32, dst *[7]path.Point) {
	// Use binary search approach
	t := monoCubicClosestT([4]float32{src[0].X, src[1].X, src[2].X, src[3].X}, x)
	if nt, ok := normalized.NewNormalizedF32Exclusive(float32(t)); ok {
		path.ChopCubicAt2(src, nt, dst)
	} else {
		// Fallback to endpoints
		dst[0] = src[0]
		dst[1] = src[1]
		dst[2] = src[2]
		dst[3] = src[3]
		dst[4] = src[3]
		dst[5] = src[3]
		dst[6] = src[3]
	}
}

func chopMonoCubicAtY(src [4]path.Point, y float32, dst *[7]path.Point) {
	// Use binary search approach
	t := monoCubicClosestT([4]float32{src[0].Y, src[1].Y, src[2].Y, src[3].Y}, y)
	if nt, ok := normalized.NewNormalizedF32Exclusive(float32(t)); ok {
		path.ChopCubicAt2(src, nt, dst)
	} else {
		// Fallback to endpoints
		dst[0] = src[0]
		dst[1] = src[1]
		dst[2] = src[2]
		dst[3] = src[3]
		dst[4] = src[3]
		dst[5] = src[3]
		dst[6] = src[3]
	}
}

func monoCubicClosestT(src [4]float32, x float32) normalized.NormalizedF32Exclusive {
	t := float32(0.5)
	bestT := t
	step := float32(0.25)
	d := src[0]
	a := src[3] + 3.0*(src[1]-src[2]) - d
	b := 3.0 * (src[2] - src[1] - src[1] + d)
	c := 3.0 * (src[1] - d)
	x -= d
	closest := float32(math.MaxFloat32)
	for {
		loc := ((a*t+b)*t + c) * t
		dist := float32(math.Abs(float64(loc - x)))
		if closest > dist {
			closest = dist
			bestT = t
		}

		lastT := t
		if loc < x {
			t += step
		} else {
			t -= step
		}
		step *= 0.5

		if !(closest > 0.25 && lastT != t) {
			break
		}
	}
	// Convert bestT to NormalizedF32Exclusive
	result, _ := normalized.NewNormalizedF32Exclusive(bestT)
	return result
}
