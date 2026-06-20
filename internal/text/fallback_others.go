//go:build !darwin && !linux && !windows
// +build !darwin,!linux,!windows

// Copyright 2012 The Chromium Authors
// Copyright 2026 Mozilla
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package text

func GetFallbackFonts(script Script) []string {
	return nil
}
