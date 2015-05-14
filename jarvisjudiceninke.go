package dithering

import (
	"image"
	"image/color"
	"image/draw"
)

// Jarvis, Judice, and Ninke algorithm
var JarvisJudiceNinke draw.Drawer = jarvisJudiceNinke{}

type jarvisJudiceNinke struct{}

func (jarvisJudiceNinke) Draw(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point) {
	quantError0 := make([][3]int32, r.Dx()+4)
	quantError1 := make([][3]int32, r.Dx()+4)
	quantError2 := make([][3]int32, r.Dx()+4)

	out := color.RGBA64{A: 0xffff}
	for y := 0; y != r.Dy(); y++ {
		for x := 0; x != r.Dx(); x++ {
			sr, sg, sb, _ := src.At(sp.X+x, sp.Y+y).RGBA()
			er, eg, eb := int32(sr), int32(sg), int32(sb)

			er = clamp(er + quantError0[x+2][0]/48)
			eg = clamp(eg + quantError0[x+2][1]/48)
			eb = clamp(eb + quantError0[x+2][2]/48)

			out.R = uint16(er)
			out.G = uint16(eg)
			out.B = uint16(eb)
			dst.Set(r.Min.X+x, r.Min.Y+y, &out)

			sr, sg, sb, _ = dst.At(r.Min.X+x, r.Min.Y+y).RGBA()
			er -= int32(sr)
			eg -= int32(sg)
			eb -= int32(sb)

			quantError0[x+3][0] += er * 7
			quantError0[x+3][1] += eg * 7
			quantError0[x+3][2] += eb * 7
			quantError0[x+4][0] += er * 5
			quantError0[x+4][1] += eg * 5
			quantError0[x+4][2] += eb * 5
			quantError1[x+0][0] += er * 3
			quantError1[x+0][1] += eg * 3
			quantError1[x+0][2] += eb * 3
			quantError1[x+1][0] += er * 5
			quantError1[x+1][1] += eg * 5
			quantError1[x+1][2] += eb * 5
			quantError1[x+2][0] += er * 7
			quantError1[x+2][1] += eg * 7
			quantError1[x+2][2] += eb * 7
			quantError1[x+3][0] += er * 5
			quantError1[x+3][1] += eg * 5
			quantError1[x+3][2] += eb * 5
			quantError1[x+4][0] += er * 3
			quantError1[x+4][1] += eg * 3
			quantError1[x+4][2] += eb * 3
			quantError2[x+0][0] += er * 1
			quantError2[x+0][1] += eg * 1
			quantError2[x+0][2] += eb * 1
			quantError2[x+1][0] += er * 3
			quantError2[x+1][1] += eg * 3
			quantError2[x+1][2] += eb * 3
			quantError2[x+2][0] += er * 5
			quantError2[x+2][1] += eg * 5
			quantError2[x+2][2] += eb * 5
			quantError2[x+3][0] += er * 3
			quantError2[x+3][1] += eg * 3
			quantError2[x+3][2] += eb * 3
			quantError2[x+4][0] += er * 1
			quantError2[x+4][1] += eg * 1
			quantError2[x+4][2] += eb * 1
		}

		// Recycle the quantization error buffers.
		quantError0, quantError1, quantError2 = quantError1, quantError2, quantError0
		for i := range quantError2 {
			quantError2[i] = [3]int32{}
		}
	}
}
