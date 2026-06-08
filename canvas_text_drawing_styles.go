// Copyright 2006 The Android Open Source Project
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package tinyskia

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"slices"
	"sync"
	"sync/atomic"
	"unicode"

	"golang.org/x/exp/mmap"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/text/language"
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
	dc.fontChain0 = make(map[string][]typeface)
	dc.fontChain1 = make(map[language.Script][]typeface)
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
	dc.fontChain0 = make(map[string][]typeface)
	dc.fontChain1 = make(map[language.Script][]typeface)
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

type typeface struct {
	font  *sfnt.Font
	file  string
	index int

	family  string
	generic FontGeneric
	weight  FontWeight
	style   FontStyle
}

type registry struct {
	sync.RWMutex
	version uint64
	assets  map[string][]typeface
	locals  map[language.Script][]typeface
}

var fonts = &registry{
	assets: make(map[string][]typeface),
	locals: make(map[language.Script][]typeface),
}

func RegisterFont(path string, fontFace FontFace) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return RegisterFontWithData(data, fontFace)
}

func RegisterFontWithData(data []byte, fontFace FontFace) error {
	if fontFace.Family == "" {
		return errors.New("family name cannot be empty")
	}
	if fontFace.Weight < 0 {
		fontFace.Weight = FontWeightNormal
	}
	if fontFace.Style == "" {
		fontFace.Style = FontStyleNormal
	}

	collection, err := sfnt.ParseCollection(data)
	if err != nil {
		return err
	}
	ttf, err := collection.Font(0)
	if err != nil {
		return err
	}
	typefaces := typeface{
		font:   ttf,
		family: fontFace.Family,
		weight: fontFace.Weight,
		style:  fontFace.Style,
	}

	fonts.Lock()
	defer fonts.Unlock()
	atomic.AddUint64(&fonts.version, 1)
	if _, ok := fonts.assets[fontFace.Family]; ok {
		fonts.assets[fontFace.Family] = append(fonts.assets[fontFace.Family], typefaces)
	} else {
		fonts.assets[fontFace.Family] = []typeface{typefaces}
	}
	return nil
}

func init() {
	fonts.Lock()
	defer fonts.Unlock()
	switch os := runtime.GOOS; os {
	case "windows":
		var base = "C:\\Windows\\Fonts\\"
		var (
			seguisb = typeface{family: "Segoe UI", generic: FontGenericSansSerif, weight: FontWeightBold, style: FontStyleNormal, file: base + "seguisb.ttf"}
			arial   = typeface{family: "Arial", generic: FontGenericSansSerif, weight: FontWeightNormal, style: FontStyleNormal, file: base + "arial.ttf"}
			ariali  = typeface{family: "Arial", generic: FontGenericSansSerif, weight: FontWeightNormal, style: FontStyleItalic, file: base + "ariali.ttf"}
			msyh    = typeface{family: "Microsoft YaHei", generic: FontGenericSansSerif, weight: FontWeightNormal, style: FontStyleNormal, file: base + "msyh.ttc"}
			msyhbd  = typeface{family: "Microsoft YaHei", generic: FontGenericSansSerif, weight: FontWeightBold, style: FontStyleNormal, file: base + "msyhbd.ttc"}
			simsun  = typeface{family: "SimSun", generic: FontGenericSerif, weight: FontWeightNormal, style: FontStyleNormal, file: base + "simsun.ttc"}
		)
		fonts.locals = map[language.Script][]typeface{
			LangScriptLatin: {
				seguisb, arial, ariali},
			LangScriptHans: {msyh, msyhbd, simsun},
		}
	}
}

var (
	LangBaseZH      = language.MustParseBase("zh")
	LangBaseJA      = language.MustParseBase("ja")
	LangRegionCN    = language.MustParseRegion("CN")
	LangRegionTW    = language.MustParseRegion("TW")
	LangRegionHK    = language.MustParseRegion("HK")
	LangRegionMO    = language.MustParseRegion("MO")
	LangRegionSG    = language.MustParseRegion("SG")
	LangRegionMY    = language.MustParseRegion("MY")
	LangRegionJP    = language.MustParseRegion("JP")
	LangRegionKR    = language.MustParseRegion("KR")
	LangScriptLatin = language.MustParseScript("Latn")
	LangScriptHans  = language.MustParseScript("Hans")
	LangScriptHant  = language.MustParseScript("Hant")
	LangScriptKore  = language.MustParseScript("Kore")
	LangScriptJpan  = language.MustParseScript("Jpan")
	LangScriptZzzz  = language.MustParseScript("Zzzz")
)

