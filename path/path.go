// Copyright 2006 The Android Open Source Project
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package path

import (
	"github.com/chewxy/math32"

	"github.com/lumifloat/tinyskia/internal/normalized"
)

// PathVerb defines a path verb.
type PathVerb int

const (
	PathVerbMove PathVerb = iota
	PathVerbLine
	PathVerbQuad
	PathVerbCubic
	PathVerbClose
)

// Path is a Bezier path.
type Path struct {
	verbs  []PathVerb
	points []Point
	bounds Rect
}

// Copy creates a copy of the path.
func (p *Path) Copy() *Path {
	newPath := &Path{
		verbs:  make([]PathVerb, len(p.verbs)),
		points: make([]Point, len(p.points)),
		bounds: p.bounds,
	}
	copy(newPath.verbs, p.verbs)
	copy(newPath.points, p.points)
	return newPath
}

// Len returns the number of segments in the path.
func (p *Path) Len() int {
	return len(p.verbs)
}

// IsEmpty returns if the path is empty.
func (p *Path) IsEmpty() bool {
	return len(p.verbs) == 0
}

// Bounds returns the bounds of the path's points.
func (p *Path) Bounds() Rect {
	return p.bounds
}

// ComputeTightBounds calculates path's tight bounds.
func (p *Path) ComputeTightBounds() (Rect, bool) {
	if p.IsEmpty() {
		return Rect{}, false
	}

	extremas := [5]Point{}
	min := p.points[0]
	max := p.points[0]

	iter := p.Segments()
	lastPoint := Point{}
	for segment := iter.Next(); segment != nil; segment = iter.Next() {
		count := 0
		switch s := segment.(type) {
		case PathSegmentMoveTo:
			extremas[0] = Point(s)
			count = 1
		case PathSegmentLineTo:
			extremas[0] = Point(s)
			count = 1
		case PathSegmentQuadTo:
			count = computeQuadExtremas(lastPoint, s.P0, s.P1, &extremas)
		case PathSegmentCubicTo:
			count = computeCubicExtremas(lastPoint, s.P0, s.P1, s.P2, &extremas)
		case PathSegmentClose:
			// do nothing
		}

		lastPoint = iter.lastPoint

		for i := 0; i < count; i++ {
			tmp := extremas[i]
			min.X = math32.Min(min.X, tmp.X)
			min.Y = math32.Min(min.Y, tmp.Y)
			max.X = math32.Max(max.X, tmp.X)
			max.Y = math32.Max(max.Y, tmp.Y)
		}
	}

	return NewRectFromLTRB(min.X, min.Y, max.X, max.Y)
}

// Verbs returns an internal slice of verbs.
func (p *Path) Verbs() []PathVerb {
	return p.verbs
}

// Points returns an internal slice of points.
func (p *Path) Points() []Point {
	return p.points
}

// Transform returns a transformed in-place path.
func (p *Path) Transform(ts Transform) *Path {
	if ts.IsIdentity() {
		return p
	}

	ts.MapPoints(p.points)

	bounds, ok := NewRectFromPoints(p.points)
	if !ok {
		return nil
	}
	p.bounds = bounds

	return p
}

// Segments returns an iterator over path's segments.
func (p *Path) Segments() *PathSegmentsIter {
	return &PathSegmentsIter{
		path:        p,
		verbIndex:   0,
		pointsIndex: 0,
		isAutoClose: false,
		lastMoveTo:  Point{},
		lastPoint:   Point{},
	}
}

// Clear clears the path and returns a PathBuilder that will reuse an allocated memory.
func (p Path) Clear() *PathBuilder {
	verbs := p.verbs[:0]
	points := p.points[:0]

	return &PathBuilder{
		verbs:           verbs,
		points:          points,
		lastMoveToIndex: 0,
		moveToRequired:  true,
	}
}

func computeQuadExtremas(p0, p1, p2 Point, extremas *[5]Point) int {
	src := [3]Point{p0, p1, p2}
	idx := 0
	if t, ok := findQuadExtrema(p0.X, p1.X, p2.X); ok {
		extremas[idx] = evalQuadAt(src, t.ToNormalized())
		idx++
	}
	if t, ok := findQuadExtrema(p0.Y, p1.Y, p2.Y); ok {
		extremas[idx] = evalQuadAt(src, t.ToNormalized())
		idx++
	}
	extremas[idx] = p2
	return idx + 1
}

