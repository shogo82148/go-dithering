package main

import (
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"log"
	"os"

	"github.com/shogo82148/go-dithering"
)

func main() {
	reader, err := os.Open("david.png")
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	img, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}

	p := []color.Color{color.Black, color.White}
	paletted := image.NewPaletted(img.Bounds(), p)

	out := func(d draw.Drawer, name string) {
		d.Draw(paletted, img.Bounds(), img, image.ZP)
		f, err := os.Create(name)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		png.Encode(f, paletted)
	}

	out(draw.Src, "src.png")
	out(draw.FloydSteinberg, "floyd-steinberg.png")
	out(dithering.FalseFloydSteinberg, "false-floyd-steinberg.png")
	out(dithering.JarvisJudiceNinke, "jarvis-judice-ninke.png")
	out(dithering.Stucki, "stucki.png")
	out(dithering.Atkinson, "atkinson.png")
	out(dithering.Burkes, "burkes.png")
}
