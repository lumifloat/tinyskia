// Copyright 2006 The Android Open Source Project
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// NOTE: this is not SkPathBuilder, but rather a reimplementation of SkPath.
package path

import (
	"github.com/chewxy/math32"

	"github.com/lumifloat/tinyskia/internal/numeric/scalar"
)

type pathDirection int

const (
	// Clockwise direction for adding closed contours.
	pathDirectionCW pathDirection = iota
	// Counter-clockwise direction for adding closed contours.
	pathDirectionCCW
)

// PathBuilder is a path builder.
type PathBuilder struct {
	verbs           []PathVerb
	points          []Point
	lastMoveToIndex int
	moveToRequired  bool
}

// New creates a new builder.
func NewPathBuilder() *PathBuilder {
	return &PathBuilder{
		verbs:           make([]PathVerb, 0),
		points:          make([]Point, 0),
		lastMoveToIndex: 0,
		moveToRequired:  true,
	}
}

// NewPathFromRect creates a new Path from Rect.
// Never fails since `Rect` is always valid.
// Segments are created clockwise: TopLeft -> TopRight -> BottomRight -> BottomLeft
// The contour is closed.
func NewPathFromRect(rect Rect) *Path {
	verbs := []PathVerb{
		PathVerbMove,
		PathVerbLine,
		PathVerbLine,
		PathVerbLine,
		PathVerbClose,
	}

	points := []Point{
		{rect.Left(), rect.Top()},
		{rect.Right(), rect.Top()},
		{rect.Right(), rect.Bottom()},
		{rect.Left(), rect.Bottom()},
	}

	return &Path{
		bounds: rect,
		verbs:  verbs,
		points: points,
	}
}

// NewPathFromCircle creates a new Path from a circle.
func NewPathFromCircle(cx, cy, radius float32) *Path {
	b := NewPathBuilder()
	b.PushCircle(cx, cy, radius)
	return b.Finish()
}

// NewPathFromOval creates a new Path from an oval.
func NewPathFromOval(oval Rect) *Path {
	b := NewPathBuilder()
	b.PushOval(oval)
	return b.Finish()
}

// Len returns the current number of segments in the builder.
func (b *PathBuilder) Len() int {
	return len(b.verbs)
}

// IsEmpty checks if the builder has any segments added.
func (b *PathBuilder) IsEmpty() bool {
	return len(b.verbs) == 0
}

// MoveTo adds beginning of a contour.
func (b *PathBuilder) MoveTo(x, y float32) {
	if len(b.verbs) > 0 && b.verbs[len(b.verbs)-1] == PathVerbMove {
		lastIdx := len(b.points) - 1
		b.points[lastIdx] = Point{x, y}
	} else {
		b.lastMoveToIndex = len(b.points)
		b.moveToRequired = false

		b.verbs = append(b.verbs, PathVerbMove)
		b.points = append(b.points, Point{x, y})
	}
}

func (b *PathBuilder) injectMoveToIfNeeded() {
	if b.moveToRequired {
		if b.lastMoveToIndex < len(b.points) {
			p := b.points[b.lastMoveToIndex]
			b.MoveTo(p.X, p.Y)
		} else {
			b.MoveTo(0.0, 0.0)
		}
	}
}

// LineTo adds a line from the last point.
func (b *PathBuilder) LineTo(x, y float32) {
	b.injectMoveToIfNeeded()
	b.verbs = append(b.verbs, PathVerbLine)
	b.points = append(b.points, Point{x, y})
}

// QuadTo adds a quad curve from the last point to x, y.
func (b *PathBuilder) QuadTo(x1, y1, x, y float32) {
	b.injectMoveToIfNeeded()

	b.verbs = append(b.verbs, PathVerbQuad)
	b.points = append(b.points, Point{x1, y1})
	b.points = append(b.points, Point{x, y})
}

func (b *PathBuilder) quadToPt(p1, p Point) {
	b.QuadTo(p1.X, p1.Y, p.X, p.Y)
}

func (b *PathBuilder) conicTo(x1, y1, x, y, weight float32) {
	if !(weight > 0.0) {
		b.LineTo(x, y)
	} else if math32.IsInf(weight, 0) || math32.IsNaN(weight) {
		b.LineTo(x1, y1)
		b.LineTo(x, y)
	} else if weight == 1.0 {
		b.QuadTo(x1, y1, x, y)
	} else {
		b.injectMoveToIfNeeded()
		last, _ := b.LastPoint()
		quadder, flag := ComputeAutoConicToQuads(
			last,
			Point{x1, y1},
			Point{x, y},
			weight,
		)
		if flag {
			offset := 1
			for i := 0; i < quadder.Len; i++ {
				pt1 := quadder.Points[offset+0]
				pt2 := quadder.Points[offset+1]
				b.QuadTo(pt1.X, pt1.Y, pt2.X, pt2.Y)
				offset += 2
			}
		}
	}
}

func (b *PathBuilder) conicPointsTo(pt1, pt2 Point, weight float32) {
	b.conicTo(pt1.X, pt1.Y, pt2.X, pt2.Y, weight)
}

// CubicTo adds a cubic curve from the last point to x, y.
func (b *PathBuilder) CubicTo(x1, y1, x2, y2, x, y float32) {
	b.injectMoveToIfNeeded()
	b.verbs = append(b.verbs, PathVerbCubic)
	b.points = append(b.points, Point{x1, y1}, Point{x2, y2}, Point{x, y})
}

