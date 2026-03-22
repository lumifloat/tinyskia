// Copyright 2014 Google Inc.
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// This module is a mix of SkDashPath, SkDashPathEffect, SkcontourMeasure and SkPathMeasure.
package path

import (
	"github.com/chewxy/math32"
	"github.com/lumifloat/tinyskia/internal/normalized"
)

const (
	maxTValue uint32 = 0x3FFFFFFF
)

// StrokeDash a stroke dashing properties.
//
// Contains an array of pairs, where the first number indicates an "on" interval
// and the second one indicates an "off" interval;
// a dash offset value and internal properties.
//
// # Guarantees
//
// - The dash array always have an even number of values.
// - All dash array values are finite and >= 0.
// - There is at least two dash array values.
// - The sum of all dash array values is positive and finite.
// - Dash offset is finite.
type StrokeDash struct {
	array       []float32
	offset      float32
	intervalLen float32
	firstLen    float32
	firstIndex  int
}

// NewStrokeDash creates a new stroke dashing object.
func NewStrokeDash(array []float32, offest float32) *StrokeDash {
	if math32.IsInf(offest, 0) || math32.IsNaN(offest) {
		return nil
	}

	if len(array) < 2 || len(array)%2 != 0 {
		return nil
	}

	var intervalLen float32
	for _, n := range array {
		if n < 0.0 || math32.IsInf(n, 0) || math32.IsNaN(n) {
			return nil
		}
		intervalLen += n
	}

	if intervalLen <= 0 {
		return nil
	}

	adjustedOffset := adjustDashOffset(offest, intervalLen)
	firstLen, firstIndex := findFirstInterval(array, adjustedOffset)

	return &StrokeDash{
		array:       array,
		offset:      adjustedOffset,
		intervalLen: intervalLen,
		firstLen:    firstLen,
		firstIndex:  firstIndex,
	}
}

// Adjust phase to be between 0 and len, "flipping" phase if negative.
// e.g., if len is 100, then phase of -20 (or -120) is equivalent to 80.
func adjustDashOffset(offset, length float32) float32 {
	if offset < 0.0 {
		offset = -offset
		if offset > length {
			offset = math32.Mod(offset, length)
		}
		offset = length - offset

		if offset == length {
			offset = 0.0
		}
		return offset
	} else if offset >= length {
		return math32.Mod(offset, length)
	}
	return offset
}

func findFirstInterval(array []float32, offset float32) (float32, int) {
	for i, gap := range array {
		if offset > gap || (offset == gap && gap != 0.0) {
			offset -= gap
		} else {
			return gap - offset, i
		}
	}
	return array[0], 0
}

// Dash converts the current path into a dashed one.
//
// Returns nil when more than 1_000_000 dashes had to be produced
// or when the final path has an invalid bounding box.
func (p *Path) Dash(dash *StrokeDash, resScale float32) *Path {
	// Since the path length / dash length ratio may be arbitrarily large, we can exert
	// significant memory pressure while attempting to build the filtered path. To avoid this,
	// we simply give up dashing beyond a certain threshold.
	//
	// The original bug report (http://crbug.com/165432) is based on a path yielding more than
	// 90 million dash segments and crashing the memory allocator. A limit of 1 million
	// segments seems reasonable: at 2 verbs per segment * 9 bytes per verb, this caps the
	// maximum dash memory overhead at roughly 17MB per path.
	const maxDashCount = 1000000

	pb := NewPathBuilder()
	var dashCount float32

	iter := newcontourMeasureIter(p, resScale)
	for {
		contour, ok := iter.next()
		if !ok {
			break
		}

		skipFirstSegment := contour.isClosed
		addedSegment := false
		length := contour.length
		index := dash.firstIndex

		dashCount += length * float32(len(dash.array)>>1) / dash.intervalLen
		if dashCount > maxDashCount {
			return nil
		}

		var distance float32 = 0.0
		dLen := dash.firstLen

		for distance < length {
			addedSegment = false
			if index%2 == 0 && !skipFirstSegment {
				addedSegment = true
				contour.pushSegment(distance, distance+dLen, true, pb)
			}

			distance += dLen
			skipFirstSegment = false

			index++
			if index == len(dash.array) {
				index = 0
			}
			dLen = dash.array[index]
		}

		if contour.isClosed && dash.firstIndex%2 == 0 && dash.firstLen >= 0.0 {
			contour.pushSegment(0.0, dash.firstLen, !addedSegment, pb)
		}
	}

	return pb.Finish()
}

