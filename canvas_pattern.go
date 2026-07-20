// Copyright 2006 The Android Open Source Project
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package tinyskia

import "errors"

var ErrInvalidPatternTransform = errors.New("pattern transform is invalid")

type Pattern interface {
	Style
	SetTransform() error
	SetTransformWithMatrix(transform *Matrix) error
}

// SetTransform resets the pattern's transformation matrix to the identity matrix.
func (p *ImagePattern) SetTransform() error {
	p.transform = NewMatrixIdentity()
	return nil
}

// SetTransformWithMatrix sets the pattern's transform matrix to the specified matrix.
func (p *ImagePattern) SetTransformWithMatrix(transform *Matrix) error {
	if transform == nil {
		return ErrInvalidPatternTransform
	}
	p.transform = transform
	return nil
}
