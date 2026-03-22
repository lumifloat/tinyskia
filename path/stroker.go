// Copyright 2008 The Android Open Source Project
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Based on SkStroke.cpp
package path

import (
	"github.com/chewxy/math32"

	"github.com/lumifloat/tinyskia/internal/normalized"
	"github.com/lumifloat/tinyskia/internal/scalar"
)

type swappableBuilders struct {
	inner *PathBuilder
	outer *PathBuilder
}

func (s *swappableBuilders) swap() {
	// Skia swaps pointers to inner and outer builders during joining,
	// but not builders itself. So a simple `core::mem::swap` will produce invalid results.
	// And if we try to use use `core::mem::swap` on references, like below,
	// borrow checker will be unhappy.
	// That's why we need this wrapper. Maybe there is a better solution.
	s.inner, s.outer = s.outer, s.inner
}

// Stroke properties.
type Stroke struct {
	// A stroke thickness.
	// Must be >= 0.
	// When set to 0, a hairline stroking will be used.
	// Default: 1.0
	Width float32

	// The limit at which a sharp corner is drawn beveled.
	// Default: 4.0
	MiterLimit float32

	// A stroke line cap.
	// Default: LineCapButt
	LineCap LineCap

	// A stroke line join.
	// Default: LineJoinMiter
	LineJoin LineJoin

	// A stroke dashing properties.
	// Default: nil (no dashing)
	Dash interface{}
}

func defaultStroke() Stroke {
	return Stroke{
		Width:      1.0,
		MiterLimit: 4.0,
		LineCap:    LineCapButt,
		LineJoin:   LineJoinMiter,
		Dash:       nil,
	}
}

// LineCap draws at the beginning and end of an open path contour.
type LineCap int

const (
	// No stroke extension.
	LineCapButt LineCap = iota
	// Adds circle.
	LineCapRound
	// Adds square.
	LineCapSquare
)

// LineJoin specifies how corners are drawn when a shape is stroked.
// Join affects the four corners of a stroked rectangle, and the connected segments in a
// stroked path.
// Choose miter join to draw sharp corners. Choose round join to draw a circle with a
// radius equal to the stroke width on top of the corner. Choose bevel join to minimally
// connect the thick strokes.
// The fill path constructed to describe the stroked path respects the join setting but may
// not contain the actual join. For instance, a fill path constructed with round joins does
// not necessarily include circles at each connected segment.
type LineJoin int

const (
	// Extends to miter limit, then switches to bevel.
	LineJoinMiter LineJoin = iota
	// Extends to miter limit, then clips the corner.
	LineJoinMiterClip
	// Adds circle.
	LineJoinRound
	// Connects outside edges.
	LineJoinBevel
)

const quadRecursiveLimit = 3

// quads with extreme widths (e.g. (0,1) (1,6) (0,3) width=5e7) recurse to point of failure
// largest seen for normal cubics: 5, 26
// largest seen for normal quads: 11
var recursiveLimits = [4]int{5 * 3, 26 * 3, 11 * 3, 11 * 3} // 3x limits seen in practice

type capProc func(
	pivot Point,
	normal Point,
	stop Point,
	otherPath *PathBuilder,
	path *PathBuilder,
)

type joinProc func(
	beforeUnitNormal Point,
	pivot Point,
	afterUnitNormal Point,
	radius float32,
	invMiterLimit float32,
	prevIsLine bool,
	currIsLine bool,
	builders *swappableBuilders,
)

type reductionType int

const (
	reductionTypePoint       reductionType = iota // all curve points are practically identical
	reductionTypeLine                             // the control point is on the line between the ends
	reductionTypeQuad                             // the control point is outside the line between the ends
	reductionTypeDegenerate                       // the control point is on the line but outside the ends
	reductionTypeDegenerate2                      // two control points are on the line but outside ends (cubic)
	reductionTypeDegenerate3                      // three areas of max curvature found (for cubic)
)

type strokeType int

const (
	strokeTypeOuter strokeType = 1 // use sign-opposite values later to flip perpendicular axis
	strokeTypeInner strokeType = -1
)

type resultType int

const (
	resultTypeSplit      resultType = iota // the caller should split the quad stroke in two
	resultTypeDegenerate                   // the caller should add a line
	resultTypeQuad                         // the caller should (continue to try to) add a quad stroke
)

type intersectRayType int

const (
	intersectRayTypeCtrlPt intersectRayType = iota
	intersectRayTypeResultType
)

// Stroke
// If you plan stroking multiple paths, you can try using [`PathStroker`]
// which will preserve temporary allocations required during stroking.
// This might improve performance a bit.
func (p *Path) Stroke(stroke Stroke, resolutionScale float32) *Path {
	return newPathStroker().Stroke(p, stroke, resolutionScale)
}

// PathStroker a path stroker.
type PathStroker struct {
	radius             float32
	invMiterLimit      float32
	resScale           float32
	invResScale        float32
	invResScaleSquared float32

	firstNormal     Point
	prevNormal      Point
	firstUnitNormal Point
	prevUnitNormal  Point

	// on original path
	firstPt Point
	prevPt  Point

	firstOuterPt               Point
	firstOuterPtIndexInContour int
	segmentCount               int
	prevIsLine                 bool

	capper capProc
	joiner joinProc

	// outer is our working answer, inner is temp
	inner  *PathBuilder
	outer  *PathBuilder
	cusper *PathBuilder

	strokeType strokeType

	recursionDepth int  // track stack depth to abort if numerics run amok
	foundTangents  bool // do less work until tangents meet (cubic)
	joinCompleted  bool // previous join was not degenerate
}

func newPathStroker() *PathStroker {
	return &PathStroker{
		radius:             0.0,
		invMiterLimit:      0.0,
		resScale:           1.0,
		invResScale:        1.0,
		invResScaleSquared: 1.0,

		firstNormal:     Point{0, 0},
		prevNormal:      Point{0, 0},
		firstUnitNormal: Point{0, 0},
		prevUnitNormal:  Point{0, 0},

		firstPt: Point{0, 0},
		prevPt:  Point{0, 0},

		firstOuterPt:               Point{0, 0},
		firstOuterPtIndexInContour: 0,
		segmentCount:               -1,
		prevIsLine:                 false,

		capper: buttCapper,
		joiner: miterJoiner,

		inner:  NewPathBuilder(),
		outer:  NewPathBuilder(),
		cusper: NewPathBuilder(),

		strokeType: strokeTypeOuter,

		recursionDepth: 0,
		foundTangents:  false,
		joinCompleted:  false,
	}
}

// ComputeResolutionScale
// Computes a resolution scale.
// Resolution scale is the "intended" resolution for the output. Default is 1.0.
// Larger values (res > 1) indicate that the result should be more precise, since it will
// be zoomed up, and small errors will be magnified.
// Smaller values (0 < res < 1) indicate that the result can be less precise, since it will
// be zoomed down, and small errors may be invisible.
func ComputeResolutionScale(ts Transform) float32 {
	sx := Point{ts.SX, ts.KX}.Length()
	sy := Point{ts.KY, ts.SY}.Length()
	if !math32.IsNaN(sx) && !math32.IsInf(sx, 0) &&
		!math32.IsNaN(sy) && !math32.IsInf(sy, 0) {
		scale := math32.Max(sx, sy)
		if scale > 0.0 {
			return scale
		}
	}

	return 1.0
}

