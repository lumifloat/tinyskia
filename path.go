// Copyright 2016 Michael Fogleman
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package tinyskia

import (
	"math"

	"github.com/lumifloat/tinyskia/internal/path"
)

type Path2D struct {
	builder *path.PathBuilder
}

func NewPath2D() *Path2D {
	return &Path2D{builder: path.NewPathBuilder()}
}

// AddPath adds to the path the path given by the argument.
func (p *Path2D) AddPath(p0 *Path2D) {
	pp := p0.builder.Finish()
	if pp != nil {
		p.builder.PushPath(pp)
	}
}

// AddPathWithTransform adds to the path the path given by the argument, transformed by the given transform.
func (p *Path2D) AddPathWithTransform(p0 *Path2D, transform *matrix) {
	pp := p0.builder.Finish()
	if pp != nil {
		p.builder.PushPathWithTransform(pp, transform.transform)
	}
}

// MoveTo creates a new subpath with the given point.
func (p *Path2D) MoveTo(x, y float64) {
	p1 := path.Point{X: float32(x), Y: float32(y)}
	p.builder.MoveTo(p1.X, p1.Y)
}

// ClosePath Marks the current subpath as closed,
// and starts a new subpath with a point the same as the start and end of the newly closed subpath.
func (p *Path2D) ClosePath() {
	p.builder.Close()
}

// LineTo adds the given point to the current subpath, connected to the previous one by a straight line.
func (p *Path2D) LineTo(x, y float64) {
	_, hasCurrent := p.builder.LastPoint()
	if !hasCurrent {
		p.MoveTo(x, y)
		return
	}
	p1 := path.Point{X: float32(x), Y: float32(y)}
	p.builder.LineTo(p1.X, p1.Y)
}

// QuadraticCurveTo adds the given point to the current subpath, connected to the previous one
// by a quadratic Bézier curve with the given control point.
func (p *Path2D) QuadraticCurveTo(x1, y1, x2, y2 float64) {
	_, hasCurrent := p.builder.LastPoint()
	if !hasCurrent {
		p.MoveTo(x1, y1)
	}
	p1 := path.Point{X: float32(x1), Y: float32(y1)}
	p2 := path.Point{X: float32(x2), Y: float32(y2)}
	p.builder.QuadTo(p1.X, p1.Y, p2.X, p2.Y)
}

// BezierCurveTo adds the given point to the current subpath, connected to the previous one
// by a cubic Bézier curve with the given control points.
func (p *Path2D) BezierCurveTo(cp1x, cp1y, cp2x, cp2y, x, y float64) {
	_, hasCurrent := p.builder.LastPoint()
	if !hasCurrent {
		p.MoveTo(cp1x, cp1y)
	}
	p1 := path.Point{X: float32(cp1x), Y: float32(cp1y)}
	p2 := path.Point{X: float32(cp2x), Y: float32(cp2y)}
	p3 := path.Point{X: float32(x), Y: float32(y)}
	p.builder.CubicTo(p1.X, p1.Y, p2.X, p2.Y, p3.X, p3.Y)
}

// ArcTo adds an arc with the given control points and radius to the current subpath,
// connected to the previous point by a straight line.
func (p *Path2D) ArcTo(x1, y1, x2, y2, radius float64) {
	p.builder.ArcTo(
		float32(x1), float32(y1),
		float32(x2), float32(y2),
		float32(radius),
	)
}

// Rect adds a new closed subpath to the path, representing the given rectangle.
func (p *Path2D) Rect(x, y, w, h float64) {
	if math.IsInf(x, 0) || math.IsNaN(x) ||
		math.IsInf(y, 0) || math.IsNaN(y) ||
		math.IsInf(w, 0) || math.IsNaN(w) ||
		math.IsInf(h, 0) || math.IsNaN(h) {
		return
	}

	p.MoveTo(x, y)
	p.LineTo(x+w, y)
	p.LineTo(x+w, y+h)
	p.LineTo(x, y+h)
	p.ClosePath()
}

