// Copyright 2006 The Android Open Source Project
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package tinyskia

import "math"

// Scale changes the current transformation matrix to apply a scaling transformation with the given characteristics.
func (dc *Context) Scale(x, y float64) {
	dc.matrix = dc.matrix.Scale(x, y)
}

// Rotate changes the current transformation matrix to apply a scaling transformation with the given characteristics.
func (dc *Context) Rotate(angle float64) {
	angleDegrees := angle * 180.0 / math.Pi
	dc.matrix = dc.matrix.Rotate(angleDegrees)
}

// Translate changes the current transformation matrix to apply a translation transformation with the given characteristics.
func (dc *Context) Translate(x, y float64) {
	dc.matrix = dc.matrix.Translate(x, y)
}

// Transform hanges the current transformation matrix to apply the matrix given by the arguments as described below.
func (dc *Context) Transform(a, b, c, d, e, f float64) {
	dc.matrix = dc.matrix.Multiply(NewMatrix(a, b, c, d, e, f))
}

// GetMatrix returns a copy of the current transformation matrix, as a newly created DOMMatrix object.
func (dc *Context) GetMatrix() *matrix {
	return &matrix{transform: dc.matrix.transform}
}

// SetTranform hanges the current transformation matrix to the matrix given by the arguments as described below.
func (dc *Context) SetTransform(a, b, c, d, e, f float64) {
	dc.matrix = NewMatrix(a, b, c, d, e, f)
}

// SetTransformWithMatrix sets the current transformation matrix to the given matrix.
func (dc *Context) SetTransformWithMatrix(matrix *matrix) {
	dc.matrix = matrix
}

// ResetTransform resets the current transformation matrix to the identity matrix.
func (dc *Context) ResetTransform() {
	dc.matrix = NewMatrixIdentity()
}

// RotateAbout rotates the current transformation matrix around a point (x, y).
// This is equivalent to: translate(x, y); rotate(angle); translate(-x, -y)
func (dc *Context) RotateAbout(angle, x, y float64) {
	dc.Translate(x, y)
	dc.Rotate(angle)
	dc.Translate(-x, -y)
}
