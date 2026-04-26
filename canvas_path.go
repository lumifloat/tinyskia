// Copyright 2006 The Android Open Source Project
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package tinyskia

// MoveTo creates a new subpath with the given point.
func (dc *Context) MoveTo(x, y float64) {
	dc.path2d.MoveTo(x, y)
}

// ClosePath Marks the current subpath as closed,
// and starts a new subpath with a point the same as the start and end of the newly closed subpath.
func (dc *Context) ClosePath() {
	dc.path2d.ClosePath()
}

// LineTo adds the given point to the current subpath, connected to the previous one by a straight line.
func (dc *Context) LineTo(x, y float64) {
	dc.path2d.LineTo(x, y)
}

// QuadraticCurveTo adds the given point to the current subpath, connected to the previous one
// by a quadratic Bézier curve with the given control point.
func (dc *Context) QuadraticCurveTo(x1, y1, x2, y2 float64) {
	dc.path2d.QuadraticCurveTo(x1, y1, x2, y2)
}

// BezierCurveTo adds the given point to the current subpath, connected to the previous one
// by a cubic Bézier curve with the given control points.
func (dc *Context) BezierCurveTo(cp1x, cp1y, cp2x, cp2y, x, y float64) {
	dc.path2d.BezierCurveTo(cp1x, cp1y, cp2x, cp2y, x, y)
}

// ArcTo adds an arc with the given control points and radius to the current subpath,
// connected to the previous point by a straight line.
func (dc *Context) ArcTo(x1, y1, x2, y2, radius float64) {
	dc.path2d.ArcTo(x1, y1, x2, y2, radius)
}

// Rect adds a new closed subpath to the path, representing the given rectangle.
func (dc *Context) Rect(x, y, w, h float64) {
	dc.path2d.Rect(x, y, w, h)
}

// RoundRect adds a new closed subpath to the path representing the given rounded rectangle.
func (dc *Context) RoundRect(x, y, w, h float64, radii []float64) {
	dc.path2d.RoundRect(x, y, w, h, radii)
}

// Arc adds points to the subpath such that the arc described by the circumference of the circle
// described by the arguments, starting at the given start angle and ending at the given end angle,
// going in the given direction (defaulting to clockwise), is added to the path, connected to
// the previous point by a straight line.
func (dc *Context) Arc(x, y, radius, startAngle, endAngle float64, counterclockwise bool) {
	dc.path2d.Arc(x, y, radius, startAngle, endAngle, counterclockwise)
}

// Ellipse adds points to the subpath such that the arc described by the circumference of the ellipse
// described by the arguments, starting at the given start angle and ending at the given end angle,
// going in the given direction (defaulting to clockwise), is added to the path, connected to
// the previous point by a straight line.
func (dc *Context) Ellipse(x, y, radiusX, radiusY, rotation, startAngle, endAngle float64, counterclockwise bool) {
	dc.path2d.Ellipse(x, y, radiusX, radiusY, rotation, startAngle, endAngle, counterclockwise)
}
