// Copyright 2016 Michael Fogleman
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package examples

import (
	"image/color"
	"math"
	"math/rand"
	"os"
	"testing"

	gg "github.com/lumifloat/tinyskia"
)

type ScatterPoint struct {
	X, Y float64
}

func createScatterPoints(n int, size float64) []ScatterPoint {
	points := make([]ScatterPoint, n)
	for i := 0; i < n; i++ {
		x := 0.5 + rand.NormFloat64()*0.1
		y := x + rand.NormFloat64()*0.1
		points[i] = ScatterPoint{x * size, y * size}
	}
	return points
}

// drawCircle draws a circle centered at (x, y) with radius r
func drawScatterCircle(dc *gg.Context, x, y, r float64) {
	dc.BeginPath()
	dc.Arc(x, y, r, 0, 2*math.Pi, false)
	dc.ClosePath()
}

func TestGGScatter(t *testing.T) {
	const S = 1024
	const P = 64

	c := gg.NewCanvas(S, S)
	dc := c.GetContext()

	// Set white background
	dc.SetFillStyleSolidColor(color.RGBA{255, 255, 255, 255})
	dc.FillRect(0, 0, S, S)

	size := float64(S - P*2)
	points := createScatterPoints(1000, size)

	// Apply transformations
	//dc.Scale(S-P*2, S-P*2)
	dc.Translate(P, P)

	dc.BeginPath()
	// Draw minor grid
	for i := 1; i <= 10; i++ {
		x := float64(i) / 10
		dc.MoveTo(x*size, 0)
		dc.LineTo(x*size, size)
		dc.MoveTo(0, x*size)
		dc.LineTo(size, x*size)
	}
	dc.SetStrokeStyleSolidColor(color.RGBA{0, 0, 0, 64}) // 0.25 alpha = 64/255
	dc.SetLineWidth(1)
	dc.Stroke()

	// Draw axes
	dc.BeginPath()
	dc.MoveTo(0, size)
	dc.LineTo(size, size)
	dc.MoveTo(0, 0)
	dc.LineTo(0, size)
	dc.SetStrokeStyleSolidColor(color.RGBA{0, 0, 0, 255})
	dc.SetLineWidth(4)
	dc.Stroke()

	// Draw points
	dc.SetFillStyleSolidColor(color.RGBA{0, 0, 255, 128}) // 0.5 alpha = 128/255
	for _, p := range points {
		drawScatterCircle(dc, size-p.X, p.Y, 3.0/S*size)
		dc.Fill()
	}

	// Reset transform for text
	dc.ResetTransform()

	dc.SetFillStyleSolidColor(color.RGBA{0, 0, 0, 255})

	dc.SetFont(gg.FontAttr{Family: []string{"arial"}, Weight: gg.FontWeightBold, Size: 24})
	text := "Chart Title"
	metrics := dc.MeasureText(text)
	textWidth := metrics.Width
	textHeight := metrics.FontBoundingBoxAscent - metrics.FontBoundingBoxDescent

	x := S/2 - textWidth/2
	y := (P / 2) + textHeight/2
	dc.FillText(text, x, y)

	// Draw X axis title
	dc.SetFont(gg.FontAttr{Family: []string{"arial"}, Weight: gg.FontWeightNormal, Size: 18})
	text = "X Axis Title"
	metrics = dc.MeasureText(text)
	textWidth = metrics.Width
	textHeight = metrics.FontBoundingBoxAscent - metrics.FontBoundingBoxDescent

	x = S/2 - textWidth/2
	y = (S - P/2) + textHeight/2
	dc.FillText(text, x, y)

	fi, err := os.Create("gg_out.png")
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	c.WritePNG(fi, nil)
}
