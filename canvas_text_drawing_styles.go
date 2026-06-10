// Copyright 2006 The Android Open Source Project
// Copyright 2012 The Chromium Authors
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

var (
	_LANG_BASE_ZH = language.MustParseBase("zh")
	_LANG_BASE_JA = language.MustParseBase("ja")

	_LANG_REGION_CN = language.MustParseRegion("CN")
	_LANG_REGION_TW = language.MustParseRegion("TW")
	_LANG_REGION_HK = language.MustParseRegion("HK")
	_LANG_REGION_MO = language.MustParseRegion("MO")
	_LANG_REGION_SG = language.MustParseRegion("SG")
	_LANG_REGION_MY = language.MustParseRegion("MY")
	_LANG_REGION_JP = language.MustParseRegion("JP")
	_LANG_REGION_KR = language.MustParseRegion("KR")

	_LANG_SCRIPT_LATIN           = language.MustParseScript("Latn") // 拉丁/英文
	_LANG_SCRIPT_CYRILLIC        = language.MustParseScript("Cyrl") // 西里尔文
	_LANG_SCRIPT_ARABIC          = language.MustParseScript("Arab") // 阿拉伯文
	_LANG_SCRIPT_DEVANAGARI      = language.MustParseScript("Deva") // 天城文/印度因纽特
	_LANG_SCRIPT_GREEK           = language.MustParseScript("Grek") // 希腊文
	_LANG_SCRIPT_JAPANESE        = language.MustParseScript("Jpan") // 日文
	_LANG_SCRIPT_KOREAN          = language.MustParseScript("Kore") // 韩文
	_LANG_SCRIPT_SIMPLIFIED_HAN  = language.MustParseScript("Hans") // 简体中文
	_LANG_SCRIPT_TRADITIONAL_HAN = language.MustParseScript("Hant") // 繁体中文
	_LANG_SCRIPT_ZZZZ            = language.MustParseScript("Zzzz") // 未知符号

	_IDS_STANDARD_FONT_FAMILY                   []string
	_IDS_FIXED_FONT_FAMILY                      []string
	_IDS_SERIF_FONT_FAMILY                      []string
	_IDS_SANS_SERIF_FONT_FAMILY                 []string
	_IDS_CURSIVE_FONT_FAMILY                    []string
	_IDS_FANTASY_FONT_FAMILY                    []string
	_IDS_MATH_FONT_FAMILY                       []string
	_IDS_FIXED_FONT_FAMILY_ARABIC               []string
	_IDS_SANS_SERIF_FONT_FAMILY_ARABIC          []string
	_IDS_STANDARD_FONT_FAMILY_CYRILLIC          []string
	_IDS_FIXED_FONT_FAMILY_CYRILLIC             []string
	_IDS_SERIF_FONT_FAMILY_CYRILLIC             []string
	_IDS_SANS_SERIF_FONT_FAMILY_CYRILLIC        []string
	_IDS_STANDARD_FONT_FAMILY_GREEK             []string
	_IDS_FIXED_FONT_FAMILY_GREEK                []string
	_IDS_SERIF_FONT_FAMILY_GREEK                []string
	_IDS_SANS_SERIF_FONT_FAMILY_GREEK           []string
	_IDS_STANDARD_FONT_FAMILY_JAPANESE          []string
	_IDS_FIXED_FONT_FAMILY_JAPANESE             []string
	_IDS_SERIF_FONT_FAMILY_JAPANESE             []string
	_IDS_SANS_SERIF_FONT_FAMILY_JAPANESE        []string
	_IDS_STANDARD_FONT_FAMILY_KOREAN            []string
	_IDS_FIXED_FONT_FAMILY_KOREAN               []string
	_IDS_SERIF_FONT_FAMILY_KOREAN               []string
	_IDS_SANS_SERIF_FONT_FAMILY_KOREAN          []string
	_IDS_CURSIVE_FONT_FAMILY_KOREAN             []string
	_IDS_STANDARD_FONT_FAMILY_SIMPLIFIED_HAN    []string
	_IDS_FIXED_FONT_FAMILY_SIMPLIFIED_HAN       []string
	_IDS_SERIF_FONT_FAMILY_SIMPLIFIED_HAN       []string
	_IDS_SANS_SERIF_FONT_FAMILY_SIMPLIFIED_HAN  []string
	_IDS_CURSIVE_FONT_FAMILY_SIMPLIFIED_HAN     []string
	_IDS_STANDARD_FONT_FAMILY_TRADITIONAL_HAN   []string
	_IDS_FIXED_FONT_FAMILY_TRADITIONAL_HAN      []string
	_IDS_SERIF_FONT_FAMILY_TRADITIONAL_HAN      []string
	_IDS_SANS_SERIF_FONT_FAMILY_TRADITIONAL_HAN []string
	_IDS_CURSIVE_FONT_FAMILY_TRADITIONAL_HAN    []string
	_IDS_STANDARD_FONT_FAMILY_DEVANAGARI        []string
	_IDS_FIXED_FONT_FAMILY_DEVANAGARI           []string
	_IDS_SERIF_FONT_FAMILY_DEVANAGARI           []string
	_IDS_SANS_SERIF_FONT_FAMILY_DEVANAGARI      []string
)

