package examples

import (
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"math"
	"os"
	"sort"
	"testing"

	"github.com/lumifloat/tinyskia"
)

var (
	base00 = color.RGBA{0x00, 0x00, 0x00, 0x00}
	base01 = color.RGBA{0x7c, 0x8c, 0x98, 0xff}
	base02 = color.RGBA{0x8e, 0x9e, 0xaa, 0xff}
	base03 = color.RGBA{0xae, 0xbe, 0xca, 0xff}
	base04 = color.RGBA{0x57, 0x57, 0x57, 0xff}
	base05 = color.RGBA{0x97, 0x97, 0x97, 0xff}
	base06 = color.RGBA{0x33, 0x33, 0x33, 0xff}
	base07 = color.RGBA{0x07, 0x07, 0x07, 0xff}
	base08 = color.RGBA{0x32, 0x32, 0x32, 0xff}
	base09 = color.RGBA{0x1d, 0x1d, 0x1d, 0xff}
	base10 = color.RGBA{0x04, 0x04, 0x04, 0xff}
	base11 = color.RGBA{0x66, 0x76, 0x82, 0xff}
	base12 = color.RGBA{0x39, 0x39, 0x39, 0xff}
	base13 = color.RGBA{0x94, 0xa4, 0xb0, 0xff}
)

var base0 = [][]color.RGBA{
	{base01, base01, base01, base02, base02, base01, base01, base01, base01, base01, base01, base01, base03, base03, base01, base01, base01, base01, base01, base02, base02, base01, base01, base01},
	{base02, base04, base05, base05, base04, base05, base06, base05, base07, base04, base04, base06, base06, base06, base07, base06, base04, base04, base05, base05, base04, base05, base06, base03},
	{base02, base06, base06, base04, base04, base04, base04, base06, base06, base06, base06, base06, base04, base04, base05, base05, base04, base06, base06, base04, base04, base04, base04, base01},
	{base01, base06, base04, base08, base09, base09, base09, base09, base10, base04, base04, base08, base08, base08, base08, base09, base09, base09, base08, base08, base09, base06, base06, base01},
	{base01, base06, base06, base08, base08, base09, base09, base09, base09, base08, base04, base04, base10, base04, base08, base09, base09, base09, base09, base08, base08, base06, base06, base01},
	{base11, base04, base04, base09, base08, base08, base04, base04, base04, base04, base04, base04, base09, base04, base04, base04, base10, base08, base08, base09, base08, base04, base05, base03},
	{base02, base05, base04, base08, base08, base10, base09, base09, base10, base09, base09, base09, base04, base09, base08, base08, base08, base04, base08, base08, base08, base07, base06, base01},
	{base01, base06, base06, base08, base08, base04, base08, base08, base08, base10, base08, base04, base08, base08, base09, base09, base09, base09, base09, base08, base08, base05, base04, base02},
	{base03, base05, base05, base10, base08, base04, base04, base09, base09, base09, base09, base09, base09, base09, base04, base04, base04, base04, base04, base10, base08, base05, base05, base01},
	{base01, base06, base07, base08, base09, base09, base08, base08, base08, base09, base09, base09, base08, base10, base08, base08, base09, base09, base10, base08, base09, base06, base04, base02},
	{base03, base05, base05, base10, base04, base04, base08, base12, base08, base09, base04, base04, base08, base08, base08, base08, base04, base04, base04, base10, base04, base05, base04, base13},
	{base01, base06, base04, base08, base08, base08, base08, base09, base09, base09, base09, base09, base09, base10, base09, base09, base09, base09, base08, base08, base08, base04, base04, base01},
	{base03, base05, base05, base04, base08, base04, base04, base10, base04, base09, base08, base08, base08, base04, base04, base04, base04, base04, base04, base04, base08, base05, base05, base11},
	{base01, base06, base05, base08, base08, base09, base09, base09, base09, base09, base09, base08, base10, base08, base08, base09, base09, base09, base04, base08, base08, base06, base06, base01},
	{base03, base04, base04, base08, base08, base08, base04, base09, base09, base04, base04, base09, base10, base09, base04, base04, base04, base08, base08, base08, base08, base04, base05, base01},
	{base02, base05, base05, base10, base10, base08, base08, base08, base09, base08, base08, base08, base04, base04, base08, base08, base08, base04, base04, base10, base10, base04, base04, base02},
	{base01, base06, base06, base08, base08, base09, base09, base09, base09, base09, base09, base09, base04, base04, base09, base09, base09, base09, base09, base08, base08, base06, base06, base01},
	{base02, base04, base05, base04, base08, base04, base09, base04, base10, base08, base08, base09, base09, base09, base10, base09, base08, base08, base04, base04, base08, base05, base06, base03},
	{base02, base06, base06, base08, base08, base08, base08, base09, base09, base09, base09, base09, base08, base08, base04, base04, base08, base09, base09, base08, base08, base04, base04, base01},
	{base01, base06, base06, base08, base08, base09, base09, base09, base09, base08, base04, base04, base10, base04, base08, base09, base09, base09, base09, base08, base08, base06, base06, base01},
	{base11, base04, base04, base09, base08, base08, base04, base04, base04, base04, base04, base04, base09, base04, base04, base04, base10, base08, base08, base09, base08, base04, base05, base03},
	{base02, base05, base04, base04, base04, base07, base06, base06, base07, base06, base06, base06, base05, base06, base04, base04, base04, base05, base04, base04, base04, base07, base06, base01},
	{base01, base06, base06, base04, base04, base05, base04, base04, base04, base07, base04, base05, base04, base04, base06, base06, base06, base06, base06, base04, base04, base05, base04, base02},
	{base03, base03, base03, base11, base02, base03, base03, base01, base01, base01, base01, base01, base01, base01, base03, base03, base03, base03, base03, base11, base02, base03, base03, base01},
}

