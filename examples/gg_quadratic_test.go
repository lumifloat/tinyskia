// Copyright 2016 Michael Fogleman
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package examples

import (
	"image/color"
	"math"
	"os"
	"runtime"
	"testing"

	gg "github.com/lumifloat/tinyskia"
)

func TestGGQuadratic(t *testing.T) {
	const S = 1000
	c := gg.NewCanvas(S, S)
	dc := c.GetContext()

	switch runtime.GOOS {
	case "windows":
		gg.RegisterFont("C:/Windows/Fonts/arial.ttf", gg.FontFace{Family: "arial"})
	case "darwin":
		gg.RegisterFont("/Library/Fonts/Arial.ttf", gg.FontFace{Family: "arial"})
	default:
		// pass
	}

	dc.SetFillStyleSolidColor(color.White)
	dc.FillRect(0, 0, float64(S), float64(S))

	dc.Translate(float64(S/2), float64(S/2))

	var x0, y0, x1, y1, x2, y2, x3, y3, x4, y4 float64
	x0, y0 = -400, 0
	x1, y1 = -200, -400
	x2, y2 = 0, 0
	x3, y3 = 200, 400
	x4, y4 = 400, 0

	dc.BeginPath()
	dc.MoveTo(x0, y0)
	dc.LineTo(x1, y1)
	dc.LineTo(x2, y2)
	dc.LineTo(x3, y3)
	dc.LineTo(x4, y4)
	dc.SetStrokeStyleSolidColor(color.RGBA{255, 45, 0, 255}) // FF2D00
	dc.SetLineWidth(8)
	dc.Stroke()

	dc.BeginPath()
	dc.MoveTo(x0, y0)
	dc.QuadraticCurveTo(x1, y1, x2, y2)
	dc.QuadraticCurveTo(x3, y3, x4, y4)
	dc.SetFillStyleSolidColor(color.RGBA{62, 96, 111, 255}) // 3E606F
	dc.Fill()
	dc.SetLineWidth(16)
	dc.SetStrokeStyleSolidColor(color.Black)
	dc.Stroke()

	dc.SetFillStyleSolidColor(color.White)
	dc.SetLineWidth(4)
	dc.SetStrokeStyleSolidColor(color.Black)

	dc.BeginPath()
	dc.Arc(x0, y0, 20, 0, 2*math.Pi, false)
	dc.Fill()
	dc.Stroke()

	dc.BeginPath()
	dc.Arc(x1, y1, 20, 0, 2*math.Pi, false)
	dc.Fill()
	dc.Stroke()

	dc.BeginPath()
	dc.Arc(x2, y2, 20, 0, 2*math.Pi, false)
	dc.Fill()
	dc.Stroke()

	dc.BeginPath()
	dc.Arc(x3, y3, 20, 0, 2*math.Pi, false)
	dc.Fill()
	dc.Stroke()

	dc.BeginPath()
	dc.Arc(x4, y4, 20, 0, 2*math.Pi, false)
	dc.Fill()
	dc.Stroke()

	dc.SetFillStyleSolidColor(color.Black)
	dc.SetFont(gg.FontAttr{Family: "arial", Size: 200})

	var text string
	var metrics gg.TextMetrics

	text = "g"
	metrics = dc.MeasureText(text)
	dc.FillText(text,
		-200-metrics.Width/2,
		200+(metrics.FontBoundingBoxAscent-metrics.FontBoundingBoxDescent)/2)

	text = "G"
	metrics = dc.MeasureText(text)
	dc.FillText(text,
		200-metrics.Width/2,
		-200+(metrics.FontBoundingBoxAscent-metrics.FontBoundingBoxDescent)/2)

	fi, err := os.Create("gg_out.png")
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	c.WritePNG(fi, nil)
}
