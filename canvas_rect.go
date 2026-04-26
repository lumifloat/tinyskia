// Copyright 2006 The Android Open Source Project
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package tinyskia

import "image/color"

// ClearRect clears all pixels on the bitmap in the given rectangle to transparent black.
func (dc *Context) ClearRect(x, y, w, h float64) {
	pd := NewPath2D()
	pd.Rect(x, y, w, h)
	style := dc.GetFillStyle()
	dc.SetFillStyleSolidColor(color.RGBA{0, 0, 0, 0})
	dc.FillPath(pd)
	dc.SetFillStyle(style)
}

// FillRect paints the given rectangle onto the bitmap, using the current fill style.
func (dc *Context) FillRect(x, y, w, h float64) {
	pd := NewPath2D()
	pd.Rect(x, y, w, h)
	dc.FillPath(pd)
}

// StrokeRect paints the box that outlines the given rectangle onto the bitmap, using the current stroke style.
func (dc *Context) StrokeRect(x, y, w, h float64) {
	pd := NewPath2D()
	pd.Rect(x, y, w, h)
	dc.StrokePath(pd)
}
