package core

import (
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"log"
	"math"
	"os"
	"strings"

	"golang.org/x/image/tiff"
)

type MipMaper interface {
	GetUV(x, y float64)
	GetRGBA(x, y int)
}
type MipMap struct {
	Num        int
	LevelImage []image.Image
}

func NewMipMap(path string) (mipmap *MipMap) {
	var (
		imag  image.Image
		err   error
		level int
		last  image.Image
		// tmep   image.Image
	)
	imag, err = readImage(path)
	mipmap = &MipMap{}
	if err != nil {
		log.Fatalf("err :%v", err)
		return
	}
	level = int(math.Log2(float64(imag.Bounds().Dx())) + 1)
	mipmap.Num = level
	mipmap.LevelImage = make([]image.Image, level)
	mipmap.LevelImage[0] = imag
	for i := 1; i < level; i++ {
		last = mipmap.LevelImage[i-1]
		temp := image.NewRGBA(image.Rect(0, 0, last.Bounds().Dx()/2, last.Bounds().Dy()/2))
		for k := 0; k < temp.Rect.Dx(); k += 2 {
			for h := 0; h < temp.Rect.Dx(); h += 2 {
				r, g, b, a := last.At(k, h).RGBA()
				r2, g2, b2, a2 := last.At(k+1, h).RGBA()
				r3, g3, b3, a3 := last.At(k+1, h+1).RGBA()
				r4, g4, b4, a4 := last.At(k, h+1).RGBA()
				// fmt.Print(r)
				r = (r&255 + r2&255 + r3&255 + r4&255) / 4
				g = (g&255 + g2&255 + g3&255 + g4&255) / 4
				b = (b&255 + b2&255 + b3&255 + b4&255) / 4
				a = (a&255 + a2&255 + a3&255 + a4&255) / 4

				temp.Set(k/2, h/2, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
			}
		}
		mipmap.LevelImage[i] = temp
	}
	return
}
func readImage(path string) (images image.Image, err error) {
	var (
		file *os.File
	)
	file, err = os.Open(path)
	if strings.HasSuffix(path, "png") {
		images, err = png.Decode(file)
	} else if strings.HasSuffix(path, "jpeg") {
		images, err = jpeg.Decode(file)

	} else if strings.HasSuffix(path, "tif") {
		images, err = tiff.Decode(file)
	}
	return
}