var base1 = [][]color.RGBA{
	{base02, base04, base05, base04, base04, base06, base06, base06, base06, base06, base04, base04, base05, base04, base07, base04, base04, base04, base05, base04, base04, base06, base06, base01},
	{base01, base06, base07, base04, base04, base04, base05, base04, base04, base04, base06, base05, base06, base06, base06, base07, base06, base06, base07, base04, base04, base04, base05, base02},
	{base03, base05, base04, base04, base06, base04, base04, base07, base05, base05, base05, base06, base05, base05, base05, base05, base05, base05, base04, base04, base06, base04, base04, base11},
	{base01, base06, base06, base04, base04, base06, base06, base06, base06, base04, base05, base07, base05, base05, base04, base06, base06, base06, base06, base04, base04, base06, base06, base01},
	{base01, base06, base06, base06, base04, base04, base06, base06, base06, base04, base04, base04, base04, base05, base05, base07, base06, base06, base06, base06, base04, base04, base06, base01},
	{base01, base04, base04, base04, base04, base06, base06, base04, base05, base05, base04, base04, base06, base06, base06, base06, base06, base04, base04, base04, base04, base06, base06, base02},
	{base03, base06, base05, base04, base05, base05, base04, base04, base06, base07, base06, base06, base06, base04, base04, base07, base05, base06, base05, base04, base05, base05, base04, base02},
	{base01, base01, base01, base02, base02, base01, base01, base01, base01, base01, base03, base03, base01, base01, base01, base01, base01, base01, base01, base02, base02, base01, base01, base01},
}

var (
	crystal00 = color.RGBA{0x00, 0x00, 0x00, 0x00}
	crystal01 = color.RGBA{0xe0, 0xf2, 0xff, 0xff}
	crystal02 = color.RGBA{0xe0, 0xf2, 0xff, 0xff}
	crystal03 = color.RGBA{0xc6, 0xe7, 0xff, 0xff}
	crystal04 = color.RGBA{0x93, 0xd2, 0xff, 0xff}
)