// Stroke the path.
// Can be called multiple times to reuse allocated buffers.
// `resolution_scale` can be obtained via
func (s *PathStroker) Stroke(path *Path, stroke Stroke, resolutionScale float32) *Path {
	if stroke.Width <= 0 {
		return nil
	}
	return s.strokeInner(
		path,
		stroke.Width,
		stroke.MiterLimit,
		stroke.LineCap,
		stroke.LineJoin,
		resolutionScale,
	)
}

func (s *PathStroker) strokeInner(
	path *Path,
	width float32,
	miterLimit float32,
	lineCap LineCap,
	lineJoin LineJoin,
	resScale float32,
) *Path {
	// TODO: stroke_rect optimization

	var invMiterLimit float32 = 0.0

	if lineJoin == LineJoinMiter {
		if miterLimit <= 1.0 {
			lineJoin = LineJoinBevel
		} else {
			invMiterLimit = 1.0 / miterLimit
		}
	}

	if lineJoin == LineJoinMiterClip {
		invMiterLimit = 1.0 / miterLimit
	}

	s.resScale = resScale
	// The '4' below matches the fill scan converter's error term.
	s.invResScale = 1.0 / (resScale * 4.0)
	s.invResScaleSquared = s.invResScale * s.invResScale

	s.radius = width * 0.5
	s.invMiterLimit = invMiterLimit

	s.firstNormal = Point{0, 0}
	s.prevNormal = Point{0, 0}
	s.firstUnitNormal = Point{0, 0}
	s.prevUnitNormal = Point{0, 0}

	s.firstPt = Point{0, 0}
	s.prevPt = Point{0, 0}

	s.firstOuterPt = Point{0, 0}
	s.firstOuterPtIndexInContour = 0
	s.segmentCount = -1
	s.prevIsLine = false

	s.capper = capFactory(lineCap)
	s.joiner = joinFactory(lineJoin)

	s.inner.Clear()
	// reserve logic simplified for Go slices
	s.outer.Clear()
	s.cusper.Clear()

	s.strokeType = strokeTypeOuter
	s.recursionDepth = 0
	s.foundTangents = false
	s.joinCompleted = false

	lastSegmentIsLine := false
	iter := path.Segments()
	iter.SetAutoClose(true)

	for segment := iter.Next(); segment != nil; segment = iter.Next() {
		switch seg := segment.(type) {
		case PathSegmentMoveTo:
			s.moveTo(Point(seg))
		case PathSegmentLineTo:
			s.lineTo(Point(seg), iter)
			lastSegmentIsLine = true
		case PathSegmentQuadTo:
			s.quadTo(Point(seg.P0), Point(seg.P1))
			lastSegmentIsLine = false
		case PathSegmentCubicTo:
			s.cubicTo(Point(seg.P0), Point(seg.P1), Point(seg.P2))
			lastSegmentIsLine = false
		case PathSegmentClose:
			if lineCap != LineCapButt {
				if s.hasOnlyMoveTo() {
					s.lineTo(s.moveToPt(), nil)
					lastSegmentIsLine = true
					continue
				}

				if s.isCurrentContourEmpty() {
					lastSegmentIsLine = true
					continue
				}
			}
			s.close(lastSegmentIsLine)
		}
	}

	return s.finish(lastSegmentIsLine)
}

func (s *PathStroker) builders() *swappableBuilders {
	return &swappableBuilders{
		inner: s.inner,
		outer: s.outer,
	}
}

func (s *PathStroker) moveToPt() Point {
	return s.firstPt
}

func (s *PathStroker) moveTo(p Point) {
	if s.segmentCount > 0 {
		s.finishContour(false, false)
	}

	s.segmentCount = 0
	s.firstPt = p
	s.prevPt = p
	s.joinCompleted = false
}

func (s *PathStroker) lineTo(p Point, iter *PathSegmentsIter) {
	teenyLine := s.prevPt.EqualsWithinTolerance(p, scalar.SCALAR_NEARLY_ZERO*s.invResScale)

	// Note: In Go we compare function pointers directly if possible,
	// or use a flag if Capper is complex.
	if s.capper != nil && teenyLine {
		return
	}

	if teenyLine && (s.joinCompleted || (iter != nil && iter.hasValidTangent())) {
		return
	}

	var normal, unitNormal, ok = s.preJoinTo(p, true)
	if !ok {
		return
	}

	s.outer.LineTo(p.X+normal.X, p.Y+normal.Y)
	s.inner.LineTo(p.X-normal.X, p.Y-normal.Y)

	s.postJoinTo(p, normal, unitNormal)
}

func (s *PathStroker) quadTo(p1, p2 Point) {
	quad := [3]Point{s.prevPt, p1, p2}
	reduction, reductionType := checkQuadLinear(quad)

	if reductionType == reductionTypePoint || reductionType == reductionTypeLine {
		s.lineTo(p2, nil)
		return
	}

	if reductionType == reductionTypeDegenerate {
		s.lineTo(reduction, nil)
		saveJoiner := s.joiner
		s.joiner = roundJoiner
		s.lineTo(p2, nil)
		s.joiner = saveJoiner
		return
	}

	var normalAB, unitAB, ok1 = s.preJoinTo(p1, false)
	if !ok1 {
		s.lineTo(p2, nil)
		return
	}

	var quadPoints quadConstruct
	s.initQuad(strokeTypeOuter, 0.0, 1.0, &quadPoints)
	s.quadStroke(quad, &quadPoints)
	s.initQuad(strokeTypeInner, 0.0, 1.0, &quadPoints)
	s.quadStroke(quad, &quadPoints)

	var normalBC, unitBC, ok2 = withNormalUnitNormal(quad[1], quad[2], s.resScale, s.radius)
	if !ok2 {
		normalBC = normalAB
		unitBC = unitAB
	}

	s.postJoinTo(p2, normalBC, unitBC)
}

