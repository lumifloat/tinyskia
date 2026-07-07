package text

import (
	"github.com/go-text/typesetting/font"
	"github.com/go-text/typesetting/font/opentype"
	"github.com/go-text/typesetting/fontscan"
	"github.com/go-text/typesetting/shaping"
	"github.com/lumifloat/tinyskia/internal/path"
	"golang.org/x/image/math/fixed"
)

// TODO ADD ITALIC AND BOLD
func Outline(outline font.GlyphOutline, scale float32, x, y float32) *path.Path {
	var pb = path.NewPathBuilder()
	var hasPath = false
	for _, s := range outline.Segments {
		switch s.Op {
		case opentype.SegmentOpMoveTo:
			if hasPath {
				pb.Close()
			}
			pb.MoveTo(
				s.Args[0].X*scale+x,
				-s.Args[0].Y*scale+y,
			)
			hasPath = true
		case opentype.SegmentOpLineTo:
			pb.LineTo(
				s.Args[0].X*scale+x,
				-s.Args[0].Y*scale+y,
			)
		case opentype.SegmentOpQuadTo:
			pb.QuadTo(
				s.Args[0].X*scale+x,
				-s.Args[0].Y*scale+y,
				s.Args[1].X*scale+x,
				-s.Args[1].Y*scale+y,
			)
		case opentype.SegmentOpCubeTo:
			pb.CubicTo(
				s.Args[0].X*scale+x,
				-s.Args[0].Y*scale+y,
				s.Args[1].X*scale+x,
				-s.Args[1].Y*scale+y,
				s.Args[2].X*scale+x,
				-s.Args[2].Y*scale+y,
			)
		}
	}
	if hasPath {
		pb.Close()
	}
	return pb.Finish()
}

func Shape(input shaping.Input, query fontscan.Query) (shapes []shaping.Output, width fixed.Int26_6) {
	var segmenter shaping.Segmenter

	FontLock.Lock()
	outputs := segmenter.Split(input, FontMap)
	FontLock.Unlock()

	shaper := shaping.HarfbuzzShaper{}
	for _, output := range outputs {
		shape := shaper.Shape(output)
		shapes = append(shapes, shape)
		width += shape.Advance
	}
	return shapes, width
}