var crystal = [][]color.RGBA{
	{crystal01, crystal02, crystal02, crystal03, crystal03, crystal03, crystal04, crystal04, crystal04, crystal04, crystal03, crystal03, crystal03, crystal02, crystal02, crystal01},
	{crystal02, crystal03, crystal00, crystal03, crystal00, crystal03, crystal00, crystal03, crystal00, crystal03, crystal00, crystal03, crystal00, crystal03, crystal00, crystal02},
	{crystal02, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal03, crystal02},
	{crystal03, crystal03, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal03},
	{crystal03, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal03, crystal03},
	{crystal03, crystal03, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal03},
	{crystal04, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal03, crystal04},
	{crystal04, crystal03, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal04},
	{crystal04, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal03, crystal04},
	{crystal04, crystal03, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal04},
	{crystal03, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal03, crystal03},
	{crystal03, crystal03, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal03},
	{crystal03, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal03, crystal03},
	{crystal02, crystal03, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal00, crystal02},
	{crystal02, crystal00, crystal03, crystal00, crystal03, crystal00, crystal03, crystal00, crystal03, crystal00, crystal03, crystal00, crystal03, crystal00, crystal03, crystal02},
	{crystal01, crystal02, crystal02, crystal03, crystal03, crystal03, crystal04, crystal04, crystal04, crystal04, crystal03, crystal03, crystal03, crystal02, crystal02, crystal01},
}

var (
	core00 = color.RGBA{0x00, 0x00, 0x00, 0x00}
	core01 = color.RGBA{0xc7, 0x43, 0xff, 0xff}
	core02 = color.RGBA{0xd1, 0x61, 0xff, 0xff}
	core03 = color.RGBA{0xd5, 0x6f, 0xff, 0xff}
	core04 = color.RGBA{0xe3, 0x81, 0xff, 0xff}
	core05 = color.RGBA{0xd8, 0x77, 0xff, 0xff}
	core06 = color.RGBA{0xf7, 0x6a, 0xd2, 0xff}
	core07 = color.RGBA{0xc6, 0x11, 0x7c, 0xff}
	core08 = color.RGBA{0xe4, 0x83, 0xff, 0xff}
	core09 = color.RGBA{0xe4, 0x22, 0x81, 0xff}
)
var core0 = [][]color.RGBA{
	{core01, core02, core02, core02, core02, core03, core03, core02, core02, core02, core02, core03, core03, core02, core02, core01},
	{core02, core04, core05, core05, core04, core05, core04, core05, core05, core05, core05, core04, core05, core05, core04, core02},
	{core02, core04, core04, core04, core05, core05, core04, core04, core04, core04, core04, core05, core04, core04, core05, core03},
	{core03, core05, core04, core04, core06, core06, core06, core04, core05, core05, core04, core06, core04, core04, core04, core03},
	{core03, core05, core05, core06, core07, core07, core07, core06, core05, core05, core06, core07, core06, core04, core04, core02},
	{core03, core04, core05, core04, core06, core06, core07, core06, core08, core06, core06, core09, core06, core04, core04, core02},
	{core02, core04, core04, core04, core05, core06, core07, core06, core06, core07, core09, core09, core06, core05, core05, core03},
	{core02, core04, core04, core05, core06, core09, core06, core04, core06, core06, core06, core09, core06, core05, core04, core03},
	{core02, core04, core05, core06, core09, core06, core05, core05, core06, core07, core09, core07, core06, core04, core04, core03},
	{core02, core05, core04, core05, core06, core09, core06, core05, core06, core09, core06, core06, core04, core05, core05, core03},
	{core03, core04, core04, core04, core04, core06, core09, core06, core06, core09, core06, core06, core04, core04, core05, core02},
	{core02, core05, core04, core04, core04, core05, core06, core09, core09, core09, core09, core07, core06, core04, core04, core02},
	{core03, core05, core04, core05, core05, core05, core04, core06, core06, core06, core06, core06, core04, core04, core04, core03},
	{core03, core04, core04, core04, core04, core05, core04, core05, core05, core04, core04, core05, core04, core04, core05, core03},
	{core02, core05, core04, core04, core05, core05, core05, core04, core05, core05, core04, core04, core05, core05, core05, core02},
	{core01, core02, core02, core03, core03, core02, core02, core02, core02, core03, core03, core02, core02, core02, core02, core01},
}