func (s *PathStroker) cubicTo(pt1, pt2, pt3 Point) {
	cubic := [4]Point{s.prevPt, pt1, pt2, pt3}
	reductionType, reduction, tangentPt := checkCubicLinear(cubic)

	if reductionType == reductionTypePoint {
		// If the stroke consists of a moveTo followed by a degenerate curve, treat it
		// as if it were followed by a zero-length line. Lines without length
		// can have square and round end caps.
		s.lineTo(pt3, nil)
		return
	}

	if reductionType == reductionTypeLine {
		s.lineTo(pt3, nil)
		return
	}

	if reductionType >= reductionTypeDegenerate &&
		reductionType <= reductionTypeDegenerate3 {
		s.lineTo(reduction[0], nil)
		saveJoiner := s.joiner
		s.joiner = roundJoiner
		if reductionType >= reductionTypeDegenerate2 {
			s.lineTo(reduction[1], nil)
		}
		if reductionType == reductionTypeDegenerate3 {
			s.lineTo(reduction[2], nil)
		}
		s.lineTo(pt3, nil)
		s.joiner = saveJoiner
		return
	}

	var normalAB, unitAB, normalCD, unitCD Point
	var ok bool
	normalAB, unitAB, ok = s.preJoinTo(tangentPt, false)
	if !ok {
		s.lineTo(pt3, nil)
		return
	}

	tmp := [3]normalized.NormalizedF32Exclusive{}
	tValues := findCubicInflections(cubic, &tmp)
	var lastT = normalized.NormalizedF32Zero

	for i := 0; i <= len(tValues); i++ {
		var nextT normalized.NormalizedF32
		if i < len(tValues) {
			nextT = tValues[i].ToNormalized()
		} else {
			nextT = normalized.NormalizedF32One
		}

		var quadPoints quadConstruct
		s.initQuad(strokeTypeOuter, lastT, nextT, &quadPoints)
		s.cubicStroke(cubic, &quadPoints)
		s.initQuad(strokeTypeInner, lastT, nextT, &quadPoints)
		s.cubicStroke(cubic, &quadPoints)
		lastT = nextT
	}

	if cusp, found := findCubicCusp(cubic); found {
		cuspLoc := evalCubicPosAt(cubic, cusp.ToNormalized())
		s.cusper.PushCircle(cuspLoc.X, cuspLoc.Y, s.radius)
	}

	normalCD, unitCD = s.withCubicEndNormal(cubic, normalAB, unitAB)
	s.postJoinTo(pt3, normalCD, unitCD)
}

func (s *PathStroker) cubicStroke(cubic [4]Point, quadPoints *quadConstruct) bool {
	if !s.foundTangents {
		resultType := s.tangentsMeet(cubic, quadPoints)
		if resultType != resultTypeQuad {
			ok := pointsWithinDist(quadPoints.quad[0], quadPoints.quad[2], s.invResScale)
			if (resultType == resultTypeDegenerate || ok) && s.cubicMidOnLine(cubic, quadPoints) {
				s.addDegenerateLine(quadPoints)
				return true
			}
		} else {
			s.foundTangents = true
		}
	}

	if s.foundTangents {
		resultType := s.compareQuadCubic(cubic, quadPoints)
		if resultType == resultTypeQuad {
			stroke := quadPoints.quad
			if s.strokeType == strokeTypeOuter {
				s.outer.QuadTo(stroke[1].X, stroke[1].Y, stroke[2].X, stroke[2].Y)
			} else {
				s.inner.QuadTo(stroke[1].X, stroke[1].Y, stroke[2].X, stroke[2].Y)
			}
			return true
		}

		if resultType == resultTypeDegenerate {
			if !quadPoints.oppositeTangents {
				s.addDegenerateLine(quadPoints)
				return true
			}
		}
	}

	if math32.IsNaN(quadPoints.quad[2].X) || math32.IsInf(quadPoints.quad[2].X, 0) {
		return false
	}

	s.recursionDepth++
	limitIdx := 0
	if s.foundTangents {
		limitIdx = 1
	}
	if s.recursionDepth > recursiveLimits[limitIdx] {
		return false
	}

	var half quadConstruct
	if !half.initWithStart(quadPoints) {
		s.addDegenerateLine(quadPoints)
		s.recursionDepth--
		return true
	}

	if !s.cubicStroke(cubic, &half) {
		return false
	}

	if !half.initWithEnd(quadPoints) {
		s.addDegenerateLine(quadPoints)
		s.recursionDepth--
		return true
	}

	if !s.cubicStroke(cubic, &half) {
		return false
	}

	s.recursionDepth--
	return true
}

func (s *PathStroker) cubicMidOnLine(cubic [4]Point, quadPoints *quadConstruct) bool {
	var strokeMid Point
	s.cubicQuadMid(cubic, quadPoints, &strokeMid)
	dist := ptToLine(strokeMid, quadPoints.quad[0], quadPoints.quad[2])
	return dist < s.invResScaleSquared
}

func (s *PathStroker) cubicQuadMid(cubic [4]Point, quadPoints *quadConstruct, mid *Point) {
	var cubicMidPt Point
	s.cubicPerpRay(cubic, quadPoints.midT, &cubicMidPt, mid, nil)
}

func (s *PathStroker) cubicPerpRay(
	cubic [4]Point,
	t normalized.NormalizedF32,
	tPt *Point,
	onPt *Point,
	tangent *Point,
) {
	*tPt = evalCubicPosAt(cubic, t)
	dxy := evalCubicTangentAt(cubic, t)

	if dxy.X == 0.0 && dxy.Y == 0.0 {
		var cPoints *[4]Point
		if t.Get() <= scalar.SCALAR_NEARLY_ZERO {
			dxy = cubic[2].Sub(cubic[0])
		} else if t.Get() >= 1.0-scalar.SCALAR_NEARLY_ZERO {
			dxy = cubic[3].Sub(cubic[1])
		} else {
			// If the cubic inflection falls on the cusp, subdivide the cubic
			// to find the tangent at that point.
			var chopped [7]Point
			ChopCubicAt2(cubic, normalized.NormalizedF32Exclusive(t), &chopped)
			dxy = chopped[3].Sub(chopped[2])
			if dxy.X == 0.0 && dxy.Y == 0.0 {
				dxy = chopped[3].Sub(chopped[1])
				// Use first 4 points of chopped as fallback
				temp := [4]Point{chopped[0], chopped[1], chopped[2], chopped[3]}
				cPoints = &temp
			}
		}

		if dxy.X == 0.0 && dxy.Y == 0.0 {
			if cPoints != nil {
				dxy = cPoints[3].Sub(cPoints[0])
			} else {
				dxy = cubic[3].Sub(cubic[0])
			}
		}
	}

	s.setRayPoints(*tPt, &dxy, onPt, tangent)
}

func (s *PathStroker) withCubicEndNormal(
	cubic [4]Point,
	normalAB Point,
	unitNormalAB Point,
) (normalCD Point, unitNormalCD Point) {
	ab := cubic[1].Sub(cubic[0])
	cd := cubic[3].Sub(cubic[2])

	degenerateAB := degenerateVector(ab)
	degenerateCB := degenerateVector(cd)

	if degenerateAB && degenerateCB {
		normalCD = normalAB
		unitNormalCD = unitNormalAB
		return
	}

	if degenerateAB {
		ab = cubic[2].Sub(cubic[0])
		degenerateAB = degenerateVector(ab)
	}

	if degenerateCB {
		cd = cubic[3].Sub(cubic[1])
		degenerateCB = degenerateVector(cd)
	}

	if degenerateAB || degenerateCB {
		normalCD = normalAB
		unitNormalCD = unitNormalAB
		return
	}

	normalCD, unitNormalCD, _ = withNormalUnitNormal2(cd, s.radius)
	return
}