func (b *PathBuilder) cubicToPt(p1, p2, p Point) {
	b.CubicTo(p1.X, p1.Y, p2.X, p2.Y, p.X, p.Y)
}

// Close closes the current contour.
func (b *PathBuilder) Close() {
	if len(b.verbs) > 0 {
		if b.verbs[len(b.verbs)-1] != PathVerbClose {
			b.verbs = append(b.verbs, PathVerbClose)
		}
	}
	b.moveToRequired = true
}

// LastPoint returns the last point if any.
func (b *PathBuilder) LastPoint() (Point, bool) {
	if len(b.points) == 0 {
		return Point{}, false
	}
	return b.points[len(b.points)-1], true
}

func (b *PathBuilder) setLastPoint(pt Point) {
	if len(b.points) > 0 {
		b.points[len(b.points)-1] = pt
	} else {
		b.MoveTo(pt.X, pt.Y)
	}
}

func (b *PathBuilder) isZeroLengthSincePoint(startPtIndex int) bool {
	count := len(b.points) - startPtIndex
	if count < 2 {
		return true
	}

	first := b.points[startPtIndex]
	for i := 1; i < count; i++ {
		if first != b.points[startPtIndex+i] {
			return false
		}
	}
	return true
}

// PushRect adds a rectangle contour.
func (b *PathBuilder) PushRect(rect Rect) {
	b.MoveTo(rect.Left(), rect.Top())
	b.LineTo(rect.Right(), rect.Top())
	b.LineTo(rect.Right(), rect.Bottom())
	b.LineTo(rect.Left(), rect.Bottom())
	b.Close()
}

// PushOval adds an oval contour bounded by the provided rectangle.
func (b *PathBuilder) PushOval(oval Rect) {
	cx := oval.Left()*0.5 + oval.Right()*0.5
	cy := oval.Top()*0.5 + oval.Bottom()*0.5

	ovalPoints := [4]Point{
		{cx, oval.Bottom()},
		{oval.Left(), cy},
		{cx, oval.Top()},
		{oval.Right(), cy},
	}

	rectPoints := [4]Point{
		{oval.Right(), oval.Bottom()},
		{oval.Left(), oval.Bottom()},
		{oval.Left(), oval.Top()},
		{oval.Right(), oval.Top()},
	}

	weight := scalar.SCALAR_ROOT_2_OVER_2
	b.MoveTo(ovalPoints[3].X, ovalPoints[3].Y)
	for i := 0; i < 4; i++ {
		b.conicPointsTo(rectPoints[i], ovalPoints[i], weight)
	}
	b.Close()
}

// PushCircle adds a circle contour.
func (b *PathBuilder) PushCircle(x, y, r float32) {
	if rect, ok := NewRectFromXYWH(x-r, y-r, r+r, r+r); ok {
		b.PushOval(rect)
	}
}

// PushPath adds a path.
func (b *PathBuilder) PushPath(other *Path) {
	b.lastMoveToIndex = len(b.points)
	b.verbs = append(b.verbs, other.verbs...)
	b.points = append(b.points, other.points...)
}

func (b *PathBuilder) pushPathBuilder(other *PathBuilder) {
	if other.IsEmpty() {
		return
	}
	if b.lastMoveToIndex != 0 {
		b.lastMoveToIndex = len(b.points) + other.lastMoveToIndex
	}
	b.verbs = append(b.verbs, other.verbs...)
	b.points = append(b.points, other.points...)
}

func (b *PathBuilder) reversePathTo(other *PathBuilder) {
	if other.IsEmpty() {
		return
	}

	pointsOffset := len(other.points) - 1
	for i := len(other.verbs) - 1; i >= 0; i-- {
		verb := other.verbs[i]
		switch verb {
		case PathVerbMove:
			return
		case PathVerbLine:
			pt := other.points[pointsOffset-1]
			pointsOffset -= 1
			b.LineTo(pt.X, pt.Y)
		case PathVerbQuad:
			pt1 := other.points[pointsOffset-1]
			pt2 := other.points[pointsOffset-2]
			pointsOffset -= 2
			b.QuadTo(pt1.X, pt1.Y, pt2.X, pt2.Y)
		case PathVerbCubic:
			pt1 := other.points[pointsOffset-1]
			pt2 := other.points[pointsOffset-2]
			pt3 := other.points[pointsOffset-3]
			pointsOffset -= 3
			b.CubicTo(pt1.X, pt1.Y, pt2.X, pt2.Y, pt3.X, pt3.Y)
		case PathVerbClose:
			// skip
		}
	}
}

// Clear reset the builder.
func (b *PathBuilder) Clear() {
	b.verbs = b.verbs[:0]
	b.points = b.points[:0]
	b.lastMoveToIndex = 0
	b.moveToRequired = true
}

// Finish finishes the builder and returns a Path.
func (b *PathBuilder) Finish() *Path {
	if b.IsEmpty() || len(b.verbs) == 1 {
		return nil
	}

	bounds, ok := NewRectFromPoints(b.points)
	if !ok {
		return nil
	}

	return &Path{
		bounds: bounds,
		verbs:  append([]PathVerb(nil), b.verbs...),
		points: append([]Point(nil), b.points...),
	}
}