var core1 = [][]color.RGBA{
	{core01, core02, core02, core02, core02, core03, core03, core02, core02, core02, core02, core03, core03, core02, core02, core01},
	{core02, core04, core05, core05, core04, core05, core04, core05, core05, core05, core05, core04, core05, core05, core04, core02},
	{core02, core04, core04, core04, core05, core05, core04, core04, core04, core04, core04, core05, core04, core04, core05, core03},
	{core03, core05, core04, core04, core06, core06, core06, core04, core05, core05, core04, core05, core04, core04, core04, core03},
	{core03, core05, core05, core06, core07, core09, core09, core06, core05, core05, core05, core05, core05, core04, core04, core02},
	{core03, core04, core05, core06, core09, core06, core09, core06, core08, core05, core08, core08, core05, core04, core04, core02},
	{core02, core04, core04, core06, core09, core06, core07, core06, core04, core05, core05, core05, core04, core05, core05, core03},
	{core02, core04, core04, core06, core07, core06, core06, core06, core06, core06, core06, core06, core04, core05, core04, core03},
	{core02, core04, core05, core06, core07, core07, core09, core09, core09, core09, core07, core07, core06, core04, core04, core03},
	{core02, core05, core04, core05, core06, core06, core06, core09, core06, core06, core06, core06, core04, core05, core05, core03},
	{core03, core04, core04, core04, core04, core05, core06, core09, core06, core05, core05, core05, core04, core04, core05, core02},
	{core02, core05, core04, core04, core04, core05, core06, core07, core06, core05, core05, core05, core04, core04, core04, core02},
	{core03, core05, core04, core05, core05, core05, core04, core06, core04, core04, core05, core05, core04, core04, core04, core03},
	{core03, core04, core04, core04, core04, core05, core04, core05, core05, core04, core04, core05, core04, core04, core05, core03},
	{core02, core05, core04, core04, core05, core05, core05, core04, core05, core05, core04, core04, core05, core05, core05, core02},
	{core01, core02, core02, core03, core03, core02, core02, core02, core02, core03, core03, core02, core02, core02, core02, core01},
}

var core2 = [][]color.RGBA{
	{core01, core02, core02, core02, core02, core03, core03, core02, core02, core02, core02, core03, core03, core02, core02, core01},
	{core02, core04, core05, core05, core04, core05, core04, core05, core05, core05, core05, core04, core05, core05, core04, core02},
	{core02, core04, core04, core04, core05, core05, core04, core04, core04, core04, core04, core05, core04, core04, core05, core03},
	{core03, core05, core04, core04, core04, core04, core04, core06, core06, core06, core06, core06, core04, core04, core04, core03},
	{core03, core05, core05, core05, core04, core05, core06, core07, core09, core09, core09, core07, core06, core04, core04, core02},
	{core03, core04, core05, core04, core05, core06, core07, core06, core06, core09, core06, core06, core05, core04, core04, core02},
	{core02, core04, core04, core04, core06, core09, core06, core04, core06, core09, core06, core05, core04, core05, core05, core03},
	{core02, core04, core04, core06, core07, core06, core04, core04, core04, core06, core09, core06, core04, core05, core04, core03},
	{core02, core04, core05, core06, core07, core07, core09, core09, core09, core09, core09, core06, core04, core04, core04, core03},
	{core02, core05, core04, core05, core06, core06, core09, core06, core06, core06, core07, core06, core04, core05, core05, core03},
	{core03, core04, core04, core04, core04, core06, core07, core06, core04, core05, core06, core07, core06, core04, core05, core02},
	{core02, core05, core04, core04, core06, core07, core06, core04, core04, core05, core06, core07, core06, core04, core04, core02},
	{core03, core05, core04, core05, core05, core06, core04, core04, core04, core04, core05, core06, core04, core04, core04, core03},
	{core03, core04, core04, core04, core04, core05, core04, core05, core05, core04, core04, core05, core04, core04, core05, core03},
	{core02, core05, core04, core04, core05, core05, core05, core04, core05, core05, core04, core04, core05, core05, core05, core02},
	{core01, core02, core02, core03, core03, core02, core02, core02, core02, core03, core03, core02, core02, core02, core02, core01},
}