func loadFont(font typeface) (*sfnt.Font, error) {
	if font.font != nil {
		return font.font, nil
	}
	reader, err := mmap.Open(font.file)
	if err != nil {
		return nil, err
	}

	ft, err := sfnt.ParseReaderAt(reader)
	if err != nil {
		_ = reader.Close()
		return nil, err
	}
	font.font = ft
	return ft, nil
}

func (ctx *Context) fmatch0(family string) (candidates []typeface) {
	if family == "" {
		return
	}

	if fonts.version == ctx.fontVersion {
		if candidates, hit := ctx.fontChain0[family]; hit {
			return candidates
		}
	} else {
		ctx.fontVersion = fonts.version
		ctx.fontChain0 = make(map[string][]typeface)
	}

	fonts.RLock()
	if assets, ok := fonts.assets[family]; ok {
		candidates = append(candidates, assets...)
	}
	fonts.RUnlock()

	if family == FontGenericFantasy || family == FontGenericMath || family == FontGenericEmoji ||
		family == FontGenericSerif || family == FontGenericSansSerif ||
		family == FontGenericCursive || family == FontGenericMonospace {
		var script language.Script
		b, s, r := ctx.lang.Raw()
		if s == LangScriptZzzz {
			switch {
			case b == LangBaseZH, r == LangRegionCN:
				script = LangScriptHans
			case r == LangRegionTW, r == LangRegionHK, r == LangRegionMO,
				r == LangRegionSG, r == LangRegionMY:
				script = LangScriptHant
			default:
				script = LangScriptLatin
			}
		}
		for _, typeface := range fonts.locals[script] {
			if typeface.generic == family {
				candidates = append(candidates, typeface)
			}
		}
	} else {
		for _, typefaces := range fonts.locals {
			for _, typeface := range typefaces {
				if typeface.family == family {
					candidates = append(candidates, typeface)
				}
			}
		}
	}

	slices.SortFunc(candidates, func(a, b typeface) int {
		if (a.style == ctx.font.Style) != (b.style == ctx.font.Style) {
			if a.style == ctx.font.Style {
				return -1
			}
			return 1
		}

		if (a.weight == ctx.font.Weight) != (b.weight == ctx.font.Weight) {
			if a.weight == ctx.font.Weight {
				return -1
			}
			return 1
		}

		return int(fabs(a.weight-ctx.font.Weight) - fabs(b.weight-ctx.font.Weight))
	})

	ctx.fontChain0[family] = candidates
	return candidates
}

func (ctx *Context) fmatch1(r rune) (candidates []typeface) {

	var script language.Script
	switch {
	case unicode.Is(unicode.Latin, r):
		script = LangScriptLatin
	case unicode.Is(unicode.Han, r):
		_, s, r := ctx.lang.Raw()
		switch {
		case s == LangScriptHans, r == LangRegionCN:
			script = LangScriptHans
		case s == LangScriptHant,
			r == LangRegionTW, r == LangRegionHK:
			script = LangScriptHant
		default:
			script = LangScriptHant
		}
	case unicode.Is(unicode.Hangul, r):
		script = LangScriptKore
	case unicode.Is(unicode.Hiragana, r) || unicode.Is(unicode.Katakana, r):
		script = LangScriptJpan
	default:
		script = LangScriptLatin
	}

	if fonts.version == ctx.fontVersion {
		if candidates, hit := ctx.fontChain1[script]; hit {
			return candidates
		}
	} else {
		ctx.fontVersion = fonts.version
		ctx.fontChain0 = make(map[string][]typeface)
	}

	if typefaces, ok := fonts.locals[script]; ok {
		candidates = append(candidates, typefaces...)
	}

	slices.SortStableFunc(candidates, func(a, b typeface) int {
		if (a.style == ctx.font.Style) != (b.style == ctx.font.Style) {
			if a.style == ctx.font.Style {
				return -1
			}
			return 1
		}

		if (a.weight == ctx.font.Weight) != (b.weight == ctx.font.Weight) {
			if a.weight == ctx.font.Weight {
				return -1
			}
			return 1
		}

		return int(fabs(a.weight-ctx.font.Weight) - fabs(b.weight-ctx.font.Weight))
	})
	ctx.fontChain1[script] = candidates
	fmt.Println(script, candidates)
	return candidates
}

func fabs(x FontWeight) FontWeight {
	if x < 0 {
		return -x
	}
	return x
}
