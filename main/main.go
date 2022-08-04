package main

import (
	"Matrix/vector"
	"core"
	"image"
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
	m, mm := core.Parse("..\\model\\cube\\", "cube.obj")
	// m, mm := core.ParseNoMtl("..\\model\\spot\\", "spot_control_mesh.obj")
	// m, mm := core.ParseNoMtl("E:\\MyPro\\C\\tinyrenderer\\obj\\", "floor.obj")
	// m, mm := core.ParseNoMtl("..\\model\\rock\\", "rock.obj")
	r, _ := core.NewRenderObj(m)
	imag := image.NewRGBA(image.Rect(0, 0, w, h))
	r.MakeModelMat(50, 50, 50)
	r.MakeRoation(140, 0, 45)
	r.MakeViewMat(
		vector.NewVector3D(0, 0, 200),
		vector.NewVector3D(0, 0, 0),
		vector.NewVector3D(0, 1, 0))
	r.MakePerspectMat(90, 1, -1, -1000)
	r.MakeViewPort(float64(w), float64(h))
	core.SetLight(vector.NewVector3D(1, 1, 1))
	r.Render(imag, mm)
	f, _ := os.OpenFile("cube.png", os.O_CREATE, 0777)
	png.Encode(f, imag)
}