type contourMeasureIter struct {
	iter      *PathSegmentsIter
	tolerance float32
}

func newcontourMeasureIter(path *Path, resScale float32) *contourMeasureIter {
	// can't use tangents, since we need [0..1..................2] to be seen
	// as definitely not a line (it is when drawn, but not parametrically)
	// so we compare midpoints
	const cheapDistLimit = 0.5 // just made this value up
	return &contourMeasureIter{
		iter:      path.Segments(),
		tolerance: cheapDistLimit / resScale,
	}
}

func (iter *contourMeasureIter) next() (*contourMeasure, bool) {
	contour := &contourMeasure{}
	var pointIndex int
	var distance float32
	haveSeenClose := false
	prevP := Point{0, 0}

	for seg := iter.iter.Next(); seg != nil; seg = iter.iter.Next() {

		switch s := seg.(type) {
		case PathSegmentMoveTo:
			contour.points = append(contour.points, Point(s))
			prevP = Point(s)
		case PathSegmentLineTo:
			prevD := distance
			distance = contour.computeLineSeg(
				prevP, Point(s), distance, pointIndex,
			)
			if distance > prevD {
				contour.points =
					append(contour.points, Point(s))
				pointIndex++
			}
			prevP = Point(s)
		case PathSegmentQuadTo:
			prevD := distance
			distance = contour.computeQuadSegs(
				prevP, s.P0, s.P1, distance, 0,
				maxTValue, pointIndex, iter.tolerance,
			)
			if distance > prevD {
				contour.points =
					append(contour.points, s.P0, s.P1)
				pointIndex += 2
			}
			prevP = s.P1
		case PathSegmentCubicTo:
			prevD := distance
			distance = contour.computeCubicSegs(
				prevP, s.P0, s.P1, s.P2, distance, 0,
				maxTValue, pointIndex, iter.tolerance,
			)
			if distance > prevD {
				contour.points =
					append(contour.points, s.P0, s.P1, s.P2)
				pointIndex += 3
			}
			prevP = s.P2
		case PathSegmentClose:
			haveSeenClose = true
		}

		if iter.iter.NextVerb() == PathVerbMove {
			break
		}
	}

	if math32.IsNaN(distance) || math32.IsInf(distance, 0) {
		return nil, false
	}

	if haveSeenClose {
		prevD := distance
		firstPt := contour.points[0]
		distance = contour.computeLineSeg(contour.points[pointIndex], firstPt, distance, pointIndex)
		if distance > prevD {
			contour.points = append(contour.points, firstPt)
		}
	}

	contour.length = distance
	contour.isClosed = haveSeenClose

	if len(contour.points) == 0 {
		return nil, false
	}
	return contour, true
}

type segmentType int

const (
	segmentTypeLine segmentType = iota
	segmentTypeQuad
	segmentTypeCubic
)

type segment struct {
	distance   float32 // total distance up to this point
	pointIndex int     // index into the contourMeasure::points array
	tValue     uint32
	kind       segmentType
}

func (s segment) scalarT() float32 {
	return float32(s.tValue) / float32(maxTValue)
}

type contourMeasure struct {
	segments []segment
	points   []Point
	length   float32
	isClosed bool
}

