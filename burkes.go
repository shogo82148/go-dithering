package dithering

import (
	"image"
	"image/color"
	"image/draw"
)

// Burkes
var Burkes draw.Drawer = burkes{}

type burkes struct{}

func (burkes) Draw(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point) {
	quantError0 := make([][3]int32, r.Dx()+4)
	quantError1 := make([][3]int32, r.Dx()+4)

	out := color.RGBA64{A: 0xffff}
	for y := 0; y != r.Dy(); y++ {
		for x := 0; x != r.Dx(); x++ {
			sr, sg, sb, _ := src.At(sp.X+x, sp.Y+y).RGBA()
			er, eg, eb := int32(sr), int32(sg), int32(sb)

			er = clamp(er + quantError0[x+2][0]/32)
			eg = clamp(eg + quantError0[x+2][1]/32)
			eb = clamp(eb + quantError0[x+2][2]/32)

			out.R = uint16(er)
			out.G = uint16(eg)
			out.B = uint16(eb)
			dst.Set(r.Min.X+x, r.Min.Y+y, &out)

			sr, sg, sb, _ = dst.At(r.Min.X+x, r.Min.Y+y).RGBA()
			er -= int32(sr)
			eg -= int32(sg)
			eb -= int32(sb)

			quantError0[x+3][0] += er * 8
			quantError0[x+3][1] += eg * 8
			quantError0[x+3][2] += eb * 8
			quantError0[x+4][0] += er * 4
			quantError0[x+4][1] += eg * 4
			quantError0[x+4][2] += eb * 4
			quantError1[x+0][0] += er * 2
			quantError1[x+0][1] += eg * 2
			quantError1[x+0][2] += eb * 2
			quantError1[x+1][0] += er * 4
			quantError1[x+1][1] += eg * 4
			quantError1[x+1][2] += eb * 4
			quantError1[x+2][0] += er * 8
			quantError1[x+2][1] += eg * 8
			quantError1[x+2][2] += eb * 8
			quantError1[x+3][0] += er * 4
			quantError1[x+3][1] += eg * 4
			quantError1[x+3][2] += eb * 4
			quantError1[x+4][0] += er * 2
			quantError1[x+4][1] += eg * 2
			quantError1[x+4][2] += eb * 2
		}

		// Recycle the quantization error buffers.
		quantError0, quantError1 = quantError1, quantError0
		for i := range quantError1 {
			quantError1[i] = [3]int32{}
		}
	}
}