func init() {
	fonts.Lock()
	defer fonts.Unlock()
	switch os := runtime.GOOS; os {
	case "windows":
		_IDS_STANDARD_FONT_FAMILY = []string{"Times New Roman"}
		_IDS_FIXED_FONT_FAMILY = []string{"Courier New", "Consolas"}
		_IDS_SERIF_FONT_FAMILY = []string{"Times New Roman"}
		_IDS_SANS_SERIF_FONT_FAMILY = []string{"Arial"}
		_IDS_CURSIVE_FONT_FAMILY = []string{"Comic Sans MS"}
		_IDS_FANTASY_FONT_FAMILY = []string{"Impact"}
		_IDS_MATH_FONT_FAMILY = []string{"Cambria Math"}

		_IDS_FIXED_FONT_FAMILY_ARABIC = []string{"Courier New"}
		_IDS_SANS_SERIF_FONT_FAMILY_ARABIC = []string{"Segoe UI"}

		_IDS_STANDARD_FONT_FAMILY_CYRILLIC = []string{"Times New Roman"}
		_IDS_FIXED_FONT_FAMILY_CYRILLIC = []string{"Courier New"}
		_IDS_SERIF_FONT_FAMILY_CYRILLIC = []string{"Times New Roman"}
		_IDS_SANS_SERIF_FONT_FAMILY_CYRILLIC = []string{"Arial"}

		_IDS_STANDARD_FONT_FAMILY_GREEK = []string{"Times New Roman"}
		_IDS_FIXED_FONT_FAMILY_GREEK = []string{"Courier New"}
		_IDS_SERIF_FONT_FAMILY_GREEK = []string{"Times New Roman"}
		_IDS_SANS_SERIF_FONT_FAMILY_GREEK = []string{"Arial"}

		_IDS_STANDARD_FONT_FAMILY_JAPANESE = []string{"Noto Sans JP", "Noto Sans CJK JP", "Meiryo", "Yu Gothic"}
		_IDS_FIXED_FONT_FAMILY_JAPANESE = []string{"BIZ UDGothic", "MS Gothic"}
		_IDS_SERIF_FONT_FAMILY_JAPANESE = []string{"Noto Serif JP", "Noto Serif CJK JP", "Yu Mincho", "MS PMincho"}
		_IDS_SANS_SERIF_FONT_FAMILY_JAPANESE = []string{"Noto Sans JP", "Noto Sans CJK JP", "Meiryo", "Yu Gothic"}

		_IDS_STANDARD_FONT_FAMILY_KOREAN = []string{"Noto Sans KR", "Noto Sans CJK KR", "Malgun Gothic"}
		_IDS_FIXED_FONT_FAMILY_KOREAN = []string{"Gulimche"}
		_IDS_SERIF_FONT_FAMILY_KOREAN = []string{"Noto Serif KR", "Noto Serif CJK KR", "Batang"}
		_IDS_SANS_SERIF_FONT_FAMILY_KOREAN = []string{"Noto Sans KR", "Noto Sans CJK KR", "Malgun Gothic"}
		_IDS_CURSIVE_FONT_FAMILY_KOREAN = []string{"Gungsuh"}

		_IDS_STANDARD_FONT_FAMILY_SIMPLIFIED_HAN = []string{"Noto Sans SC", "Noto Sans CJK SC", "Microsoft YaHei"}
		_IDS_FIXED_FONT_FAMILY_SIMPLIFIED_HAN = []string{"NSimsun"}                                                  // 新宋体
		_IDS_SERIF_FONT_FAMILY_SIMPLIFIED_HAN = []string{"Noto Serif SC", "Noto Serif CJK SC", "Simsun"}             // 宋体
		_IDS_SANS_SERIF_FONT_FAMILY_SIMPLIFIED_HAN = []string{"Noto Sans SC", "Noto Sans CJK SC", "Microsoft YaHei"} // 微软雅黑
		_IDS_CURSIVE_FONT_FAMILY_SIMPLIFIED_HAN = []string{"KaiTi"}                                                  // 楷体

		_IDS_STANDARD_FONT_FAMILY_TRADITIONAL_HAN = []string{"Noto Sans TC", "Noto Sans CJK TC", "Microsoft JhengHei"}
		_IDS_FIXED_FONT_FAMILY_TRADITIONAL_HAN = []string{"MingLiU"}                                                     // 明体/细明体
		_IDS_SERIF_FONT_FAMILY_TRADITIONAL_HAN = []string{"Noto Serif TC", "Noto Serif CJK TC", "PMingLiU"}              // 新细明体
		_IDS_SANS_SERIF_FONT_FAMILY_TRADITIONAL_HAN = []string{"Noto Sans TC", "Noto Sans CJK TC", "Microsoft JhengHei"} // 微軟正黑體
		_IDS_CURSIVE_FONT_FAMILY_TRADITIONAL_HAN = []string{"DFKai-SB"}                                                  // 标楷体

		_IDS_STANDARD_FONT_FAMILY_DEVANAGARI = []string{"Nirmala UI"}
		_IDS_FIXED_FONT_FAMILY_DEVANAGARI = []string{"Consolas"}
		_IDS_SERIF_FONT_FAMILY_DEVANAGARI = []string{"Nirmala UI"}
		_IDS_SANS_SERIF_FONT_FAMILY_DEVANAGARI = []string{"Nirmala UI"}
	case "linux":
		_IDS_STANDARD_FONT_FAMILY = []string{"Times New Roman"}
		_IDS_FIXED_FONT_FAMILY = []string{"Monospace"}
		_IDS_SERIF_FONT_FAMILY = []string{"Times New Roman"}
		_IDS_SANS_SERIF_FONT_FAMILY = []string{"Arial"}
		_IDS_CURSIVE_FONT_FAMILY = []string{"Comic Sans MS"}
		_IDS_FANTASY_FONT_FAMILY = []string{"Impact"}
		_IDS_MATH_FONT_FAMILY = []string{"Latin Modern Math"}

		_IDS_FIXED_FONT_FAMILY_ARABIC = []string{}
		_IDS_SANS_SERIF_FONT_FAMILY_ARABIC = []string{}

		_IDS_STANDARD_FONT_FAMILY_CYRILLIC = []string{}
		_IDS_FIXED_FONT_FAMILY_CYRILLIC = []string{}
		_IDS_SERIF_FONT_FAMILY_CYRILLIC = []string{}
		_IDS_SANS_SERIF_FONT_FAMILY_CYRILLIC = []string{}

		_IDS_STANDARD_FONT_FAMILY_GREEK = []string{}
		_IDS_FIXED_FONT_FAMILY_GREEK = []string{}
		_IDS_SERIF_FONT_FAMILY_GREEK = []string{}
		_IDS_SANS_SERIF_FONT_FAMILY_GREEK = []string{}

		_IDS_STANDARD_FONT_FAMILY_JAPANESE = []string{"Noto Sans JP", "Noto Sans CJK JP", "Times New Roman"}
		_IDS_FIXED_FONT_FAMILY_JAPANESE = []string{"Noto Sans Mono CJK JP"}
		_IDS_SERIF_FONT_FAMILY_JAPANESE = []string{"Noto Serif JP", "Noto Serif CJK JP", "Times New Roman"}
		_IDS_SANS_SERIF_FONT_FAMILY_JAPANESE = []string{"Noto Sans JP", "Noto Sans CJK JP", "Arial"}

		_IDS_STANDARD_FONT_FAMILY_KOREAN = []string{"Noto Sans KR", "Noto Sans CJK KR", "Times New Roman"}
		_IDS_FIXED_FONT_FAMILY_KOREAN = []string{}
		_IDS_SERIF_FONT_FAMILY_KOREAN = []string{"Noto Serif KR", "Noto Serif CJK KR", "Times New Roman"}
		_IDS_SANS_SERIF_FONT_FAMILY_KOREAN = []string{"Noto Sans KR", "Noto Sans CJK KR", "Arial"}
		_IDS_CURSIVE_FONT_FAMILY_KOREAN = []string{}

		_IDS_STANDARD_FONT_FAMILY_SIMPLIFIED_HAN = []string{"Noto Sans SC", "Noto Sans CJK SC", "Times New Roman"}
		_IDS_FIXED_FONT_FAMILY_SIMPLIFIED_HAN = []string{}
		_IDS_SERIF_FONT_FAMILY_SIMPLIFIED_HAN = []string{"Noto Serif SC", "Noto Serif CJK SC", "Times New Roman"}
		_IDS_SANS_SERIF_FONT_FAMILY_SIMPLIFIED_HAN = []string{"Noto Sans SC", "Noto Sans CJK SC", "Arial"}
		_IDS_CURSIVE_FONT_FAMILY_SIMPLIFIED_HAN = []string{}

		_IDS_STANDARD_FONT_FAMILY_TRADITIONAL_HAN = []string{"Noto Sans TC", "Noto Sans CJK TC", "Times New Roman"}
		_IDS_FIXED_FONT_FAMILY_TRADITIONAL_HAN = []string{}
		_IDS_SERIF_FONT_FAMILY_TRADITIONAL_HAN = []string{"Noto Serif TC", "Noto Serif CJK TC", "Times New Roman"}
		_IDS_SANS_SERIF_FONT_FAMILY_TRADITIONAL_HAN = []string{"Noto Sans TC", "Noto Sans CJK TC", "Arial"}
		_IDS_CURSIVE_FONT_FAMILY_TRADITIONAL_HAN = []string{}

		_IDS_STANDARD_FONT_FAMILY_DEVANAGARI = []string{"Noto Sans Devanagari"}
		_IDS_FIXED_FONT_FAMILY_DEVANAGARI = []string{"Noto Sans Mono"}
		_IDS_SERIF_FONT_FAMILY_DEVANAGARI = []string{"Noto Serif Devanagari"}
		_IDS_SANS_SERIF_FONT_FAMILY_DEVANAGARI = []string{"Noto Sans Devanagari"}
	case "macos":
		_IDS_STANDARD_FONT_FAMILY = []string{"Times"}
		_IDS_FIXED_FONT_FAMILY = []string{"Menlo"}
		_IDS_SERIF_FONT_FAMILY = []string{"Times"}
		_IDS_SANS_SERIF_FONT_FAMILY = []string{"Helvetica"}
		_IDS_CURSIVE_FONT_FAMILY = []string{"Apple Chancery"}
		_IDS_FANTASY_FONT_FAMILY = []string{"Papyrus"}
		_IDS_MATH_FONT_FAMILY = []string{"STIX Two Math"}

		_IDS_FIXED_FONT_FAMILY_ARABIC = []string{}
		_IDS_SANS_SERIF_FONT_FAMILY_ARABIC = []string{}

		_IDS_STANDARD_FONT_FAMILY_CYRILLIC = []string{}
		_IDS_FIXED_FONT_FAMILY_CYRILLIC = []string{}
		_IDS_SERIF_FONT_FAMILY_CYRILLIC = []string{}
		_IDS_SANS_SERIF_FONT_FAMILY_CYRILLIC = []string{}

		_IDS_STANDARD_FONT_FAMILY_GREEK = []string{}
		_IDS_FIXED_FONT_FAMILY_GREEK = []string{}
		_IDS_SERIF_FONT_FAMILY_GREEK = []string{}
		_IDS_SANS_SERIF_FONT_FAMILY_GREEK = []string{}

		_IDS_STANDARD_FONT_FAMILY_JAPANESE = []string{"Hiragino Kaku Gothic ProN"}
		_IDS_FIXED_FONT_FAMILY_JAPANESE = []string{"Osaka", "BIZ UDGothic", "Menlo"}
		_IDS_SERIF_FONT_FAMILY_JAPANESE = []string{"Hiragino Mincho ProN"}
		_IDS_SANS_SERIF_FONT_FAMILY_JAPANESE = []string{"Hiragino Kaku Gothic ProN"}

		_IDS_STANDARD_FONT_FAMILY_KOREAN = []string{"Apple SD Gothic Neo"}
		_IDS_FIXED_FONT_FAMILY_KOREAN = []string{}
		_IDS_SERIF_FONT_FAMILY_KOREAN = []string{"AppleMyungjo"}
		_IDS_SANS_SERIF_FONT_FAMILY_KOREAN = []string{"Apple SD Gothic Neo"}
		_IDS_CURSIVE_FONT_FAMILY_KOREAN = []string{}

		_IDS_STANDARD_FONT_FAMILY_SIMPLIFIED_HAN = []string{"PingFang SC", "STHeiti"}
		_IDS_FIXED_FONT_FAMILY_SIMPLIFIED_HAN = []string{}
		_IDS_SERIF_FONT_FAMILY_SIMPLIFIED_HAN = []string{"Songti SC"}
		_IDS_SANS_SERIF_FONT_FAMILY_SIMPLIFIED_HAN = []string{"PingFang SC", "STHeiti"}
		_IDS_CURSIVE_FONT_FAMILY_SIMPLIFIED_HAN = []string{"Kaiti SC"}

		_IDS_STANDARD_FONT_FAMILY_TRADITIONAL_HAN = []string{"PingFang TC", "Heiti TC"}
		_IDS_FIXED_FONT_FAMILY_TRADITIONAL_HAN = []string{}
		_IDS_SERIF_FONT_FAMILY_TRADITIONAL_HAN = []string{"Songti TC"}
		_IDS_SANS_SERIF_FONT_FAMILY_TRADITIONAL_HAN = []string{"PingFang TC", "Heiti TC"}
		_IDS_CURSIVE_FONT_FAMILY_TRADITIONAL_HAN = []string{"Kaiti TC"}

		_IDS_STANDARD_FONT_FAMILY_DEVANAGARI = []string{"Devanagari MT"}
		_IDS_FIXED_FONT_FAMILY_DEVANAGARI = []string{"Menlo"}
		_IDS_SERIF_FONT_FAMILY_DEVANAGARI = []string{"Devanagari MT"}
		_IDS_SANS_SERIF_FONT_FAMILY_DEVANAGARI = []string{"Kohinoor Devanagari"}
	}
}