func (cm *contourMeasure) pushSegment(startD, stopD float32, startWithMoveTo bool, pb *PathBuilder) {
	if startD < 0.0 {
		startD = 0.0
	}
	if stopD > cm.length {
		stopD = cm.length
	}
	if !(startD <= stopD) {
		return
	}
	if len(cm.segments) == 0 {
		return
	}

	segIndex, startT, ok1 := cm.distanceToSegment(startD)
	if !ok1 {
		return
	}
	seg := cm.segments[segIndex]

	stopSegIndex, stopT, ok2 := cm.distanceToSegment(stopD)
	if !ok2 {
		return
	}
	stopSeg := cm.segments[stopSegIndex]

	if startWithMoveTo {
		pos, _ := computePosTan(
			cm.points[seg.pointIndex:],
			seg.kind,
			startT,
			Point{0, 0},
			Point{math32.NaN(), math32.NaN()},
		)
		pb.MoveTo(pos.X, pos.Y)
	}

	if seg.pointIndex == stopSeg.pointIndex {
		segmentTo(
			cm.points[seg.pointIndex:],
			seg.kind, startT, stopT, pb,
		)
	} else {
		currSegIndex := segIndex
		for {
			segmentTo(
				cm.points[seg.pointIndex:],
				seg.kind, startT, 1.0, pb,
			)

			oldPointIndex := seg.pointIndex
			for {
				currSegIndex++
				if cm.segments[currSegIndex].pointIndex != oldPointIndex {
					break
				}
			}
			seg = cm.segments[currSegIndex]
			startT = 0.0

			if seg.pointIndex >= stopSeg.pointIndex {
				break
			}
		}
		segmentTo(
			cm.points[seg.pointIndex:],
			seg.kind, 0.0, stopT, pb,
		)
	}
}

func (cm *contourMeasure) distanceToSegment(distance float32) (int, normalized.NormalizedF32, bool) {
	index := findSegment(cm.segments, distance)
	// don't care if we hit an exact match or not, so we xor index if it is negative
	index ^= index >> 31
	seg := cm.segments[index]

	// now interpolate t-values with the prev segment (if possible)
	var startT, startD float32 = 0.0, 0.0
	// check if the prev segment is legal, and references the same set of points
	if index > 0 {
		startD = cm.segments[index-1].distance
		if cm.segments[index-1].pointIndex == seg.pointIndex {
			startT = cm.segments[index-1].scalarT()
		}
	}

	t := startT + (seg.scalarT()-startT)*(distance-startD)/(seg.distance-startD)
	tt, flag := normalized.NewNormalizedF32(t)
	return index, tt, flag
}

func (cm *contourMeasure) computeLineSeg(p0, p1 Point, distance float32, pointIndex int) float32 {
	d := p0.Distance(p1)
	prevD := distance
	distance += d
	if distance > prevD {
		cm.segments = append(cm.segments, segment{
			distance:   distance,
			pointIndex: pointIndex,
			tValue:     maxTValue,
			kind:       segmentTypeLine,
		})
	}
	return distance
}

func (cm *contourMeasure) computeQuadSegs(p0, p1, p2 Point, distance float32, minT, maxT uint32, pointIndex int, tolerance float32) float32 {
	if tSpanBigEnough(maxT-minT) != 0 && quadTooCurvy(p0, p1, p2, tolerance) {
		tmp := [5]Point{}
		halfT := (minT + maxT) >> 1

		ChopQuadAt([3]Point{p0, p1, p2}, 0.5, &tmp)
		distance = cm.computeQuadSegs(
			tmp[0],
			tmp[1],
			tmp[2],
			distance,
			minT,
			halfT,
			pointIndex,
			tolerance,
		)
		distance = cm.computeQuadSegs(
			tmp[2],
			tmp[3],
			tmp[4],
			distance,
			halfT,
			maxT,
			pointIndex,
			tolerance,
		)
	} else {
		d := p0.Distance(p2)
		prevD := distance
		distance += d
		if distance > prevD {
			cm.segments = append(cm.segments, segment{
				distance:   distance,
				pointIndex: pointIndex,
				tValue:     maxT,
				kind:       segmentTypeQuad,
			})
		}
	}
	return distance
}

