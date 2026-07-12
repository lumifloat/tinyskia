// Copyright 2006 The Android Open Source Project
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package tinyskia

import "math"

// Scale changes the current transformation matrix to apply a scaling transformation with the given characteristics.
func (ctx *Context) Scale(x, y float64) {
	ctx.matrix = ctx.matrix.Scale(x, y)
}

// Rotate changes the current transformation matrix to apply a scaling transformation with the given characteristics.
func (ctx *Context) Rotate(angle float64) {
	angle = angle * 180.0 / math.Pi
	ctx.matrix = ctx.matrix.Rotate(angle)
}

// Translate changes the current transformation matrix to apply a translation transformation with the given characteristics.
func (ctx *Context) Translate(x, y float64) {
	ctx.matrix = ctx.matrix.Translate(x, y)
}

// Transform changes the current transformation matrix to apply the matrix given by the arguments as described below.
func (ctx *Context) Transform(a, b, c, d, e, f float64) {
	ctx.matrix = ctx.matrix.Multiply(NewMatrix(a, b, c, d, e, f))
}

// GetTransform returns a copy of the current transformation matrix, as a newly created DOMMatrix object.
func (ctx *Context) GetTransform() *Matrix {
	return &Matrix{transform: ctx.matrix.transform}
}

// SetTransform changes the current transformation matrix to the matrix given by the arguments as described below.
func (ctx *Context) SetTransform(matrix *Matrix) {
	if matrix == nil {
		ctx.matrix = NewMatrixIdentity()
	} else {
		ctx.matrix = &Matrix{matrix.transform}
	}
}

// SetTransformValues changes the current transformation matrix to the matrix given by the arguments as described below.
func (ctx *Context) SetTransformValues(a, b, c, d, e, f float64) {
	ctx.matrix = NewMatrix(a, b, c, d, e, f)
}

// ResetTransform resets the current transformation matrix to the identity matrix.
func (ctx *Context) ResetTransform() {
	ctx.matrix = NewMatrixIdentity()
}