func fallback(family string, script language.Script) []string {
	switch family {
	case FontGenericSansSerif:
		switch script {
		case _LANG_SCRIPT_ARABIC:
			return _IDS_SANS_SERIF_FONT_FAMILY_ARABIC
		case _LANG_SCRIPT_CYRILLIC:
			return _IDS_SANS_SERIF_FONT_FAMILY_CYRILLIC
		case _LANG_SCRIPT_GREEK:
			return _IDS_SANS_SERIF_FONT_FAMILY_GREEK
		case _LANG_SCRIPT_JAPANESE:
			return _IDS_SANS_SERIF_FONT_FAMILY_JAPANESE
		case _LANG_SCRIPT_KOREAN:
			return _IDS_SANS_SERIF_FONT_FAMILY_KOREAN
		case _LANG_SCRIPT_SIMPLIFIED_HAN:
			return _IDS_SANS_SERIF_FONT_FAMILY_SIMPLIFIED_HAN
		case _LANG_SCRIPT_TRADITIONAL_HAN:
			return _IDS_SANS_SERIF_FONT_FAMILY_TRADITIONAL_HAN
		case _LANG_SCRIPT_DEVANAGARI:
			return _IDS_SANS_SERIF_FONT_FAMILY_DEVANAGARI
		default:
			return _IDS_SANS_SERIF_FONT_FAMILY
		}
	case FontGenericSerif:
		switch script {
		case _LANG_SCRIPT_CYRILLIC:
			return _IDS_SERIF_FONT_FAMILY_CYRILLIC
		case _LANG_SCRIPT_GREEK:
			return _IDS_SERIF_FONT_FAMILY_GREEK
		case _LANG_SCRIPT_JAPANESE:
			return _IDS_SERIF_FONT_FAMILY_JAPANESE
		case _LANG_SCRIPT_KOREAN:
			return _IDS_SERIF_FONT_FAMILY_KOREAN
		case _LANG_SCRIPT_SIMPLIFIED_HAN:
			return _IDS_SERIF_FONT_FAMILY_SIMPLIFIED_HAN
		case _LANG_SCRIPT_TRADITIONAL_HAN:
			return _IDS_SERIF_FONT_FAMILY_TRADITIONAL_HAN
		case _LANG_SCRIPT_DEVANAGARI:
			return _IDS_SERIF_FONT_FAMILY_DEVANAGARI
		default:
			return _IDS_SERIF_FONT_FAMILY
		}
	case FontGenericMonospace:
		switch script {
		case _LANG_SCRIPT_ARABIC:
			return _IDS_FIXED_FONT_FAMILY_ARABIC
		case _LANG_SCRIPT_CYRILLIC:
			return _IDS_FIXED_FONT_FAMILY_CYRILLIC
		case _LANG_SCRIPT_GREEK:
			return _IDS_FIXED_FONT_FAMILY_GREEK
		case _LANG_SCRIPT_JAPANESE:
			return _IDS_FIXED_FONT_FAMILY_JAPANESE
		case _LANG_SCRIPT_KOREAN:
			return _IDS_FIXED_FONT_FAMILY_KOREAN
		case _LANG_SCRIPT_SIMPLIFIED_HAN:
			return _IDS_FIXED_FONT_FAMILY_SIMPLIFIED_HAN
		case _LANG_SCRIPT_TRADITIONAL_HAN:
			return _IDS_FIXED_FONT_FAMILY_TRADITIONAL_HAN
		case _LANG_SCRIPT_DEVANAGARI:
			return _IDS_FIXED_FONT_FAMILY_DEVANAGARI
		default:
			return _IDS_FIXED_FONT_FAMILY
		}
	case FontGenericCursive:
		switch script {
		case _LANG_SCRIPT_KOREAN:
			return _IDS_CURSIVE_FONT_FAMILY_KOREAN
		case _LANG_SCRIPT_SIMPLIFIED_HAN:
			return _IDS_CURSIVE_FONT_FAMILY_SIMPLIFIED_HAN
		case _LANG_SCRIPT_TRADITIONAL_HAN:
			return _IDS_CURSIVE_FONT_FAMILY_TRADITIONAL_HAN
		default:
			return _IDS_CURSIVE_FONT_FAMILY
		}
	case FontGenericFantasy:
		return _IDS_FANTASY_FONT_FAMILY
	case FontGenericMath:
		return _IDS_MATH_FONT_FAMILY
	default:
		switch script {
		case _LANG_SCRIPT_CYRILLIC:
			return _IDS_STANDARD_FONT_FAMILY_CYRILLIC
		case _LANG_SCRIPT_GREEK:
			return _IDS_STANDARD_FONT_FAMILY_GREEK
		case _LANG_SCRIPT_JAPANESE:
			return _IDS_STANDARD_FONT_FAMILY_JAPANESE
		case _LANG_SCRIPT_KOREAN:
			return _IDS_STANDARD_FONT_FAMILY_KOREAN
		case _LANG_SCRIPT_SIMPLIFIED_HAN:
			return _IDS_STANDARD_FONT_FAMILY_SIMPLIFIED_HAN
		case _LANG_SCRIPT_TRADITIONAL_HAN:
			return _IDS_STANDARD_FONT_FAMILY_TRADITIONAL_HAN
		case _LANG_SCRIPT_DEVANAGARI:
			return _IDS_STANDARD_FONT_FAMILY_DEVANAGARI
		default:
			return _IDS_STANDARD_FONT_FAMILY
		}
	}
}

