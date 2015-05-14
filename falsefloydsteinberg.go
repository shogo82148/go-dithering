package dithering

import (
	"image"
	"image/color"
	"image/draw"
)

// False Floyd-Steinberg algorithm
var FalseFloydSteinberg draw.Drawer = falseFloydSteinberg{}

type falseFloydSteinberg struct{}

func (falseFloydSteinberg) Draw(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point) {
	quantErrorNext := make([][3]int32, r.Dx()+1)

	out := color.RGBA64{A: 0xffff}
	for y := 0; y != r.Dy(); y++ {
		quantError := [3]int32{}
		quantErrorNext[0] = [3]int32{}
		for x := 0; x != r.Dx(); x++ {
			sr, sg, sb, _ := src.At(sp.X+x, sp.Y+y).RGBA()
			er, eg, eb := int32(sr), int32(sg), int32(sb)

			er = clamp(er + (quantErrorNext[x][0]+quantError[0])/16)
			eg = clamp(eg + (quantErrorNext[x][1]+quantError[1])/16)
			eb = clamp(eb + (quantErrorNext[x][2]+quantError[2])/16)

			out.R = uint16(er)
			out.G = uint16(eg)
			out.B = uint16(eb)
			dst.Set(r.Min.X+x, r.Min.Y+y, &out)

			sr, sg, sb, _ = dst.At(r.Min.X+x, r.Min.Y+y).RGBA()
			er -= int32(sr)
			eg -= int32(sg)
			eb -= int32(sb)
			quantError[0] = er * 3
			quantError[1] = eg * 3
			quantError[2] = eb * 3
			quantErrorNext[x+0][0] += er * 3
			quantErrorNext[x+0][1] += eg * 3
			quantErrorNext[x+0][2] += eb * 3
			quantErrorNext[x+1][0] = er * 2
			quantErrorNext[x+1][1] = eg * 2
			quantErrorNext[x+1][2] = eb * 2
		}
	}
}
