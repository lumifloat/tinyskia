// Copyright 2006 The Android Open Source Project
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package pixmap

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io"
	"os"

	"github.com/chewxy/math32"
	"github.com/lumifloat/tinyskia/color"
	"github.com/lumifloat/tinyskia/path"
)

// Number of bytes per pixel.
const BYTES_PER_PIXEL = 4

// A container that owns premultiplied RGBA pixels.
//
// The data is not aligned, therefore width == stride.
type Pixmap struct {
	data []uint8
	size path.IntSize
}

// Allocates a new pixmap.
//
// A pixmap is filled with transparent black by default, aka (0, 0, 0, 0).
//
// Zero size in an error.
func NewPixmap(width, height uint32) *Pixmap {
	size, ok := path.NewIntSize(width, height)
	if !ok {
		return nil
	}
	dataLen, ok := dataLenForSize(size)
	if !ok {
		return nil
	}

	return &Pixmap{
		data: make([]uint8, dataLen),
		size: size,
	}
}

// Creates a new pixmap by taking ownership over an image buffer
// (premultiplied RGBA pixels).
//
// The size needs to match the data provided.
//
// Pixmap's width is limited by i32::MAX/4.
func PixmapFromVec(data []uint8, size path.IntSize) *Pixmap {
	dataLen, ok := dataLenForSize(size)
	if !ok || int(len(data)) != dataLen {
		return nil
	}

	return &Pixmap{data: data, size: size}
}

// Decodes a PNG data into a `Pixmap`.
//
// Only 8-bit images are supported.
// Index PNGs are not supported.
func DecodePNG(r io.Reader) (*Pixmap, error) {
	img, err := png.Decode(r)
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Fastest path: if the image is already *image.RGBA, directly use its data
	// Both tiny-skia and Go's image.RGBA use identical premultiplied RGBA layout,
	// so we can share the memory without any copying
	if rgba, ok := img.(*image.RGBA); ok {
		size, _ := path.NewIntSize(uint32(width), uint32(height))
		pixmap := &Pixmap{
			data: rgba.Pix,
			size: size,
		}
		return pixmap, nil
	}

	// Fallback for other image types: allocate new memory and convert pixel by pixel
	pixmap := NewPixmap(uint32(width), uint32(height))
	if pixmap == nil {
		return nil, fmt.Errorf("failed to create a pixmap")
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c := img.At(x+bounds.Min.X, y+bounds.Min.Y)
			r, g, b, a := c.RGBA()
			idx := (y*width + x) * BYTES_PER_PIXEL
			pixmap.data[idx+0] = uint8(r >> 8)
			pixmap.data[idx+1] = uint8(g >> 8)
			pixmap.data[idx+2] = uint8(b >> 8)
			pixmap.data[idx+3] = uint8(a >> 8)
		}
	}

	return pixmap, nil
}

// Loads a PNG file into a `Pixmap`.
func LoadPNG(file string) (*Pixmap, error) {
	fi, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer fi.Close()
	return DecodePNG(fi)
}

// Encodes pixmap into a PNG data.
func (p *Pixmap) EncodePNG() (io.Reader, error) {
	// Direct conversion: both tiny-skia and Go's image.RGBA use premultiplied alpha
	// with identical memory layout (RGBA order, 4 bytes per pixel).
	img := &image.RGBA{
		Pix:    p.data,
		Stride: int(p.Width()) * BYTES_PER_PIXEL,
		Rect:   image.Rect(0, 0, int(p.Width()), int(p.Height())),
	}

	buf := bytes.NewBuffer(make([]byte, 0))
	return buf, png.Encode(buf, img)
}

// Saves pixmap as a PNG file.
func (p *Pixmap) SavePng(file string) error {
	w, err := p.EncodePNG()
	if err != nil {
		return err
	}
	fi, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	_, err = io.Copy(fi, w)
	return err
}

func (p *Pixmap) Width() uint32  { return p.size.Width() }
func (p *Pixmap) Height() uint32 { return p.size.Height() }

func (p *Pixmap) Size() path.IntSize { return p.size }

// Fills the entire pixmap with a specified color.
func (p *Pixmap) Fill(color color.ColorU8) {
	c := color.Premultiply()
	for i := 0; i < len(p.data); i += BYTES_PER_PIXEL {
		p.data[i+0] = c.Red()
		p.data[i+1] = c.Green()
		p.data[i+2] = c.Blue()
		p.data[i+3] = c.Alpha()
	}
}

// Data returns the internal data as bytes.
func (p *Pixmap) Data() []uint8 { return p.data }

// DataMut returns the mutable internal data as bytes.
func (p *Pixmap) DataMut() []uint8 { return p.data }

// TakeDemultiplied consumes the pixmap and returns the internal data as demultiplied RGBA bytes.
func (p *Pixmap) TakeDemultiplied() []byte {
	// Demultiply alpha.
	for i := 0; i < len(p.data); i += BYTES_PER_PIXEL {
		a := p.data[i+3]
		if a != 0 && a != 255 {
			p.data[i+0] = byte((uint32(p.data[i+0]) * 255) / uint32(a))
			p.data[i+1] = byte((uint32(p.data[i+1]) * 255) / uint32(a))
			p.data[i+2] = byte((uint32(p.data[i+2]) * 255) / uint32(a))
		}
	}
	return p.data
}

