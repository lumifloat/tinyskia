// Copyright 2006 The Android Open Source Project
// Copyright 2012 The Chromium Authors
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package tinyskia

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/go-text/typesetting/di"
	"github.com/go-text/typesetting/font"
	"github.com/go-text/typesetting/fontscan"
)

var (
	fonts *fontscan.FontMap
	flock sync.Mutex
)

func init() {
	// TODO
	fonts = fontscan.NewFontMap(log.New(io.Discard, "", 0))
	cacheDir := ""
	if runtime.GOOS == "android" {
		parent := os.Getenv("FILESDIR")
		cacheDir = filepath.Join(parent, "fontcache")
	}

	fonts.UseSystemFonts(cacheDir)
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

type FontStyle = uint8

var (
	FontStyleNormal  FontStyle = FontStyle(font.StyleNormal)
	FontStyleItalic  FontStyle = FontStyle(font.StyleItalic)
	FontStyleOblique FontStyle = FontStyle(font.StyleItalic)
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

type Direction uint8

const (
	DirectionLTR = Direction(di.DirectionLTR)
	DirectionRTL = Direction(di.DirectionRTL)
)

type FontStretch float32

const (
	FontStretchUltraCondensed = FontStretch(font.StretchUltraCondensed)
	FontStretchExtraCondensed = FontStretch(font.StretchExtraCondensed)
	FontStretchCondensed      = FontStretch(font.StretchCondensed)
	FontStretchSemiCondensed  = FontStretch(font.StretchSemiCondensed)
	FontStretchNormal         = FontStretch(font.StretchNormal)
	FontStretchSemiExpanded   = FontStretch(font.StretchSemiExpanded)
	FontStretchExpanded       = FontStretch(font.StretchExpanded)
	FontStretchExtraExpanded  = FontStretch(font.StretchExtraExpanded)
	FontStretchUltraExpanded  = FontStretch(font.StretchUltraExpanded)
)

type FontVariant string

const (
	FontVariantNormal        FontVariant = "normal"
	FontVariantSmallCaps     FontVariant = "small-caps"
	FontVariantAllSmallCaps  FontVariant = "all-small-caps"
	FontVariantPetiteCaps    FontVariant = "petite-caps"
	FontVariantAllPetiteCaps FontVariant = "all-petite-caps"
	FontVariantUnicase       FontVariant = "unicase"
	FontVariantTitlingCaps   FontVariant = "titling-caps"
)

type FontKerning uint32

const (
	// FontKerningAuto
	FontKerningAuto FontKerning = 1
	// FontKerningNormal
	FontKerningNormal FontKerning = 1
	// FontKerningNone
	FontKerningNone FontKerning = 0
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

// SetDirection
func (dc *Context) SetDirection(direction Direction) {
	dc.direction = direction
}

// GetDirection
func (dc *Context) GetDirection() Direction {
	return dc.direction
}

// SetFontStretch
func (dc *Context) SetFontStretch(stretch FontStretch) {
	dc.fontStretch = stretch
}

// GetFontStretch
func (dc *Context) GetFontStretch() FontStretch {
	return dc.fontStretch
}

// SetFontVariant
func (dc *Context) SetFontVariant(variant FontVariant) {
	dc.fontVariant = variant
}

// GetFontVariant
func (dc *Context) GetFontVariant() FontVariant {
	return dc.fontVariant
}

// SetFontKerning
func (dc *Context) SetFontKerning(kerning FontKerning) {
	dc.fontKerning = kerning
}

// GetFontKerning
func (dc *Context) GetFontKerning() FontKerning {
	return dc.fontKerning
}

func RegisterFont(file string, face FontFace) error {
	if face.Weight <= 0 {
		face.Weight = FontWeightNormal
	}
	if face.Style == 0 {
		face.Style = FontStyleNormal
	}
	flock.Lock()
	defer flock.Unlock()
	fi, err := os.Open(file)
	if err != nil {
		return err
	}
	faces, err := font.ParseTTC(fi)
	if err != nil {
		return fmt.Errorf("unsupported font resource: %s", err)
	}
	loc := fontscan.Location{File: file}
	desc := font.Description{
		Family: face.Family,
		Aspect: font.Aspect{
			Weight: font.Weight(face.Weight),
			Style:  font.Style(face.Style),
		},
	}
	fonts.AddFace(faces[0], loc, desc)
	return nil
}

func RegisterFontWithResource(file font.Resource, location string, face FontFace) error {
	if face.Weight <= 0 {
		face.Weight = FontWeightNormal
	}
	if face.Style == 0 {
		face.Style = FontStyleNormal
	}
	flock.Lock()
	defer flock.Unlock()
	faces, err := font.ParseTTC(file)
	if err != nil {
		return fmt.Errorf("unsupported font resource: %s", err)
	}
	loc := fontscan.Location{File: location}
	desc := font.Description{
		Family: face.Family,
		Aspect: font.Aspect{
			Weight: font.Weight(face.Weight),
			Style:  font.Style(face.Style),
		},
	}
	fonts.AddFace(faces[0], loc, desc)
	return nil
}
