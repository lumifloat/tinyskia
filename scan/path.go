// Copyright 2006 The Android Open Source Project
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package scan

import (
	"math"
	"sort"

	"github.com/chewxy/math32"
	"github.com/lumifloat/tinyskia/blitter"
	"github.com/lumifloat/tinyskia/edge"
	"github.com/lumifloat/tinyskia/internal/fixed"
	"github.com/lumifloat/tinyskia/path"
)

func FillPath(
	p *path.Path,
	fillRule int,
	clip path.ScreenIntRect,
	blitter blitter.Blitter,
) {
	// Check for nil or empty path
	if p == nil || p.IsEmpty() {
		return
	}

	ir, ok := conservativeRoundToInt(p.Bounds())
	if !ok {
		return
	}

	var pathContainedInClip bool
	if bounds, ok := ir.ToScreenIntRect(); ok {
		pathContainedInClip = clip.Contains(bounds)
	}

	FillPathImpl(
		p,
		fillRule,
		clip,
		ir.Top(),
		ir.Bottom(),
		0,
		pathContainedInClip,
		blitter,
	)
}

func conservativeRoundToInt(src path.Rect) (path.IntRect, bool) {
	return path.NewIntRectFromLTRB(
		roundDownToInt(src.Left()),
		roundDownToInt(src.Top()),
		roundUpToInt(src.Right()),
		roundUpToInt(src.Bottom()),
	)
}

const conservativeRoundBias float64 = 0.5 + 1.5/fixed.FDot6One

func roundDownToInt(x float32) int32 {
	var xx float64 = float64(x)
	xx = xx - conservativeRoundBias
	return int32(math.Ceil(xx))
}

func roundUpToInt(x float32) int32 {
	var xx float64 = float64(x)
	xx = xx + conservativeRoundBias
	return int32(math.Floor(xx))
}

func FillPathImpl(
	path *path.Path,
	fillRule int,
	clipRect path.ScreenIntRect,
	startY int32,
	stopY int32,
	shiftEdgesUp int32,
	pathContainedInClip bool,
	blitter blitter.Blitter,
) {
	var hasClip bool
	var clip edge.ShiftedIntRect
	shiftedClip, ok := edge.NewShiftedIntRect(clipRect, shiftEdgesUp)
	if !ok {
		return
	}
	if !pathContainedInClip {
		hasClip = true
		clip = shiftedClip
	}

	edges := edge.BuildEdges(path, hasClip, clip, shiftEdgesUp)
	if len(edges) == 0 {
		return
	}

	sort.Slice(edges, func(i, j int) bool {
		valueA := edges[i].AsLine().FirstY
		valueB := edges[j].AsLine().FirstY
		if valueA == valueB {
			return edges[i].AsLine().X < edges[j].AsLine().X
		}
		return valueA < valueB
	})

	for i := range edges {
		edges[i].AsLine().Prev = i
		edges[i].AsLine().Next = i + 2
	}

	const EDGE_HEAD_Y = math32.MinInt32
	const EDGE_TAIL_Y = math32.MaxInt32

	// Insert Head
	head := edge.Edge{
		Type: edge.EdgeTypeLine,
		Line: &edge.LineEdge{
			Prev:   -1,
			Next:   1,
			X:      math32.MinInt32,
			DX:     0,
			FirstY: EDGE_HEAD_Y,
		},
	}
	edges = append([]edge.Edge{head}, edges...)

	// Append Tail
	tail := edge.Edge{
		Type: edge.EdgeTypeLine,
		Line: &edge.LineEdge{
			Prev:   len(edges) - 1,
			Next:   -1,
			X:      0,
			DX:     0,
			FirstY: EDGE_TAIL_Y,
		},
	}
	edges = append(edges, tail)

	startY <<= shiftEdgesUp
	stopY <<= shiftEdgesUp

	top := int32(shiftedClip.Shifted().Y())
	if !pathContainedInClip && startY < top {
		startY = top
	}

	bottom := int32(shiftedClip.Shifted().Bottom())
	if !pathContainedInClip && stopY > bottom {
		stopY = bottom
	}

	// Check for negative values (equivalent to Rust's u32::try_from error handling)
	if startY < 0 || stopY < 0 {
		return
	}

	walkEdges(
		fillRule,
		uint32(startY),
		uint32(stopY),
		uint32(shiftedClip.Shifted().Right()),
		edges,
		blitter,
	)
}