// Pixel returns a pixel color at the specified position.
// Returns nil when position is out of bounds.
func (p *Pixmap) Pixel(x, y int) (color.PremultipliedColorU8, bool) {
	if x < 0 || y < 0 || x >= int(p.Width()) || y >= int(p.Height()) {
		return color.PremultipliedColorU8{}, false
	}
	idx := (y*int(p.Width()) + x) * BYTES_PER_PIXEL
	c := color.PremultipliedColorU8FromRGBAUnchecked(
		p.data[idx+0],
		p.data[idx+1],
		p.data[idx+2],
		p.data[idx+3],
	)
	return c, true
}

// Pixels returns a slice of all pixels.
func (p *Pixmap) Pixels() []color.PremultipliedColorU8 {
	pixels := make([]color.PremultipliedColorU8, p.Width()*p.Height())
	for i := 0; i < len(p.data); i += BYTES_PER_PIXEL {
		pixels[i/BYTES_PER_PIXEL] = color.PremultipliedColorU8FromRGBAUnchecked(
			p.data[i+0],
			p.data[i+1],
			p.data[i+2],
			p.data[i+3],
		)
	}
	return pixels
}

func (p *Pixmap) AsSubPixmap(x, y, width, height int) *SubPixmap {
	return &SubPixmap{
		Data:      p.data,
		Size:      p.size,
		RealWidth: int(p.size.Width()),
	}
}

func minRowBytes(size path.IntSize) (int, bool) {
	w := size.Width()
	res := int(w) * BYTES_PER_PIXEL
	if res > math32.MaxUint32 {
		return 0, false
	}
	return res, true
}

func computeDataLen(size path.IntSize, rowBytes int) (int, bool) {
	if size.Height() == 0 {
		return 0, false
	}
	h := size.Height() - 1
	hLen := int(h) * rowBytes
	wLen := int(size.Width()) * BYTES_PER_PIXEL

	return hLen + wLen, true
}

func dataLenForSize(size path.IntSize) (int, bool) {
	rowBytes, ok := minRowBytes(size)
	if !ok {
		return 0, false
	}
	return computeDataLen(size, rowBytes)
}

// SubPixmap represents a subregion of a Pixmap.
type SubPixmap struct {
	Data      []uint8
	Size      path.IntSize
	RealWidth int
}

// NewSubPixmap creates a sub-pixmap from the given pixmap at the specified region.
func NewSubPixmap(p *Pixmap, x, y, width, height int) (*SubPixmap, error) {
	if x < 0 || y < 0 || width <= 0 || height <= 0 {
		return nil, fmt.Errorf("invalid subregion")
	}
	if x+width > int(p.Width()) || y+height > int(p.Height()) {
		return nil, fmt.Errorf("subregion out of bounds")
	}

	rowBytes := int(p.Width()) * BYTES_PER_PIXEL
	offset := y*rowBytes + x*BYTES_PER_PIXEL
	dataLen := height*rowBytes - (int(p.Width())-width)*BYTES_PER_PIXEL

	size, ok := path.NewIntSize(uint32(width), uint32(height))
	if !ok {
		return nil, fmt.Errorf("invalid size")
	}

	return &SubPixmap{
		Data:      p.data[offset : offset+dataLen],
		Size:      size,
		RealWidth: int(p.Width()),
	}, nil
}

// SubPixmap creates a sub-pixmap from the given pixmap at the specified rect.
// Returns nil if the rect is invalid or out of bounds.
func (p *Pixmap) SubPixmap(rect path.IntRect) (*SubPixmap, bool) {
	x := int(rect.Left())
	y := int(rect.Top())
	width := int(rect.Width())
	height := int(rect.Height())

	if x < 0 || y < 0 || width <= 0 || height <= 0 {
		return nil, false
	}
	if x+width > int(p.Width()) || y+height > int(p.Height()) {
		return nil, false
	}

	rowBytes := int(p.Width()) * BYTES_PER_PIXEL
	offset := y*rowBytes + x*BYTES_PER_PIXEL
	dataLen := height*rowBytes - (int(p.Width())-width)*BYTES_PER_PIXEL

	size, ok := path.NewIntSize(uint32(width), uint32(height))
	if !ok {
		return nil, false
	}

	return &SubPixmap{
		Data:      p.data[offset : offset+dataLen],
		Size:      size,
		RealWidth: int(p.Width()),
	}, true
}

// ToRGBA converts the pixmap to an image.RGBA.
// Since both use premultiplied RGBA format, this is a zero-copy operation.
func (p *Pixmap) ToRGBA() *image.RGBA {
	width := int(p.Width())
	height := int(p.Height())

	// Create image.RGBA that shares the same underlying data
	// Both pixmap and image.RGBA use premultiplied RGBA format
	return &image.RGBA{
		Pix:    p.data,
		Stride: width * BYTES_PER_PIXEL,
		Rect:   image.Rect(0, 0, width, height),
	}
}