func (cm *contourMeasure) computeCubicSegs(p0, p1, p2, p3 Point, distance float32, minT, maxT uint32, pointIndex int, tolerance float32) float32 {
	if tSpanBigEnough(maxT-minT) != 0 && cubicTooCurvy(p0, p1, p2, p3, tolerance) {
		tmp := [7]Point{}
		halfT := (minT + maxT) >> 1

		ChopCubicAt2([4]Point{p0, p1, p2, p3}, 0.5, &tmp)
		distance = cm.computeCubicSegs(
			tmp[0],
			tmp[1],
			tmp[2],
			tmp[3],
			distance,
			minT,
			halfT,
			pointIndex,
			tolerance,
		)
		distance = cm.computeCubicSegs(
			tmp[3],
			tmp[4],
			tmp[5],
			tmp[6],
			distance,
			halfT,
			maxT,
			pointIndex,
			tolerance,
		)
	} else {
		d := p0.Distance(p3)
		prevD := distance
		distance += d
		if distance > prevD {
			cm.segments = append(cm.segments, segment{
				distance:   distance,
				pointIndex: pointIndex,
				tValue:     maxT,
				kind:       segmentTypeCubic,
			})
		}
	}
	return distance
}

func findSegment(base []segment, key float32) int {
	lo, hi := 0, len(base)-1
	for lo < hi {
		mid := (lo + hi) >> 1
		if base[mid].distance < key {
			lo = mid + 1
		} else {
			hi = mid
		}
	}

	if base[hi].distance < key {
		hi++
		hi = ^hi
	} else if key < base[hi].distance {
		hi = ^hi
	}
	return hi
}

func computePosTan(points []Point, kind segmentType, t normalized.NormalizedF32, pos, tangent Point) (Point, Point) {
	switch kind {
	case segmentTypeLine:
		if len(points) < 2 {
			return Point{math32.NaN(), math32.NaN()}, Point{math32.NaN(), math32.NaN()}
		}
		if pos.IsFinite() {
			pos.X = interpF32(points[0].X, points[1].X, t)
			pos.Y = interpF32(points[0].Y, points[1].Y, t)
		}
		if tangent.IsFinite() {
			tangent = Point{
				points[1].X - points[0].X,
				points[1].Y - points[0].Y,
			}
			tangent, _ = tangent.WithNormalizeFrom()
		}
	case segmentTypeQuad:
		if len(points) < 3 {
			return Point{math32.NaN(), math32.NaN()}, Point{math32.NaN(), math32.NaN()}
		}
		if pos.IsFinite() {
			pos = evalQuadAt([3]Point{points[0], points[1], points[2]}, normalized.NormalizedF32(t))
		}
		if tangent.IsFinite() {
			tangent, _ = evalQuadTangentAt([3]Point{points[0], points[1], points[2]}, t).WithNormalizeFrom()
		}
	case segmentTypeCubic:
		if len(points) < 4 {
			return Point{math32.NaN(), math32.NaN()}, Point{math32.NaN(), math32.NaN()}
		}
		if pos.IsFinite() {
			pos = evalCubicPosAt([4]Point{points[0], points[1], points[2], points[3]}, t)
		}
		if tangent.IsFinite() {
			tangent, _ = evalCubicTangentAt([4]Point{points[0], points[1], points[2], points[3]}, t).WithNormalizeFrom()
		}
	}
	return pos, tangent
}

