// Copyright 2006 The Android Open Source Project
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package scan

import (
	"math/bits"

	"github.com/chewxy/math32"
	"github.com/lumifloat/tinyskia/blitter"
	"github.com/lumifloat/tinyskia/internal/fixed"
	"github.com/lumifloat/tinyskia/internal/normalized"
	"github.com/lumifloat/tinyskia/internal/wide"
	"github.com/lumifloat/tinyskia/path"
)

type LineProc func(points []path.Point, clip path.ScreenIntRect, hasClip bool, blitter blitter.Blitter)

const maxCubicSubdivideLevel uint8 = 9
const maxQuadSubdivideLevel uint8 = 5

func StrokePath(
	p *path.Path,
	lineCap path.LineCap,
	clip path.ScreenIntRect,
	hasClip bool,
	blitter blitter.Blitter,
) {
	strokePathImpl(p, lineCap, clip, hasClip, hairLineRgn, blitter)
}

func hairLineRgn(points []path.Point, clip path.ScreenIntRect, hasClip bool, blitter blitter.Blitter) {
	max := float32(32767.0)
	fixedBounds, _ := path.NewRectFromLTRB(-max, -max, max, max)

	var clipBounds path.Rect
	if hasClip {
		clipBounds = clip.ToRect()
	}

	for i := 0; i < len(points)-1; i++ {
		pts := [2]path.Point{points[i], points[i+1]}

		// We have to pre-clip the line to fit in a Fixed, so we just chop the line.
		if !path.Intersect(pts, fixedBounds, &pts) {
			continue
		}

		if hasClip {
			tmp := pts
			// Perform a clip in scalar space, so we catch huge values which might
			// be missed after we convert to FDot6 (overflow).
			if !path.Intersect(tmp, clipBounds, &pts) {
				continue
			}
		}

		x0 := fixed.NewFDot6FromF32(float32(pts[0].X))
		y0 := fixed.NewFDot6FromF32(float32(pts[0].Y))
		x1 := fixed.NewFDot6FromF32(float32(pts[1].X))
		y1 := fixed.NewFDot6FromF32(float32(pts[1].Y))

		dx := x1 - x0
		dy := y1 - y0

		if fixed.FDot6Abs(dx) > fixed.FDot6Abs(dy) {
			// mostly horizontal
			if x0 > x1 {
				// we want to go left-to-right
				x0, x1 = x1, x0
				y0, y1 = y1, y0
			}

			ix0 := fixed.FDot6Round(x0)
			ix1 := fixed.FDot6Round(x1)
			if ix0 == ix1 {
				// too short to blitter
				continue
			}

			slope := fixed.FDot6DivToFDot16(dy, dx)
			startY := fixed.FDot6ToFDot16(y0) + (slope * ((32 - x0) & 63) >> 6)

			var maxY fixed.FDot16 = math32.MaxInt32
			if hasClip {
				maxY = fixed.NewFDot16FromF32(float32(clipBounds.Bottom()))
			}

			for {
				if ix0 >= 0 && startY >= 0 && startY < maxY {
					blitter.BlitH(uint32(ix0), uint32(startY>>16), 1)
				}

				startY += slope
				ix0 += 1
				if ix0 >= ix1 {
					break
				}
			}
		} else {
			// mostly vertical
			if y0 > y1 {
				// we want to go top-to-bottom
				x0, x1 = x1, x0
				y0, y1 = y1, y0
			}

			iy0 := fixed.FDot6Round(y0)
			iy1 := fixed.FDot6Round(y1)
			if iy0 == iy1 {
				// too short to blitter
				continue
			}

			slope := fixed.FDot6DivToFDot16(dx, dy)
			startX := fixed.FDot6ToFDot16(x0) + (slope * ((32 - y0) & 63) >> 6)

			for {
				if startX >= 0 && iy0 >= 0 {
					blitter.BlitH(uint32(startX>>16), uint32(iy0), 1)
				}

				startX += slope
				iy0 += 1
				if iy0 >= iy1 {
					break
				}
			}
		}
	}
}

