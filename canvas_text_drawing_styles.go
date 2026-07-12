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
	"os"

	"github.com/go-text/typesetting/di"
	"github.com/go-text/typesetting/font"
	"github.com/go-text/typesetting/font/opentype"
	"github.com/go-text/typesetting/fontscan"
	"github.com/go-text/typesetting/shaping"
	"github.com/lumifloat/tinyskia/internal/text"
)

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

type CanvasTextAlign string

const (
	CanvasTextAlignStart  CanvasTextAlign = "start"
	CanvasTextAlignEnd    CanvasTextAlign = "end"
	CanvasTextAlignLeft   CanvasTextAlign = "left"
	CanvasTextAlignCenter CanvasTextAlign = "center"
	CanvasTextAlignRight  CanvasTextAlign = "right"
)

type CanvasTextBaseline string

const (
	CanvasTextBaselineTop         CanvasTextBaseline = "top"
	CanvasTextBaselineHanging     CanvasTextBaseline = "hanging"
	CanvasTextBaselineMiddle      CanvasTextBaseline = "middle"
	CanvasTextBaselineAlphabetic  CanvasTextBaseline = "alphabetic"
	CanvasTextBaselineIdeographic CanvasTextBaseline = "ideographic"
	CanvasTextBaselineBottom      CanvasTextBaseline = "bottom"
)

type CanvasDirection string

const (
	CanvasDirectionLTR     CanvasDirection = "ltr"
	CanvasDirectionRTL     CanvasDirection = "rtl"
	CanvasDirectionInherit CanvasDirection = "inherit"
)

type CanvasFontKerning string

const (
	CanvasFontKerningAuto   CanvasFontKerning = "auto"
	CanvasFontKerningNormal CanvasFontKerning = "normal"
	CanvasFontKerningNone   CanvasFontKerning = "none"
)

type CanvasFontStretch string

const (
	CanvasFontStretchUltraCondensed CanvasFontStretch = "ultra-condensed"
	CanvasFontStretchExtraCondensed CanvasFontStretch = "extra-condensed"
	CanvasFontStretchCondensed      CanvasFontStretch = "condensed"
	CanvasFontStretchSemiCondensed  CanvasFontStretch = "semi-condensed"
	CanvasFontStretchNormal         CanvasFontStretch = "normal"
	CanvasFontStretchSemiExpanded   CanvasFontStretch = "semi-expanded"
	CanvasFontStretchExpanded       CanvasFontStretch = "expanded"
	CanvasFontStretchExtraExpanded  CanvasFontStretch = "extra-expanded"
	CanvasFontStretchUltraExpanded  CanvasFontStretch = "ultra-expanded"
)

type CanvasFontVariantCaps string

const (
	CanvasFontVariantCapsNormal        CanvasFontVariantCaps = "normal"
	CanvasFontVariantCapsSmallCaps     CanvasFontVariantCaps = "small-caps"
	CanvasFontVariantCapsAllSmallCaps  CanvasFontVariantCaps = "all-small-caps"
	CanvasFontVariantCapsPetiteCaps    CanvasFontVariantCaps = "petite-caps"
	CanvasFontVariantCapsAllPetiteCaps CanvasFontVariantCaps = "all-petite-caps"
	CanvasFontVariantCapsUnicase       CanvasFontVariantCaps = "unicase"
	CanvasFontVariantCapsTitlingCaps   CanvasFontVariantCaps = "titling-caps"
)

type CanvasTextRendering string

const (
	CanvasTextRenderingAuto               CanvasTextRendering = "auto"
	CanvasTextRenderingOptimizeSpeed      CanvasTextRendering = "optimizeSpeed"
	CanvasTextRenderingOptimizeLegibility CanvasTextRendering = "optimizeLegibility"
	CanvasTextRenderingGeometricPrecision CanvasTextRendering = "geometricPrecision"
)

// SetLang
func (ctx *Context) SetLang(lang string) {
	ctx.lang = lang
}

// GetLang
func (ctx *Context) GetLang() string {
	return ctx.lang
}

// SetFont
func (ctx *Context) SetFont(font FontAttr) {
	if len(font.Family) == 0 {
		return
	}
	for _, family := range font.Family {
		if family == "" {
			return
		}
	}
	ctx.font = font
}

// GetFont
func (ctx *Context) GetFont() FontAttr {
	return ctx.font
}

// SetTextAlign
func (ctx *Context) SetTextAlign(align CanvasTextAlign) {
	ctx.textAlign = align
}

// GetTextAlign
func (ctx *Context) GetTextAlign() CanvasTextAlign {
	return ctx.textAlign
}

// SetTextBaseline
func (ctx *Context) SetTextBaseline(baseline CanvasTextBaseline) {
	ctx.textBaseline = baseline
}

// GetTextBaseline
func (ctx *Context) GetTextBaseline() CanvasTextBaseline {
	return ctx.textBaseline
}

// SetDirection
func (ctx *Context) SetDirection(direction CanvasDirection) {
	ctx.direction = direction
}

// GetDirection
func (ctx *Context) GetDirection() CanvasDirection {
	return ctx.direction
}

// SetFontKerning
func (ctx *Context) SetFontKerning(kerning CanvasFontKerning) {
	ctx.fontKerning = kerning
}

