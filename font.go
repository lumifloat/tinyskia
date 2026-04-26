// Copyright 2016 Michael Fogleman
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package tinyskia

import (
	"fmt"
	"io/ioutil"

	"golang.org/x/image/font"
	"golang.org/x/image/font/sfnt"
)

// Font holds font-related information
type Font struct {
	ttf     *sfnt.Font  // Keep original SFNT font for glyph outlines
	buf     sfnt.Buffer // Buffer for glyph operations
	size    float64
	dpi     float64 // dots per inch
	hinting font.Hinting
}

// NewFont creates a new Font with the specified SFNT font, size in points, and DPI.
func NewFont(ttf *sfnt.Font, size, dpi float64) *Font {
	return &Font{
		ttf:     ttf,
		size:    size,
		dpi:     dpi,
		hinting: font.HintingFull,
	}
}

// Size returns the font size in points
func (f *Font) Size() float64 {
	return f.size
}

// DPI returns the dots per inch
func (f *Font) DPI() float64 {
	return f.dpi
}

// LoadFontFace loads a TrueType font from a file.
func (dc *Context) LoadFontFace(path string, points float64) error {
	return dc.LoadFontFaceWithDPI(path, points, 72.0)
}

// LoadFontFaceWithDPI loads a TrueType font from a file with custom DPI.
// Supports both single fonts (TTF/OTF) and font collections (TTC/OTC).
// For collections, it loads the first font (index 0).
func (dc *Context) LoadFontFaceWithDPI(path string, points, dpi float64) error {
	// Load the font file
	fontBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	// Try to parse as a collection first (TTC/OTC)
	collection, err := sfnt.ParseCollection(fontBytes)
	if err == nil && collection.NumFonts() > 0 {
		// Successfully parsed as collection, use the first font
		ttf, err := collection.Font(0)
		if err != nil {
			return fmt.Errorf("failed to get first font from collection: %v", err)
		}
		// Create font object with raw data
		dc.font = NewFont(ttf, points, dpi)
		return nil
	}

	// Not a collection or failed to parse as collection, try as single font
	ttf, err := sfnt.Parse(fontBytes)
	if err != nil {
		return err
	}

	// Create font object with raw data
	dc.font = NewFont(ttf, points, dpi)
	return nil
}

// LoadFontFaceFromData loads a TrueType font from byte data.
func (dc *Context) LoadFontFaceFromData(data []byte, points float64) error {
	return dc.LoadFontFaceFromDataWithDPI(data, points, 72.0)
}

// LoadFontFaceFromDataWithDPI loads a TrueType font from byte data with custom DPI.
// Supports both single fonts (TTF/OTF) and font collections (TTC/OTC).
// For collections, it loads the first font (index 0).
func (dc *Context) LoadFontFaceFromDataWithDPI(data []byte, points, dpi float64) error {
	// Try to parse as a collection first (TTC/OTC)
	collection, err := sfnt.ParseCollection(data)
	if err == nil && collection.NumFonts() > 0 {
		// Successfully parsed as collection, use the first font
		ttf, err := collection.Font(0)
		if err != nil {
			return fmt.Errorf("failed to get first font from collection: %v", err)
		}
		dc.font = NewFont(ttf, points, dpi)
		return nil
	}

	// Not a collection or failed to parse as collection, try as single font
	ttf, err := sfnt.Parse(data)
	if err != nil {
		return err
	}

	dc.font = NewFont(ttf, points, dpi)
	return nil
}

// SetFontFace sets the current font face.
// Accepts either *Font (for vector rendering with outlines) or font.Face (for bitmap rendering).
func (dc *Context) SetFontFace(f *Font) {
	dc.font = f
}