func strokePathImpl(
	p *path.Path,
	lineCap path.LineCap,
	clip path.ScreenIntRect,
	hasClip bool,
	lineProc LineProc,
	blitter blitter.Blitter,
) {
	var insetClip path.IntRect
	var outsetClip path.IntRect
	hasInsetClip := false
	hasOutsetClip := false

	capOut := float32(1.0)
	if lineCap != path.LineCapButt {
		capOut = float32(2.0)
	}

	bounds := p.Bounds()
	outsetBounds, ok := bounds.Outset(capOut, capOut)
	if !ok {
		return
	}
	ibounds, ok := outsetBounds.RoundOut()
	if !ok {
		return
	}

	if hasClip {
		_, ok := clip.ToIntRect().Intersect(ibounds)
		if !ok {
			return
		}

		if !clip.ToIntRect().Contains(ibounds) {
			outset, ok := clip.ToIntRect().Outset(1, 1)
			if !ok {
				return
			}
			outsetClip = outset
			hasOutsetClip = true
			inset, ok := clip.ToIntRect().Inset(1, 1)
			if !ok {
				return
			}
			insetClip = inset
			hasInsetClip = true
		}
	}

	var prevVerb path.PathVerb = path.PathVerbMove
	var firstPt path.Point
	var lastPt path.Point

	iter := p.Segments()
	for {
		segment := iter.Next()
		if segment == nil {
			break
		}

		// Get next verb for extendPts (similar to Rust's iter.next_verb())
		nextVerb := iter.NextVerb()
		// hasNextVerb is true if we're not at the end of the path
		hasNextVerb := iter.HasNext()

		var lastPt2 path.Point

		switch s := segment.(type) {
		case path.PathSegmentMoveTo:
			firstPt = path.Point(s)
			lastPt = path.Point(s)
			lastPt2 = path.Point(s)
		case path.PathSegmentLineTo:
			points := [2]path.Point{lastPt, path.Point(s)}
			if lineCap != path.LineCapButt {
				extendPts(lineCap, prevVerb, nextVerb, hasNextVerb, points[:])
			}
			lineProc(points[:], clip, hasClip, blitter)
			lastPt = path.Point(s)
			lastPt2 = points[0]
		case path.PathSegmentQuadTo:
			points := [3]path.Point{lastPt, s.P0, s.P1}
			if lineCap != path.LineCapButt {
				extendPts(lineCap, prevVerb, nextVerb, hasNextVerb, points[:])
			}
			hairQuad(
				points,
				clip,
				hasClip,
				insetClip,
				hasInsetClip,
				outsetClip,
				hasOutsetClip,
				computeQuadLevel(points),
				lineProc,
				blitter,
			)
			lastPt = s.P1
			lastPt2 = points[0]
		case path.PathSegmentCubicTo:
			points := [4]path.Point{lastPt, s.P0, s.P1, s.P2}
			if lineCap != path.LineCapButt {
				extendPts(lineCap, prevVerb, nextVerb, hasNextVerb, points[:])
			}
			hairCubic(
				points,
				clip,
				hasClip,
				insetClip,
				hasInsetClip,
				outsetClip,
				hasOutsetClip,
				lineProc,
				blitter,
			)
			lastPt = s.P2
			lastPt2 = points[0]
		case path.PathSegmentClose:
			points := [2]path.Point{lastPt, firstPt}
			if lineCap != path.LineCapButt && prevVerb == path.PathVerbMove {
				extendPts(lineCap, prevVerb, nextVerb, hasNextVerb, points[:])
			}
			lineProc(points[:], clip, hasClip, blitter)
			lastPt2 = points[0]
		}

		if lineCap != path.LineCapButt {
			var currVerb path.PathVerb
			switch segment.(type) {
			case path.PathSegmentLineTo:
				currVerb = path.PathVerbLine
			case path.PathSegmentQuadTo:
				currVerb = path.PathVerbQuad
			case path.PathSegmentCubicTo:
				currVerb = path.PathVerbCubic
			}
			if prevVerb == path.PathVerbMove &&
				(currVerb == path.PathVerbLine || currVerb == path.PathVerbQuad || currVerb == path.PathVerbCubic) {
				firstPt = lastPt2
			}
			prevVerb = currVerb
		}
	}
}