func (s *PathStroker) compareQuadCubic(
	cubic [4]Point,
	quadPoints *quadConstruct,
) resultType {
	// get the quadratic approximation of the stroke
	s.cubicQuadEnds(cubic, quadPoints)
	resultType := s.intersectRay(intersectRayTypeCtrlPt, quadPoints)
	if resultType != resultTypeQuad {
		return resultType
	}

	// project a ray from the curve to the stroke
	// points near midpoint on quad, midpoint on cubic
	var ray0, ray1 Point
	s.cubicPerpRay(cubic, quadPoints.midT, &ray1, &ray0, nil)

	// In Go, slices are already references, and Point is a value type, so clone() is just a copy
	quadCopy := quadPoints.quad
	return s.strokeCloseEnough(quadCopy[:], []Point{ray0, ray1}, quadPoints)
}

// Given a cubic and a t range, find the start and end if they haven't been found already.
func (s *PathStroker) cubicQuadEnds(cubic [4]Point, quadPoints *quadConstruct) {
	if !quadPoints.startSet {
		var cubicStartPt Point
		s.cubicPerpRay(
			cubic,
			quadPoints.startT,
			&cubicStartPt,
			&quadPoints.quad[0],
			&quadPoints.tangentStart,
		)
		quadPoints.startSet = true
	}

	if !quadPoints.endSet {
		var cubicEndPt Point
		s.cubicPerpRay(
			cubic,
			quadPoints.endT,
			&cubicEndPt,
			&quadPoints.quad[2],
			&quadPoints.tangentEnd,
		)
		quadPoints.endSet = true
	}
}

func (s *PathStroker) close(isLine bool) {
	s.finishContour(true, isLine)
}

func (s *PathStroker) finishContour(close bool, currIsLine bool) {
	if s.segmentCount > 0 {
		if close {
			s.joiner(
				s.prevUnitNormal,
				s.prevPt,
				s.firstUnitNormal,
				s.radius,
				s.invMiterLimit,
				s.prevIsLine,
				currIsLine,
				s.builders(),
			)
			s.outer.Close()

			// now add inner as its own contour
			pt, _ := s.inner.LastPoint()
			s.outer.MoveTo(pt.X, pt.Y)
			s.outer.reversePathTo(s.inner)
			s.outer.Close()
		} else {
			// add caps to start and end

			// cap the end
			pt, _ := s.inner.LastPoint()
			var otherPath *PathBuilder
			if currIsLine {
				otherPath = s.inner
			}
			s.capper(
				s.prevPt,
				s.prevNormal,
				pt,
				otherPath,
				s.outer,
			)
			s.outer.reversePathTo(s.inner)

			// cap the start
			if s.prevIsLine {
				otherPath = s.inner
			} else {
				otherPath = nil
			}
			s.capper(
				s.firstPt,
				s.firstNormal.Neg(),
				s.firstOuterPt,
				otherPath,
				s.outer,
			)
			s.outer.Close()
		}

		if !s.cusper.IsEmpty() {
			s.outer.pushPathBuilder(s.cusper)
			s.cusper.Clear()
		}
	}

	// since we may re-use `inner`, we rewind instead of reset, to save on
	// reallocating its internal storage.
	s.inner.Clear()
	s.segmentCount = -1
	s.firstOuterPtIndexInContour = len(s.outer.points)
}

func (s *PathStroker) preJoinTo(
	p Point,
	currIsLine bool,
) (normal Point, unitNormal Point, ok bool) {
	prevX := s.prevPt.X
	prevY := s.prevPt.Y

	normal, unitNormal, ok = withNormalUnitNormal(
		s.prevPt,
		p,
		s.resScale,
		s.radius,
	)

	if !ok {
		if s.capper != nil {
			return
		}

		// Square caps and round caps draw even if the segment length is zero.
		// Since the zero length segment has no direction, set the orientation
		// to upright as the default orientation.
		normal = Point{s.radius, 0.0}
		unitNormal = Point{1.0, 0.0}
	}

	if s.segmentCount == 0 {
		s.firstNormal = normal
		s.firstUnitNormal = unitNormal
		s.firstOuterPt = Point{prevX + normal.X, prevY + normal.Y}

		s.outer.MoveTo(s.firstOuterPt.X, s.firstOuterPt.Y)
		s.inner.MoveTo(prevX-normal.X, prevY-normal.Y)
	} else {
		// we have a previous segment
		s.joiner(
			s.prevUnitNormal,
			s.prevPt,
			unitNormal,
			s.radius,
			s.invMiterLimit,
			s.prevIsLine,
			currIsLine,
			s.builders(),
		)
	}
	s.prevIsLine = currIsLine
	return normal, unitNormal, true
}

func (s *PathStroker) postJoinTo(p Point, normal Point, unitNormal Point) {
	s.joinCompleted = true
	s.prevPt = p
	s.prevUnitNormal = unitNormal
	s.prevNormal = normal
	s.segmentCount += 1
}

func (s *PathStroker) initQuad(
	strokeType strokeType,
	start normalized.NormalizedF32,
	end normalized.NormalizedF32,
	quadPoints *quadConstruct,
) {
	s.strokeType = strokeType
	s.foundTangents = false
	quadPoints.init(start, end)
}

func (s *PathStroker) quadStroke(quad [3]Point, quadPoints *quadConstruct) bool {
	resultType := s.compareQuadQuad(quad, quadPoints)
	if resultType == resultTypeQuad {
		var path *PathBuilder
		if s.strokeType == strokeTypeOuter {
			path = s.outer
		} else {
			path = s.inner
		}

		path.QuadTo(
			quadPoints.quad[1].X,
			quadPoints.quad[1].Y,
			quadPoints.quad[2].X,
			quadPoints.quad[2].Y,
		)

		return true
	}

	if resultType == resultTypeDegenerate {
		s.addDegenerateLine(quadPoints)
		return true
	}

	s.recursionDepth += 1
	if s.recursionDepth > recursiveLimits[quadRecursiveLimit] {
		return false // just abort if projected quad isn't representable
	}

	var half quadConstruct
	half.initWithStart(quadPoints)
	if !s.quadStroke(quad, &half) {
		return false
	}

	half.initWithEnd(quadPoints)
	if !s.quadStroke(quad, &half) {
		return false
	}

	s.recursionDepth -= 1
	return true
}

func (s *PathStroker) compareQuadQuad(
	quad [3]Point,
	quadPoints *quadConstruct,
) resultType {
	// get the quadratic approximation of the stroke
	if !quadPoints.startSet {
		var quadStartPt Point
		s.quadPerpRay(
			quad[:],
			quadPoints.startT,
			&quadStartPt,
			&quadPoints.quad[0],
			&quadPoints.tangentStart,
		)
		quadPoints.startSet = true
	}

	if !quadPoints.endSet {
		var quadEndPt Point
		s.quadPerpRay(
			quad[:],
			quadPoints.endT,
			&quadEndPt,
			&quadPoints.quad[2],
			&quadPoints.tangentEnd,
		)
		quadPoints.endSet = true
	}

	resultType := s.intersectRay(intersectRayTypeCtrlPt, quadPoints)
	if resultType != resultTypeQuad {
		return resultType
	}

	// project a ray from the curve to the stroke
	var ray0, ray1 Point
	s.quadPerpRay(quad[:], quadPoints.midT, &ray1, &ray0, nil)
	quadCopy := quadPoints.quad
	return s.strokeCloseEnough(quadCopy[:], []Point{ray0, ray1}, quadPoints)
}

