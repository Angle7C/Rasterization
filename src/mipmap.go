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

	"github.com/ftrvxmtrx/tga"
	"golang.org/x/image/tiff"
)

type MipMaper interface {
	UV(l, u, v float64)
	GetRGBA()
}
type MipMap struct {
	Num        int
	w          int
	h          int
	LevelImage []image.Image
}

func (mipmap *MipMap) UV(u, v, u1, v1, u2, v2 float64) color.RGBA {
	ut, vt := u, v
	u = u * float64(mipmap.w)
	v = v * float64(mipmap.h)
	u1 = u1 * float64(mipmap.w)
	v1 = v1 * float64(mipmap.h)
	u2 = u2 * float64(mipmap.w)
	v2 = v2 * float64(mipmap.h)
	l := math.Max(math.Sqrt(math.Pow(u1-u, 2)+math.Pow(v1-v, 2)), math.Sqrt(math.Pow(u2-u, 2)+math.Pow(v2-v, 2)))
	l = math.Log2(l)
	l = math.Max(l, 0)
	l = math.Min(l, float64(len(mipmap.LevelImage)))
	floor := math.Floor(l)
	ceil := math.Ceil(l)
	c1 := mipmap.LevelImage[int(floor)].At(int(ut*float64(mipmap.LevelImage[int(floor)].Bounds().Dx()-1)), int(vt*float64(mipmap.LevelImage[int(floor)].Bounds().Dx()-1)))
	c2 := mipmap.LevelImage[int(ceil)].At(int(ut*float64(mipmap.LevelImage[int(ceil)].Bounds().Dx()-1)), int(vt*float64(mipmap.LevelImage[int(ceil)].Bounds().Dx()-1)))
	r1, g1, b1, a1 := c1.RGBA()
	r1 = r1 & 255
	g1 = g1 & 255
	b1 = b1 & 255
	a1 = a1 & 255
	r2, g2, b2, a2 := c2.RGBA()
	r2 = r2 & 255
	g2 = g2 & 255
	b2 = b2 & 255
	a2 = a2 & 255
	r := r1 + r2*(uint32(ceil-l))
	g := g1 + g2*(uint32(ceil-l))
	b := b1 + b2*(uint32(ceil-l))
	a := a1 + a2*(uint32(ceil-l))
	return color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
}
func (mipmap *MipMap) NearUV(u, v float64) color.RGBA {
	u = u * float64(mipmap.w-1)
	v = v * float64(mipmap.h-1)
	r, g, b, a := mipmap.LevelImage[0].At(int(u), int(v)).RGBA()
	return color.RGBA{uint8(r & 255), uint8(g & 255), uint8(b & 255), uint8(a & 255)}
}
func (mipmap *MipMap) BilinearUV(u, v float64) color.RGBA {
	u = u * float64(mipmap.w-1)
	v = v * float64(mipmap.h-1)

	r, g, b, a := mipmap.LevelImage[0].At(int(u), int(v)).RGBA()
	return color.RGBA{uint8(r & 255), uint8(g & 255), uint8(b & 255), uint8(a & 255)}
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
	mipmap.w = imag.Bounds().Dx()
	mipmap.h = imag.Bounds().Dy()
	level = int(math.Log2(float64(imag.Bounds().Dx())) + 1)
	mipmap.Num = level
	mipmap.LevelImage = make([]image.Image, level)
	mipmap.LevelImage[0] = imag
	for i := 1; i < level; i++ {
		last = mipmap.LevelImage[i-1]
		temp := image.NewRGBA(image.Rect(0, 0, last.Bounds().Dx()/2, last.Bounds().Dy()/2))
		for k := 0; k < last.Bounds().Size().X; k += 2 {
			for h := 0; h < last.Bounds().Size().Y; h += 2 {
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
	} else if strings.HasSuffix(path, "tga") {
		images, err = tga.Decode(file)
	} else if strings.HasSuffix(path, "jpg") {
		images, err = jpeg.Decode(file)
	}
	return
}