func extendPts(
	lineCap path.LineCap,
	prevVerb path.PathVerb,
	nextVerb path.PathVerb,
	hasNextVerb bool,
	points []path.Point,
) {
	capOutset := float32(math32.Pi / 8.0)
	if lineCap == path.LineCapSquare {
		capOutset = 0.5
	}

	if prevVerb == path.PathVerbMove {
		first := points[0]
		offset := 0
		controls := len(points) - 1
		var tangent path.Point
		for {
			offset++
			tangent = first.Sub(points[offset])
			if !tangent.IsZero() {
				break
			}
			controls--
			if controls == 0 {
				break
			}
		}

		if tangent.IsZero() {
			tangent = path.Point{X: 1.0, Y: 0.0}
			controls = len(points) - 1
		} else {
			tangent, _ = tangent.WithNormalizeFrom()
		}

		offset = 0
		for {
			points[offset].X += tangent.X * capOutset
			points[offset].Y += tangent.Y * capOutset
			offset++
			controls++
			if controls >= len(points) {
				break
			}
		}
	}

	isEndOfContour := !hasNextVerb || nextVerb == path.PathVerbMove || nextVerb == path.PathVerbClose
	if isEndOfContour {
		last := points[len(points)-1]
		offset := len(points) - 1
		controls := len(points) - 1
		var tangent path.Point
		for {
			offset--
			tangent = last.Sub(points[offset])
			if !tangent.IsZero() {
				break
			}
			controls--
			if controls == 0 {
				break
			}
		}

		if tangent.IsZero() {
			tangent = path.Point{X: -1.0, Y: 0.0}
			controls = len(points) - 1
		} else {
			tangent, _ = tangent.WithNormalizeFrom()
		}

		offset = len(points) - 1
		for {
			points[offset].X += tangent.X * capOutset
			points[offset].Y += tangent.Y * capOutset
			offset--
			controls++
			if controls >= len(points) {
				break
			}
		}
	}
}

func hairQuad(
	points [3]path.Point,
	clip path.ScreenIntRect,
	hasClip bool,
	insetClip path.IntRect,
	hasInsetClip bool,
	outsetClip path.IntRect,
	hasOutsetClip bool,
	level uint8,
	lineProc LineProc,
	blitter blitter.Blitter,
) {
	if hasInsetClip {
		outsetRect := outsetClip.ToRect()
		insetRect := insetClip.ToRect()
		bounds, ok := computeNocheckQuadBounds(points)
		if !ok {
			return
		}
		if !geometricOverlap(outsetRect, bounds) {
			return
		} else if geometricContains(insetRect, bounds) {
			hasClip = false
		}
	}
	hairQuad2(points, clip, hasClip, level, lineProc, blitter)
}

func computeNocheckQuadBounds(points [3]path.Point) (path.Rect, bool) {
	minX, minY := points[0].X, points[0].Y
	maxX, maxY := minX, minY
	for i := 1; i < 3; i++ {
		if points[i].X < minX {
			minX = points[i].X
		}
		if points[i].Y < minY {
			minY = points[i].Y
		}
		if points[i].X > maxX {
			maxX = points[i].X
		}
		if points[i].Y > maxY {
			maxY = points[i].Y
		}
	}
	return path.NewRectFromLTRB(minX, minY, maxX, maxY)
}

func geometricOverlap(a, b path.Rect) bool {
	return a.Left() < b.Right() && b.Left() < a.Right() && a.Top() < b.Bottom() && b.Top() < a.Bottom()
}

func geometricContains(outer, inner path.Rect) bool {
	return inner.Right() <= outer.Right() && inner.Left() >= outer.Left() &&
		inner.Bottom() <= outer.Bottom() && inner.Top() >= outer.Top()
}

func hairQuad2(
	points [3]path.Point,
	clip path.ScreenIntRect,
	hasClip bool,
	level uint8,
	lineProc LineProc,
	blitter blitter.Blitter,
) {
	coeff := path.NewQuadCoeffFromPoints([3]path.Point{points[0], points[1], points[2]})
	lines := 1 << level
	tmp := make([]path.Point, lines+1)
	tmp[0] = points[0]

	dt := 1.0 / float64(lines)
	t := 0.0
	for i := 1; i < lines; i++ {
		t += dt
		v := coeff.Eval(wide.Splat(float32(t)))
		tmp[i] = path.Point{X: float32(v[0]), Y: float32(v[1])}
	}
	tmp[lines] = points[2]
	lineProc(tmp, clip, hasClip, blitter)
}

func computeQuadLevel(points [3]path.Point) uint8 {
	d := computeIntQuadDist(points)
	if d == 0 {
		return 0
	}
	// Use 33 instead of 64 to match Rust's u32-based calculation
	level := (33 - bits.LeadingZeros32(uint32(d))) >> 1
	if level > int(maxQuadSubdivideLevel) {
		level = int(maxQuadSubdivideLevel)
	}
	return uint8(level)
}

