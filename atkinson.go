package dithering

import (
	"image"
	"image/color"
	"image/draw"
)

// Atkinson
var Atkinson draw.Drawer = atkinson{}

type atkinson struct{}

func (atkinson) Draw(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point) {
	quantError0 := make([][3]int32, r.Dx()+3)
	quantError1 := make([][3]int32, r.Dx()+3)
	quantError2 := make([][3]int32, r.Dx()+3)

	out := color.RGBA64{A: 0xffff}
	for y := 0; y != r.Dy(); y++ {
		for x := 0; x != r.Dx(); x++ {
			sr, sg, sb, _ := src.At(sp.X+x, sp.Y+y).RGBA()
			er, eg, eb := int32(sr), int32(sg), int32(sb)

			er = clamp(er + quantError0[x+1][0]/8)
			eg = clamp(eg + quantError0[x+1][1]/8)
			eb = clamp(eb + quantError0[x+1][2]/8)

			out.R = uint16(er)
			out.G = uint16(eg)
			out.B = uint16(eb)
			dst.Set(r.Min.X+x, r.Min.Y+y, &out)

			sr, sg, sb, _ = dst.At(r.Min.X+x, r.Min.Y+y).RGBA()
			er -= int32(sr)
			eg -= int32(sg)
			eb -= int32(sb)

			quantError0[x+2][0] += er
			quantError0[x+2][1] += eg
			quantError0[x+2][2] += eb
			quantError0[x+3][0] += er
			quantError0[x+3][1] += eg
			quantError0[x+3][2] += eb
			quantError1[x+0][0] += er
			quantError1[x+0][1] += eg
			quantError1[x+0][2] += eb
			quantError1[x+1][0] += er
			quantError1[x+1][1] += eg
			quantError1[x+1][2] += eb
			quantError1[x+2][0] += er
			quantError1[x+2][1] += eg
			quantError1[x+2][2] += eb
			quantError2[x+1][0] += er
			quantError2[x+1][1] += eg
			quantError2[x+1][2] += eb
		}

		// Recycle the quantization error buffers.
		quantError0, quantError1, quantError2 = quantError1, quantError2, quantError0
		for i := range quantError2 {
			quantError2[i] = [3]int32{}
		}
	}
}