// RoundRect adds a new closed subpath to the path representing the given rounded rectangle.
func (p *Path2D) RoundRect(x, y, w, h float64, radii []float64) {
	if math.IsInf(x, 0) || math.IsNaN(x) ||
		math.IsInf(y, 0) || math.IsNaN(y) ||
		math.IsInf(w, 0) || math.IsNaN(w) ||
		math.IsInf(h, 0) || math.IsNaN(h) {
		return
	}

	if radii == nil {
		p.Rect(x, y, w, h)
		return
	}

	if (len(radii) == 0) || len(radii) > 4 {
		return
	}

	var rx, ry [4]float64

	switch len(radii) {
	case 1:
		rx = [4]float64{radii[0], radii[0], radii[0], radii[0]}
		ry = [4]float64{radii[0], radii[0], radii[0], radii[0]}
	case 2:
		rx = [4]float64{radii[0], radii[1], radii[0], radii[1]}
		ry = [4]float64{radii[0], radii[1], radii[0], radii[1]}
	case 3:
		rx = [4]float64{radii[0], radii[1], radii[2], radii[1]}
		ry = [4]float64{radii[0], radii[1], radii[2], radii[1]}
	case 4:
		rx = [4]float64{radii[0], radii[1], radii[2], radii[3]}
		ry = [4]float64{radii[0], radii[1], radii[2], radii[3]}
	default:
		// unreachable
	}

	clockwise := true
	if w < 0 {
		clockwise = !clockwise
		w = -w
		x -= w
		rx[0], rx[1] = rx[1], rx[0]
		rx[2], rx[3] = rx[3], rx[2]
		ry[0], ry[1] = ry[1], ry[0]
		ry[2], ry[3] = ry[3], ry[2]
	}
	if h < 0 {
		clockwise = !clockwise
		h = -h
		y -= h
		rx[0], rx[3] = rx[3], rx[0]
		rx[1], rx[2] = rx[2], rx[1]
		ry[0], ry[3] = ry[3], ry[0]
		ry[1], ry[2] = ry[2], ry[1]
	}

	top := rx[0] + rx[1]
	right := ry[1] + ry[2]
	bottom := rx[2] + rx[3]
	left := ry[0] + ry[3]

	scale := 1.0
	if top > 0 {
		scale = math.Min(scale, w/top)
	}
	if right > 0 {
		scale = math.Min(scale, h/right)
	}
	if bottom > 0 {
		scale = math.Min(scale, w/bottom)
	}
	if left > 0 {
		scale = math.Min(scale, h/left)
	}

	if scale < 1 {
		for i := 0; i < 4; i++ {
			rx[i] *= scale
			ry[i] *= scale
		}
	}

	p.MoveTo(x+rx[0], y)

	if clockwise {
		p.LineTo(x+w-rx[1], y)
		if rx[1] > 0 || ry[1] > 0 {
			radius := math.Max(rx[1], ry[1])
			p.Arc(x+w-rx[1], y+ry[1], radius, -math.Pi/2, 0, false)
		}
		p.LineTo(x+w, y+h-ry[2])
		if rx[2] > 0 || ry[2] > 0 {
			radius := math.Max(rx[2], ry[2])
			p.Arc(x+w-rx[2], y+h-ry[2], radius, 0, math.Pi/2, false)
		}
		p.LineTo(x+rx[3], y+h)
		if rx[3] > 0 || ry[3] > 0 {
			radius := math.Max(rx[3], ry[3])
			p.Arc(x+rx[3], y+h-ry[3], radius, math.Pi/2, math.Pi, false)
		}
		p.LineTo(x, y+ry[0])
		if rx[0] > 0 || ry[0] > 0 {
			radius := math.Max(rx[0], ry[0])
			p.Arc(x+rx[0], y+ry[0], radius, math.Pi, 3*math.Pi/2, false)
		}
	} else {
		p.LineTo(x, y+ry[0])
		if rx[0] > 0 || ry[0] > 0 {
			radius := math.Max(rx[0], ry[0])
			p.Arc(x+rx[0], y+ry[0], radius, 3*math.Pi/2, math.Pi, false)
		}
		p.LineTo(x+w-rx[1], y)
		if rx[1] > 0 || ry[1] > 0 {
			radius := math.Max(rx[1], ry[1])
			p.Arc(x+w-rx[1], y+ry[1], radius, math.Pi, -math.Pi/2, false)
		}
		p.LineTo(x+w, y+h-ry[2])
		if rx[2] > 0 || ry[2] > 0 {
			radius := math.Max(rx[2], ry[2])
			p.Arc(x+w-rx[2], y+h-ry[2], radius, -math.Pi/2, 0, false)
		}
		p.LineTo(x+rx[3], y+h)
		if rx[3] > 0 || ry[3] > 0 {
			radius := math.Max(rx[3], ry[3])
			p.Arc(x+rx[3], y+h-ry[3], radius, 0, math.Pi/2, false)
		}
		p.LineTo(x, y+ry[0])
	}

	p.ClosePath()
}

// Arc adds points to the subpath such that the arc described by the circumference of the circle
// described by the arguments, starting at the given start angle and ending at the given end angle,
// going in the given direction (defaulting to clockwise), is added to the path, connected to
// the previous point by a straight line.
func (p *Path2D) Arc(x, y, radius, startAngle, endAngle float64, counterclockwise bool) {
	p.Ellipse(x, y, radius, radius, 0, startAngle, endAngle, counterclockwise)
}

// Ellipse adds points to the subpath such that the arc described by the circumference of the ellipse
// described by the arguments, starting at the given start angle and ending at the given end angle,
// going in the given direction (defaulting to clockwise), is added to the path, connected to
// the previous point by a straight line.
func (p *Path2D) Ellipse(x, y, radiusX, radiusY, rotation, startAngle, endAngle float64, counterclockwise bool) {
	sweepAngle := endAngle - startAngle
	if counterclockwise {
		sweepAngle = -sweepAngle
	}

	oval, ok := path.NewRectFromXYWH(float32(x-radiusX), float32(y-radiusY), float32(radiusX*2), float32(radiusY*2))
	if !ok {
		return
	}

	startDeg := float32(startAngle * 180.0 / math.Pi)
	sweepDeg := float32(sweepAngle * 180.0 / math.Pi)

	if math.Abs(float64(sweepDeg)) >= 360 {
		if sweepDeg > 0 {
			sweepDeg = 359.9999
		} else {
			sweepDeg = -359.9999
		}
	}

	if rotation == 0 {
		p.builder.ArcToOval(oval, startDeg, sweepDeg)
		return
	}

	pb := path.NewPathBuilder()
	pb.ArcToOval(oval, startDeg, sweepDeg)

	if pb.IsEmpty() {
		return
	}

	rotationDeg := float32(rotation * 180.0 / math.Pi)
	ts := path.NewTransformDefault().
		PostRotateAt(rotationDeg, float32(x), float32(y))

	arcPath := pb.Finish()
	if arcPath == nil {
		return
	}

	p.builder.PushPathWithTransform(arcPath, ts)
}
