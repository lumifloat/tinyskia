// Copyright 2009 The Android Open Source Project
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package path

const MAX_POINTS = 4

// Clip clips the line pts[0]...pts[1] against clip, ignoring segments that
// lie completely above or below the clip. For portions to the left or
// right, turn those into vertical line segments that are aligned to the
// edge of the clip.
//
// The clipped points are stored in points buffer. Returns the number of points written.
func Clip(src [2]Point, clip Rect, canCullToTheRight bool, points *[MAX_POINTS]Point) int {
	index0, index1 := 0, 1
	if src[0].Y < src[1].Y {
		index0, index1 = 0, 1
	} else {
		index0, index1 = 1, 0
	}

	// Check if we're completely clipped out in Y (above or below)
	if src[index1].Y <= clip.Top() {
		// we're above the clip
		return 0
	}

	if src[index0].Y >= clip.Bottom() {
		// we're below the clip
		return 0
	}

	// Chop in Y to produce a single segment, stored in tmp[0..1]
	tmp := src

	// now compute intersections
	if src[index0].Y < clip.Top() {
		tmp[index0] = Point{
			X: sectWithHorizontal(tmp, clip.Top()),
			Y: clip.Top(),
		}
	}

	if tmp[index1].Y > clip.Bottom() {
		tmp[index1] = Point{
			X: sectWithHorizontal(tmp, clip.Bottom()),
			Y: clip.Bottom(),
		}
	}

	// Chop it into 1..3 segments that are wholly within the clip in X.
	// temp storage for up to 3 segments
	var resultStorage [MAX_POINTS]Point
	lineCount := 1
	var reverse bool

	if src[0].X < src[1].X {
		index0, index1 = 0, 1
		reverse = false
	} else if src[0].X > src[1].X {
		index0, index1 = 1, 0
		reverse = true
	} else {
		// Vertical line: X coordinates are equal, use Y ordering
		if src[0].Y <= src[1].Y {
			index0, index1 = 0, 1
		} else {
			index0, index1 = 1, 0
		}
		reverse = false // Don't reverse for vertical lines
	}

	var result []Point
	if tmp[index1].X <= clip.Left() {
		// wholly to the left
		tmp[0].X = clip.Left()
		tmp[1].X = clip.Left()
		reverse = false
		result = tmp[:]
	} else if tmp[index0].X >= clip.Right() {
		// wholly to the right
		if canCullToTheRight {
			return 0
		}

		tmp[0].X = clip.Right()
		tmp[1].X = clip.Right()
		reverse = false
		result = tmp[:]
	} else {
		offset := 0

		if tmp[index0].X < clip.Left() {
			resultStorage[offset] = Point{X: clip.Left(), Y: tmp[index0].Y}
			offset++
			resultStorage[offset] = Point{
				X: clip.Left(),
				Y: sectClampWithVertical(tmp, clip.Left()),
			}
		} else {
			resultStorage[offset] = tmp[index0]
		}
		offset++

		if tmp[index1].X > clip.Right() {
			resultStorage[offset] = Point{
				X: clip.Right(),
				Y: sectClampWithVertical(tmp, clip.Right()),
			}
			offset++
			resultStorage[offset] = Point{X: clip.Right(), Y: tmp[index1].Y}
		} else {
			resultStorage[offset] = tmp[index1]
		}

		lineCount = offset
		result = resultStorage[:lineCount+1]
	}

	// Now copy the results into the caller's points[] parameter
	if reverse {
		// copy the pts in reverse order to maintain winding order
		for i := 0; i <= lineCount; i++ {
			points[lineCount-i] = result[i]
		}
	} else {
		copy(points[:lineCount+1], result[:lineCount+1])
	}

	return lineCount + 1
}

func sectWithHorizontal(pts [2]Point, y float32) float32 {
	dy := pts[1].Y - pts[0].Y
	if dy == 0 {
		// Use average when line is horizontal (matches Rust's .ave() method)
		return (pts[0].X + pts[1].X) / 2
	}
	// Use f64 for better precision during calculation
	x0 := float64(pts[0].X)
	y0 := float64(pts[0].Y)
	x1 := float64(pts[1].X)
	y1 := float64(pts[1].Y)
	result := x0 + (float64(y)-y0)*(x1-x0)/(y1-y0)
	// Pin the result to handle floating point precision issues
	return float32(pinUnsortedF64(result, float64(pts[0].X), float64(pts[1].X)))
}

func sectClampWithVertical(pts [2]Point, x float32) float32 {
	y := sectWithVertical(pts, x)
	// Our caller expects y to be between pts[0].Y and pts[1].Y (unsorted), but due to the
	// numerics of floats/doubles, we might have computed a value slightly outside of that,
	// so we have to manually clamp afterwards.
	// See skbug.com/7491
	return pinUnsortedF32(y, pts[0].Y, pts[1].Y)
}