// Given a point on the curve and its derivative, scale the derivative by the radius, and
// compute the perpendicular point and its tangent.
func (s *PathStroker) setRayPoints(
	tp Point,
	dxy *Point,
	onP *Point,
	tangent *Point,
) {
	var ok bool
	*dxy, ok = dxy.WithLengthFrom(s.radius)
	if !ok {
		*dxy = Point{X: s.radius, Y: 0.0}
	}

	axisFlip := float32(s.strokeType) // go opposite ways for outer, inner
	onP.X = tp.X + axisFlip*dxy.Y
	onP.Y = tp.Y - axisFlip*dxy.X

	if tangent != nil {
		tangent.X = onP.X + dxy.X
		tangent.Y = onP.Y + dxy.Y
	}
}

// Given a quad and t, return the point on curve,
// its perpendicular, and the perpendicular tangent.
func (s *PathStroker) quadPerpRay(
	quad []Point,
	t normalized.NormalizedF32,
	tp *Point,
	onP *Point,
	tangent *Point,
) {
	*tp = evalQuadAt([3]Point{quad[0], quad[1], quad[2]}, t)
	dxy := evalQuadTangentAt([3]Point{quad[0], quad[1], quad[2]}, t)

	if dxy.IsZero() {
		dxy = quad[2].Sub(quad[0])
	}

	s.setRayPoints(*tp, &dxy, onP, tangent)
}

func (s *PathStroker) addDegenerateLine(quadPoints *quadConstruct) {
	if s.strokeType == strokeTypeOuter {
		s.outer.LineTo(quadPoints.quad[2].X, quadPoints.quad[2].Y)
	} else {
		s.inner.LineTo(quadPoints.quad[2].X, quadPoints.quad[2].Y)
	}
}

func (s *PathStroker) strokeCloseEnough(
	stroke []Point,
	ray []Point,
	quadPoints *quadConstruct,
) resultType {
	strokeMid := evalQuadAt([3]Point{stroke[0], stroke[1], stroke[2]}, 0.5)
	// measure the distance from the curve to the quad-stroke midpoint, compare to radius
	if pointsWithinDist(ray[0], strokeMid, s.invResScale) {
		// if the difference is small
		if sharpAngle([3]Point{quadPoints.quad[0], quadPoints.quad[1], quadPoints.quad[2]}) {
			return resultTypeSplit
		}

		return resultTypeQuad
	}

	// measure the distance to quad's bounds (quick reject)
	if !ptInQuadBounds([3]Point{stroke[0], stroke[1], stroke[2]}, ray[0], s.invResScale) {
		// if far, subdivide
		return resultTypeSplit
	}

	// measure the curve ray distance to the quad-stroke
	var tmp [3]normalized.NormalizedF32Exclusive
	roots := intersectQuadRay([2]Point{ray[0], ray[1]}, [3]Point{stroke[0], stroke[1], stroke[2]}, &tmp)
	if len(roots) != 1 {
		return resultTypeSplit
	}

	quadPt := evalQuadAt([3]Point{stroke[0], stroke[1], stroke[2]}, normalized.NormalizedF32(roots[0]))
	errorVal := s.invResScale * (1.0 - math32.Abs(roots[0].Get()-0.5)*2.0)
	if pointsWithinDist(ray[0], quadPt, errorVal) {
		// if the difference is small, we're done
		if sharpAngle([3]Point{quadPoints.quad[0], quadPoints.quad[1], quadPoints.quad[2]}) {
			return resultTypeSplit
		}

		return resultTypeQuad
	}

	// otherwise, subdivide
	return resultTypeSplit
}

// Find the intersection of the stroke tangents to construct a stroke quad.
func (s *PathStroker) intersectRay(
	intersectRayType intersectRayType,
	quadPoints *quadConstruct,
) resultType {
	start := quadPoints.quad[0]
	end := quadPoints.quad[2]
	aLen := quadPoints.tangentStart.Sub(start)
	bLen := quadPoints.tangentEnd.Sub(end)

	denom := aLen.Cross(bLen)
	if denom == 0.0 || math32.IsNaN(denom) || math32.IsInf(denom, 0) {
		quadPoints.oppositeTangents = aLen.Dot(bLen) < 0.0
		return resultTypeDegenerate
	}

	quadPoints.oppositeTangents = false
	ab0 := start.Sub(end)
	numerA := bLen.Cross(ab0)
	numerB := aLen.Cross(ab0)
	if (numerA >= 0.0) == (numerB >= 0.0) {
		// if the control point is outside the quad ends
		dist1 := ptToLine(start, end, quadPoints.tangentEnd)
		dist2 := ptToLine(end, start, quadPoints.tangentStart)
		if math32.Max(dist1, dist2) <= s.invResScaleSquared {
			return resultTypeDegenerate
		}

		return resultTypeSplit
	}

	// check to see if the denominator is teeny relative to the numerator
	numerA /= denom
	validDivide := numerA > numerA-1.0
	if validDivide {
		if intersectRayType == intersectRayTypeCtrlPt {
			quadPoints.quad[1].X = start.X*(1.0-numerA) + quadPoints.tangentStart.X*numerA
			quadPoints.quad[1].Y = start.Y*(1.0-numerA) + quadPoints.tangentStart.Y*numerA
		}

		return resultTypeQuad
	}

	quadPoints.oppositeTangents = aLen.Dot(bLen) < 0.0
	return resultTypeDegenerate
}

// Given a cubic and a t-range, determine if the stroke can be described by a quadratic.
func (s *PathStroker) tangentsMeet(cubic [4]Point, quadPoints *quadConstruct) resultType {
	s.cubicQuadEnds(cubic, quadPoints)
	return s.intersectRay(intersectRayTypeResultType, quadPoints)
}

func (s *PathStroker) finish(isLine bool) *Path {
	s.finishContour(false, isLine)

	// Swap out the outer builder.
	buf := s.outer
	s.outer = NewPathBuilder()

	return buf.Finish()
}

func (s *PathStroker) hasOnlyMoveTo() bool {
	return s.segmentCount == 0
}

func (s *PathStroker) isCurrentContourEmpty() bool {
	return s.inner.isZeroLengthSincePoint(0) &&
		s.outer.isZeroLengthSincePoint(int(s.firstOuterPtIndexInContour))
}

func capFactory(cap LineCap) capProc {
	switch cap {
	case LineCapButt:
		return buttCapper
	case LineCapRound:
		return roundCapper
	case LineCapSquare:
		return squareCapper
	default:
		return buttCapper
	}
}

func buttCapper(_ Point, _ Point, stop Point, _ *PathBuilder, path *PathBuilder) {
	path.LineTo(stop.X, stop.Y)
}

