package main

import (
	"Matrix/vector"
	"core"
	"image"
	"image/color"
	"image/png"
	"os"
)

func main() {
	var (
		w int = 1024
		h int = 1024
	// 	err error
	// 	f   *os.File
	// 	m   *core.ModelPly
	// 	i   *image.RGBA
	// 	r   *core.Render
	)
	// m, mm := core.Parse("..\\model\\cube\\", "cube.obj")
	m, mm := core.ParseNoMtl("..\\model\\spot\\", "spot_triangulated_good.obj")
	// m, mm := core.ParseNoMtl("E:\\MyPro\\C\\tinyrenderer\\obj\\", "floor.obj")
	// m, mm := core.ParseNoMtl("..\\model\\rock\\", "rock.obj")
	eye := vector.NewVector3D(00, 200, 250)
	r, _ := core.NewRenderObj(m)
	imag := image.NewRGBA(image.Rect(0, 0, w, h))
	core.MakeLight(
		*vector.NewVector3D(150, 150, 150),
		*vector.NewVector3D(900, 900, 900),
		*eye)
	r.MakeModelMat(200, 200, 200)
	r.MakeRoation(45, 45, 40)
	r.MakeTranslation(0, 00, 00)
	r.MakeViewMat(
		eye,
		vector.NewVector3D(0, 0, 0),
		vector.NewVector3D(0, 1, 0))
	r.MakePerspectMat(90, 1, -1, -1000)
	r.MakeViewPort(float64(w), float64(h))
	for i := 0; i < imag.Rect.Dx(); i++ {
		for j := 0; j < imag.Rect.Dy(); j++ {
			imag.Set(i, j, color.Black)
		}
	}
	r.Render(imag, mm)
	f, _ := os.OpenFile("cube.png", os.O_CREATE, 0777)
	png.Encode(f, imag)
}
