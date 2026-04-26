// Copyright 2006 The Android Open Source Project
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package tinyskia

type TextAlign int

const (
	// Align to the start edge of the text (left side in left-to-right text, right side in right-to-left text).
	TextAlignStart TextAlign = iota
	// Align to the end edge of the text (right side in left-to-right text, left side in right-to-left text).
	TextAlignEnd
	// Align to the left.
	TextAlignLeft
	// Align to the left.
	TextAlignCenter
	// Align to the center.
	TextAlignRight
)

type FontKerning int

const (
	// FontKerningAuto
	FontKerningAuto FontKerning = iota
	// FontKerningNormal
	FontKerningNormal
	// FontKerningNone
	FontKerningNone
)

// SetLang
func (dc *Context) SetLang(lang string) {
	dc.lang = lang
}

// GetLang
func (dc *Context) GetLang() string {
	return dc.lang
}

// SetFont
func (dc *Context) SetFont(font *Font) {
	dc.font = font
}

// GetFont
func (dc *Context) GetFont() *Font {
	return dc.font
}

// SetTextAlign
func (dc *Context) SetTextAlign(align TextAlign) {
	dc.textAlign = align
}

// GetTextAlign
func (dc *Context) GetTextAlign() TextAlign {
	return dc.textAlign
}

// SetFontKerning
func (dc *Context) SetFontKerning(kerning FontKerning) {
	dc.fontKerning = kerning
}

// GetFontKerning
func (dc *Context) GetFontKerning() FontKerning {
	return dc.fontKerning
}