func roundCapper(
	pivot Point,
	normal Point,
	stop Point,
	_ *PathBuilder,
	path *PathBuilder,
) {
	parallel := normal.WithRotateCWFrom()

	projectedCenter := pivot.Add(parallel)

	path.conicPointsTo(
		projectedCenter.Add(normal),
		projectedCenter,
		scalar.SCALAR_ROOT_2_OVER_2,
	)
	path.conicPointsTo(projectedCenter.Sub(normal), stop, scalar.SCALAR_ROOT_2_OVER_2)
}

func squareCapper(
	pivot Point,
	normal Point,
	stop Point,
	otherPath *PathBuilder,
	path *PathBuilder,
) {
	parallel := normal.WithRotateCWFrom()

	if otherPath != nil {
		path.setLastPoint(Point{
			pivot.X + normal.X + parallel.X,
			pivot.Y + normal.Y + parallel.Y,
		})
		path.LineTo(
			pivot.X-normal.X+parallel.X,
			pivot.Y-normal.Y+parallel.Y,
		)
	} else {
		path.LineTo(
			pivot.X+normal.X+parallel.X,
			pivot.Y+normal.Y+parallel.Y,
		)
		path.LineTo(
			pivot.X-normal.X+parallel.X,
			pivot.Y-normal.Y+parallel.Y,
		)
		path.LineTo(stop.X, stop.Y)
	}
}

func joinFactory(join LineJoin) joinProc {
	switch join {
	case LineJoinMiter:
		return miterJoiner
	case LineJoinMiterClip:
		return miterClipJoiner
	case LineJoinRound:
		return roundJoiner
	case LineJoinBevel:
		return bevelJoiner
	default:
		return bevelJoiner
	}
}

func isClockwise(before Point, after Point) bool {
	return before.X*after.Y > before.Y*after.X
}

type angleType int

const (
	angleTypeNearly180 angleType = iota
	angleTypeSharp
	angleTypeShallow
	angleTypeNearlyLine
)

func dotToAngleType(dot float32) angleType {
	if dot >= 0.0 {
		// shallow or line
		if scalar.IsNearlyZero(1.0 - dot) {
			return angleTypeNearlyLine
		} else {
			return angleTypeShallow
		}
	} else {
		// sharp or 180
		if scalar.IsNearlyZero(1.0 + dot) {
			return angleTypeNearly180
		} else {
			return angleTypeSharp
		}
	}
}

func handleInnerJoin(pivot Point, after Point, inner *PathBuilder) {
	// In the degenerate case that the stroke radius is larger than our segments
	// just connecting the two inner segments may "show through" as a funny
	// diagonal. To pseudo-fix this, we go through the pivot point. This adds
	// an extra point/edge, but I can't see a cheap way to know when this is
	// not needed :(
	inner.LineTo(pivot.X, pivot.Y)

	inner.LineTo(pivot.X-after.X, pivot.Y-after.Y)
}

func bevelJoiner(
	beforeUnitNormal Point,
	pivot Point,
	afterUnitNormal Point,
	radius float32,
	_ float32,
	_ bool,
	_ bool,
	builders *swappableBuilders,
) {
	after := afterUnitNormal.WithScaleFrom(radius)

	if !isClockwise(beforeUnitNormal, afterUnitNormal) {
		builders.swap()
		after = after.Neg()
	}

	builders.outer.LineTo(pivot.X+after.X, pivot.Y+after.Y)
	handleInnerJoin(pivot, after, builders.inner)
}

func roundJoiner(
	beforeUnitNormal Point,
	pivot Point,
	afterUnitNormal Point,
	radius float32,
	_ float32,
	_ bool,
	_ bool,
	builders *swappableBuilders,
) {
	dotProd := beforeUnitNormal.Dot(afterUnitNormal)
	angleType := dotToAngleType(dotProd)

	if angleType == angleTypeNearlyLine {
		return
	}

	before := beforeUnitNormal
	after := afterUnitNormal
	dir := pathDirectionCW

	if !isClockwise(before, after) {
		builders.swap()
		before = before.Neg()
		after = after.Neg()
		dir = pathDirectionCCW
	}

	ts := TransformFromRow(radius, 0.0, 0.0, radius, pivot.X, pivot.Y)

	var conics [5]conic
	conicsSlice := BuildUnitArc(before, after, dir, ts, &conics)
	if conicsSlice != nil {
		for _, conic := range conicsSlice {
			builders.outer.conicPointsTo(conic.Points[1], conic.Points[2], conic.Weight)
		}

		scaledAfter := after.WithScaleFrom(radius)
		handleInnerJoin(pivot, scaledAfter, builders.inner)
	}
}

func miterJoiner(
	beforeUnitNormal Point,
	pivot Point,
	afterUnitNormal Point,
	radius float32,
	invMiterLimit float32,
	prevIsLine bool,
	currIsLine bool,
	builders *swappableBuilders,
) {
	miterJoinerInner(
		beforeUnitNormal,
		pivot,
		afterUnitNormal,
		radius,
		invMiterLimit,
		false,
		prevIsLine,
		currIsLine,
		*builders,
	)
}

func miterClipJoiner(
	beforeUnitNormal Point,
	pivot Point,
	afterUnitNormal Point,
	radius float32,
	invMiterLimit float32,
	prevIsLine bool,
	currIsLine bool,
	builders *swappableBuilders,
) {
	miterJoinerInner(
		beforeUnitNormal,
		pivot,
		afterUnitNormal,
		radius,
		invMiterLimit,
		true,
		prevIsLine,
		currIsLine,
		*builders,
	)
}

