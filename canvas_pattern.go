// Copyright 2006 The Android Open Source Project
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package tinyskia

type Pattern interface {
	Style
	SetTransform()
	SetTransformWithMatrix(transform *Matrix)
}

// SetTransform resets the pattern's transformation matrix to the identity matrix.
func (p *ImagePattern) SetTransform() {
	p.transform = NewMatrixIdentity()
}

// SetTransformWithMatrix sets the pattern's transform matrix to the specified matrix.
func (p *ImagePattern) SetTransformWithMatrix(transform *Matrix) {
	if transform == nil {
		transform = NewMatrixIdentity()
	}
	p.transform = transform
}
