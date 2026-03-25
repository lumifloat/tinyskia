// Copyright 2011 Google Inc.
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package edge

import (
	"github.com/lumifloat/tinyskia/path"
)

type combine int

const (
	combineNo combine = iota
	combinePartial
	combineTotal
)

type ShiftedIntRect struct {
	shifted path.ScreenIntRect
	shift   int32
}

// NewShiftedIntRect creates a new ShiftedIntRect by shifting the rect coordinates.
func NewShiftedIntRect(rect path.ScreenIntRect, shift int32) (ShiftedIntRect, bool) {
	shifted, ok := path.NewScreenIntRectFromXYWH(
		rect.X()<<shift,
		rect.Y()<<shift,
		rect.Width()<<shift,
		rect.Height()<<shift,
	)
	if !ok {
		return ShiftedIntRect{}, false
	}
	return ShiftedIntRect{shifted: shifted, shift: shift}, true
}

func (s ShiftedIntRect) Shifted() path.ScreenIntRect {
	return s.shifted
}

func (s ShiftedIntRect) Recover() path.ScreenIntRect {
	// cannot fail, because the original rect was valid
	r, _ := path.NewScreenIntRectFromXYWH(
		s.shifted.X()>>s.shift,
		s.shifted.Y()>>s.shift,
		s.shifted.Width()>>s.shift,
		s.shifted.Height()>>s.shift,
	)
	return r
}

type BasicEdgeBuilder struct {
	edges     []Edge
	clipShift int32
}

// NewBasicEdgeBuilder creates a new BasicEdgeBuilder.
func NewBasicEdgeBuilder(clipShift int32) *BasicEdgeBuilder {
	return &BasicEdgeBuilder{
		edges:     make([]Edge, 0, 64), // TODO: stack array + fallback
		clipShift: clipShift,
	}
}

// BuildEdges mimics Skia's edge building logic.
// Returns nil if infinite or NaN segments are detected.
func BuildEdges(p *path.Path, hasClip bool, clip ShiftedIntRect, clipShift int32) []Edge {
	// If we're convex, then we need both edges, even if the right edge is past the clip.
	canCullToTheRight := false // TODO: implement isConvex logic

	builder := NewBasicEdgeBuilder(clipShift)
	if !builder.Build(p, hasClip, clip, canCullToTheRight) {
		// infinite or NaN segments detected during edges building
		return nil
	}

	if len(builder.edges) < 2 {
		return nil
	}

	return builder.edges
}

func (b *BasicEdgeBuilder) Build(p *path.Path, hasClip bool, clip ShiftedIntRect, canCullToTheRight bool) bool {
	if hasClip {
		clipRect := clip.Recover().ToRect()
		iter := NewEdgeClipperIter(p, clipRect, canCullToTheRight)
		for {
			edges, ok := iter.next()
			if !ok {
				break
			}
			for _, e := range edges {
				switch v := e.(type) {
				case PathEdgeLineTo:
					// Check finite (manually or via path.Point method)
					b.pushLine([2]path.Point{v.P0, v.P1})
				case PathEdgeQuadTo:
					b.pushQuad([3]path.Point{v.P0, v.P1, v.P2})
				case PathEdgeCubicTo:
					b.pushCubic([4]path.Point{v.P0, v.P1, v.P2, v.P3})
				}
			}
		}
	} else {
		iter := NewPathEdgeIter(p)
		for {
			e, ok := iter.next()
			if !ok {
				break
			}
			switch v := e.(type) {
			case PathEdgeLineTo:
				b.pushLine([2]path.Point{v.P0, v.P1})
			case PathEdgeQuadTo:
				points := [3]path.Point{v.P0, v.P1, v.P2}
				monoX := [5]path.Point{}
				n := path.ChopQuadAtYExtrema(points, &monoX)
				for i := 0; i <= n; i++ {
					b.pushQuad([3]path.Point{
						monoX[i*2+0],
						monoX[i*2+1],
						monoX[i*2+2],
					})
				}
			case PathEdgeCubicTo:
				points := [4]path.Point{v.P0, v.P1, v.P2, v.P3}
				monoY := [10]path.Point{}
				n := path.ChopCubicAtYExtrema(points, &monoY)
				for i := 0; i <= n; i++ {
					b.pushCubic([4]path.Point{
						monoY[i*3+0],
						monoY[i*3+1],
						monoY[i*3+2],
						monoY[i*3+3],
					})
				}
			}
		}
	}
	return true
}

func (b *BasicEdgeBuilder) pushLine(points [2]path.Point) {
	e := NewLineEdge(points[0], points[1], b.clipShift)
	if e == nil {
		return
	}

	combine := combineNo
	if e.IsVertical() && len(b.edges) > 0 {
		lastEdge := &b.edges[len(b.edges)-1]
		if lastEdge.Type == EdgeTypeLine && lastEdge.Line != nil {
			combine = combineVertical(e, lastEdge.Line)
		}
	}

	switch combine {
	case combineTotal:
		b.edges = b.edges[:len(b.edges)-1]
	case combinePartial:
		// already modified 'last' in combineVertical
	case combineNo:
		b.edges = append(b.edges, Edge{Type: EdgeTypeLine, Line: e})
	}
}

