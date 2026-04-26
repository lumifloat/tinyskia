// Copyright 2006 The Android Open Source Project
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package tinyskia

// TextMetrics holds text measurement information, similar to HTML5 Canvas TextMetrics.
type TextMetrics struct {
	// Width is the advance width of the text (inline box width).
	Width float64

	// ActualBoundingBoxLeft is the distance from the alignment point to the left side
	// of the bounding rectangle. Positive values indicate left direction.
	ActualBoundingBoxLeft float64

	// ActualBoundingBoxRight is the distance from the alignment point to the right side
	// of the bounding rectangle. Positive values indicate right direction.
	ActualBoundingBoxRight float64

	// FontBoundingBoxAscent is the distance from the baseline to the font's ascent metric,
	// in CSS pixels. Positive values indicate upward direction.
	FontBoundingBoxAscent float64

	// FontBoundingBoxDescent is the distance from the baseline to the font's descent metric,
	// in CSS pixels. Positive values indicate downward direction.
	FontBoundingBoxDescent float64

	// ActualBoundingBoxAscent is the distance from the baseline to the top of the
	// bounding rectangle of the given text. Positive values indicate upward.
	ActualBoundingBoxAscent float64

	// ActualBoundingBoxDescent is the distance from the baseline to the bottom of the
	// bounding rectangle of the given text. Positive values indicate downward.
	ActualBoundingBoxDescent float64
}
