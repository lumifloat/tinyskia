// Copyright 2006 The Android Open Source Project
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package tinyskia

type Pattern interface {
	style
	SetTransform(transform *matrix)
}

func (p *imagePattern) SetTransform(transform *matrix) {
	p.transform = transform
}