func (b *BasicEdgeBuilder) pushQuad(points [3]path.Point) {
	if e := NewQuadraticEdge(points, b.clipShift); e != nil {
		b.edges = append(b.edges, Edge{Type: EdgeTypeQuadratic, Quadratic: e})
	}
}

func (b *BasicEdgeBuilder) pushCubic(points [4]path.Point) {
	if e := NewCubicEdge(points, b.clipShift); e != nil {
		b.edges = append(b.edges, Edge{Type: EdgeTypeCubic, Cubic: e})
	}
}

func combineVertical(edge *LineEdge, last *LineEdge) combine {
	if last.DX != 0 || edge.X != last.X {
		return combineNo
	}

	if edge.Winding == last.Winding {
		if edge.LastY+1 == last.FirstY {
			last.FirstY = edge.FirstY
			return combinePartial
		} else if edge.FirstY == last.LastY+1 {
			last.LastY = edge.LastY
			return combinePartial
		} else {
			return combineNo
		}
	}

	if edge.FirstY == last.FirstY {
		if edge.LastY == last.LastY {
			return combineTotal
		} else if edge.LastY < last.LastY {
			last.FirstY = edge.LastY + 1
			return combinePartial
		} else {
			last.FirstY = last.LastY + 1
			last.LastY = edge.LastY
			last.Winding = edge.Winding
			return combinePartial
		}
	}

	if edge.LastY == last.LastY {
		if edge.FirstY > last.FirstY {
			last.LastY = edge.FirstY - 1
		} else {
			last.LastY = last.FirstY - 1
			last.FirstY = edge.FirstY
			last.Winding = edge.Winding
		}
		return combinePartial
	}

	return combineNo
}

// PathEdge represents a path edge segment.
type PathEdge interface{}

// PathEdgeLine represents a line edge.
type PathEdgeLineTo struct {
	P0 path.Point
	P1 path.Point
}

// PathEdgeQuad represents a quadratic Bezier edge.
type PathEdgeQuadTo struct {
	P0 path.Point
	P1 path.Point
	P2 path.Point
}

// PathEdgeCubic represents a cubic Bezier edge.
type PathEdgeCubicTo struct {
	P0 path.Point
	P1 path.Point
	P2 path.Point
	P3 path.Point
}

type PathEdgeIter struct {
	path           *path.Path
	verbIdx        int
	pointIdx       int
	moveTo         path.Point
	needsCloseLine bool
}

// NewPathEdgeIter creates a new PathEdgeIter.
func NewPathEdgeIter(p *path.Path) *PathEdgeIter {
	return &PathEdgeIter{
		path:   p,
		moveTo: path.Point{X: 0, Y: 0},
	}
}

func (iter *PathEdgeIter) closeLine() PathEdge {
	iter.needsCloseLine = false
	return PathEdgeLineTo{
		P0: iter.path.Points()[iter.pointIdx-1],
		P1: iter.moveTo,
	}
}

func (iter *PathEdgeIter) next() (PathEdge, bool) {
	verbs := iter.path.Verbs()
	points := iter.path.Points()

	if iter.verbIdx < len(verbs) {
		verb := verbs[iter.verbIdx]
		iter.verbIdx++

		switch verb {
		case path.PathVerbMove:
			if iter.needsCloseLine {
				res := iter.closeLine()
				iter.moveTo = points[iter.pointIdx]
				iter.pointIdx++
				return res, true
			}
			iter.moveTo = points[iter.pointIdx]
			iter.pointIdx++
			return iter.next()

		case path.PathVerbClose:
			if iter.needsCloseLine {
				return iter.closeLine(), true
			}
			return iter.next()

		case path.PathVerbLine:
			iter.needsCloseLine = true
			edge := PathEdgeLineTo{P0: points[iter.pointIdx-1], P1: points[iter.pointIdx]}
			iter.pointIdx++
			return edge, true

		case path.PathVerbQuad:
			iter.needsCloseLine = true
			edge := PathEdgeQuadTo{
				P0: points[iter.pointIdx-1],
				P1: points[iter.pointIdx],
				P2: points[iter.pointIdx+1],
			}
			iter.pointIdx += 2
			return edge, true

		case path.PathVerbCubic:
			iter.needsCloseLine = true
			edge := PathEdgeCubicTo{
				P0: points[iter.pointIdx-1],
				P1: points[iter.pointIdx],
				P2: points[iter.pointIdx+1],
				P3: points[iter.pointIdx+2],
			}
			iter.pointIdx += 3
			return edge, true
		}
	} else if iter.needsCloseLine {
		return iter.closeLine(), true
	}

	return nil, false
}