var core3 = [][]color.RGBA{
	{core01, core02, core02, core02, core02, core03, core03, core02, core02, core02, core02, core03, core03, core02, core02, core01},
	{core02, core04, core05, core05, core04, core05, core04, core05, core05, core05, core05, core04, core05, core05, core04, core02},
	{core02, core04, core04, core04, core05, core05, core04, core04, core04, core04, core04, core05, core04, core04, core05, core03},
	{core03, core05, core04, core04, core06, core04, core04, core04, core05, core06, core04, core05, core04, core04, core04, core03},
	{core03, core05, core05, core06, core09, core06, core05, core05, core06, core07, core06, core05, core05, core04, core04, core02},
	{core03, core04, core05, core04, core06, core09, core06, core06, core06, core09, core06, core08, core05, core04, core04, core02},
	{core02, core04, core04, core04, core05, core06, core07, core07, core09, core09, core06, core05, core04, core05, core05, core03},
	{core02, core04, core04, core05, core05, core05, core06, core06, core06, core09, core06, core04, core04, core05, core04, core03},
	{core02, core04, core05, core05, core06, core06, core06, core06, core06, core09, core06, core05, core04, core04, core04, core03},
	{core02, core05, core04, core06, core07, core09, core09, core09, core09, core07, core06, core05, core04, core05, core05, core03},
	{core03, core04, core04, core04, core06, core06, core09, core06, core06, core06, core07, core06, core04, core04, core05, core02},
	{core02, core05, core04, core04, core04, core06, core07, core06, core04, core05, core06, core07, core06, core04, core04, core02},
	{core03, core05, core04, core05, core05, core05, core06, core04, core04, core04, core05, core06, core04, core04, core04, core03},
	{core03, core04, core04, core04, core04, core05, core04, core05, core05, core04, core04, core05, core04, core04, core05, core03},
	{core02, core05, core04, core04, core05, core05, core05, core04, core05, core05, core04, core04, core05, core05, core05, core02},
	{core01, core02, core02, core03, core03, core02, core02, core02, core02, core03, core03, core02, core02, core02, core02, core01},
}

var core4 = [][]color.RGBA{
	{core01, core02, core02, core02, core02, core03, core03, core02, core02, core02, core02, core03, core03, core02, core02, core01},
	{core02, core04, core05, core05, core04, core05, core04, core05, core05, core05, core05, core04, core05, core05, core04, core02},
	{core02, core04, core04, core04, core05, core05, core04, core04, core04, core04, core04, core05, core04, core04, core05, core03},
	{core03, core05, core04, core04, core06, core04, core04, core06, core06, core06, core04, core05, core04, core04, core04, core03},
	{core03, core05, core05, core06, core09, core06, core06, core07, core07, core09, core06, core05, core05, core04, core04, core02},
	{core03, core04, core05, core06, core07, core06, core06, core07, core06, core06, core09, core06, core05, core04, core04, core02},
	{core02, core04, core04, core04, core06, core05, core06, core09, core06, core05, core06, core07, core06, core05, core05, core03},
	{core02, core04, core04, core05, core05, core06, core09, core07, core06, core05, core05, core06, core04, core05, core04, core03},
	{core02, core04, core05, core05, core06, core09, core06, core06, core04, core05, core05, core06, core04, core04, core04, core03},
	{core02, core05, core04, core06, core09, core06, core05, core05, core05, core05, core06, core07, core06, core05, core05, core03},
	{core03, core04, core04, core06, core07, core06, core06, core06, core06, core06, core06, core09, core06, core04, core05, core02},
	{core02, core05, core04, core06, core07, core07, core09, core09, core09, core09, core07, core07, core06, core04, core04, core02},
	{core03, core05, core04, core05, core06, core06, core06, core06, core06, core06, core06, core06, core04, core04, core04, core03},
	{core03, core04, core04, core04, core04, core05, core04, core05, core05, core04, core04, core05, core04, core04, core05, core03},
	{core02, core05, core04, core04, core05, core05, core05, core04, core05, core05, core04, core04, core05, core05, core05, core02},
	{core01, core02, core02, core03, core03, core02, core02, core02, core02, core03, core03, core02, core02, core02, core02, core01},
}