// GetFontKerning
func (ctx *Context) GetFontKerning() CanvasFontKerning {
	return ctx.fontKerning
}

// SetFontStretch
func (ctx *Context) SetFontStretch(stretch CanvasFontStretch) {
	ctx.fontStretch = stretch
}

// GetFontStretch
func (ctx *Context) GetFontStretch() CanvasFontStretch {
	return ctx.fontStretch
}

// SetFontVariant
func (ctx *Context) SetFontVariantCaps(variant CanvasFontVariantCaps) {
	ctx.fontVariantCaps = variant
}

// GetFontVariant
func (ctx *Context) GetFontVariantCaps() CanvasFontVariantCaps {
	return ctx.fontVariantCaps
}

// SetTextRendering
func (ctx *Context) SetTextRendering(rendering CanvasTextRendering) {
	ctx.textRendering = rendering
}

// GetTextRendering
func (ctx *Context) GetTextRendering() CanvasTextRendering {
	return ctx.textRendering
}

// SetWordSpacing
func (ctx *Context) SetWordSpacing(spacing float64) {
	// TODO
	ctx.wordSpacing = spacing
}

// GetWordSpacing
func (ctx *Context) GetWordSpacing() float64 {
	return ctx.wordSpacing
}

func direction(d CanvasDirection) di.Direction {
	switch d {
	case CanvasDirectionLTR:
		return di.DirectionLTR
	case CanvasDirectionRTL:
		return di.DirectionRTL
	default:
		return di.DirectionLTR
	}
}

func kerning(k CanvasFontKerning, features []shaping.FontFeature) {
	switch k {
	case CanvasFontKerningAuto:
		// pass
	case CanvasFontKerningNormal:
		features = append(features, shaping.FontFeature{Tag: opentype.NewTag('k', 'e', 'r', 'n'), Value: 1})
	case CanvasFontKerningNone:
		features = append(features, shaping.FontFeature{Tag: opentype.NewTag('k', 'e', 'r', 'n'), Value: 0})
	default:
		// pass
	}
}

func stretch(s CanvasFontStretch) font.Stretch {
	switch s {
	case CanvasFontStretchUltraCondensed:
		return font.StretchUltraCondensed
	case CanvasFontStretchExtraCondensed:
		return font.StretchExtraCondensed
	case CanvasFontStretchCondensed:
		return font.StretchCondensed
	case CanvasFontStretchSemiCondensed:
		return font.StretchSemiCondensed
	case CanvasFontStretchNormal:
		return font.StretchNormal
	case CanvasFontStretchSemiExpanded:
		return font.StretchSemiExpanded
	case CanvasFontStretchExpanded:
		return font.StretchExpanded
	case CanvasFontStretchExtraExpanded:
		return font.StretchExtraExpanded
	case CanvasFontStretchUltraExpanded:
		return font.StretchUltraExpanded
	default:
		return font.StretchNormal
	}
}

func variant(v CanvasFontVariantCaps, features []shaping.FontFeature) {
	switch v {
	case CanvasFontVariantCapsNormal:
		break
	case CanvasFontVariantCapsSmallCaps:
		features = append(features, shaping.FontFeature{
			Tag:   opentype.NewTag('s', 'm', 'c', 'p'),
			Value: 1,
		})
	case CanvasFontVariantCapsAllSmallCaps:
		features = append(features, shaping.FontFeature{
			Tag:   opentype.NewTag('c', '2', 's', 'c'),
			Value: 1,
		}, shaping.FontFeature{
			Tag:   opentype.NewTag('s', 'm', 'c', 'p'),
			Value: 1,
		})
	case CanvasFontVariantCapsPetiteCaps:
		features = append(features, shaping.FontFeature{
			Tag:   opentype.NewTag('p', 'c', 'a', 'p'),
			Value: 1,
		})
	case CanvasFontVariantCapsAllPetiteCaps:
		features = append(features, shaping.FontFeature{
			Tag:   opentype.NewTag('c', '2', 'p', 'c'),
			Value: 1,
		}, shaping.FontFeature{
			Tag:   opentype.NewTag('p', 'c', 'a', 'p'),
			Value: 1,
		})
	case CanvasFontVariantCapsUnicase:
		features = append(features, shaping.FontFeature{
			Tag:   opentype.NewTag('u', 'n', 'i', 'c'),
			Value: 1,
		})
	case CanvasFontVariantCapsTitlingCaps:
		features = append(features, shaping.FontFeature{
			Tag:   opentype.NewTag('t', 'i', 't', 'l'),
			Value: 1,
		})
	default:
		// pass
	}
}

func RegisterFont(file string, face FontFace) error {
	if face.Weight <= 0 {
		face.Weight = FontWeightNormal
	}
	if face.Style == 0 {
		face.Style = FontStyleNormal
	}
	text.FontLock.Lock()
	defer text.FontLock.Unlock()
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
	text.FontMap.AddFace(faces[0], loc, desc)
	return nil
}

func RegisterFontWithResource(file font.Resource, location string, face FontFace) error {
	if face.Weight <= 0 {
		face.Weight = FontWeightNormal
	}
	if face.Style == 0 {
		face.Style = FontStyleNormal
	}
	text.FontLock.Lock()
	defer text.FontLock.Unlock()
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
	text.FontMap.AddFace(faces[0], loc, desc)
	return nil
}