func computeCubicExtremas(p0, p1, p2, p3 Point, extremas *[5]Point) int {

	ts0 := [3]normalized.NormalizedF32Exclusive{}
	ts1 := [3]normalized.NormalizedF32Exclusive{}

	n0 := findCubicExtrema(p0.X, p1.X, p2.X, p3.X, &ts0)
	n1 := findCubicExtrema(p0.Y, p1.Y, p2.Y, p3.Y, &ts1)

	src := [4]Point{p0, p1, p2, p3}
	idx := 0
	for i := 0; i < n0; i++ {
		extremas[idx] = evalCubicPosAt(src, ts0[i].ToNormalized())
		idx++
	}
	for i := 0; i < n1; i++ {
		extremas[idx] = evalCubicPosAt(src, ts1[i].ToNormalized())
		idx++
	}
	extremas[idx] = p3
	return idx + 1
}

type (
	PathSegmentMoveTo  Point
	PathSegmentLineTo  Point
	PathSegmentQuadTo  struct{ P0, P1 Point }
	PathSegmentCubicTo struct{ P0, P1, P2 Point }
	PathSegmentClose   struct{}
)

// PathSegmentsIter a path segments iterator.
type PathSegmentsIter struct {
	path        *Path
	verbIndex   int
	pointsIndex int

	isAutoClose bool
	lastMoveTo  Point
	lastPoint   Point
}

func (iter *PathSegmentsIter) SetAutoClose(flag bool) {
	iter.isAutoClose = flag
}

func (iter *PathSegmentsIter) autoClose() interface{} {
	if iter.isAutoClose && iter.lastPoint != iter.lastMoveTo {
		iter.verbIndex--
		return PathSegmentLineTo(iter.lastMoveTo)
	}
	return PathSegmentClose{}
}

func (iter *PathSegmentsIter) hasValidTangent() bool {
	it := *iter
	for {
		segment := it.Next()
		if segment == nil {
			break
		}
		switch s := segment.(type) {
		case PathSegmentMoveTo:
			return false
		case PathSegmentLineTo:
			p := Point(s)
			if it.lastPoint == p {
				continue
			}
			return true
		case PathSegmentQuadTo:
			if it.lastPoint == s.P0 && it.lastPoint == s.P1 {
				continue
			}
			return true
		case PathSegmentCubicTo:
			if it.lastPoint == s.P0 && it.lastPoint == s.P1 && it.lastPoint == s.P2 {
				continue
			}
			return true
		case PathSegmentClose:
			return false
		}
	}
	return false
}

// CurrVerb returns the current verb.
func (iter *PathSegmentsIter) CurrVerb() PathVerb {
	if iter.verbIndex == 0 {
		panic("CurrVerb called before first segment")
	}
	return iter.path.verbs[iter.verbIndex-1]
}

// NextVerb returns the next verb or PathVerbClose if at end.
func (iter *PathSegmentsIter) NextVerb() PathVerb {
	if iter.verbIndex >= len(iter.path.verbs) {
		return PathVerbClose
	}
	return iter.path.verbs[iter.verbIndex]
}

// HasNext returns true if there are more verbs to iterate.
func (iter *PathSegmentsIter) HasNext() bool {
	return iter.verbIndex < len(iter.path.verbs)
}

func (iter *PathSegmentsIter) Next() interface{} {
	if iter.verbIndex >= len(iter.path.verbs) {
		return nil
	}

	verb := iter.path.verbs[iter.verbIndex]
	iter.verbIndex++

	switch verb {
	case PathVerbMove:
		iter.pointsIndex++
		iter.lastMoveTo = iter.path.points[iter.pointsIndex-1]
		iter.lastPoint = iter.lastMoveTo
		return PathSegmentMoveTo(iter.lastMoveTo)
	case PathVerbLine:
		iter.pointsIndex++
		iter.lastPoint = iter.path.points[iter.pointsIndex-1]
		return PathSegmentLineTo(iter.lastPoint)
	case PathVerbQuad:
		iter.pointsIndex += 2
		iter.lastPoint = iter.path.points[iter.pointsIndex-1]
		return PathSegmentQuadTo{
			iter.path.points[iter.pointsIndex-2],
			iter.lastPoint,
		}
	case PathVerbCubic:
		iter.pointsIndex += 3
		iter.lastPoint = iter.path.points[iter.pointsIndex-1]
		return PathSegmentCubicTo{
			iter.path.points[iter.pointsIndex-3],
			iter.path.points[iter.pointsIndex-2],
			iter.lastPoint,
		}
	case PathVerbClose:
		seg := iter.autoClose()
		iter.lastPoint = iter.lastMoveTo
		return seg
	}
	return nil
}