func walkEdges(
	fillRule int,
	startY uint32,
	stopY uint32,
	rightClip uint32,
	edges []edge.Edge,
	blitter blitter.Blitter,
) {
	currY := startY
	var windingMask int32 = -1
	if fillRule == fillRuleEvenOdd {
		windingMask = 1
	}

	for {
		var w int32 = 0
		var left uint32 = 0
		prevX := edges[0].AsLine().X

		currIdx := int(edges[0].AsLine().Next)
		for edges[currIdx].AsLine().FirstY <= int32(currY) {
			x := uint32(fixed.FDot16RoundToI32(edges[currIdx].AsLine().X))

			if (w & windingMask) == 0 {
				left = x
			}

			w += int32(edges[currIdx].AsLine().Winding)

			if (w & windingMask) == 0 {
				if left < x {
					blitter.BlitH(left, currY, x-left)
				}
			}

			nextIdx := int(edges[currIdx].AsLine().Next)
			var newX fixed.FDot16

			if edges[currIdx].AsLine().LastY == int32(currY) {
				if edges[currIdx].Type == edge.EdgeTypeLine {
					removeEdge(currIdx, edges)
				} else if edges[currIdx].Update() {
					newX = edges[currIdx].AsLine().X
					if newX < prevX {
						backwardInsertEdgeBasedOnX(currIdx, edges)
					} else {
						prevX = newX
					}
				} else {
					removeEdge(currIdx, edges)
				}
			} else {
				newX = edges[currIdx].AsLine().X + edges[currIdx].AsLine().DX
				edges[currIdx].AsLine().X = newX
				if newX < prevX {
					backwardInsertEdgeBasedOnX(currIdx, edges)
				} else {
					prevX = newX
				}
			}
			currIdx = nextIdx
		}

		if (w & windingMask) != 0 {
			if left < rightClip {
				blitter.BlitH(left, currY, rightClip-left)
			}
		}

		currY++
		if currY >= stopY {
			break
		}

		insertNewEdges(currIdx, int32(currY), edges)
	}
}

func removeEdge(currIdx int, edges []edge.Edge) {
	prev := edges[currIdx].AsLine().Prev
	next := edges[currIdx].AsLine().Next
	edges[prev].AsLine().Next = next
	edges[next].AsLine().Prev = prev
}

func backwardInsertEdgeBasedOnX(currIdx int, edges []edge.Edge) {
	x := edges[currIdx].AsLine().X
	prevIdx := int(edges[currIdx].AsLine().Prev)
	for prevIdx != 0 {
		if edges[prevIdx].AsLine().X > x {
			prevIdx = int(edges[prevIdx].AsLine().Prev)
		} else {
			break
		}
	}

	nextIdx := int(edges[prevIdx].AsLine().Next)
	if nextIdx != currIdx {
		removeEdge(currIdx, edges)
		insertEdgeAfter(currIdx, prevIdx, edges)
	}
}

func insertEdgeAfter(currIdx int, afterIdx int, edges []edge.Edge) {
	edges[currIdx].AsLine().Prev = afterIdx
	edges[currIdx].AsLine().Next = edges[afterIdx].AsLine().Next

	afterNextIdx := edges[afterIdx].AsLine().Next
	edges[afterNextIdx].AsLine().Prev = currIdx
	edges[afterIdx].AsLine().Next = currIdx
}

func backwardInsertStart(prevIdx int, x fixed.FDot16, edges []edge.Edge) int {
	for {
		prev := edges[prevIdx].AsLine().Prev
		if prev == 0 && edges[0].AsLine().X > x {
			break
		}
		prevIdx = int(prev)
		if edges[prevIdx].AsLine().X <= x {
			break
		}
		if prev == 0 {
			break
		}
	}
	return prevIdx
}

func insertNewEdges(newIdx int, currY int32, edges []edge.Edge) {
	if edges[newIdx].AsLine().FirstY != currY {
		return
	}

	prevIdx := int(edges[newIdx].AsLine().Prev)
	if edges[prevIdx].AsLine().X <= edges[newIdx].AsLine().X {
		return
	}

	startIdx := backwardInsertStart(prevIdx, edges[newIdx].AsLine().X, edges)
	for {
		nextIdx := int(edges[newIdx].AsLine().Next)
		keepEdge := false
		for {
			afterIdx := int(edges[startIdx].AsLine().Next)
			if afterIdx == newIdx {
				keepEdge = true
				break
			}
			if edges[afterIdx].AsLine().X >= edges[newIdx].AsLine().X {
				break
			}
			startIdx = afterIdx
		}

		if !keepEdge {
			removeEdge(newIdx, edges)
			insertEdgeAfter(newIdx, startIdx, edges)
		}

		startIdx = newIdx
		newIdx = nextIdx
		if edges[newIdx].AsLine().FirstY != currY {
			break
		}
	}
}
