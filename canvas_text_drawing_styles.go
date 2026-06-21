// Copyright 2006 The Android Open Source Project
// Copyright 2012 The Chromium Authors
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package tinyskia

import (
	"os"

	"github.com/lumifloat/tinyskia/internal/text"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/text/language"
)

var script = text.ScriptCache{}
var fonts = text.NewRegistry()

func init() {
	fonts.ScanSystemFonts()
}

type FontAttr struct {
	Family []string
	Weight FontWeight
	Style  FontStyle
	Size   float64
}

type FontFace struct {
	Family string
	Weight FontWeight
	Style  FontStyle
}

type FontStyle = string

var (
	FontStyleNormal  FontStyle = FontStyle("normal")
	FontStyleItalic  FontStyle = FontStyle("italic")
	FontStyleOblique FontStyle = FontStyle("oblique")
)

type FontWeight = float64

const (
	FontWeightThin       FontWeight = 100
	FontWeightExtraLight FontWeight = 200
	FontWeightUltraLight FontWeight = 200
	FontWeightLight      FontWeight = 300
	FontWeightRegular    FontWeight = 400
	FontWeightNormal     FontWeight = 400
	FontWeightMedium     FontWeight = 500
	FontWeightSemiBold   FontWeight = 600
	FontWeightDemiBold   FontWeight = 600
	FontWeightBold       FontWeight = 700
	FontWeightExtraBold  FontWeight = 800
	FontWeightUltraBold  FontWeight = 800
	FontWeightBlack      FontWeight = 900
	FontWeightHeavy      FontWeight = 900
	FontWeightExtraBlack FontWeight = 950
	FontWeightUltraBlack FontWeight = 950
)

type FontGeneric = string

const (
	FontGenericFantasy   FontGeneric = "fantasy"
	FontGenericMath      FontGeneric = "math"
	FontGenericEmoji     FontGeneric = "emoji"
	FontGenericSerif     FontGeneric = "serif"
	FontGenericSansSerif FontGeneric = "sans-serif"
	FontGenericCursive   FontGeneric = "cursive"
	FontGenericMonospace FontGeneric = "monospace"
)

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
	tag, err := language.Parse(lang)
	if err != nil {
		return
	}
	dc.lang = tag
}

// GetLang
func (dc *Context) GetLang() string {
	return dc.lang.String()
}

// SetFont
func (dc *Context) SetFont(font FontAttr) {
	if len(font.Family) == 0 {
		return
	}
	for _, family := range font.Family {
		if family == "" {
			return
		}
	}
	dc.font = font
}

// GetFont
func (dc *Context) GetFont() FontAttr {
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

func RegisterFont(font *sfnt.Font, face FontFace) error {
	return fonts.RegisterFont(font, face.Family, face.Weight, face.Style)
}

func RegisterFontWithFile(file string, face FontFace) error {
	fi, err := os.Open(file)
	if err != nil {
		return err
	}
	collection, err := sfnt.ParseCollectionReaderAt(fi)
	if err != nil {
		return err
	}
	if collection.NumFonts() == 0 {
		return sfnt.ErrNotFound
	}
	ttf, err := collection.Font(0)
	if err != nil {
		return err
	}
	return fonts.RegisterFont(ttf, face.Family, face.Weight, face.Style)
}

func RegisterFontWithData(data []byte, face FontFace) error {
	collection, err := sfnt.ParseCollection(data)
	if err != nil {
		return err
	}
	if collection.NumFonts() == 0 {
		return sfnt.ErrNotFound
	}
	ttf, err := collection.Font(0)
	if err != nil {
		return err
	}
	return fonts.RegisterFont(ttf, face.Family, face.Weight, face.Style)
}