func miterJoinerInner(
	beforeUnitNormal Point,
	pivot Point,
	afterUnitNormal Point,
	radius float32,
	invMiterLimit float32,
	miterClip bool,
	prevIsLine bool,
	currIsLine bool,
	builders swappableBuilders,
) {
	doBluntOrClipped := func(
		builders swappableBuilders,
		pivot Point,
		radius float32,
		prevIsLine bool,
		currIsLine bool,
		before Point,
		mid Point,
		after Point,
		invMiterLimit float32,
		miterClip bool,
	) {
		after = after.WithScaleFrom(radius)

		if miterClip {
			mid, _ = mid.WithNormalizeFrom()

			cosBeta := before.Dot(mid)
			sinBeta := before.Cross(mid)

			var x float32
			if math32.Abs(sinBeta) <= scalar.SCALAR_NEARLY_ZERO {
				x = 1.0 / invMiterLimit
			} else {
				x = ((1.0 / invMiterLimit) - cosBeta) / sinBeta
			}

			before = before.WithScaleFrom(radius)

			beforeTangent := before.WithRotateCWFrom()

			afterTangent := after.WithRotateCCWFrom()

			c1 := pivot.Add(before).Add(beforeTangent.WithScaleFrom(x))
			c2 := pivot.Add(after).Add(afterTangent.WithScaleFrom(x))

			if prevIsLine {
				builders.outer.setLastPoint(c1)
			} else {
				builders.outer.LineTo(c1.X, c1.Y)
			}

			builders.outer.LineTo(c2.X, c2.Y)
		}

		if !currIsLine {
			builders.outer.LineTo(pivot.X+after.X, pivot.Y+after.Y)
		}

		handleInnerJoin(pivot, after, builders.inner)
	}

	doMiter := func(
		builders swappableBuilders,
		pivot Point,
		radius float32,
		prevIsLine bool,
		currIsLine bool,
		mid Point,
		after Point,
	) {
		after = after.WithScaleFrom(radius)

		if prevIsLine {
			builders.outer.setLastPoint(Point{pivot.X + mid.X, pivot.Y + mid.Y})
		} else {
			builders.outer.LineTo(pivot.X+mid.X, pivot.Y+mid.Y)
		}

		if !currIsLine {
			builders.outer.LineTo(pivot.X+after.X, pivot.Y+after.Y)
		}

		handleInnerJoin(pivot, after, builders.inner)
	}

	dotProd := beforeUnitNormal.Dot(afterUnitNormal)
	angleType := dotToAngleType(dotProd)
	before := beforeUnitNormal
	after := afterUnitNormal
	var mid Point

	if angleType == angleTypeNearlyLine {
		return
	}

	if angleType == angleTypeNearly180 {
		currIsLine = false
		mid = after.Sub(before).WithScaleFrom(radius / 2.0)
		doBluntOrClipped(
			builders,
			pivot,
			radius,
			prevIsLine,
			currIsLine,
			before,
			mid,
			after,
			invMiterLimit,
			miterClip,
		)
		return
	}

	ccw := !isClockwise(before, after)
	if ccw {
		builders.swap()
		before = before.Neg()
		after = after.Neg()
	}

	if dotProd == 0.0 && invMiterLimit <= scalar.SCALAR_ROOT_2_OVER_2 {
		mid = before.Add(after).WithScaleFrom(radius)
		doMiter(
			builders,
			pivot,
			radius,
			prevIsLine,
			currIsLine,
			mid,
			after,
		)
		return
	}

	if angleType == angleTypeSharp {
		mid = Point{after.Y - before.Y, before.X - after.X}
		if ccw {
			mid = mid.Neg()
		}
	} else {
		mid = Point{before.X + after.X, before.Y + after.Y}
	}

	sinHalfAngle := math32.Sqrt((1.0 + dotProd) * 0.5)
	if sinHalfAngle < invMiterLimit {
		currIsLine = false
		doBluntOrClipped(
			builders,
			pivot,
			radius,
			prevIsLine,
			currIsLine,
			before,
			mid,
			after,
			invMiterLimit,
			miterClip,
		)
		return
	}

	mid, _ = mid.WithLengthFrom(radius / sinHalfAngle)
	doMiter(
		builders,
		pivot,
		radius,
		prevIsLine,
		currIsLine,
		mid,
		after,
	)
}

func withNormalUnitNormal(
	before Point,
	after Point,
	scale float32,
	radius float32,
) (normal Point, unitNormal Point, ok bool) {
	unitNormal = Point{
		(after.X - before.X) * scale,
		(after.Y - before.Y) * scale,
	}
	unitNormal, ok = unitNormal.WithNormalizeFrom()
	if !ok {
		return
	}
	unitNormal = unitNormal.WithRotateCCWFrom()
	normal = unitNormal.WithScaleFrom(radius)
	return normal, unitNormal, true
}

func withNormalUnitNormal2(
	vec Point,
	radius float32,
) (normal Point, unitNormal Point, ok bool) {
	unitNormal, ok = vec.WithNormalizeFrom()
	if !ok {
		return
	}

	unitNormal = unitNormal.WithRotateCCWFrom()
	normal = unitNormal.WithScaleFrom(radius)
	return normal, unitNormal, true
}

type quadConstruct struct {
	// The state of the quad stroke under construction.
	quad             [3]Point                 // the stroked quad parallel to the original curve
	tangentStart     Point                    // a point tangent to quad[0]
	tangentEnd       Point                    // a point tangent to quad[2]
	startT           normalized.NormalizedF32 // a segment of the original curve
	midT             normalized.NormalizedF32
	endT             normalized.NormalizedF32
	startSet         bool // state to share common points across structs
	endSet           bool
	oppositeTangents bool // set if coincident tangents have opposite directions
}

// init return false if start and end are too close to have a unique middle
func (qc *quadConstruct) init(start normalized.NormalizedF32, end normalized.NormalizedF32) bool {
	qc.startT = start
	qc.midT = normalized.NewNormalizedF32WithClamped((start.Get() + end.Get()) * 0.5)
	qc.endT = end
	qc.startSet = false
	qc.endSet = false
	return qc.startT < qc.midT && qc.midT < qc.endT
}

func (qc *quadConstruct) initWithStart(parent *quadConstruct) bool {
	if !qc.init(parent.startT, parent.midT) {
		return false
	}

	qc.quad[0] = parent.quad[0]
	qc.tangentStart = parent.tangentStart
	qc.startSet = true
	return true
}

func (qc *quadConstruct) initWithEnd(parent *quadConstruct) bool {
	if !qc.init(parent.midT, parent.endT) {
		return false
	}

	qc.quad[2] = parent.quad[2]
	qc.tangentEnd = parent.tangentEnd
	qc.endSet = true
	return true
}

func checkQuadLinear(quad [3]Point) (Point, reductionType) {
	degenerateAB := degenerateVector(quad[1].Sub(quad[0]))
	degenerateBC := degenerateVector(quad[2].Sub(quad[1]))
	if degenerateAB && degenerateBC {
		return Point{}, reductionTypePoint
	}

	if degenerateAB || degenerateBC {
		return Point{}, reductionTypeLine
	}

	if !quadInLine(quad) {
		return Point{}, reductionTypeQuad
	}

	t := findQuadMaxCurvature(quad)
	if t == normalized.NormalizedF32Zero || t == normalized.NormalizedF32One {
		return Point{}, reductionTypeLine
	}

	return evalQuadAt(quad, t), reductionTypeDegenerate
}

func degenerateVector(v Point) bool {
	return !v.canNormalize()
}

// quadInLine given quad, see if all there points are in a line.
// Return true if the inside point is close to a line connecting the outermost points.
// /
// Find the outermost point by looking for the largest difference in X or Y.
// Since the XOR of the indices is 3  (0 ^ 1 ^ 2)
// the missing index equals: outer_1 ^ outer_2 ^ 3.
func quadInLine(quad [3]Point) bool {
	var ptMax float32 = -1.0
	outer1 := 0
	outer2 := 0
	for index := 0; index < 2; index++ {
		for inner := index + 1; inner < 3; inner++ {
			testDiff := quad[inner].Sub(quad[index])
			testMax := math32.Max(math32.Abs(testDiff.X), math32.Abs(testDiff.Y))
			if ptMax < testMax {
				outer1 = index
				outer2 = inner
				ptMax = testMax
			}
		}
	}

	mid := outer1 ^ outer2 ^ 3
	const CURVATURE_SLOP float32 = 0.000005 // this multiplier is pulled out of the air
	lineSlop := ptMax * ptMax * CURVATURE_SLOP
	return ptToLine(quad[mid], quad[outer1], quad[outer2]) <= lineSlop
}

