// Copyright 2006 The Android Open Source Project
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package tinyskia

import (
	"errors"
	"os"
	"sync"

	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/sfnt"
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

type FontAttr struct {
	Family string
	Weight string
	Style  string
	Size   int
}

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
func (dc *Context) SetFont(font FontAttr) {
	if font.Family == "" {
		return
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

type FontFace struct {
	Family string
	Weight string
	Style  string
}

type fontdata struct {
	font   *sfnt.Font
	family string
	weight string
	style  string
}

type registry struct {
	sync.Mutex
	buf   sfnt.Buffer
	fonts []fontdata
}

var fonts = &registry{}

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

	collection, err := sfnt.ParseCollection(data)
	if err != nil {
		return err
	}
	ttf, err := collection.Font(0)
	if err != nil {
		return err
	}

	fonts.Lock()
	defer fonts.Unlock()

	fonts.fonts = append(fonts.fonts, fontdata{
		font:   ttf,
		family: fontFace.Family,
		weight: fontFace.Weight,
		style:  fontFace.Style,
	})

	return nil
}

func (r *registry) match(family string, weight string, style string) (*sfnt.Font, error) {
	if family == "" {
		return nil, errors.New("family name cannot be empty")
	}

	r.Lock()
	defer r.Unlock()

	if len(r.fonts) == 0 {
		return sfnt.Parse(goregular.TTF)
	}

	var current *sfnt.Font
	var best = -1

	for _, f := range fonts.fonts {
		if f.family != family {
			continue
		}

		score := 0
		if style == "" || f.style == style {
			score += 2
		}
		if weight == "" || f.weight == weight {
			score += 1
		}

		if score == 3 {
			return f.font, nil
		}

		if score > best {
			best = score
			current = f.font
		}
	}

	if best == -1 {
		return sfnt.Parse(goregular.TTF)
	}

	return current, nil
}
