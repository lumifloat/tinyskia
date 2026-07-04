// Copyright 2016 Michael Fogleman
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package examples

import (
	"image/color"
	"os"
	"strings"
	"testing"

	gg "github.com/lumifloat/tinyskia"
)

const WRAP_TEXT = "Call me Ishmael. Some years ago—never mind how long precisely—having little or no money in my purse, and nothing particular to interest me on shore, I thought I would sail about a little and see the watery part of the world. It is a way I have of driving off the spleen and regulating the circulation. Whenever I find myself growing grim about the mouth; whenever it is a damp, drizzly November in my soul; whenever I find myself involuntarily pausing before coffin warehouses, and bringing up the rear of every funeral I meet; and especially whenever my hypos get such an upper hand of me, that it requires a strong moral principle to prevent me from deliberately stepping into the street, and methodically knocking people's hats off—then, I account it high time to get to sea as soon as I can. This is my substitute for pistol and ball. With a philosophical flourish Cato throws himself upon his sword; I quietly take to the ship. There is nothing surprising in this. If they but knew it, almost all men in their degree, some time or other, cherish very nearly the same feelings towards the ocean with me."

// drawWrapLine draws a line from (x1, y1) to (x2, y2)
func drawWrapLine(dc *gg.Context, x1, y1, x2, y2 float64) {
	dc.MoveTo(x1, y1)
	dc.LineTo(x2, y2)
}

// drawWrapRectangle draws a rectangle at (x, y) with width w and height h
func drawWrapRectangle(dc *gg.Context, x, y, w, h float64) {
	dc.MoveTo(x, y)
	dc.LineTo(x+w, y)
	dc.LineTo(x+w, y+h)
	dc.LineTo(x, y+h)
	dc.ClosePath()
}

// wrapText wraps text to fit within maxWidth and returns lines
func wrapText(dc *gg.Context, text string, maxWidth float64) []string {
	var result []string
	for _, line := range strings.Split(text, "\n") {
		fields := splitOnSpace(line)
		if len(fields)%2 == 1 {
			fields = append(fields, "")
		}

		x := ""
		for i := 0; i < len(fields); i += 2 {
			metrics := dc.MeasureText(x + fields[i])
			if metrics.Width > maxWidth {
				if x == "" {
					result = append(result, fields[i])
					x = ""
					continue
				} else {
					result = append(result, x)
					x = ""
				}
			}
			x += fields[i] + fields[i+1]
		}
		if x != "" {
			result = append(result, x)
		}
	}

	// Trim whitespace from each line
	for i, line := range result {
		result[i] = strings.TrimSpace(line)
	}

	return result
}