var core5 = [][]color.RGBA{
	{core01, core02, core02, core02, core02, core03, core03, core02, core02, core02, core02, core03, core03, core02, core02, core01},
	{core02, core04, core05, core05, core04, core05, core04, core05, core05, core05, core05, core04, core05, core05, core04, core02},
	{core02, core04, core04, core04, core05, core05, core04, core04, core04, core04, core04, core05, core04, core04, core05, core03},
	{core03, core05, core04, core04, core06, core06, core06, core06, core06, core06, core04, core05, core04, core04, core04, core03},
	{core03, core05, core05, core06, core07, core07, core09, core09, core09, core07, core06, core05, core05, core04, core04, core02},
	{core03, core04, core05, core04, core06, core06, core06, core06, core06, core09, core06, core06, core05, core04, core04, core02},
	{core02, core04, core04, core04, core05, core06, core07, core07, core09, core09, core07, core07, core06, core05, core05, core03},
	{core02, core04, core04, core05, core05, core06, core07, core06, core06, core09, core06, core07, core06, core05, core04, core03},
	{core02, core04, core05, core05, core05, core06, core09, core06, core06, core09, core06, core09, core06, core04, core04, core03},
	{core02, core05, core04, core05, core04, core06, core07, core06, core06, core07, core06, core09, core06, core05, core05, core03},
	{core03, core04, core04, core04, core06, core06, core06, core06, core06, core06, core06, core07, core06, core04, core05, core02},
	{core02, core05, core04, core06, core07, core09, core09, core09, core09, core09, core07, core07, core06, core04, core04, core02},
	{core03, core05, core04, core05, core06, core06, core06, core06, core06, core06, core06, core06, core04, core04, core04, core03},
	{core03, core04, core04, core04, core04, core05, core04, core05, core05, core04, core04, core05, core04, core04, core05, core03},
	{core02, core05, core04, core04, core05, core05, core05, core04, core05, core05, core04, core04, core05, core05, core05, core02},
	{core01, core02, core02, core03, core03, core02, core02, core02, core02, core03, core03, core02, core02, core02, core02, core01},
}

// Camera represents a camera with position and orientation
type Camera struct {
	Scale float64
	Pitch float64
	Yaw   float64
}

// TestEndCrystal renders a Minecraft End Crystal animation and saves as GIF
func TestEndCrystal(t *testing.T) {
	width := 250
	height := 450
	numFrames := 100

	camera := Camera{
		Scale: 200,
		Pitch: 0.48,
		Yaw:   0.785,
	}

	var palette = color.Palette([]color.Color{
		base00, base01, base02, base03, base04, base05, base06, base07, base08, base09,
		base10, base11, base12, base13, crystal01, crystal02, crystal03, crystal04,
		core00, core01, core02, core03, core04, core05, core06, core07, core08, core09,
	})

	var frames []*image.Paletted
	var delays []int
	var disposal []byte

	c := tinyskia.NewCanvas(width, height)
	ctx := c.GetContext()
	ctx.SetAntiAlias(false)

	for frame := 0; frame < numFrames; frame++ {
		ctx.ResetTransform()
		ctx.ClearRect(0, 0, float64(width), float64(height))

		centerX := float64(width) / 2
		centerY := float64(height) / 2

		ctx.SetTransformWithMatrix(tinyskia.NewMatrixFromTranslate(centerX, centerY))

		// Animation parameters - designed for 100-frame loop
		angle := float64(frame) * (2 * math.Pi / 100)
		bobbing := math.Sin(float64(frame)*(6*math.Pi/100)) * 0.25

		var collector = &PixelQuadsCollector{
			Camera:     camera,
			PixelQuads: []PixelQuad{},
		}

		// Base
		collector.Cube(Point3D{0.75, 0.25, 0.75}, Point3D{0, 0.68, 0}, Point3D{0, 0, 0}, [6][][]color.RGBA{base0, base1, base1, base1, base1, base0})

		// Core
		collector.Cube(Point3D{0.35, 0.35, 0.35}, Point3D{0, -0.35 + bobbing, 0}, Point3D{-angle, -angle, angle}, [6][][]color.RGBA{core0, core1, core2, core3, core4, core5})

		// Crystal 0
		collector.Cube(Point3D{0.45, 0.45, 0.45}, Point3D{0, -0.35 + bobbing, 0}, Point3D{-0.75, angle, -angle}, [6][][]color.RGBA{crystal, crystal, crystal, crystal, crystal, crystal})

		// Crystal 1
		collector.Cube(Point3D{0.5, 0.5, 0.5}, Point3D{0, -0.35 + bobbing, 0}, Point3D{0.25, -angle, -0.25}, [6][][]color.RGBA{crystal, crystal, crystal, crystal, crystal, crystal})

		sort.SliceStable(collector.PixelQuads, func(i, j int) bool {
			return collector.PixelQuads[i].Depth > collector.PixelQuads[j].Depth
		})

		// Render sorted quads
		for _, quad := range collector.PixelQuads {
			ctx.SetFillStyleSolidColor(quad.Color)
			ctx.BeginPath()
			ctx.MoveTo(quad.Points[0].X, quad.Points[0].Y)
			ctx.LineTo(quad.Points[1].X, quad.Points[1].Y)
			ctx.LineTo(quad.Points[2].X, quad.Points[2].Y)
			ctx.LineTo(quad.Points[3].X, quad.Points[3].Y)
			ctx.ClosePath()
			ctx.Fill()
		}

		src := c.Image()
		bounds := src.Bounds()
		dst := image.NewPaletted(bounds, palette)
		draw.Draw(dst, bounds, src, bounds.Min, draw.Src)
		frames = append(frames, dst)
		delays = append(delays, 5)
		disposal = append(disposal, gif.DisposalBackground)
		t.Logf("Frame %d: Added %d quads", frame, len(collector.PixelQuads))
	}

	fi, err := os.Create("end_crystal.gif")
	if err != nil {
		t.Fatalf("Failed to create GIF file: %v", err)
	}
	defer fi.Close()
	if err := gif.EncodeAll(fi, &gif.GIF{
		Image:    frames,
		Delay:    delays,
		Disposal: disposal,
	}); err != nil {
		t.Fatalf("Failed to encode GIF: %v", err)
	}
}