func loadFont(font typeface) (*sfnt.Font, error) {
	if font.font != nil {
		fmt.Println("cache hit")
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
		if s == _LANG_SCRIPT_ZZZZ {
			switch {
			case b == _LANG_BASE_ZH, r == _LANG_REGION_CN:
				script = _LANG_SCRIPT_SIMPLIFIED_HAN
			case r == _LANG_REGION_TW, r == _LANG_REGION_HK, r == _LANG_REGION_MO,
				r == _LANG_REGION_SG, r == _LANG_REGION_MY:
				script = _LANG_SCRIPT_TRADITIONAL_HAN
			default:
				script = _LANG_SCRIPT_LATIN
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
		script = _LANG_SCRIPT_LATIN
	case unicode.Is(unicode.Han, r):
		_, s, r := ctx.lang.Raw()
		switch {
		case s == _LANG_SCRIPT_SIMPLIFIED_HAN, r == _LANG_REGION_CN:
			script = _LANG_SCRIPT_SIMPLIFIED_HAN
		case s == _LANG_SCRIPT_TRADITIONAL_HAN,
			r == _LANG_REGION_TW, r == _LANG_REGION_HK:
			script = _LANG_SCRIPT_TRADITIONAL_HAN
		default:
			script = _LANG_SCRIPT_TRADITIONAL_HAN
		}
	case unicode.Is(unicode.Hangul, r):
		script = _LANG_SCRIPT_KOREAN
	case unicode.Is(unicode.Hiragana, r) || unicode.Is(unicode.Katakana, r):
		script = _LANG_SCRIPT_JAPANESE
	default:
		script = _LANG_SCRIPT_LATIN
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