// sectWithVertical returns Y coordinate of intersection with vertical line at X.
func sectWithVertical(pts [2]Point, x float32) float32 {
	dx := pts[1].X - pts[0].X
	if dx == 0 {
		// Vertical line: return average of Y values (matches Rust's .ave() method)
		return (pts[0].Y + pts[1].Y) / 2
	}
	// Use f64 for better precision during calculation
	x0 := float64(pts[0].X)
	y0 := float64(pts[0].Y)
	x1 := float64(pts[1].X)
	y1 := float64(pts[1].Y)
	result := y0 + (float64(x)-x0)*(y1-y0)/(x1-x0)
	// Pin the result to handle floating point precision issues
	return float32(pinUnsortedF64(result, float64(pts[0].Y), float64(pts[1].Y)))
}

// pinUnsortedF32 clamps value to be within [limit0, limit1],
// handling the case where limits may not be sorted.
func pinUnsortedF32(value, limit0, limit1 float32) float32 {
	l0, l1 := limit0, limit1
	if l1 < l0 {
		l0, l1 = l1, l0
	}
	// Now limits are sorted
	if value < l0 {
		return l0
	} else if value > l1 {
		return l1
	}
	return value
}

// pinUnsortedF64 clamps value to be within [limit0, limit1],
// handling the case where limits may not be sorted (f64 version).
func pinUnsortedF64(value, limit0, limit1 float64) float64 {
	l0, l1 := limit0, limit1
	if l1 < l0 {
		l0, l1 = l1, l0
	}
	// Now limits are sorted
	if value < l0 {
		return l0
	} else if value > l1 {
		return l1
	}
	return value
}

// Intersect intersects the line segment against the rect. If there is a non-empty
// resulting segment, return true and set dst[] to that segment. If not,
// return false and ignore dst[].
//
// Clip is specialized for scan-conversion, as it adds vertical
// segments on the sides to show where the line extended beyond the
// left or right sides. Intersect does not.
func Intersect(src [2]Point, clip Rect, dst *[2]Point) bool {
	// Compute bounds of the source line
	left := src[0].X
	if src[1].X < left {
		left = src[1].X
	}
	right := src[0].X
	if src[1].X > right {
		right = src[1].X
	}
	top := src[0].Y
	if src[1].Y < top {
		top = src[1].Y
	}
	bottom := src[0].Y
	if src[1].Y > bottom {
		bottom = src[1].Y
	}

	bounds, ok := NewRectFromLTRB(left, top, right, bottom)
	if ok {
		if containsNoEmptyCheck(clip, bounds) {
			dst[0] = src[0]
			dst[1] = src[1]
			return true
		}

		// Check for no overlap, and only permit coincident edges if the line
		// and the edge are collinear
		if nestedLt(bounds.Right(), clip.Left(), bounds.Width()) ||
			nestedLt(clip.Right(), bounds.Left(), bounds.Width()) ||
			nestedLt(bounds.Bottom(), clip.Top(), bounds.Height()) ||
			nestedLt(clip.Bottom(), bounds.Top(), bounds.Height()) {
			return false
		}
	}

	index0, index1 := 0, 1
	if src[0].Y >= src[1].Y {
		index0, index1 = 1, 0
	}

	tmp := src

	// Now compute Y intersections
	if tmp[index0].Y < clip.Top() {
		tmp[index0] = Point{
			X: sectWithHorizontal(tmp, clip.Top()),
			Y: clip.Top(),
		}
	}

	if tmp[index1].Y > clip.Bottom() {
		tmp[index1] = Point{
			X: sectWithHorizontal(tmp, clip.Bottom()),
			Y: clip.Bottom(),
		}
	}

	if tmp[0].X < tmp[1].X {
		index0, index1 = 0, 1
	} else {
		index0, index1 = 1, 0
	}

	// Check for quick-reject in X again, now that we may have been chopped
	if tmp[index1].X <= clip.Left() || tmp[index0].X >= clip.Right() {
		// Usually we will return false, but we don't if the line is vertical and coincident
		// with the clip.
		if tmp[0].X != tmp[1].X || tmp[0].X < clip.Left() || tmp[0].X > clip.Right() {
			return false
		}
	}

	if tmp[index0].X < clip.Left() {
		tmp[index0] = Point{
			X: clip.Left(),
			Y: sectWithVertical(tmp, clip.Left()),
		}
	}

	if tmp[index1].X > clip.Right() {
		tmp[index1] = Point{
			X: clip.Right(),
			Y: sectWithVertical(tmp, clip.Right()),
		}
	}

	dst[0] = tmp[0]
	dst[1] = tmp[1]
	return true
}

func nestedLt(a, b, dim float32) bool {
	return a <= b && (a < b || dim > 0.0)
}

// containsNoEmptyCheck returns true if outer contains inner, even if inner is empty.
func containsNoEmptyCheck(outer, inner Rect) bool {
	return outer.Left() <= inner.Left() &&
		outer.Top() <= inner.Top() &&
		outer.Right() >= inner.Right() &&
		outer.Bottom() >= inner.Bottom()
}