func computeIntQuadDist(points [3]path.Point) uint32 {
	dx := math32.Abs((points[0].X+points[2].X)*0.5 - points[1].X)
	dy := math32.Abs((points[0].Y+points[2].Y)*0.5 - points[1].Y)

	idx := uint32(math32.Ceil(dx))
	idy := uint32(math32.Ceil(dy))

	if idx > idy {
		return idx + (idy >> 1)
	}
	return idy + (idx >> 1)
}

func hairCubic(
	points [4]path.Point,
	clip path.ScreenIntRect,
	hasClip bool,
	insetClip path.IntRect,
	hasInsetClip bool,
	outsetClip path.IntRect,
	hasOutsetClip bool,
	lineProc LineProc,
	blitter blitter.Blitter,
) {
	if hasInsetClip {
		outsetRect := outsetClip.ToRect()
		insetRect := insetClip.ToRect()
		bounds, ok := computeNocheckCubicBounds(points)
		if !ok {
			return
		}
		if !geometricOverlap(outsetRect, bounds) {
			return
		} else if geometricContains(insetRect, bounds) {
			hasClip = false
		}
	}

	if quickCubicNicenessCheck(points) {
		hairCubic2(points, clip, hasClip, lineProc, blitter)
	} else {
		tmp := make([]path.Point, 13)
		tValues := [3]normalized.NormalizedF32{}
		count := path.ChopCubicAtMaxCurvature([4]path.Point{points[0], points[1], points[2], points[3]}, &tValues, tmp)
		for i := 0; i < int(count); i++ {
			offset := i * 3
			lineProc(tmp[offset:offset+4], clip, hasClip, blitter)
		}
	}
}

func computeNocheckCubicBounds(points [4]path.Point) (path.Rect, bool) {
	minX, minY := points[0].X, points[0].Y
	maxX, maxY := minX, minY
	for i := 1; i < 4; i++ {
		minX = math32.Min(minX, points[i].X)
		minY = math32.Min(minY, points[i].Y)
		maxX = math32.Max(maxX, points[i].X)
		maxY = math32.Max(maxY, points[i].Y)
	}
	return path.NewRectFromLTRB(minX, minY, maxX, maxY)
}

func quickCubicNicenessCheck(points [4]path.Point) bool {
	return lt90(points[1], points[0], points[3]) &&
		lt90(points[2], points[0], points[3]) &&
		lt90(points[1], points[3], points[0]) &&
		lt90(points[2], points[3], points[0])
}

func lt90(p0, pivot, p2 path.Point) bool {
	return p0.Sub(pivot).Dot(p2.Sub(pivot)) >= 0.0
}

func hairCubic2(
	points [4]path.Point,
	clip path.ScreenIntRect,
	hasClip bool,
	lineProc LineProc,
	blitter blitter.Blitter,
) {
	lines := computeCubicSegments(points)
	if lines == 1 {
		lineProc([]path.Point{points[0], points[3]}, clip, hasClip, blitter)
		return
	}

	coeff := path.NewCubicCoeffFromPoints([4]path.Point{points[0], points[1], points[2], points[3]})
	tmp := make([]path.Point, lines+1)
	tmp[0] = points[0]

	dt := 1.0 / float64(lines)
	t := 0.0
	for i := 1; i < int(lines); i++ {
		t += dt
		ev := coeff.Eval(wide.Splat(float32(t)))
		tmp[i] = path.Point{X: float32(ev[0]), Y: float32(ev[1])}
	}

	tmp[lines] = points[3]
	lineProc(tmp, clip, hasClip, blitter)
}

func computeCubicSegments(points [4]path.Point) uint64 {
	p0, p1, p2, p3 := points[0], points[1], points[2], points[3]

	oneThird := path.Point{X: 1.0 / 3.0, Y: 1.0 / 3.0}
	twoThird := path.Point{X: 2.0 / 3.0, Y: 2.0 / 3.0}

	p13 := p3.Mul(oneThird).Add(p0.Mul(twoThird))
	p23 := p0.Mul(oneThird).Add(p3.Mul(twoThird))

	dx := math32.Max(math32.Abs(p1.X-p13.X), math32.Abs(p2.X-p23.X))
	dy := math32.Max(math32.Abs(p1.Y-p13.Y), math32.Abs(p2.Y-p23.Y))
	diff := math32.Max(dx, dy)

	tol := float32(1.0 / 8.0)
	for i := uint8(0); i < maxCubicSubdivideLevel; i++ {
		if diff < tol {
			return 1 << i
		}
		tol *= 4.0
	}

	return 1 << maxCubicSubdivideLevel
}