// Point3D represents a 3D point in space
type Point3D struct {
	X, Y, Z float64
}

// Rotate rotates a 3D point around the X, Y, and Z axes
func (p Point3D) Rotate(ax, ay, az float64) Point3D {
	x, y, z := p.X, p.Y, p.Z
	if ax != 0 {
		s, c := math.Sin(ax), math.Cos(ax)
		y1 := y*c - z*s
		z1 := y*s + z*c
		y, z = y1, z1
	}
	if ay != 0 {
		s, c := math.Sin(ay), math.Cos(ay)
		x1 := x*c + z*s
		z1 := -x*s + z*c
		x, z = x1, z1
	}
	if az != 0 {
		s, c := math.Sin(az), math.Cos(az)
		x1 := x*c - y*s
		y1 := x*s + y*c
		x, y = x1, y1
	}
	return Point3D{X: x, Y: y, Z: z}
}

// Move moves a 3D point by the specified offsets
func (p Point3D) Move(dx, dy, dz float64) Point3D {
	return Point3D{X: p.X + dx, Y: p.Y + dy, Z: p.Z + dz}
}

// ToCamera projects a 3D point to camera space
func (p Point3D) ToCamera(scale, pitch, yaw float64) Point3D {
	sY, cY := math.Sin(yaw), math.Cos(yaw)
	xCam := p.X*cY + p.Z*sY
	zCam := -p.X*sY + p.Z*cY
	sX, cX := math.Sin(pitch), math.Cos(pitch)
	yCam := p.Y*cX - zCam*sX
	finalZ := p.Y*sX + zCam*cX
	return Point3D{X: xCam * scale, Y: yCam * scale, Z: finalZ}
}

// PixelQuad represents a projected pixel quad with depth information
type PixelQuad struct {
	Color  color.Color
	Points [4]Point3D
	Depth  float64
}

// Face represents a 3D face with texture mapping
type Face struct {
	Name    string
	Normal  Point3D
	UDir    Point3D
	VDir    Point3D
	Start   Point3D
	Texture [][]color.RGBA
}

// PixelQuadsCollector collects pixel quads for rendering
type PixelQuadsCollector struct {
	Camera     Camera
	PixelQuads []PixelQuad
}

