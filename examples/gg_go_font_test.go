// Copyright 2016 Michael Fogleman
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package examples

import (
	"bytes"
	"image/color"
	"os"
	"testing"

	gg "github.com/lumifloat/tinyskia"
	"golang.org/x/image/font/gofont/goregular"
)

func TestGoFont(t *testing.T) {
	const S = 1024
	c := gg.NewCanvas(S, S)
	dc := c.GetContext()
	dc.SetFillStyleSolidColor(color.RGBA{255, 255, 255, 255})
	dc.FillRect(0, 0, float64(S), float64(S))

	ttf := bytes.NewReader(goregular.TTF)

	if err := gg.RegisterFontWithResource(ttf, "goregular.ttf", gg.FontFace{Family: "go"}); err != nil {
		t.Errorf("Font not available: %v", err)
		return
	}
	dc.SetFont(gg.FontAttr{Family: []string{"go"}, Size: 48})
	dc.SetTextAlign(gg.CanvasTextAlignCenter)
	text := "Hello, world!"

	x := S / 2
	y := S / 2

	dc.SetFillStyleSolidColor(color.Black)
	dc.FillText(text, float64(x), float64(y))

	fi, err := os.Create("gg_out.png")
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	c.WritePNG(fi, nil)
}