// ptToLine returns the distance squared from the point to the line
func ptToLine(pt Point, lineStart Point, lineEnd Point) float32 {
	dxy := lineEnd.Sub(lineStart)
	ab0 := pt.Sub(lineStart)
	numer := dxy.Dot(ab0)
	denom := dxy.Dot(dxy)
	t := numer / denom
	if t >= 0.0 && t <= 1.0 {
		hit := Point{
			lineStart.X*(1.0-t) + lineEnd.X*t,
			lineStart.Y*(1.0-t) + lineEnd.Y*t,
		}
		return hit.DistanceToSqd(pt)
	} else {
		return pt.DistanceToSqd(lineStart)
	}
}

// intersectQuadRay intersect the line with the quad and return the t values on the quad where the line crosses.
func intersectQuadRay(line [2]Point, quad [3]Point, roots *[3]normalized.NormalizedF32Exclusive) []normalized.NormalizedF32Exclusive {
	vec := line[1].Sub(line[0])
	var r [3]float32
	for n := 0; n < 3; n++ {
		r[n] = (quad[n].Y-line[0].Y)*vec.X - (quad[n].X-line[0].X)*vec.Y
	}
	a := r[2]
	b := r[1]
	c := r[0]
	a += c - 2.0*b // A = a - 2*b + c
	b -= c         // B = -(b - c)

	length := FindUnitQuadRoots(a, 2.0*b, c, roots)
	return roots[:length]
}

func pointsWithinDist(nearPt Point, farPt Point, limit float32) bool {
	return nearPt.DistanceToSqd(farPt) <= limit*limit
}

func sharpAngle(quad [3]Point) bool {
	smaller := quad[1].Sub(quad[0])
	larger := quad[1].Sub(quad[2])
	smallerLen := smaller.LengthSqd()
	largerLen := larger.LengthSqd()

	if smallerLen > largerLen {
		smaller, larger = larger, smaller
		largerLen = smallerLen
	}

	var ok bool
	smaller, ok = smaller.WithLengthFrom(largerLen)
	if !ok {
		return false
	}

	dot := smaller.Dot(larger)
	return dot > 0.0
}

// ptInQuadBounds return true if the point is close to the bounds of the quad. This is used as a quick reject.
func ptInQuadBounds(quad [3]Point, pt Point, invResScale float32) bool {
	xMin := math32.Min(math32.Min(quad[0].X, quad[1].X), quad[2].X)
	if pt.X+invResScale < xMin {
		return false
	}

	xMax := math32.Max(math32.Max(quad[0].X, quad[1].X), quad[2].X)
	if pt.X-invResScale > xMax {
		return false
	}

	yMin := math32.Min(math32.Min(quad[0].Y, quad[1].Y), quad[2].Y)
	if pt.Y+invResScale < yMin {
		return false
	}

	yMax := math32.Max(math32.Max(quad[0].Y, quad[1].Y), quad[2].Y)
	if pt.Y-invResScale > yMax {
		return false
	}

	return true
}

func checkCubicLinear(cubic [4]Point) (reductionType, [3]Point, Point) {
	reduction := [3]Point{}
	tangentPt := Point{}

	degenerateAB := degenerateVector(cubic[1].Sub(cubic[0]))
	degenerateBC := degenerateVector(cubic[2].Sub(cubic[1]))
	degenerateCD := degenerateVector(cubic[3].Sub(cubic[2]))

	if degenerateAB && degenerateBC && degenerateCD {
		return reductionTypePoint, reduction, tangentPt
	}

	count := 0
	if degenerateAB {
		count++
	}
	if degenerateBC {
		count++
	}
	if degenerateCD {
		count++
	}

	if count == 2 {
		return reductionTypeLine, reduction, tangentPt
	}

	if !cubicInLine(cubic) {
		if degenerateAB {
			tangentPt = cubic[2]
		} else {
			tangentPt = cubic[1]
		}

		return reductionTypeQuad, reduction, tangentPt
	}

	var tValues [3]normalized.NormalizedF32
	found := FindCubicMaxCurvature(cubic, &tValues)
	rCount := 0
	// Now loop over the t-values, and reject any that evaluate to either end-point
	for i := 0; i < len(found); i++ {
		t := found[i]
		if 0.0 >= t || t >= 1.0 {
			continue
		}

		reduction[rCount] = evalCubicPosAt(cubic, tValues[i])
		if reduction[rCount] != cubic[0] && reduction[rCount] != cubic[3] {
			rCount++
		}
	}

	switch rCount {
	case 0:
		return reductionTypeLine, reduction, tangentPt
	case 1:
		return reductionTypeDegenerate, reduction, tangentPt
	case 2:
		return reductionTypeDegenerate2, reduction, tangentPt
	case 3:
		return reductionTypeDegenerate3, reduction, tangentPt
	default:
		panic("unreachable")
	}
}

// cubicInLine
// Given a cubic, determine if all four points are in a line.
// /
// Return true if the inner points is close to a line connecting the outermost points.
// /
// Find the outermost point by looking for the largest difference in X or Y.
// Given the indices of the outermost points, and that outer_1 is greater than outer_2,
// this table shows the index of the smaller of the remaining points:
// /
// ```text
//
//	                outer_2
//	            0    1    2    3
//	outer_1     ----------------
//	   0     |  -    2    1    1
//	   1     |  -    -    0    0
//	   2     |  -    -    -    0
//	   3     |  -    -    -    -
//
// ```
// /
// If outer_1 == 0 and outer_2 == 1, the smaller of the remaining indices (2 and 3) is 2.
// /
// This table can be collapsed to: (1 + (2 >> outer_2)) >> outer_1
// /
// Given three indices (outer_1 outer_2 mid_1) from 0..3, the remaining index is:
// /
// ```text
// mid_2 == (outer_1 ^ outer_2 ^ mid_1)
// ```
func cubicInLine(cubic [4]Point) bool {
	var ptMax float32 = -1.0
	outer1 := 0
	outer2 := 0
	for index := 0; index < 3; index++ {
		for inner := index + 1; inner < 4; inner++ {
			testDiff := cubic[inner].Sub(cubic[index])
			testMax := math32.Max(math32.Abs(testDiff.X), math32.Abs(testDiff.Y))
			if ptMax < testMax {
				outer1 = index
				outer2 = inner
				ptMax = testMax
			}
		}
	}

	mid1 := (1 + (2 >> outer2)) >> outer1
	mid2 := outer1 ^ outer2 ^ mid1

	lineSlop := ptMax * ptMax * 0.00001 // this multiplier is pulled out of the air

	return ptToLine(cubic[mid1], cubic[outer1], cubic[outer2]) <= lineSlop &&
		ptToLine(cubic[mid2], cubic[outer1], cubic[outer2]) <= lineSlop
}