// Cube adds a 3D cube with the given size, position, rotation, and texture to the collector
func (qc *PixelQuadsCollector) Cube(size Point3D, position Point3D, rotate Point3D, texture [6][][]color.RGBA) {

	sX := size.X / 2
	sY := size.Y / 2
	sZ := size.Z / 2

	faces := []Face{
		{Name: "top", Normal: Point3D{0, -1, 0}, UDir: Point3D{1, 0, 0}, VDir: Point3D{0, 0, 1}, Start: Point3D{-sX, -sY, -sZ}, Texture: texture[0]},
		{Name: "front", Normal: Point3D{0, 0, 1}, UDir: Point3D{1, 0, 0}, VDir: Point3D{0, 1, 0}, Start: Point3D{-sX, -sY, sZ}, Texture: texture[1]},
		{Name: "right", Normal: Point3D{1, 0, 0}, UDir: Point3D{0, 0, -1}, VDir: Point3D{0, 1, 0}, Start: Point3D{sX, -sY, sZ}, Texture: texture[2]},
		{Name: "left", Normal: Point3D{-1, 0, 0}, UDir: Point3D{0, 0, 1}, VDir: Point3D{0, 1, 0}, Start: Point3D{-sX, -sY, -sZ}, Texture: texture[3]},
		{Name: "back", Normal: Point3D{0, 0, -1}, UDir: Point3D{-1, 0, 0}, VDir: Point3D{0, 1, 0}, Start: Point3D{sX, -sY, -sZ}, Texture: texture[4]},
		{Name: "bottom", Normal: Point3D{0, 1, 0}, UDir: Point3D{1, 0, 0}, VDir: Point3D{0, 0, -1}, Start: Point3D{-sX, sY, sZ}, Texture: texture[5]},
	}

	for _, face := range faces {
		texture := face.Texture
		if texture == nil {
			continue
		}

		uResolution := len(texture[0])
		vResolution := len(texture)

		uStep := math.Abs(face.UDir.X*size.X+face.UDir.Y*size.Y+face.UDir.Z*size.Z) / float64(uResolution)
		vStep := math.Abs(face.VDir.X*size.X+face.VDir.Y*size.Y+face.VDir.Z*size.Z) / float64(vResolution)

		for u := 0; u < uResolution; u++ {
			for v := 0; v < vResolution; v++ {
				if v >= len(texture) || u >= len(texture[v]) {
					continue
				}

				var quad PixelQuad

				quad.Color = texture[v][u]

				// Calculate pixel position
				px := face.Start.X + float64(u)*face.UDir.X*uStep + float64(v)*face.VDir.X*vStep
				py := face.Start.Y + float64(u)*face.UDir.Y*uStep + float64(v)*face.VDir.Y*vStep
				pz := face.Start.Z + float64(u)*face.UDir.Z*uStep + float64(v)*face.VDir.Z*vStep

				// Create quad vertices
				p1 := Point3D{X: px, Y: py, Z: pz}
				p2 := Point3D{X: px + face.UDir.X*uStep, Y: py + face.UDir.Y*uStep, Z: pz + face.UDir.Z*uStep}
				p3 := Point3D{X: px + face.UDir.X*uStep + face.VDir.X*vStep, Y: py + face.UDir.Y*uStep + face.VDir.Y*vStep, Z: pz + face.UDir.Z*uStep + face.VDir.Z*vStep}
				p4 := Point3D{X: px + face.VDir.X*vStep, Y: py + face.VDir.Y*vStep, Z: pz + face.VDir.Z*vStep}

				quad.Points = [4]Point3D{p1, p2, p3, p4}
				for i := range quad.Points {
					quad.Points[i] = quad.Points[i].Rotate(rotate.X, rotate.Y, rotate.Z)
					quad.Points[i] = quad.Points[i].Move(position.X, position.Y, position.Z)
					quad.Points[i] = quad.Points[i].ToCamera(qc.Camera.Scale, qc.Camera.Pitch, qc.Camera.Yaw)
					quad.Depth += quad.Points[i].Z
				}
				quad.Depth = quad.Depth / 4

				qc.PixelQuads = append(qc.PixelQuads, quad)
			}
		}
	}
}