func segmentTo(points []Point, kind segmentType, startT, stopT normalized.NormalizedF32, pb *PathBuilder) {
	if startT >= stopT {
		if pt, ok := pb.LastPoint(); ok {
			pb.LineTo(pt.X, pt.Y)
		}
		return
	}

	switch kind {
	case segmentTypeLine:
		if stopT == normalized.NormalizedF32One {
			pb.LineTo(points[1].X, points[1].Y)
		} else {
			pb.LineTo(
				interpF32(points[0].X, points[1].X, stopT),
				interpF32(points[0].Y, points[1].Y, stopT),
			)
		}
	case segmentTypeQuad:
		tmp0 := [5]Point{}
		tmp1 := [5]Point{}
		if startT == normalized.NormalizedF32Zero {
			if stopT == normalized.NormalizedF32One {
				pb.quadToPt(points[1], points[2])
			} else {
				stopTt := normalized.NewNormalizedF32ExclusiveWithBounded(stopT.Get())
				ChopQuadAt([3]Point{points[0], points[1], points[2]}, stopTt, &tmp0)
				pb.quadToPt(tmp0[1], tmp0[2])
			}
		} else {
			startTt := normalized.NewNormalizedF32ExclusiveWithBounded(startT.Get())
			ChopQuadAt([3]Point{points[0], points[1], points[2]}, startTt, &tmp0)
			if stopT == normalized.NormalizedF32One {
				pb.quadToPt(tmp0[3], tmp0[4])
			} else {
				newT := (stopT - startT) / (1.0 - startT)
				newTt := normalized.NewNormalizedF32ExclusiveWithBounded(newT.Get())
				ChopQuadAt([3]Point{tmp0[2], tmp0[3], tmp0[4]}, newTt, &tmp1)
				pb.quadToPt(tmp1[1], tmp1[2])
			}
		}
	case segmentTypeCubic:
		tmp0 := [7]Point{}
		tmp1 := [7]Point{}
		if startT == normalized.NormalizedF32Zero {
			if stopT == normalized.NormalizedF32One {
				pb.cubicToPt(points[1], points[2], points[3])
			} else {
				stopTt := normalized.NewNormalizedF32ExclusiveWithBounded(stopT.Get())
				ChopCubicAt2([4]Point{points[0], points[1], points[2], points[3]}, stopTt, &tmp0)
				pb.cubicToPt(tmp0[1], tmp0[2], tmp0[3])
			}
		} else {
			startTt := normalized.NewNormalizedF32ExclusiveWithBounded(startT.Get())
			ChopCubicAt2([4]Point{points[0], points[1], points[2], points[3]}, startTt, &tmp0)
			if stopT == normalized.NormalizedF32One {
				pb.cubicToPt(tmp0[4], tmp0[5], tmp0[6])
			} else {
				newT := (stopT - startT) / (1.0 - startT)
				newTt := normalized.NewNormalizedF32ExclusiveWithBounded(newT.Get())
				ChopCubicAt2([4]Point{tmp0[3], tmp0[4], tmp0[5], tmp0[6]}, newTt, &tmp1)
				pb.cubicToPt(tmp1[1], tmp1[2], tmp1[3])
			}
		}
	}
}

func tSpanBigEnough(tSpan uint32) uint32 {
	return tSpan >> 10
}

func quadTooCurvy(p0, p1, p2 Point, tolerance float32) bool {
	// diff = (a/4 + b/2 + c/4) - (a/2 + c/2)
	// diff = -a/4 + b/2 - c/4
	dx := p1.X*0.5 - (p0.X+p2.X)*0.25
	dy := p1.Y*0.5 - (p0.Y+p2.Y)*0.25

	dist := math32.Max(math32.Abs(dx), math32.Abs(dy))
	return dist > tolerance
}

func cubicTooCurvy(p0, p1, p2, p3 Point, tolerance float32) bool {
	n0 := cheapDistExceedsLimit(
		p1,
		interpSafeF32(p0.X, p3.X, 1.0/3.0),
		interpSafeF32(p0.Y, p3.Y, 1.0/3.0),
		tolerance,
	)
	n1 := cheapDistExceedsLimit(
		p2,
		interpSafeF32(p0.X, p3.X, 2.0/3.0),
		interpSafeF32(p0.Y, p3.Y, 2.0/3.0),
		tolerance,
	)
	return n0 || n1
}

func cheapDistExceedsLimit(pt Point, x, y, tolerance float32) bool {
	dist := math32.Max(math32.Abs(x-pt.X), math32.Abs(y-pt.Y))
	// just made up the 1/2
	return dist > tolerance
}

// Linearly interpolate between A and B, based on t.
// If t is 0, return A. If t is 1, return B else interpolate.
func interpF32(a, b float32, t normalized.NormalizedF32) float32 {
	return a + (b-a)*t.Get()
}

func interpSafeF32(a, b, t float32) float32 {
	return a + (b-a)*t
}
