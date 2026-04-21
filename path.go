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

type path2d struct {
	data    *path.Path
	builder *path.PathBuilder
}

func NewPath2D() *path2d {
	return &path2d{builder: path.NewPathBuilder()}
}

// MoveTo creates a new subpath with the given point.
func (p *path2d) MoveTo(x, y float64) {
	p1 := path.Point{X: float32(x), Y: float32(y)}
	p.builder.MoveTo(p1.X, p1.Y)
}

// ClosePath Marks the current subpath as closed,
// and starts a new subpath with a point the same as the start and end of the newly closed subpath.
func (p *path2d) ClosePath() {
	p.builder.Close()
}

// LineTo adds the given point to the current subpath, connected to the previous one by a straight line.
func (p *path2d) LineTo(x, y float64) {
	if p.data != nil {
		p.data = nil
	}
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
func (p *path2d) QuadraticCurveTo(x1, y1, x2, y2 float64) {
	if p.data != nil {
		p.data = nil
	}
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
func (p *path2d) BezierCurveTo(cp1x, cp1y, cp2x, cp2y, x, y float64) {
	if p.data != nil {
		p.data = nil
	}
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
func (p *path2d) ArcTo(x1, y1, x2, y2, radius float64) {
	if p.data != nil {
		p.data = nil
	}
	currentPt, hasCurrent := p.builder.LastPoint()
	if !hasCurrent {
		return
	}

	if radius < 0 {
		return
	}

	p0 := currentPt
	p1 := path.Point{X: float32(x1), Y: float32(y1)}
	p2 := path.Point{X: float32(x2), Y: float32(y2)}

	dx1 := float64(p1.X - p0.X)
	dy1 := float64(p1.Y - p0.Y)
	dx2 := float64(p2.X - p1.X)
	dy2 := float64(p2.Y - p1.Y)

	len1 := math.Sqrt(dx1*dx1 + dy1*dy1)
	len2 := math.Sqrt(dx2*dx2 + dy2*dy2)

	if len1 < 1e-10 || len2 < 1e-10 {
		p.LineTo(x1, y1)
		return
	}

	ux1 := dx1 / len1
	uy1 := dy1 / len1
	ux2 := dx2 / len2
	uy2 := dy2 / len2

	cosAngle := ux1*ux2 + uy1*uy2
	if cosAngle > 1.0 {
		cosAngle = 1.0
	} else if cosAngle < -1.0 {
		cosAngle = -1.0
	}
	angle := math.Acos(cosAngle)

	tanHalfAngle := math.Tan(angle / 2.0)
	dist := radius / tanHalfAngle

	if dist > len1 {
		dist = len1
	}
	if dist > len2 {
		dist = len2
	}

	t1x := float64(p1.X) - ux1*dist
	t1y := float64(p1.Y) - uy1*dist
	t2x := float64(p1.X) + ux2*dist
	t2y := float64(p1.Y) + uy2*dist

	midX := (t1x + t2x) / 2.0
	midY := (t1y + t2y) / 2.0

	dirX := midX - float64(p1.X)
	dirY := midY - float64(p1.Y)
	dirLen := math.Sqrt(dirX*dirX + dirY*dirY)

	if dirLen < 1e-10 {
		p.LineTo(x1, y1)
		return
	}

	nx := dirX / dirLen
	ny := dirY / dirLen

	cpDist := radius / math.Sin(angle/2.0)
	cpx := float64(p1.X) + nx*cpDist
	cpy := float64(p1.Y) + ny*cpDist

	p.LineTo(t1x, t1y)
	p.QuadraticCurveTo(cpx, cpy, t2x, t2y)
}

// Rect adds a new closed subpath to the path, representing the given rectangle.
func (p *path2d) Rect(x, y, w, h float64) {
	if p.data != nil {
		p.data = nil
	}
	p.MoveTo(x, y)
	p.LineTo(x+w, y)
	p.LineTo(x+w, y+h)
	p.LineTo(x, y+h)
	p.ClosePath()
}

// RoundRect adds a new closed subpath to the path representing the given rounded rectangle.
func (p *path2d) RoundRect(x, y, w, h float64, radii []float64) {
	if p.data != nil {
		p.data = nil
	}
	if math.IsInf(x, 0) || math.IsNaN(x) ||
		math.IsInf(y, 0) || math.IsNaN(y) ||
		math.IsInf(w, 0) || math.IsNaN(w) ||
		math.IsInf(h, 0) || math.IsNaN(h) {
		return
	}

	if len(radii) == 0 || len(radii) > 4 {
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
func (p *path2d) Arc(x, y, radius, startAngle, endAngle float64, counterclockwise bool) {
	if p.data != nil {
		p.data = nil
	}
	if counterclockwise {
		startAngle, endAngle = endAngle, startAngle
	}

	const n = 16
	for i := 0; i < n; i++ {
		a1 := startAngle + (endAngle-startAngle)*float64(i)/n
		a2 := startAngle + (endAngle-startAngle)*float64(i+1)/n

		x0 := x + radius*math.Cos(a1)
		y0 := y + radius*math.Sin(a1)
		x1 := x + radius*math.Cos((a1+a2)/2)
		y1 := y + radius*math.Sin((a1+a2)/2)
		x2 := x + radius*math.Cos(a2)
		y2 := y + radius*math.Sin(a2)

		cx := 2*x1 - x0/2 - x2/2
		cy := 2*y1 - y0/2 - y2/2

		if i == 0 {
			_, hasCurrent := p.builder.LastPoint()
			if hasCurrent {
				p.LineTo(x0, y0)
			} else {
				p.MoveTo(x0, y0)
			}
		}
		p.QuadraticCurveTo(cx, cy, x2, y2)
	}
}

// Ellipse adds points to the subpath such that the arc described by the circumference of the ellipse
// described by the arguments, starting at the given start angle and ending at the given end angle,
// going in the given direction (defaulting to clockwise), is added to the path, connected to
// the previous point by a straight line.
func (p *path2d) Ellipse(x, y, radiusX, radiusY, rotation, startAngle, endAngle float64, counterclockwise bool) {
	if p.data != nil {
		p.data = nil
	}
	if counterclockwise {
		startAngle, endAngle = endAngle, startAngle
	}

	const n = 16
	cosRot := math.Cos(rotation)
	sinRot := math.Sin(rotation)

	for i := 0; i < n; i++ {
		a1 := startAngle + (endAngle-startAngle)*float64(i)/n
		a2 := startAngle + (endAngle-startAngle)*float64(i+1)/n

		x0u := radiusX * math.Cos(a1)
		y0u := radiusY * math.Sin(a1)
		x1u := radiusX * math.Cos((a1+a2)/2)
		y1u := radiusY * math.Sin((a1+a2)/2)
		x2u := radiusX * math.Cos(a2)
		y2u := radiusY * math.Sin(a2)

		x0 := x + x0u*cosRot - y0u*sinRot
		y0 := y + x0u*sinRot + y0u*cosRot
		x1 := x + x1u*cosRot - y1u*sinRot
		y1 := y + x1u*sinRot + y1u*cosRot
		x2 := x + x2u*cosRot - y2u*sinRot
		y2 := y + x2u*sinRot + y2u*cosRot

		cx := 2*x1 - x0/2 - x2/2
		cy := 2*y1 - y0/2 - y2/2

		if i == 0 {
			_, hasCurrent := p.builder.LastPoint()
			if hasCurrent {
				p.LineTo(x0, y0)
			} else {
				p.MoveTo(x0, y0)
			}
		}
		p.QuadraticCurveTo(cx, cy, x2, y2)
	}
}