// splitOnSpace splits a string into alternating word and space segments
func splitOnSpace(s string) []string {
	var result []string
	current := ""
	inSpace := false

	for _, c := range s {
		if c == ' ' || c == '\t' {
			if !inSpace {
				if current != "" {
					result = append(result, current)
				}
				current = ""
				inSpace = true
			}
			current += string(c)
		} else {
			if inSpace {
				if current != "" {
					result = append(result, current)
				}
				current = ""
				inSpace = false
			}
			current += string(c)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}

// drawWrappedText draws wrapped text with alignment
func drawWrappedText(dc *gg.Context, text string, x, y, anchorX, anchorY, maxWidth, lineSpacing float64, align gg.TextAlign) {
	lines := wrapText(dc, text, maxWidth)
	if len(lines) == 0 {
		return
	}

	// Get font metrics
	metrics := dc.MeasureText("Ag")
	fontHeight := metrics.FontBoundingBoxAscent + metrics.FontBoundingBoxDescent

	// Calculate total height (sync with DrawStringWrapped formula)
	totalHeight := float64(len(lines)) * fontHeight * lineSpacing
	totalHeight -= (lineSpacing - 1) * fontHeight

	// Apply anchor point to position
	x -= anchorX * maxWidth
	y -= anchorY * totalHeight

	// Adjust x based on text alignment
	var ax float64
	switch align {
	case gg.TextAlignLeft:
		ax = 0
	case gg.TextAlignCenter:
		ax = 0.5
		x += maxWidth / 2
	case gg.TextAlignRight:
		ax = 1
		x += maxWidth
	default:
		ax = 0
	}

	// Save and restore text align
	prevAlign := dc.GetTextAlign()
	dc.SetTextAlign(gg.TextAlignLeft)

	// Draw each line
	for _, line := range lines {
		// Draw line with anchor
		lineMetrics := dc.MeasureText(line)
		lineWidth := lineMetrics.Width
		lineX := x - ax*lineWidth
		lineY := y + metrics.FontBoundingBoxAscent

		dc.FillText(line, lineX, lineY)

		// Move to next line
		y += fontHeight * lineSpacing
	}

	dc.SetTextAlign(prevAlign)
}

func TestGGWrap(t *testing.T) {
	const W = 1024
	const H = 1024
	const P = 16

	c := gg.NewCanvas(W, H)
	dc := c.GetContext()

	// Set white background
	dc.SetFillStyleSolidColor(color.RGBA{255, 255, 255, 255})
	dc.FillRect(0, 0, W, H)

	// Draw center lines
	dc.BeginPath()
	drawWrapLine(dc, W/2, 0, W/2, H)
	drawWrapLine(dc, 0, H/2, W, H/2)
	drawWrapRectangle(dc, P, P, W-P-P, H-P-P)
	dc.SetStrokeStyleSolidColor(color.RGBA{0, 0, 255, 64}) // 0.25 alpha
	dc.SetLineWidth(3)
	dc.Stroke()

	dc.SetFillStyleSolidColor(color.RGBA{0, 0, 0, 255})

	// Draw corner labels
	dc.SetFont(gg.FontAttr{Family: []string{"arial"}, Weight: gg.FontWeightBold, Size: 18})
	drawWrappedText(dc, "UPPER LEFT", P, P, 0, 0, 0, 1, gg.TextAlignLeft)
	drawWrappedText(dc, "UPPER RIGHT", W-P, P, 1, 0, 0, 1, gg.TextAlignRight)
	drawWrappedText(dc, "BOTTOM LEFT", P, H-P, 0, 1, 0, 1, gg.TextAlignLeft)
	drawWrappedText(dc, "BOTTOM RIGHT", W-P, H-P, 1, 1, 0, 1, gg.TextAlignRight)
	drawWrappedText(dc, "UPPER MIDDLE", W/2, P, 0.5, 0, 0, 1, gg.TextAlignCenter)
	drawWrappedText(dc, "LOWER MIDDLE", W/2, H-P, 0.5, 1, 0, 1, gg.TextAlignCenter)
	drawWrappedText(dc, "LEFT MIDDLE", P, H/2, 0, 0.5, 0, 1, gg.TextAlignLeft)
	drawWrappedText(dc, "RIGHT MIDDLE", W-P, H/2, 1, 0.5, 0, 1, gg.TextAlignRight)

	// Load regular font for body text
	dc.SetFont(gg.FontAttr{Family: []string{"arial"}, Size: 12})
	drawWrappedText(dc, WRAP_TEXT, W/2-P, H/2-P, 1, 1, W/3, 1, gg.TextAlignLeft)
	drawWrappedText(dc, WRAP_TEXT, W/2+P, H/2-P, 0, 1, W/3, 1.2, gg.TextAlignLeft)
	drawWrappedText(dc, WRAP_TEXT, W/2-P, H/2+P, 1, 0, W/3, 1.4, gg.TextAlignLeft)
	drawWrappedText(dc, WRAP_TEXT, W/2+P, H/2+P, 0, 0, W/3, 1.6, gg.TextAlignLeft)

	fi, err := os.Create("gg_out.png")
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	c.WritePNG(fi, nil)
}
