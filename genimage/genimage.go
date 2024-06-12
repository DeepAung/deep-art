package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math/rand"
	"os"
)

func main() {
	width := flag.Int("width", 400, "width")
	height := flag.Int("height", 300, "height")
	n := flag.Int("n", 1, "times")
	flag.Parse()

	imagesDir := "../static/storage/genimages"

	err := os.RemoveAll(imagesDir + "/")
	if err != nil {
		log.Fatalf("os.RemoveAll: %v", err)
	}
	err = os.MkdirAll(imagesDir+"/", os.ModePerm)
	if err != nil {
		log.Fatalf("os.MkdirAll: %v", err)
	}

	for i := range *n {
		r := uint8(rand.Intn(256))
		g := uint8(rand.Intn(256))
		b := uint8(rand.Intn(256))
		c := color.RGBA{R: r, G: g, B: b, A: 0xff}

		img := image.NewRGBA(image.Rectangle{
			Min: image.Pt(0, 0),
			Max: image.Pt(*width, *height),
		})

		for x := 0; x < *width; x++ {
			for y := 0; y < *height; y++ {
				img.Set(x, y, c)
			}
		}

		f, err := os.Create(fmt.Sprintf("%s/%05d.png", imagesDir, i+1))
		if err != nil {
			log.Fatalf("os.Create: %v", err)
		}

		png.Encode(f, img)
	}
}
