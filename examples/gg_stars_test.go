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
	"testing"

	gg "github.com/lumifloat/tinyskia"
)

func SmallPolygon(n int) [][2]float64 {
	result := make([][2]float64, n)
	for i := 0; i < n; i++ {
		a := float64(i)*2*math.Pi/float64(n) - math.Pi/2
		result[i] = [2]float64{math.Cos(a), math.Sin(a)}
	}
	return result
}

func TestGGStars(t *testing.T) {
	const W = 1200
	const H = 120
	const S = 100
	dc := gg.NewContext(W, H)
	dc.SetFillStyleSolidColor(color.RGBA{255, 255, 255, 255})
	dc.FillRect(0, 0, float64(W), float64(H))

	n := 5
	points := SmallPolygon(n)

	source := rand.NewSource(1)
	rnd := rand.New(source)

	for x := S / 2; x < W; x += S {
		dc.Save()
		s := rnd.Float64()*S/4 + S/4
		dc.Rotate(rnd.Float64() * 2 * math.Pi)
		dc.Translate(float64(x), H/2)
		dc.BeginPath()
		for i := 0; i < n; i++ {
			index := (i * 2) % n
			p := points[index]
			dc.LineTo(p[0]*s, p[1]*s)
		}
		dc.ClosePath()
		dc.SetStrokeStyleSolidColor(color.RGBA{255, 204, 0, 255}) // #FFCC00
		dc.SetLineWidth(10)
		dc.Stroke()

		dc.SetFillStyleSolidColor(color.RGBA{255, 228, 58, 255}) // #FFE43A
		dc.Fill()
		dc.Restore()
	}

	dc.SavePNG("gg_out.png")
}
