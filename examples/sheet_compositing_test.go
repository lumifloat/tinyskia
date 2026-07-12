package examples

import (
	"image/color"
	"math"
	"os"
	"testing"

	"github.com/lumifloat/tinyskia"
)

func TestSheetCompositing(t *testing.T) {
	w := 300
	h := 300

	compositing := []struct {
		name string
		op   tinyskia.CompositeOperation
	}{
		{"source-over", tinyskia.CompositeOperationSourceOver},
		{"source-in", tinyskia.CompositeOperationSourceIn},
		{"source-out", tinyskia.CompositeOperationSourceOut},
		{"source-atop", tinyskia.CompositeOperationSourceAtop},
		{"destination-over", tinyskia.CompositeOperationDestinationOver},
		{"destination-in", tinyskia.CompositeOperationDestinationIn},
		{"destination-out", tinyskia.CompositeOperationDestinationOut},
		{"destination-atop", tinyskia.CompositeOperationDestinationAtop},
		{"lighter", tinyskia.CompositeOperationLighter},
		{"copy", tinyskia.CompositeOperationCopy},
		{"xor", tinyskia.CompositeOperationXor},
		{"multiply", tinyskia.CompositeOperationMultiply},
		{"screen", tinyskia.CompositeOperationScreen},
		{"overlay", tinyskia.CompositeOperationOverlay},
		{"darken", tinyskia.CompositeOperationDarken},
		{"lighten", tinyskia.CompositeOperationLighten},
		{"color-dodge", tinyskia.CompositeOperationColorDodge},
		{"color-burn", tinyskia.CompositeOperationColorBurn},
		{"hard-light", tinyskia.CompositeOperationHardLight},
		{"soft-light", tinyskia.CompositeOperationSoftLight},
		{"difference", tinyskia.CompositeOperationDifference},
		{"exclusion", tinyskia.CompositeOperationExclusion},
		{"hue", tinyskia.CompositeOperationHue},
		{"saturation", tinyskia.CompositeOperationSaturation},
		{"color", tinyskia.CompositeOperationColor},
		{"luminosity", tinyskia.CompositeOperationLuminosity},
	}

	len := len(compositing)

	cols := 4
	rows := (len + cols - 1) / cols
	width := w * cols
	height := h * rows

	canvas := tinyskia.NewCanvas(width, height)
	ctx := canvas.GetContext()

	for i, comp := range compositing {
		row := i / cols
		col := i % cols

		scanvas := tinyskia.NewCanvas(w, h)
		sctx := scanvas.GetContext()
		sctx.SetAntiAlias(false)
		sctx.SetForceHQPipeline(false)

		sctx.Translate(float64(w)/2, float64(h)/2)

		sctx.SetFillStyleSolidColor(color.RGBA{255, 0, 0, 255})
		sctx.BeginPath()
		sctx.Arc(-40, 20, 80, 0, 2*math.Pi)
		sctx.ClosePath()
		sctx.Fill()

		sctx.SetGlobalCompositeOperation(comp.op)

		sctx.SetFillStyleSolidColor(color.RGBA{255, 165, 0, 255})
		sctx.BeginPath()
		sctx.Arc(40, 20, 80, 0, 2*math.Pi)
		sctx.ClosePath()
		sctx.Fill()

		ctx.DrawImage(scanvas.Image(), float64(col*w), float64(row*h))
	}

	outputPath := "sheet_out.png"
	fi, err := os.Create(outputPath)
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}
	defer fi.Close()
	if err := canvas.WritePNG(fi, nil); err != nil {
		t.Fatalf("Failed to save PNG: %v", err)
	}
}
